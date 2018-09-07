package email

type Emailer interface {
	Send() error
}
