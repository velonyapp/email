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

type SendEmailFromTemplate struct {
	From           string
	To             []string
	CC             []string
	BCC            []string
	TemplateID     string
	Variables      map[string]string
	IdempotencyKey string
}

type SendEmailFromTemplateResult struct{}

func (SendEmailFromTemplate) resultType() SendEmailFromTemplateResult {
	return SendEmailFromTemplateResult{}
}

type SendEmailFromTemplateHandler Handler[SendEmailFromTemplate, SendEmailFromTemplateResult]

type sendEmailFromTemplateHandler struct {
	templateRepo repo.Template
	emailSender  port.EmailSender
}

func NewSendEmailFromTemplateHandler(
	templateRepo repo.Template,
	emailSender port.EmailSender,
) SendEmailFromTemplateHandler {
	return &sendEmailFromTemplateHandler{
		templateRepo: templateRepo,
		emailSender:  emailSender,
	}
}

func (h *sendEmailFromTemplateHandler) Handle(
	ctx context.Context,
	cmd SendEmailFromTemplate,
) (SendEmailFromTemplateResult, error) {
	templateID, err := vo.NewTemplateID(cmd.TemplateID)
	if err != nil {
		return SendEmailFromTemplateResult{}, err
	}

	template, err := h.templateRepo.FindByID(ctx, templateID)
	if err != nil {
		return SendEmailFromTemplateResult{}, err
	}
	if template == nil {
		return SendEmailFromTemplateResult{}, common.ErrTemplateNotFound
	}

	subjectTemplate, err := texttemplate.New("subject").
		Option("missingkey=error").
		Parse(template.Subject().String())
	if err != nil {
		return SendEmailFromTemplateResult{}, err
	}
	var subject bytes.Buffer
	if err := subjectTemplate.Execute(&subject, cmd.Variables); err != nil {
		return SendEmailFromTemplateResult{}, err
	}

	var html bytes.Buffer
	if template.HasHTML() {
		htmlTemplate, err := htmltemplate.New("html").
			Option("missingkey=error").
			Parse(template.HTML().String())
		if err != nil {
			return SendEmailFromTemplateResult{}, err
		}
		if err := htmlTemplate.Execute(&html, cmd.Variables); err != nil {
			return SendEmailFromTemplateResult{}, err
		}
	}

	var text bytes.Buffer
	if template.HasText() {
		textTemplate, err := texttemplate.New("text").
			Option("missingkey=error").
			Parse(template.Text().String())
		if err != nil {
			return SendEmailFromTemplateResult{}, err
		}
		if err := textTemplate.Execute(&text, cmd.Variables); err != nil {
			return SendEmailFromTemplateResult{}, err
		}
	}

	from, err := vo.NewAddress(cmd.From)
	if err != nil {
		return SendEmailFromTemplateResult{}, err
	}
	to := make([]vo.Address, len(cmd.To))
	for i, value := range cmd.To {
		address, err := vo.NewAddress(value)
		if err != nil {
			return SendEmailFromTemplateResult{}, err
		}
		to[i] = address
	}
	cc := make([]vo.Address, len(cmd.CC))
	for i, value := range cmd.CC {
		address, err := vo.NewAddress(value)
		if err != nil {
			return SendEmailFromTemplateResult{}, err
		}
		cc[i] = address
	}
	bcc := make([]vo.Address, len(cmd.BCC))
	for i, value := range cmd.BCC {
		address, err := vo.NewAddress(value)
		if err != nil {
			return SendEmailFromTemplateResult{}, err
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
		return SendEmailFromTemplateResult{}, err
	}

	return SendEmailFromTemplateResult{}, nil
}
