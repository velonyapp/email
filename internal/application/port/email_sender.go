package port

import "context"

type EmailAddress struct {
	Address string
	Name    *string
}

type Email struct {
	From    EmailAddress
	To      []EmailAddress
	CC      []EmailAddress
	BCC     []EmailAddress
	ReplyTo *EmailAddress

	Subject string
	Text    *string
	HTML    *string
}

type EmailSender interface {
	Send(ctx context.Context, email *Email, idempotencyKey *string) (messageID string, err error)
}
