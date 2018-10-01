package email

import "github.com/dwarvesf/yggdrasil/email/model"

//Emailer contain send method
type Emailer interface {
	Send(apiKey string, p *model.Payload) error
}
