package api

import (
	"context"

	v1 "github.com/velonyapp/notification/gen/api/v1"
	"github.com/velonyapp/notification/internal/application/port"
	"github.com/velonyapp/notification/internal/application/usecase"
)

type Service struct {
	v1.UnimplementedNotificationServiceServer

	sendEmailHandler *usecase.SendEmailHandler
}

func NewService(
	sendEmailHandler *usecase.SendEmailHandler,
) *Service {
	return &Service{
		sendEmailHandler: sendEmailHandler,
	}
}

func (s *Service) SendEmail(
	ctx context.Context,
	req *v1.SendEmailRequest,
) (*v1.SendEmailResponse, error) {
	email := &port.Email{
		To:      make([]port.EmailAddress, 0, len(req.Email.To)),
		Cc:      make([]port.EmailAddress, 0, len(req.Email.Cc)),
		Bcc:     make([]port.EmailAddress, 0, len(req.Email.Bcc)),
		Subject: req.Email.Subject,
		Text:    req.Email.Text,
		HTML:    req.Email.Html,
	}

	for _, address := range req.Email.To {
		if address == nil {
			continue
		}

		email.To = append(email.To, port.EmailAddress{
			Address: address.Address,
			Name:    address.Name,
		})
	}
	for _, address := range req.Email.Cc {
		if address == nil {
			continue
		}

		email.Cc = append(email.Cc, port.EmailAddress{
			Address: address.Address,
			Name:    address.Name,
		})
	}
	for _, address := range req.Email.Bcc {
		if address == nil {
			continue
		}

		email.Bcc = append(email.Bcc, port.EmailAddress{
			Address: address.Address,
			Name:    address.Name,
		})
	}

	result, err := s.sendEmailHandler.Execute(ctx, &usecase.SendEmail{
		IdempotencyKey: req.IdempotencyKey,
		Email:          email,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &v1.SendEmailResponse{
		MessageId: result.MessageID,
	}, nil
}
