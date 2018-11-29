package sendgrid

import (
	sg "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/dwarvesf/yggdrasil/email/model"
	"github.com/dwarvesf/yggdrasil/email/service"
)

//Sendgrid struct contains a Client struct
type Sendgrid struct {
	apiKey string
}

// New returns a SendGridClient struct
func New(apiKey string) email.Emailer {
	return &Sendgrid{
		apiKey: apiKey,
	}
}

// Send sends an email via sendgrid
func (sc *Sendgrid) Send(p *model.Payload) error {
	fromName := p.From.Name
	if fromName == "" {
		return ErrNameIsRequired
	}

	fromEmail := p.From.Email
	if fromEmail == "" {
		return ErrEmailIsRequired
	}

	toName := p.To.Name
	if toName == "" {
		return ErrNameIsRequired
	}

	toEmail := p.To.Email
	if toEmail == "" {
		return ErrEmailIsRequired
	}

	m := mail.NewV3Mail()

	from := mail.NewEmail(fromName, fromEmail)
	m.SetFrom(from)

	person := mail.NewPersonalization()

	if p.TemplateID == "" {
		c := mail.NewContent("text/plain", p.Content)
		m.AddContent(c)
	} else {
		m.SetTemplateID(p.TemplateID)
		person.SetDynamicTemplateData("data", p.Data)
	}

	tos := []*mail.Email{
		mail.NewEmail(toName, toEmail),
	}

	person.AddTos(tos...)
	m.AddPersonalizations(person)

	body := mail.GetRequestBody(m)
	request := sg.GetRequest(sc.apiKey, "v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = body

	_, err := sg.API(request)
	if err != nil {
		return err
	}

	return nil
}
