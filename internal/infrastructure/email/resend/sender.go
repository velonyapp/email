package resend

import (
	"context"
	"net/mail"

	"github.com/velonyapp/email/internal/application/port"
	"github.com/velonyapp/email/internal/conf"

	"github.com/resend/resend-go/v3"
)

type Sender struct {
	client *resend.Client
}

func NewSender(c *conf.Email, client *resend.Client) port.EmailSender {
	return &Sender{
		client: client,
	}
}

func (s *Sender) Send(ctx context.Context, email *port.Email, idempotencyKey *string) (string, error) {
	req := &resend.SendEmailRequest{
		To:      make([]string, 0, len(email.To)),
		Cc:      make([]string, 0, len(email.CC)),
		Bcc:     make([]string, 0, len(email.BCC)),
		Subject: email.Subject,
	}

	if email.From.Name != nil {
		req.From = (&mail.Address{
			Name:    *email.From.Name,
			Address: email.From.Address,
		}).String()
	} else {
		req.From = email.From.Address
	}
	for _, address := range email.To {
		if address.Name != nil {
			req.To = append(req.To, (&mail.Address{
				Name:    *address.Name,
				Address: address.Address,
			}).String())
		} else {
			req.To = append(req.To, address.Address)
		}
	}
	for _, address := range email.CC {
		if address.Name != nil {
			req.Cc = append(req.Cc, (&mail.Address{
				Name:    *address.Name,
				Address: address.Address,
			}).String())
		} else {
			req.Cc = append(req.Cc, address.Address)
		}
	}
	for _, address := range email.BCC {
		if address.Name != nil {
			req.Bcc = append(req.Bcc, (&mail.Address{
				Name:    *address.Name,
				Address: address.Address,
			}).String())
		} else {
			req.Bcc = append(req.Bcc, address.Address)
		}
	}

	if email.ReplyTo != nil {
		if email.ReplyTo.Name != nil {
			req.ReplyTo = (&mail.Address{
				Name:    *email.ReplyTo.Name,
				Address: email.ReplyTo.Address,
			}).String()
		} else {
			req.ReplyTo = email.ReplyTo.Address
		}
	}

	if email.Text != nil {
		req.Text = *email.Text
	}
	if email.HTML != nil {
		req.Html = *email.HTML
	}

	var result *resend.SendEmailResponse
	var err error

	if idempotencyKey != nil {
		result, err = s.client.Emails.SendWithOptions(ctx, req,
			&resend.SendEmailOptions{
				IdempotencyKey: *idempotencyKey,
			},
		)
	} else {
		result, err = s.client.Emails.SendWithContext(ctx, req)
	}
	if err != nil {
		return "", err
	}

	return result.Id, nil
}
