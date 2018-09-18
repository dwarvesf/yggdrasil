package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/dwarvesf/yggdrasil/toolkit/queue/kafka"

	"github.com/dwarvesf/yggdrasil/toolkit/queue"

	"github.com/go-kit/kit/log"
	consul "github.com/hashicorp/consul/api"

	"github.com/dwarvesf/yggdrasil/email/model"
	email "github.com/dwarvesf/yggdrasil/email/service"
	"github.com/dwarvesf/yggdrasil/email/service/sendgrid"
	"github.com/dwarvesf/yggdrasil/toolkit"
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
		defer r.Close()

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

			var emailer email.Emailer
			switch req.Provider {
			case "sendgrid":
				v := os.Getenv("SENDGRID")
				if v == "" {
					v, _ = toolkit.GetConsulValueFromKey(consulClient, "sendgrid")
				}
				emailer = sendgrid.New(v)
				emailer.Send()
			}
		}
	}()

	logger.Log("exit", <-errs)
}
