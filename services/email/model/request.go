package model

import (
	"github.com/dwarvesf/yggdrasil/toolkit"
)

// Request is a struct define request message taken from queue
type Request struct {
	Payload Payload               `json:"payload"`
	Retry   toolkit.RetryMetadata `json:"retry"`
}

// Payload is the content of request
type Payload struct {
	From       Person            `json:"from"`
	To         Person            `json:"to"`
	Provider   string            `json:"provider" validate:"nonzero"`
	TemplateID string            `json:"template_id"`
	Data       map[string]string `json:"data"`
	Content    string            `json:"content"`
	Retry      int               `json:"retry"`
	Sent       bool              `json:"sent"`
}

type Person struct {
	Name  string `validate:"nonzero"`
	Email string `validate:"nonzero"`
}

// MailgunSecret contain Mailgun secret
type MailgunSecret struct {
	APIKey    string `json:"api_key"`
	PublicKey string `json:"pub_key"`
	Domain    string `json:"domain"`
	HostName  string `json:"host"`
}
