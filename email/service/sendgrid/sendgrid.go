package sendgrid

import (
	sg "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/dwarvesf/yggdrasil/email/model"
)

//Sendgrid struct contains a Client struct
type Sendgrid struct {
	c sg.Client
}

// New returns a SendGridClient struct
func New(apiKey string) Sendgrider {
	return Sendgrid{c: *sg.NewSendClient(apiKey)}
}

// Send sends an email via sendgrid
func (sc Sendgrid) Send(apiKey string, req *model.Request) error {
	m := mail.NewV3Mail()

	fromName := req.Payload.From.Name
	if fromName == "" {
		return ErrNameIsRequired
	}

	fromEmail := req.Payload.From.Email
	if fromEmail == "" {
		return ErrEmailIsRequired
	}

	from := mail.NewEmail(fromName, fromEmail)
	m.SetFrom(from)

	p := mail.NewPersonalization()

	if req.Payload.TemplateID == "" {
		c := mail.NewContent("text/plain", req.Payload.Content)
		m.AddContent(c)
	}

	m.SetTemplateID(req.Payload.TemplateID)
	p.SetDynamicTemplateData("data", req.Payload.Data)

	toName := req.Payload.To.Name
	if toName == "" {
		return ErrNameIsRequired
	}

	toEmail := req.Payload.To.Email
	if toEmail == "" {
		return ErrEmailIsRequired
	}

	tos := []*mail.Email{
		mail.NewEmail(toName, toEmail),
	}
	p.AddTos(tos...)
	m.AddPersonalizations(p)

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
