package mailgun

import (
	"github.com/dwarvesf/yggdrasil/email/model"
	"github.com/dwarvesf/yggdrasil/email/service"
	mailgun "github.com/mailgun/mailgun-go"
)

//Mailgun contain mailgun client
type Mailgun struct {
	m mailgun.Mailgun
}

//New crete new mailgun client
func New(domain, apiKey, pubKey string) email.Emailer {
	return &Mailgun{
		m: mailgun.NewMailgun(domain, apiKey, pubKey),
	}
}

//Send mail via mailgun service
func (mg Mailgun) Send(apiKey string, p *model.Payload) error {
	message := mg.m.NewMessage(p.From.Email, "", p.Content, p.To.Email)
	_, _, err := mg.m.Send(message)
	return err
}
