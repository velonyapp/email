package port

import (
	"context"

	"github.com/velonyapp/email/internal/domain/vo"
)

type Email struct {
	From    vo.Address
	To      []vo.Address
	CC      []vo.Address
	BCC     []vo.Address
	Subject vo.Subject
	Text    vo.Text
	HTML    vo.HTML
}

type EmailSender interface {
	Send(ctx context.Context, email Email, idempotencyKey string) error
}
