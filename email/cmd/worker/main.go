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

	"github.com/dwarvesf/yggdrasil/email/model"
	email "github.com/dwarvesf/yggdrasil/email/service"
	mailgun "github.com/dwarvesf/yggdrasil/email/service/mailgun"
	sendgrid "github.com/dwarvesf/yggdrasil/email/service/sendgrid"
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
		logger.Log("worker", "email")
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

	svcName := "email"
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

			if err := sendEmail(req, consulClient); err != nil {
				logger.Log("error", err.Error())

				message, err := toolkit.CreateRetryMessage("email", req.Payload, req.Retry)
				if err != nil {
					logger.Log("error", err.Error())
					continue
				}

				w.Write("email", message)
				logger.Log("info", "retry sent")
			}
		}
	}()

	logger.Log("exit", <-errs)
}

func sendEmail(r model.Request, consulClient *consul.Client) error {
	var emailer email.Emailer

	switch r.Payload.Provider {
	case "sendgrid":
		v := os.Getenv("SENDGRID")
		if v == "" {
			v, _ = toolkit.GetConsulValueFromKey(consulClient, "sendgrid")
		}

		sendgrid.New(v, &r.Payload)
		return emailer.Send(v, &r.Payload)

	case "mailgun":
		v := os.Getenv("MAILGUN")
		if v == "" {
			v, _ = toolkit.GetConsulValueFromKey(consulClient, "mailgun")
		}

		value := model.MailgunSecret{}
		err := json.Unmarshal([]byte(v), &value)
		if err != nil {
			return err
		}

		mailgun.New(value.Domain, value.APIKey, value.PublicKey)
		return emailer.Send(value.APIKey, &r.Payload)

	default:
		return errors.New("INVALID_PROVIDER")
	}
}
