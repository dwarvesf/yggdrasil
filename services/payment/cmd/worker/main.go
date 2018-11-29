package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-kit/kit/log"
	consul "github.com/hashicorp/consul/api"
	validator "gopkg.in/validator.v2"

	"github.com/dwarvesf/yggdrasil/services/payment/model"
	payment "github.com/dwarvesf/yggdrasil/services/payment/service"
	"github.com/dwarvesf/yggdrasil/services/payment/service/stripe"
	"github.com/dwarvesf/yggdrasil/toolkit"
	"github.com/dwarvesf/yggdrasil/toolkit/queue"
	"github.com/dwarvesf/yggdrasil/toolkit/queue/kafka"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	errs := make(chan error)
	go func() {
		logger.Log("worker", "payment")
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	consulClient, err := consul.NewClient(&consul.Config{
		Address: fmt.Sprintf("consul-server:8500"),
	})
	if err != nil {
		panic(err)
	}

	svcName := "payment"
	go func() {
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			panic(err)
		}
		logger.Log("consul", "registering", "name", svcName)

		if err := toolkit.RegisterService(consulClient, svcName, port); err != nil {
			panic(err)
		}
	}()

	go func() {
		var q queue.Queue
		q = kafka.New(consulClient)
		r := q.NewReader(svcName)
		w := q.NewWriter("scheduler")
		defer r.Close()
		defer w.Close()

		for {
			b, err := r.Read()
			if err != nil {
				logger.Log("error", err.Error())
				continue
			}

			var req model.Request
			if err = json.Unmarshal(b, &req); err != nil {
				logger.Log("error", err.Error())
				continue
			}
			if err := validator.Validate; err != nil {
				logger.Log("error", err)
				continue
			}

			if err := sendPayment(req.Payload, consulClient); err != nil {
				logger.Log("error", err.Error())

				message, err := toolkit.CreateRetryMessage("payment", req.Payload, req.Retry)
				if err != nil {
					logger.Log("error", err.Error())
					continue
				}

				w.Write("payment", message)
				logger.Log("info", "retry sent")
			}
		}
	}()

	logger.Log("exit", <-errs)
}

func sendPayment(p model.Payload, consulClient *consul.Client) error {
	var paymentClient payment.Payer
	switch p.Provider {
	case "stripe":
		v := os.Getenv("STRIPE")
		if v == "" {
			v, _ = toolkit.GetConsulValueFromKey(consulClient, "stripe")
		}

		paymentClient = stripe.New(v)
		return paymentClient.Pay(p)
	default:
		return errors.New("INVALID_PROVIDER")
	}
}
