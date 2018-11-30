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

	"github.com/dwarvesf/yggdrasil/logger"
	"github.com/dwarvesf/yggdrasil/services/email/model"
	email "github.com/dwarvesf/yggdrasil/services/email/service"
	mailgun "github.com/dwarvesf/yggdrasil/services/email/service/mailgun"
	sendgrid "github.com/dwarvesf/yggdrasil/services/email/service/sendgrid"
	"github.com/dwarvesf/yggdrasil/toolkit"
	"github.com/dwarvesf/yggdrasil/toolkit/queue"
	"github.com/dwarvesf/yggdrasil/toolkit/queue/kafka"
)

func main() {
	logger := logger.NewLogger()

	errs := make(chan error)
	go func() {
		logger.Info("starting email worker")
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
			logger.Error("unable to get port %s", err.Error())
			panic(err)
		}
		logger.Info("registering %s to consul", svcName)

		if err := toolkit.RegisterService(consulClient, svcName, port); err != nil {
			logger.Error("unable to register to consul %s", err.Error())
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

			if err := sendEmail(req, consulClient); err != nil {
				logger.Info("sending email")
				message, err := toolkit.CreateRetryMessage("email", req.Payload, req.Retry)
				if err != nil {
					logger.Error("unable to send an email %s", err.Error())
					continue
				}
				w.Write("email", message)
				logger.Info("retry payload")
			}
		}
	}()

	logger.Error("exit", <-errs)
}

func sendEmail(r model.Request, consulClient *consul.Client) error {
	var emailer email.Emailer

	switch r.Payload.Provider {
	case "sendgrid":
		v := os.Getenv("SENDGRID")
		if v == "" {
			v, _ = toolkit.GetConsulValueFromKey(consulClient, "sendgrid")
		}
		emailer = sendgrid.New(v)
		return emailer.Send(&r.Payload)

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

		emailer = mailgun.New(value.Domain, value.APIKey, value.PublicKey)
		return emailer.Send(&r.Payload)

	default:
		return errors.New("INVALID_PROVIDER")
	}
}
