package sendgrid

import (
	sendgrid "github.com/sendgrid/sendgrid-go"
)

type SendGridClient struct {
	c *sendgrid.Client
}

func New(apiKey string) *SendGridClient {
	return &SendGridClient{c: sendgrid.NewSendClient(apiKey)}
}

func (sc SendGridClient) Send() error {
	return nil
}
