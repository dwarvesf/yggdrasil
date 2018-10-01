package email

import (
	"github.com/dwarvesf/yggdrasil/email/model"
	"github.com/dwarvesf/yggdrasil/email/service/mailgun"
	"github.com/dwarvesf/yggdrasil/email/service/sendgrid"
)

//Emailer contain send method
type Emailer interface {
	NewSendgrid(apiKey string) *sendgrid.Client
	NewMailgun(domain, privateAPIKey, publicValidationKey string) *mailgun.Client
	SendSendgrid(apiKey string, r *model.Request) error
	SendMailgun(sender, body, recipient string) error
}
