package model

// Request is a struct define request message taken from queue
type Request struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Provider string `json:"provider"`
	Content  string `json:"content"`
}
