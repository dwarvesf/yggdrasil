package model

import "github.com/dwarvesf/yggdrasil/toolkit"

// Request is a struct define request message taken from queue
type Request struct {
	Payload Payload               `json:"payload"`
	Retry   toolkit.RetryMetadata `json:"retry"`
}

// Payload is the content of request
type Payload struct {
	DeviceTokens []string               `json:"device_tokens"`
	Body         string                 `json:"body"`
	Title        string                 `json:"title"`
	Provider     string                 `json:"provider"`
	DeviceType   string                 `json:"device_type"`
	Data         map[string]interface{} `json:"data"`
}
