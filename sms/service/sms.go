package sms

type SMS interface {
	Send() error
}
