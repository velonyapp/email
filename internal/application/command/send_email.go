package command

import (
	"bytes"
	"context"
	htmltemplate "html/template"
	texttemplate "text/template"

	"github.com/velonyapp/email/internal/application/common"
	"github.com/velonyapp/email/internal/application/port"
	"github.com/velonyapp/email/internal/domain/repo"
	"github.com/velonyapp/email/internal/domain/vo"
)

type SendEmail struct {
	From           string
	To             []string
	CC             []string
	BCC            []string
	TemplateAlias  string
	Variables      map[string]string
	IdempotencyKey string
}

type SendEmailResult struct{}

func (SendEmail) resultType() SendEmailResult {
	return SendEmailResult{}
}

type SendEmailHandler Handler[SendEmail, SendEmailResult]

type sendEmailHandler struct {
	templateRepo repo.Template
	emailSender  port.EmailSender
}

func NewSendEmailHandler(
	templateRepo repo.Template,
	emailSender port.EmailSender,
) SendEmailHandler {
	return &sendEmailHandler{
		templateRepo: templateRepo,
		emailSender:  emailSender,
	}
}

func (h *sendEmailHandler) Handle(
	ctx context.Context,
	cmd SendEmail,
) (SendEmailResult, error) {
	templateAlias, err := vo.NewAlias(cmd.TemplateAlias)
	if err != nil {
		return SendEmailResult{}, err
	}

	template, err := h.templateRepo.FindByAlias(ctx, templateAlias)
	if err != nil {
		return SendEmailResult{}, err
	}
	if template == nil {
		return SendEmailResult{}, common.ErrTemplateNotFound
	}

	subjectTemplate, err := texttemplate.New("subject").
		Option("missingkey=error").
		Parse(template.Subject().String())
	if err != nil {
		return SendEmailResult{}, err
	}
	var subject bytes.Buffer
	if err := subjectTemplate.Execute(&subject, cmd.Variables); err != nil {
		return SendEmailResult{}, err
	}

	var html bytes.Buffer
	if template.HasHTML() {
		htmlTemplate, err := htmltemplate.New("html").
			Option("missingkey=error").
			Parse(template.HTML().String())
		if err != nil {
			return SendEmailResult{}, err
		}
		if err := htmlTemplate.Execute(&html, cmd.Variables); err != nil {
			return SendEmailResult{}, err
		}
	}

	var text bytes.Buffer
	if template.HasText() {
		textTemplate, err := texttemplate.New("text").
			Option("missingkey=error").
			Parse(template.Text().String())
		if err != nil {
			return SendEmailResult{}, err
		}
		if err := textTemplate.Execute(&text, cmd.Variables); err != nil {
			return SendEmailResult{}, err
		}
	}

	from, err := vo.NewAddress(cmd.From)
	if err != nil {
		return SendEmailResult{}, err
	}
	to := make([]vo.Address, len(cmd.To))
	for i, value := range cmd.To {
		address, err := vo.NewAddress(value)
		if err != nil {
			return SendEmailResult{}, err
		}
		to[i] = address
	}
	cc := make([]vo.Address, len(cmd.CC))
	for i, value := range cmd.CC {
		address, err := vo.NewAddress(value)
		if err != nil {
			return SendEmailResult{}, err
		}
		cc[i] = address
	}
	bcc := make([]vo.Address, len(cmd.BCC))
	for i, value := range cmd.BCC {
		address, err := vo.NewAddress(value)
		if err != nil {
			return SendEmailResult{}, err
		}
		bcc[i] = address
	}

	if err = h.emailSender.Send(ctx, port.Email{
		From:    from,
		To:      to,
		CC:      cc,
		BCC:     bcc,
		Subject: vo.NewSubject(subject.String()),
		Text:    vo.NewText(text.String()),
		HTML:    vo.NewHTML(html.String()),
	}, cmd.IdempotencyKey); err != nil {
		return SendEmailResult{}, err
	}

	return SendEmailResult{}, nil
}
