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
	templateAlias, err := parseEmailTemplateName(payload.Template)
	if err != nil {
		return err
	}

	if _, err := command.Send(ctx, s.commandBus, command.SendEmail{
		From:           payload.From,
		To:             payload.To,
		CC:             payload.Cc,
		BCC:            payload.Bcc,
		TemplateAlias:  templateAlias,
		Variables:      payload.Variables,
		IdempotencyKey: payload.IdempotencyKey,
	}); err != nil {
		return err
	}

	return nil
}
