package sendgrid

import email "github.com/dwarvesf/yggdrasil/email/service"

type SendgridClient struct {
	email.Emailer
}

func NewSendgridClient() *SendgridClient {
	return &SendgridClient{}
}

func (sc SendgridClient) Send() error {
	return nil
}
