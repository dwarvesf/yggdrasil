package scheduler

import (
	"encoding/json"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/dwarvesf/yggdrasil/scheduler/model"
	"github.com/dwarvesf/yggdrasil/scheduler/service"
	"github.com/dwarvesf/yggdrasil/scheduler/validator"
)

// Scheduler to represent scheduler functions: producer and consumer
type Scheduler interface {
	HandleRequests(d time.Duration)
	ListenMessages()
}

// NewScheduler to create a new worker
func NewScheduler(service service.Service, logger log.Logger) Scheduler {
	return &schedulerImpl{
		Service: service,
		Logger:  logger,
	}
}

type schedulerImpl struct {
	Service service.Service
	Logger  log.Logger
}

// HandleRequests will periodically check for request in db and broadcast it to kafka
func (s *schedulerImpl) HandleRequests(d time.Duration) {
	for t := range time.Tick(d) {
		s.Logger.Log("start", t)

		requests, err := s.Service.RequestService.GetRequests()
		if err != nil {
			s.Logger.Log("error", err.Error())
			continue
		}
		if len(requests) == 0 {
			s.Logger.Log("info", "skipping")
			continue
		}

		for _, entity := range requests {
			r, err := entity.ToRequest()
			if err != nil {
				s.Logger.Log("error", err)
				continue
			}

			r.Retry.CurrenyRetry++

			s.logRequest("sending", r)
			err = s.sendResponse(r)
			if err != nil {
				s.Logger.Log("error", err)
				continue
			}
		}

		s.Logger.Log("sending", "success")
		s.deleteRequests(requests)
	}
}

func (s *schedulerImpl) logRequest(name string, r model.Request) {
	s.Logger.Log(name, r.Service, "time", r.Timestamp)
	s.Logger.Log("retry", r.Retry.CurrenyRetry, "maxRetry", r.Retry.MaxRetry, "retryAfter", r.Retry.RetryAfter)
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
		s.Logger.Log("error", err.Error())
	}
}

// ListenMessages will continuosly check for messages in kafka and save to db
func (s *schedulerImpl) ListenMessages() {
	r := s.Service.QueueService.NewReader("scheduler")
	defer r.Close()

	for {
		value, err := r.Read()
		if err != nil {
			s.Logger.Log("error", err.Error())
			continue
		}

		req, err := s.parseRequest(value)
		if err != nil {
			s.Logger.Log("error", err.Error())
			continue
		}

		s.logRequest("saving", req)
		err = s.saveRequest(req)
		if err != nil {
			s.Logger.Log("error", err.Error())
			continue
		}

		s.Logger.Log("saving", "success")
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
