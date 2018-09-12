package sms

type SMS interface {
	Send(from, to, content, appSid string) error
}
