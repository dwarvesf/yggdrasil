package model

import (
	"time"

	"github.com/dwarvesf/yggdrasil/toolkit"
)

// Request is a struct define request message taken from queue
type Request struct {
	Payload Payload               `json:"payload"`
	Retry   toolkit.RetryMetadata `json:"retry"`
}

// Payload is the content of request
type Payload struct {
	Currency   string     `json:"currency,omitempty"`
	Desc       string     `json:"desc,omitempty"`
	Provider   string     `json:"provider,omitempty"`
	CreditCard CreditCard `json:"credit_card,omitempty"`
	Amount     int64      `json:"amount,omitempty"`
	Token      string     `json:"token,omitempty"`
}

// CreditCard is a struct to store info of customer card
type CreditCard struct {
	Name       string    `json:"name,omitempty"`
	Number     string    `json:"number,omitempty"`
	Cvc        string    `json:"cvc,omitempty"`
	ExpireDate time.Time `json:"expire_date,omitempty"`
}
