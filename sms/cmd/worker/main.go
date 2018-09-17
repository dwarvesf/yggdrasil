package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-kit/kit/log"
	consul "github.com/hashicorp/consul/api"
	validator "gopkg.in/validator.v2"

	"github.com/dwarvesf/yggdrasil/sms/model"
	sms "github.com/dwarvesf/yggdrasil/sms/service"
	"github.com/dwarvesf/yggdrasil/sms/service/twilio"
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
		logger.Log("worker", "sms")
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

	go func() {
		name := "sms"
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			panic(err)
		}
		logger.Log("consul", "registering", "name", name)

		if err := toolkit.RegisterService(consulClient, name, port); err != nil {
			panic(err)
		}
	}()

	go func() {
		var q queue.Queue
		q = kafka.New(consulClient, "sms")
		defer q.Close()

		for {
			b := q.Read()

			var req model.Request

			if err = json.Unmarshal(b, &req); err != nil {
				logger.Log("error", err.Error())
				continue
			}
			if err := validator.Validate; err != nil {
				logger.Log("error", err)
				continue
			}

			var smsClient sms.SMS
			switch req.Provider {
			case "twilio":
				v := os.Getenv("TWILIO")
				if v == "" {
					v, _ = toolkit.GetConsulValueFromKey(consulClient, "twilio")
				}
				value := model.TwilioSecret{}
				if err = json.Unmarshal([]byte(v), &value); err != nil {
					logger.Log("error", err.Error())
				}
				smsClient = twilio.New(value.AppSid, value.AuthToken)
				smsClient.Send(value.AppNumber, req.To, req.Content, value.AppSid)
			}
		}
	}()

	logger.Log("exit", <-errs)
}
