package service

import (
	"context"

	v1 "github.com/velonyapp/email/gen/job/v1"
	"github.com/velonyapp/email/internal/application/command"
)

type Job struct {
	commandBus *command.Bus
}

func NewJob(
	commandBus *command.Bus,
) *Job {
	return &Job{
		commandBus: commandBus,
	}
}

func (s *Job) SendEmail(ctx context.Context, payload *v1.SendEmailPayload) error {
	if content := payload.GetCustom(); content != nil {
		if _, err := command.Send(ctx, s.commandBus, command.SendEmail{
			From:           payload.From,
			To:             payload.To,
			CC:             payload.Cc,
			BCC:            payload.Bcc,
			Subject:        content.Subject,
			HTML:           content.Html,
			Text:           content.Text,
			IdempotencyKey: payload.IdempotencyKey,
		}); err != nil {
			return err
		}
	}
	if content := payload.GetTemplate(); content != nil {
		templateID, err := parseEmailTemplateName(content.Name)
		if err != nil {
			return err
		}

		if _, err := command.Send(ctx, s.commandBus, command.SendEmailFromTemplate{
			From:           payload.From,
			To:             payload.To,
			CC:             payload.Cc,
			BCC:            payload.Bcc,
			TemplateID:     templateID,
			Variables:      content.Variables,
			IdempotencyKey: payload.IdempotencyKey,
		}); err != nil {
			return err
		}
	}

	return nil
}
