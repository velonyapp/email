package resend

import (
	"context"

	"github.com/velonyapp/email/internal/application/port"

	"github.com/resend/resend-go/v3"
)

type sender struct {
	client *resend.Client
}

func NewSender(client *resend.Client) port.EmailSender {
	return &sender{
		client: client,
	}
}

func (s *sender) Send(ctx context.Context, email port.Email, idempotencyKey string) error {
	req := &resend.SendEmailRequest{
		From:    email.From.String(),
		To:      make([]string, 0, len(email.To)),
		Cc:      make([]string, 0, len(email.CC)),
		Bcc:     make([]string, 0, len(email.BCC)),
		Subject: email.Subject.String(),
		Html:    email.HTML.String(),
		Text:    email.Text.String(),
	}

	for _, address := range email.To {
		req.To = append(req.To, address.String())
	}
	for _, address := range email.CC {
		req.Cc = append(req.Cc, address.String())
	}
	for _, address := range email.BCC {
		req.Bcc = append(req.Bcc, address.String())
	}

	if idempotencyKey == "" {
		if _, err := s.client.Emails.SendWithContext(ctx, req); err != nil {
			return err
		}
	} else {
		if _, err := s.client.Emails.SendWithOptions(ctx, req,
			&resend.SendEmailOptions{
				IdempotencyKey: idempotencyKey,
			},
		); err != nil {
			return err
		}
	}

	return nil
}
