package model

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
)

// Request is a struct define request message taken from queue
type Request struct {
	Service   string
	Payload   map[string]interface{}
	Timestamp time.Time
}

// RequestEntity is a struct define request message to be saved in db
type RequestEntity struct {
	gorm.Model
	Service   string
	Payload   string
	Timestamp time.Time
}

// ToRequest convert entity to request
func (e RequestEntity) ToRequest() (Request, error) {
	var payload map[string]interface{}

	err := json.Unmarshal([]byte(e.Payload), &payload)
	if err != nil {
		return Request{}, err
	}

	return Request{
		Service:   e.Service,
		Payload:   payload,
		Timestamp: e.Timestamp,
	}, nil
}
