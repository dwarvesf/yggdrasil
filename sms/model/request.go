package model

// Request is a struct define request message taken from queue
type Request struct {
	From     string
	To       string `json:"to"`
	Provider string `json:"provider"`
	Content  string `json:"content"`
}

//TwilioSecret is a struct that define App sid, authenticate token and app number from consul
type TwilioSecret struct {
	AppSid    string `json:"sid"`
	AuthToken string `json:"token"`
	AppNumber string `json:"number"`
}
