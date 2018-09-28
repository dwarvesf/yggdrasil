package mailgun

import (
	mailgun "github.com/mailgun/mailgun-go"
)

//Mailguner contain method to send email via mailgun
type Mailguner interface {
	New(domain, privateAPIKey, publicValidationKey string) *Client
	Send(sender, subject, body, recipient string) error
}

//Client contain mailgun client
type Client struct {
	m mailgun.Mailgun
}

//New crete new mailgun client
func New(domain, privateAPIKey, publicValidationKey string) *Client {
	return &Client{m: mailgun.NewMailgun(domain, privateAPIKey, publicValidationKey)}
}

//Send mail via mailgun service
func (mg Client) Send(sender, body, recipient string) error {
	message := mg.m.NewMessage(sender, "", body, recipient)
	_, _, err := mg.m.Send(message)
	return err
}
