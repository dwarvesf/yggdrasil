package model

type Request struct {
	From       Person            `json:"from"`
	To         Person            `json:"to"`
	Type       string            `json:"type"`
	TemplateID string            `json:"template_id"`
	Data       map[string]string `json:"data"`
	Content    string            `json:"content"`
	Retry      int               `json:"retry"`
	Sent       bool              `json:"sent"`
}

type Person struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
