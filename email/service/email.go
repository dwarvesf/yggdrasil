package email

import (
	"github.com/dwarvesf/yggdrasil/email/service/mailgun"
	"github.com/dwarvesf/yggdrasil/email/service/sendgrid"
)

//Email contain send method
type Email struct {
	SendGrid sendgrid.SendGrider
	Mailgun  mailgun.Mailguner
}
