package command

import (
	"context"

	"github.com/velonyapp/email/internal/application/port"
	"github.com/velonyapp/email/internal/domain/vo"
)

type SendEmail struct {
	From           string
	To             []string
	CC             []string
	BCC            []string
	Subject        string
	HTML           string
	Text           string
	IdempotencyKey string
}

type SendEmailResult struct{}

func (SendEmail) resultType() SendEmailResult {
	return SendEmailResult{}
}

type SendEmailHandler Handler[SendEmail, SendEmailResult]

type sendEmailHandler struct {
	emailSender port.EmailSender
}

func NewSendEmailHandler(
	emailSender port.EmailSender,
) SendEmailHandler {
	return &sendEmailHandler{
		emailSender: emailSender,
	}
}

func (h *sendEmailHandler) Handle(
	ctx context.Context,
	cmd SendEmail,
) (SendEmailResult, error) {
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
		Subject: vo.NewSubject(cmd.Subject),
		Text:    vo.NewText(cmd.Text),
		HTML:    vo.NewHTML(cmd.HTML),
	}, cmd.IdempotencyKey); err != nil {
		return SendEmailResult{}, err
	}

	return SendEmailResult{}, nil
}
