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
	from   string
}

func NewSender(c *conf.Service, client *resend.Client) port.EmailSender {
	return &Sender{
		client: client,
		from: (&mail.Address{
			Name:    c.From.Name,
			Address: c.From.Address,
		}).String(),
	}
}

func (s *Sender) Send(ctx context.Context, idempotencyKey string, email *port.Email) (string, error) {
	req := &resend.SendEmailRequest{
		From:    s.from,
		To:      make([]string, 0, len(email.To)),
		Cc:      make([]string, 0, len(email.Cc)),
		Bcc:     make([]string, 0, len(email.Bcc)),
		Subject: email.Subject,
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
	for _, address := range email.Cc {
		if address.Name != nil {
			req.Cc = append(req.Cc, (&mail.Address{
				Name:    *address.Name,
				Address: address.Address,
			}).String())
		} else {
			req.Cc = append(req.Cc, address.Address)
		}
	}
	for _, address := range email.Bcc {
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

	result, err := s.client.Emails.SendWithOptions(
		ctx,
		req,
		&resend.SendEmailOptions{
			IdempotencyKey: idempotencyKey,
		},
	)
	if err != nil {
		return "", err
	}

	return result.Id, nil
}
