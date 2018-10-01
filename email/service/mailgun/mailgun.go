package mailgun

import (
	mailgun "github.com/mailgun/mailgun-go"
)

//Client contain mailgun client
type Client struct {
	m mailgun.Mailgun
}

//NewMailgun crete new mailgun client
func NewMailgun(domain, privateAPIKey, publicValidationKey string) *Client {
	return &Client{m: mailgun.NewMailgun(domain, privateAPIKey, publicValidationKey)}
}

//SendMailgun mail via mailgun service
func (mg Client) SendMailgun(sender, body, recipient string) error {
	message := mg.m.NewMessage(sender, "", body, recipient)
	_, _, err := mg.m.Send(message)
	return err
}
