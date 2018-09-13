package model

// Request is a struct define request message taken from queue
type Request struct {
	From       Person            `json:"from"`
	To         Person            `json:"to"`
	Provider   string            `json:"provider" validate:"nonzero"`
	TemplateID string            `json:"template_id"`
	Data       map[string]string `json:"data"`
	Content    string            `json:"content"`
	Retry      int               `json:"retry"`
	Sent       bool              `json:"sent"`
}

type Person struct {
	Name  string `validate:"nonzero"`
	Email string `validate:"nonzero"`
}
