package model

//Request store structer request from queue
type Request struct {
	DeviceToken string `json:"device_token"`
	Body        string `json:"body"`
	Title       string `json:"title"`
	Provider    string `json:"provider"`
	DeviceType  string `json:"device_type"`
}
