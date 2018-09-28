package sendgrid

import (
	sendgrid "github.com/sendgrid/sendgrid-go"
)

//SendGrider contain method send email via Sendgrid
type SendGrider interface {
	New(apiKey string) *Client
	Send() error
}

//Client contain Sendgrid client
type Client struct {
	c *sendgrid.Client
}

//New create new Sendgrid client
func New(apiKey string) *Client {
	return &Client{c: sendgrid.NewSendClient(apiKey)}
}

//Send send email via Sendgrid
func (sc Client) Send() error {
	return nil
}
