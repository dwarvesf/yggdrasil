package model

import "time"

// Request is a struct define request message taken from queue
type Request struct {
	Service   string
	Payload   map[string]interface{}
	Timestamp time.Time
}
