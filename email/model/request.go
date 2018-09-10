package model

// Request is a struct define request message taken from queue
type Request struct {
	From       Person            `json:"from" validate:"nonzero"`
	To         Person            `json:"to" validate:"nonzero"`
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
