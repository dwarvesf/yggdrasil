package scheduler

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/dwarvesf/yggdrasil/scheduler/model"
	"github.com/dwarvesf/yggdrasil/scheduler/service"
	"github.com/dwarvesf/yggdrasil/scheduler/service/request"
	"github.com/dwarvesf/yggdrasil/scheduler/util/testutil"
	"github.com/dwarvesf/yggdrasil/toolkit"
)

func TestListenMessagesWhenInvalidShouldNotSave(t *testing.T) {
	// Request
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	requestService := request.NewPGService(pgdb)

	// Queue
	requests := make(chan []byte)
	queueService := &testutil.MockQueueService{
		ReadData: requests,
	}

	// Scheduler
	s := service.Service{
		RequestService: requestService,
		QueueService:   queueService,
	}
	sch := NewScheduler(s, log.NewNopLogger())
	go sch.ListenMessages()

	// Given
	retry := toolkit.RetryMetadata{
		CurrenyRetry: 2,
		MaxRetry:     3,
		RetryAfter:   time.Second,
	}
	r, err := json.Marshal(model.Request{
		Service:   "invalid",
		Timestamp: time.Now().Add(5 * time.Second),
		Retry:     retry,
	})
	if err != nil {
		panic(err)
	}
	requests <- r
	time.Sleep(100 * time.Millisecond)

	// Expect
	var entities []model.RequestEntity
	err = pgdb.Find(&entities).Error
	if err != nil {
		panic(err)
	}
	if len(entities) > 0 {
		t.Errorf("Expect number of saved requests to be zero, but got %v", len(entities))
	}
}

func TestListenMessagesWhenMessageValidShouldSave(t *testing.T) {
	// Request
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	requestService := request.NewPGService(pgdb)

	// Queue
	requests := make(chan []byte)
	queueService := &testutil.MockQueueService{
		ReadData: requests,
	}

	// Scheduler
	s := service.Service{
		RequestService: requestService,
		QueueService:   queueService,
	}
	sch := NewScheduler(s, log.NewNopLogger())
	go sch.ListenMessages()

	// Given
	retry := toolkit.RetryMetadata{
		CurrenyRetry: 2,
		MaxRetry:     3,
		RetryAfter:   time.Second,
	}
	r, err := json.Marshal(model.Request{
		Service:   "email",
		Timestamp: time.Now().Add(5 * time.Second),
		Retry:     retry,
	})
	if err != nil {
		panic(err)
	}
	requests <- r
	time.Sleep(100 * time.Millisecond)

	// Expect
	var entities []model.RequestEntity
	err = pgdb.Find(&entities).Error
	if err != nil {
		panic(err)
	}
	if len(entities) != 1 {
		t.Errorf("Expect number of saved requests to be 1, but got %v", len(entities))
	}
	if entities[0].Service != "email" {
		t.Errorf("Expect request sevice is email, but got %v", entities[0].Service)
	}
}

func TestHandleRequests(t *testing.T) {
	// Request
	pgdb := testutil.GetDB()
	defer pgdb.Close()
	requestService := request.NewPGService(pgdb)

	// Given
	retry := toolkit.RetryMetadata{
		CurrenyRetry: 1,
		MaxRetry:     3,
		RetryAfter:   time.Second,
	}
	retryRaw, err := json.Marshal(retry)
	if err != nil {
		panic(err)
	}
	requestService.SaveRequest(model.RequestEntity{
		Service:   "sms",
		Payload:   testutil.MockPayload("sms"),
		Timestamp: time.Now().Add(-10 * time.Second),
		Retry:     string(retryRaw),
	})
	requestService.SaveRequest(model.RequestEntity{
		Service:   "email",
		Payload:   testutil.MockPayload("email"),
		Timestamp: time.Now().Add(10 * time.Second),
		Retry:     string(retryRaw),
	})

	// Queue
	writeMessages := make(chan testutil.Output)
	queueService := &testutil.MockQueueService{
		WriteData: writeMessages,
	}

	// Scheduler
	s := service.Service{
		RequestService: requestService,
		QueueService:   queueService,
	}
	sch := NewScheduler(s, log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout)))
	go sch.HandleRequests(100 * time.Millisecond)

	// Expect send sms
	output := <-writeMessages
	var response model.Response
	if output.Topic != "sms" {
		t.Errorf("Expect topic to be sms, but got %v", output.Topic)
	}
	if err := json.Unmarshal(output.Data, &response); err != nil {
		panic(err)
	}
	if response.Payload["content"] != "sms" {
		t.Errorf("Expect content to be sms, but got %v", response.Payload["content"])
	}
	if response.Retry.CurrenyRetry != 2 {
		t.Errorf("Expect CurrenyRetry to be 2, but got %v", response.Retry.CurrenyRetry)
	}
	if response.Retry.MaxRetry != 3 {
		t.Errorf("Expect MaxRetry to be 3, but got %v", response.Retry.MaxRetry)
	}
	if response.Retry.RetryAfter != time.Second {
		t.Errorf("Expect RetryAfter to be 1 second, but got %v", response.Retry.RetryAfter)
	}

	// Expect not send email
	time.Sleep(100 * time.Millisecond)
	var entities []model.RequestEntity
	err = pgdb.Find(&entities).Error
	if err != nil {
		panic(err)
	}
	if len(entities) != 1 {
		t.Errorf("Expect number of saved requests to be 1, but got %v", len(entities))
	}
	if entities[0].Service != "email" {
		t.Errorf("Expect request sevice is email, but got %v", entities[0].Service)
	}
}
