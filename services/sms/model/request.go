package model

import "github.com/dwarvesf/yggdrasil/toolkit"

// Request is a struct define request message taken from queue
type Request struct {
	Payload Payload               `json:"payload"`
	Retry   toolkit.RetryMetadata `json:"retry"`
}

// Payload is the content of request
type Payload struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Provider string `json:"provider"`
	Content  string `json:"content"`
}

//TwilioSecret is a struct that define App sid, authenticate token and app number from consul
type TwilioSecret struct {
	AppSid    string `json:"sid"`
	AuthToken string `json:"token"`
	AppNumber string `json:"number"`
}
