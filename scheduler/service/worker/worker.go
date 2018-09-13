package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	consul "github.com/hashicorp/consul/api"
	kafka "github.com/segmentio/kafka-go"

	"github.com/dwarvesf/yggdrasil/scheduler/model"
	"github.com/dwarvesf/yggdrasil/scheduler/service"
	"github.com/dwarvesf/yggdrasil/scheduler/service/validator"
	"github.com/dwarvesf/yggdrasil/toolkit"
)

// Worker to represent worker functions: producer and consumer
type Worker interface {
	HandleRequests(d time.Duration)
	ListenMessages()
}

// NewWorker to create a new worker
func NewWorker(service service.Service, consulClient *consul.Client, logger log.Logger) Worker {
	return &workerImpl{
		Service:      service,
		ConsulClient: consulClient,
		Logger:       logger,
	}
}

type workerImpl struct {
	Service      service.Service
	ConsulClient *consul.Client
	Logger       log.Logger
}

// HandleRequests will periodically check for request in db and broadcast it to kafka
func (w *workerImpl) HandleRequests(d time.Duration) {
	kafkaAddr, kafkaPort, err := toolkit.GetServiceAddress(w.ConsulClient, "kafka")
	if err != nil {
		panic(err)
	}

	brokers := []string{fmt.Sprintf("%v:%v", kafkaAddr, kafkaPort)}

	for t := range time.Tick(d) {
		w.Logger.Log("start", t)

		requests, err := w.Service.SchedulerService.GetRequests()
		if err != nil {
			w.Logger.Log("error", err.Error())
			continue
		}
		if len(requests) == 0 {
			w.Logger.Log("info", "skipping")
			continue
		}

		for _, entity := range requests {
			r, err := entity.ToRequest()
			if err != nil {
				w.Logger.Log("error", err)
				continue
			}

			w.Logger.Log("sending", r.Service, r.Payload, r.Timestamp)
			err = sendRequest(r, brokers)
			if err != nil {
				w.Logger.Log("error", err)
				continue
			}
		}

		w.deleteRequests(requests)
	}
}

func (w *workerImpl) deleteRequests(requests []model.RequestEntity) {
	var ids []uint

	for _, r := range requests {
		ids = append(ids, r.ID)
	}

	err := w.Service.SchedulerService.DeleteRequests(ids)
	if err != nil {
		w.Logger.Log("error", err.Error())
	}
}

// ListenMessages will continuosly check for messages in kafka and save to db
func (w *workerImpl) ListenMessages() {
	kafkaAddr, kafkaPort, err := toolkit.GetServiceAddress(w.ConsulClient, "kafka")
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
			w.Logger.Log("error", err.Error())
			continue
		}

		req, err := parseRequest(message)
		if err != nil {
			w.Logger.Log("error", err.Error())
			continue
		}

		w.Logger.Log("saving", req.Service, req.Payload, req.Timestamp)
		err = w.saveRequest(req)
		if err != nil {
			w.Logger.Log("error", err.Error())
			continue
		}
	}
}

func sendRequest(r model.Request, brokers []string) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    r.Service,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	value, err := r.ToBytes()
	if err != nil {
		return err
	}

	err = writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte("scheduler"),
			Value: value,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (w *workerImpl) saveRequest(r model.Request) error {
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

	return w.Service.SchedulerService.SaveRequest(entity)
}

func parseRequest(message kafka.Message) (model.Request, error) {
	var req model.Request

	if err := json.Unmarshal(message.Value, &req); err != nil {
		return req, err
	}

	return req, nil
}
