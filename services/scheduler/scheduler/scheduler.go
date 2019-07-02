package scheduler

import (
	"encoding/json"
	"time"

	"github.com/dwarvesf/yggdrasil/logger"
	"github.com/dwarvesf/yggdrasil/services/scheduler/model"
	"github.com/dwarvesf/yggdrasil/services/scheduler/service"
	"github.com/dwarvesf/yggdrasil/services/scheduler/validator"
)

// Scheduler to represent scheduler functions: producer and consumer
type Scheduler interface {
	HandleRequests(d time.Duration)
	ListenMessages()
}

// NewScheduler to create a new worker
func NewScheduler(service service.Service, logger logger.LoggingService) Scheduler {
	return &schedulerImpl{
		Service: service,
		Logger:  logger,
	}
}

type schedulerImpl struct {
	Service service.Service
	Logger  logger.LoggingService
}

// HandleRequests will periodically check for request in db and broadcast it to kafka
func (s *schedulerImpl) HandleRequests(d time.Duration) {
	for range time.Tick(d) {
		requests, err := s.Service.RequestService.GetRequests()
		if err != nil {
			s.Logger.Error("[HandleRequests].GetRequests %s", err.Error())
			continue
		}
		if len(requests) == 0 {
			s.Logger.Info("skipping")
			continue
		}

		for _, entity := range requests {
			r, err := entity.ToRequest()
			if err != nil {
				s.Logger.Error("[HandleRequests].ToRequest %s", err.Error())
				continue
			}

			r.Retry.CurrenyRetry++

			s.logRequest("sending", r)
			err = s.sendResponse(r)
			if err != nil {
				s.Logger.Error("[HandleRequests].sendResponse %s", err.Error())
				continue
			}
		}

		s.Logger.Info("sending successfully")
		s.deleteRequests(requests)
	}
}

func (s *schedulerImpl) logRequest(name string, r model.Request) {
	s.Logger.Info("name: %s, time: %s", r.Service, r.Timestamp)
	s.Logger.Info("retry: %s, maxRetry: %s, retryAfter: %s", r.Retry.CurrenyRetry, r.Retry.MaxRetry, r.Retry.RetryAfter)
}

func (s *schedulerImpl) sendResponse(r model.Request) error {
	w := s.Service.QueueService.NewWriter(r.Service)
	defer w.Close()

	response := model.Response{
		Payload: r.Payload,
		Retry:   r.Retry,
	}

	resp, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = w.Write("scheduler", resp)
	if err != nil {
		return err
	}

	return nil
}

func (s *schedulerImpl) deleteRequests(requests []model.RequestEntity) {
	var ids []uint

	for _, r := range requests {
		ids = append(ids, r.ID)
	}

	err := s.Service.RequestService.DeleteRequests(ids)
	if err != nil {
		s.Logger.Error("[deleteRequests] err= %s", err.Error())
	}
}

// ListenMessages will continuosly check for messages in kafka and save to db
func (s *schedulerImpl) ListenMessages() {
	r := s.Service.QueueService.NewReader("scheduler")
	defer r.Close()

	for {
		value, err := r.Read()
		if err != nil {
			s.Logger.Error("[ListenMessage]: error= %s", err.Error())
			continue
		}

		req, err := s.parseRequest(value)
		if err != nil {
			s.Logger.Error("[ListenMessage]: error= %s", err.Error())
			continue
		}

		s.logRequest("saving", req)
		err = s.saveRequest(req)
		if err != nil {
			s.Logger.Error("[ListenMessage]: error= %s", err.Error())
			continue
		}

		s.Logger.Info("saving", "success")
	}
}

func (s *schedulerImpl) saveRequest(r model.Request) error {
	err := validator.ValidateRequest(r)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(r.Payload)
	if err != nil {
		return err
	}

	retry, err := json.Marshal(r.Retry)
	if err != nil {
		return err
	}

	entity := model.RequestEntity{
		Service:   r.Service,
		Payload:   string(payload),
		Timestamp: r.Timestamp,
		Retry:     string(retry),
	}

	return s.Service.RequestService.SaveRequest(entity)
}

func (*schedulerImpl) parseRequest(value []byte) (model.Request, error) {
	var req model.Request

	if err := json.Unmarshal(value, &req); err != nil {
		return req, err
	}

	return req, nil
}
