package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	consul "github.com/hashicorp/consul/api"
	kafka "github.com/segmentio/kafka-go"

	"github.com/dwarvesf/yggdrasil/scheduler/db"
	"github.com/dwarvesf/yggdrasil/scheduler/model"
	"github.com/dwarvesf/yggdrasil/scheduler/service"
	"github.com/dwarvesf/yggdrasil/scheduler/service/scheduler"
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
		logger.Log("worker", "scheduler")
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
		name := "scheduler"
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			panic(err)
		}
		logger.Log("consul", "registering", "name", name)

		if err := toolkit.RegisterService(consulClient, name, port); err != nil {
			panic(err)
		}
	}()

	pgdb, closeDB := db.New(consulClient)
	db.Migrate(pgdb)
	s := service.Service{
		SchedulerService: scheduler.NewPGService(pgdb),
	}
	defer closeDB()

	// To check db after each X unit of time
	go checkRequests()

	// To run a forever loop to check queue
	go checkMessages(s, consulClient, logger)

	// Test send message to queue, will remove later
	sendMessages(consulClient)

	logger.Log("exit", <-errs)
}

func sendMessages(consulClient *consul.Client) {
	kafkaAddr, kafkaPort, err := toolkit.GetServiceAddress(consulClient, "kafka")
	if err != nil {
		panic(err)
	}

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{fmt.Sprintf("%v:%v", kafkaAddr, kafkaPort)},
		Topic:    "scheduler",
		Balancer: &kafka.LeastBytes{},
	})
	defer w.Close()

	payload := make(map[string]interface{})
	payload["number"] = 100
	payload["text"] = "hello world"

	w.WriteMessages(context.Background(),
		kafka.Message{
			Key: []byte("test"),
			Value: toBytes(model.Request{
				Service:   "test",
				Payload:   payload,
				Timestamp: time.Now().Add(time.Second * 10),
			}),
		},
		kafka.Message{
			Key: []byte("test"),
			Value: toBytes(model.Request{
				Service:   "sms",
				Payload:   payload,
				Timestamp: time.Now().Add(time.Second * -10),
			}),
		},
		kafka.Message{
			Key: []byte("test"),
			Value: toBytes(model.Request{
				Service:   "notification",
				Payload:   payload,
				Timestamp: time.Now().Add(time.Second * 10),
			}),
		},
	)
}

func toBytes(r model.Request) []byte {
	out, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return out
}
