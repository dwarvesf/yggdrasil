package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-kit/kit/log"
	consul "github.com/hashicorp/consul/api"
	"github.com/segmentio/kafka-go"
	validator "gopkg.in/validator.v2"

	"github.com/dwarvesf/yggdrasil/sms/model"
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
		Address: fmt.Sprintf("consul:8500"),
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
		kafkaAddr, kafkaPort, err := toolkit.GetServiceAddress(consulClient, "kafka")
		if err != nil {
			panic(err)
		}

		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{fmt.Sprintf("%v:%v", kafkaAddr, kafkaPort)},
			Topic:   "email",
		})

		defer r.Close()

		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				logger.Log("error", err.Error())
				// TODO: should break or continue if cannot read msg from queue
				break
			}
			if string(m.Value) == "" {
				continue
			}

			var req model.Request

			if err = json.Unmarshal(m.Value, &req); err != nil {
				logger.Log("error", err.Error())
				continue
			}

			if err := validator.Validate; err != nil {
				logger.Log("error", err)
				continue
			}

			// do stuff
		}
	}()

	logger.Log("exit", <-errs)
}
