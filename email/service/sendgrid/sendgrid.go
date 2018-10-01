package sendgrid

import (
	"github.com/dwarvesf/yggdrasil/email/service"
	sg "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/dwarvesf/yggdrasil/email/model"
)

//Sendgrid struct contains a Client struct
type Sendgrid struct {
	apiKey string
	p      *model.Payload
}

// New returns a SendGridClient struct
func New(apiKey string, p *model.Payload) email.Emailer {
	return &Sendgrid{
		apiKey: apiKey,
		p:      p,
	}
}

// Send sends an email via sendgrid
func (sc *Sendgrid) Send(apiKey string, p *model.Payload) error {
	m := mail.NewV3Mail()

	fromName := p.From.Name
	if fromName == "" {
		return ErrNameIsRequired
	}

	fromEmail := p.From.Email
	if fromEmail == "" {
		return ErrEmailIsRequired
	}

	from := mail.NewEmail(fromName, fromEmail)
	m.SetFrom(from)

	person := mail.NewPersonalization()

	if p.TemplateID == "" {
		c := mail.NewContent("text/plain", p.Content)
		m.AddContent(c)
	}

	m.SetTemplateID(p.TemplateID)
	person.SetDynamicTemplateData("data", p.Data)

	toName := p.To.Name
	if toName == "" {
		return ErrNameIsRequired
	}

	toEmail := p.To.Email
	if toEmail == "" {
		return ErrEmailIsRequired
	}

	tos := []*mail.Email{
		mail.NewEmail(toName, toEmail),
	}
	person.AddTos(tos...)
	m.AddPersonalizations(person)

	body := mail.GetRequestBody(m)
	request := sg.GetRequest(apiKey, "v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = body

	_, err := sg.API(request)
	if err != nil {
		return err
	}

	return nil
}
