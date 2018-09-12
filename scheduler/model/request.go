package model

import (
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
