package scheduler

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/dwarvesf/yggdrasil/scheduler/model"
	"github.com/dwarvesf/yggdrasil/scheduler/service"
	"github.com/dwarvesf/yggdrasil/scheduler/testutils"
)

func TestListenMessagesWhenMessageExpiredShouldNotSave(t *testing.T) {
	requests := make(chan []byte)
	requestService := &testutils.MockRequestService{
		Requests: make([]model.RequestEntity, 0),
	}
	messageService := &testutils.MockMessageService{
		ReadData: requests,
	}
	s := service.Service{
		RequestService: requestService,
		MessageService: messageService,
	}

	sch := NewScheduler(s, log.NewNopLogger())
	go sch.ListenMessages()

	r, err := json.Marshal(model.Request{
		Service:   "sms",
		Timestamp: time.Now().Add(-1 * time.Second),
	})
	if err != nil {
		panic(err)
	}
	requests <- r
	time.Sleep(100 * time.Millisecond)

	if len(requestService.Requests) > 0 {
		t.Errorf("Expect number of saved requests to be zero, but got %v", len(requestService.Requests))
	}
}

func TestListenMessagesWhenInvalidServiceShouldNotSave(t *testing.T) {
	requests := make(chan []byte)
	requestService := &testutils.MockRequestService{
		Requests: make([]model.RequestEntity, 0),
	}
	messageService := &testutils.MockMessageService{
		ReadData: requests,
	}
	s := service.Service{
		RequestService: requestService,
		MessageService: messageService,
	}

	sch := NewScheduler(s, log.NewNopLogger())
	go sch.ListenMessages()

	r, err := json.Marshal(model.Request{
		Service:   "invalid",
		Timestamp: time.Now().Add(5 * time.Second),
	})
	if err != nil {
		panic(err)
	}
	requests <- r
	time.Sleep(100 * time.Millisecond)

	if len(requestService.Requests) > 0 {
		t.Errorf("Expect number of saved requests to be zero, but got %v", len(requestService.Requests))
	}
}

func TestListenMessagesWhenMessageValidShouldSave(t *testing.T) {
	requests := make(chan []byte)
	requestService := &testutils.MockRequestService{
		Requests: make([]model.RequestEntity, 0),
	}
	messageService := &testutils.MockMessageService{
		ReadData: requests,
	}
	s := service.Service{
		RequestService: requestService,
		MessageService: messageService,
	}

	sch := NewScheduler(s, log.NewNopLogger())
	go sch.ListenMessages()

	r, err := json.Marshal(model.Request{
		Service:   "email",
		Timestamp: time.Now().Add(5 * time.Second),
	})
	if err != nil {
		panic(err)
	}
	requests <- r
	time.Sleep(100 * time.Millisecond)

	if len(requestService.Requests) != 1 {
		t.Errorf("Expect number of saved requests to be 1, but got %v", len(requestService.Requests))
	}

	if requestService.Requests[0].Service != "email" {
		t.Errorf("Expect request sevice is email, but got %v", requestService.Requests[0].Service)
	}
}

func TestHandleRequests(t *testing.T) {
	// Prepare request
	requests := make([]model.RequestEntity, 2)
	requests[0] = model.RequestEntity{
		Service:   "sms",
		Payload:   testutils.MockPayload("sms"),
		Timestamp: time.Now().Add(-10 * time.Second),
	}
	requests[0].ID = 1
	requests[1] = model.RequestEntity{
		Service:   "payment",
		Payload:   testutils.MockPayload("payment"),
		Timestamp: time.Now().Add(-5 * time.Second),
	}
	requests[1].ID = 2

	writeMessages := make(chan testutils.Output)
	requestService := &testutils.MockRequestService{
		Requests: requests,
	}
	messageService := &testutils.MockMessageService{
		WriteData: writeMessages,
	}
	s := service.Service{
		RequestService: requestService,
		MessageService: messageService,
	}

	sch := NewScheduler(s, log.NewNopLogger())
	go sch.HandleRequests(1 * time.Second)

	// Validate send sms
	output := <-writeMessages
	var payload map[string]interface{}
	if output.Topic != "sms" {
		t.Errorf("Expect topic to be sms, but got %v", output.Topic)
	}
	if err := json.Unmarshal(output.Data, &payload); err != nil {
		panic(err)
	}
	if payload["content"] != "sms" {
		t.Errorf("Expect content to be sms, but got %v", payload["content"])
	}

	// Validate send payment
	output = <-writeMessages
	if output.Topic != "payment" {
		t.Errorf("Expect topic to be payment, but got %v", output.Topic)
	}
	if err := json.Unmarshal(output.Data, &payload); err != nil {
		panic(err)
	}
	if payload["content"] != "payment" {
		t.Errorf("Expect content to be payment, but got %v", payload["content"])
	}

	// Validate delete requests
	time.Sleep(100 * time.Millisecond)
	if len(requestService.DeletedIds) != 2 {
		t.Errorf("Expect number of deleted ids to be 2, but got %v", len(requestService.DeletedIds))
	}
	if requestService.DeletedIds[0] != 1 {
		t.Errorf("Expect element with id 1 is deleted, but got %v", requestService.DeletedIds[0])
	}
	if requestService.DeletedIds[1] != 2 {
		t.Errorf("Expect element with id 2 is deleted, but got %v", requestService.DeletedIds[1])
	}
}
