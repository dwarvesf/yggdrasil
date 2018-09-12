package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-kit/kit/log"
	consul "github.com/hashicorp/consul/api"
	kafka "github.com/segmentio/kafka-go"

	"github.com/dwarvesf/yggdrasil/scheduler/model"
	"github.com/dwarvesf/yggdrasil/scheduler/service"
	"github.com/dwarvesf/yggdrasil/scheduler/service/validator"
	"github.com/dwarvesf/yggdrasil/toolkit"
)

func checkRequests() {
	// TODO: check db every X mins and send message to right channel in queue
}

// go routine to check message queue
func checkMessages(s service.Service, consulClient *consul.Client, logger log.Logger) {
	kafkaAddr, kafkaPort, err := toolkit.GetServiceAddress(consulClient, "kafka")
	if err != nil {
		panic(err)
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{fmt.Sprintf("%v:%v", kafkaAddr, kafkaPort)},
		Topic:   "scheduler",
	})
	defer r.Close()

	for {
		message, err := r.ReadMessage(context.Background())
		if err != nil {
			logger.Log("error", err.Error())
			continue
		}

		req, err := parseRequest(message)
		if err != nil {
			logger.Log("error", err.Error())
			continue
		}

		err = saveRequest(s, req)
		if err != nil {
			logger.Log("error", err.Error())
			continue
		}
	}
}

func saveRequest(s service.Service, r model.Request) error {
	err := validator.ValidateRequest(r)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(r.Payload)
	if err != nil {
		return err
	}

	entity := model.RequestEntity{
		Service:   r.Service,
		Payload:   string(payload),
		Timestamp: r.Timestamp,
	}

	return s.SchedulerService.SaveRequest(entity)
}

func parseRequest(message kafka.Message) (model.Request, error) {
	var req model.Request

	if err := json.Unmarshal(message.Value, &req); err != nil {
		return req, err
	}

	return req, nil
}
