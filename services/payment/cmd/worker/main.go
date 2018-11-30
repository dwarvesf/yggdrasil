package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	consul "github.com/hashicorp/consul/api"
	validator "gopkg.in/validator.v2"

	"github.com/dwarvesf/yggdrasil/logger"
	"github.com/dwarvesf/yggdrasil/services/payment/model"
	payment "github.com/dwarvesf/yggdrasil/services/payment/service"
	"github.com/dwarvesf/yggdrasil/services/payment/service/stripe"
	"github.com/dwarvesf/yggdrasil/toolkit"
	"github.com/dwarvesf/yggdrasil/toolkit/queue"
	"github.com/dwarvesf/yggdrasil/toolkit/queue/kafka"
)

func main() {
	logger := logger.NewLogger()

	errs := make(chan error)
	go func() {
		logger.Info("starting payment worker")
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	consulClient, err := consul.NewClient(&consul.Config{
		Address: fmt.Sprintf("consul-server:8500"),
	})
	if err != nil {
		logger.Error("unable to get port %s", err.Error())
		panic(err)
	}

	svcName := "payment"
	go func() {
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			logger.Error("unable to get port %s", err.Error())
			panic(err)
		}
		logger.Warn("registering %s to consul", svcName)

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
				logger.Error("unable to read from kafka %s", err.Error())
				continue
			}

			var req model.Request
			if err = json.Unmarshal(b, &req); err != nil {
				logger.Info("unable to parse request %s", err.Error())
				continue
			}
			if err := validator.Validate; err != nil {
				logger.Error("Validator error: %s", err)
				continue
			}

			if err := sendPayment(req.Payload, consulClient); err != nil {
				logger.Info("sending payment")
				message, err := toolkit.CreateRetryMessage("payment", req.Payload, req.Retry)
				if err != nil {
					logger.Error("unable to send a payment %s", err.Error())
					continue
				}
				w.Write("payment", message)
				logger.Info("retry payload")
			}
		}
	}()

	logger.Error("exit", <-errs)
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
