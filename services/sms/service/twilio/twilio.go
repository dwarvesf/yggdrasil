package twilio

import (
	twilio "github.com/sfreiberg/gotwilio"
)

//Client struct contain Twilio Client data
type Client struct {
	t *twilio.Twilio
}

//New create struct Twilio client
func New(appSid, authToken string) *Client {
	return &Client{t: twilio.NewTwilioClient(appSid, authToken)}
}

//Send send sms
func (tw *Client) Send(from, to, content, appSid string) error {
	_, _, err := tw.t.SendSMS(from, to, content, "", appSid)
	return err
}
