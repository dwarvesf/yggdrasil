package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

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

	go func() {
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			panic(err)
		}

		client, err := consul.NewClient(&consul.Config{
			Address: fmt.Sprintf("consul:8500"),
		})
		if err != nil {
			panic(err)
		}
		agent := client.Agent()

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
		r := kafka.NewReader(kafka.ReaderConfig{
			// ! move address to env
			Brokers: []string{"10.5.0.5:9092"},
			Topic:   "email",
		})

		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				break
			}
			// TODO: handle message logic here
			fmt.Println("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		}

		r.Close()
	}()

	logger.Log("exit", <-errs)
}
