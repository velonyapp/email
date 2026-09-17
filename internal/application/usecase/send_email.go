package usecase

import (
	"context"

	"github.com/velonyapp/email/internal/application/port"
)

type SendEmail struct {
	IdempotencyKey string

	Email *port.Email
}

type SendEmailResult struct {
	MessageID string
}

type SendEmailHandler struct {
	emailSender port.EmailSender
}

func NewSendEmailHandler(
	emailSender port.EmailSender,
) *SendEmailHandler {
	return &SendEmailHandler{
		emailSender: emailSender,
	}
}

func (h *SendEmailHandler) Execute(
	ctx context.Context,
	uc *SendEmail,
) (*SendEmailResult, error) {
	messageID, err := h.emailSender.Send(ctx, uc.IdempotencyKey, uc.Email)
	if err != nil {
		return nil, err
	}

	return &SendEmailResult{
		MessageID: messageID,
	}, nil
}
