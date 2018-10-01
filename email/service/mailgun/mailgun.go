package mailgun

import (
	mailgun "github.com/mailgun/mailgun-go"
)

//Mailgun contain mailgun client
type Mailgun struct {
	m mailgun.Mailgun
}

//New crete new mailgun client
func New(domain, privateAPIKey, publicValidationKey string) Mailguner {
	return &Mailgun{m: mailgun.NewMailgun(domain, privateAPIKey, publicValidationKey)}
}

//Send mail via mailgun service
func (mg Mailgun) Send(sender, body, recipient string) error {
	message := mg.m.NewMessage(sender, "", body, recipient)
	_, _, err := mg.m.Send(message)
	return err
}
