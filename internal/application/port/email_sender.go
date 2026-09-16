package port

import "context"

type EmailAddress struct {
	Address string
	Name    *string
}

type Email struct {
	To      []EmailAddress
	Cc      []EmailAddress
	Bcc     []EmailAddress
	ReplyTo *EmailAddress

	Subject string
	Text    *string
	HTML    *string
}

type EmailSender interface {
	Send(ctx context.Context, idempotencyKey string, email *Email) (messageID string, err error)
}
