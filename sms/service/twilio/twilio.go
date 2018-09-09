package twilio

type Twilio struct {
	// c *sendgrid.Client
}

func New(apiKey string) *Twilio {
	return &Twilio{}
}

func (sc Twilio) Send() error {
	return nil
}
