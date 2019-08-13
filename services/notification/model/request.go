package model

import "github.com/dwarvesf/yggdrasil/toolkit"

// Request is a struct define request message taken from queue
type Request struct {
	Payload Payload               `json:"payload"`
	Retry   toolkit.RetryMetadata `json:"retry"`
}

type DeviceToken struct {
	Token string `json:"token"`
	Badge int    `json:"badge"`
}

// Payload is the content of request
type Payload struct {
	DeviceTokens []DeviceToken          `json:"device_tokens"`
	Body         string                 `json:"body"`
	Title        string                 `json:"title"`
	Provider     string                 `json:"provider"`
	DeviceType   string                 `json:"device_type"`
	Data         map[string]interface{} `json:"data"`
}
