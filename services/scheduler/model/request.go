package model

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/dwarvesf/yggdrasil/toolkit"
)

// Request is a struct define request message taken from queue
type Request = toolkit.SchedulerRequest

// Response is a struct define message scheduler send to message queue
type Response struct {
	Payload map[string]interface{} `json:"payload"`
	Retry   toolkit.RetryMetadata  `json:"retry"`
}

// RequestEntity is a struct define request message to be saved in db
type RequestEntity struct {
	gorm.Model
	Service   string
	Payload   string
	Timestamp time.Time
	Retry     string
}

// ToRequest convert entity to request
func (e RequestEntity) ToRequest() (Request, error) {
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(e.Payload), &payload); err != nil {
		return Request{}, err
	}

	var retry toolkit.RetryMetadata
	if err := json.Unmarshal([]byte(e.Retry), &retry); err != nil {
		return Request{}, err
	}

	return Request{
		Service:   e.Service,
		Payload:   payload,
		Timestamp: e.Timestamp,
		Retry:     retry,
	}, nil
}
