package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

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
		var kafkaAddr []*consul.CatalogService
		kafkaAddr, _, err := consulClient.Catalog().Service("kafka", "", nil)
		if err != nil {
			panic(err)
		}
		type Message struct {
			Type       string            `json:"type"`
			TemplateID string            `json:"template_id"`
			Data       map[string]string `json:"data"`
			Content    string            `json:"content"`
		}

		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{fmt.Sprintf("%v:%v", kafkaAddr[0].ServiceAddress, kafkaAddr[0].ServicePort)},
			Topic:   "email",
		})
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				logger.Log("error", err.Error())
				// TODO: should break or continue if cannot read msg from kafka
				break
			}
			// fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
			if string(m.Value) == "" {
				continue
			}
			var msg Message
			if err = json.Unmarshal(m.Value, &msg); err != nil {
				logger.Log("error", err.Error())
				continue
			}
			var emailer email.Emailer
			switch msg.Type {
			case "sendgrid":
				pair, _, err := kv.Get("sendgrid", nil)
				if err != nil {
					logger.Log("error", err.Error())
					continue
				}
				emailer = sendgrid.New(string(pair.Value))
				emailer.Send()
			}
		}

		r.Close()
	}()

	logger.Log("exit", <-errs)
}
