package model

type Request struct {
	From       Person            `json:"from" validate:"nonzero"`
	To         Person            `json:"to" validate:"nonzero"`
	Type       string            `json:"type" validate:"nonzero"`
	TemplateID string            `json:"template_id"`
	Data       map[string]string `json:"data"`
	Content    string            `json:"content"`
	Retry      int               `json:"retry"`
	Sent       bool              `json:"sent"`
}

type Person struct {
	Name  string `json:"name" validate:"nonzero"`
	Email string `json:"email" validate:"nonzero"`
}
