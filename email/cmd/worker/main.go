package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/dwarvesf/yggdrasil/email/model"
	email "github.com/dwarvesf/yggdrasil/email/service"
	"github.com/dwarvesf/yggdrasil/email/service/sendgrid"
	"github.com/go-kit/kit/log"
	consul "github.com/hashicorp/consul/api"
	"github.com/segmentio/kafka-go"
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
	kv := consulClient.KV()

	go func() {
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			panic(err)
		}

		if err != nil {
			panic(err)
		}
		agent := consulClient.Agent()

		name := "email"
		if err := agent.ServiceRegister(&consul.AgentServiceRegistration{
			Name:    name,
			Port:    port,
			Address: os.Getenv("PRIVATE_IP"),
		}); err != nil {
			panic(err)
		}
		logger.Log("consul", "registered", "name", name)
	}()

	go func() {
		var queueAddr []*consul.CatalogService
		queueAddr, _, err := consulClient.Catalog().Service("kafka", "", nil)
		if err != nil {
			panic(err)
		}

		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{fmt.Sprintf("%v:%v", queueAddr[0].ServiceAddress, queueAddr[0].ServicePort)},
			Topic:   "email",
		})
		for {
			func() {
				m, err := r.ReadMessage(context.Background())
				defer r.Close()
				if err != nil {
					logger.Log("error", err.Error())
					// TODO: should break or continue if cannot read msg from queue
					return
				}
				// fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
				if string(m.Value) == "" {
					return
				}

				var req model.Request
				if err = json.Unmarshal(m.Value, &req); err != nil {
					logger.Log("error", err.Error())
					return
				}

				var emailer email.Emailer
				switch req.Type {
				case "sendgrid":
					pair, _, err := kv.Get("sendgrid", nil)
					if err != nil {
						logger.Log("error", err.Error())
						return
					}
					emailer = sendgrid.New(string(pair.Value))
					emailer.Send()
				}
			}()
		}

		r.Close()
	}()

	logger.Log("exit", <-errs)
}
