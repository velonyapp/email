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

func (s *Job) SendEmail(ctx context.Context, req *v1.SendEmail) error {
	if content := req.GetCustom(); content != nil {
		if _, err := command.Send(ctx, s.commandBus, command.SendEmail{
			From:           req.From,
			To:             req.To,
			CC:             req.Cc,
			BCC:            req.Bcc,
			Subject:        content.Subject,
			HTML:           content.Html,
			Text:           content.Text,
			IdempotencyKey: req.IdempotencyKey,
		}); err != nil {
			return err
		}
	}
	if content := req.GetTemplate(); content != nil {
		templateID, err := parseEmailTemplateName(content.Name)
		if err != nil {
			return err
		}

		if _, err := command.Send(ctx, s.commandBus, command.SendEmailFromTemplate{
			From:           req.From,
			To:             req.To,
			CC:             req.Cc,
			BCC:            req.Bcc,
			TemplateID:     templateID,
			Variables:      content.Variables,
			IdempotencyKey: req.IdempotencyKey,
		}); err != nil {
			return err
		}
	}

	return nil
}
