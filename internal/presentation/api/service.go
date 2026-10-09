package api

import (
	"context"
	"errors"

	v1 "github.com/velonyapp/email/gen/api/v1"
	"github.com/velonyapp/email/internal/application/command"
	"github.com/velonyapp/email/internal/application/query"

	"go.einride.tech/aip/resourcename"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrInvalidEmailTemplateResourceName = errors.New("invalid email template resource name")
	ErrInvalidUpdateMaskPath            = errors.New("invalid update mask path")
)

const (
	emailTemplateResourcePattern = "emailTemplates/{email_template}"
)

type Service struct {
	v1.UnimplementedEmailServiceServer

	commandBus *command.Bus
	queryBus   *query.Bus
}

func NewService(
	commandBus *command.Bus,
	queryBus *query.Bus,
) *Service {
	return &Service{
		commandBus: commandBus,
		queryBus:   queryBus,
	}
}

func (s *Service) GetEmailTemplate(ctx context.Context, req *v1.GetEmailTemplateRequest) (*v1.EmailTemplate, error) {
	var templateID string
	if err := resourcename.Sscan(req.Name, emailTemplateResourcePattern, &templateID); err != nil {
		return nil, ErrInvalidEmailTemplateResourceName
	}

	result, err := query.Send(ctx, s.queryBus, query.GetTemplate{
		TemplateID: templateID,
	})
	if err != nil {
		return nil, err
	}

	return &v1.EmailTemplate{
		Name:    resourcename.Sprint(emailTemplateResourcePattern, result.Template.ID),
		Alias:   result.Template.Alias,
		Subject: result.Template.Subject,
		Html:    result.Template.HTML,
		Text:    result.Template.Text,
	}, nil
}

func (s *Service) ListEmailTemplates(ctx context.Context, req *v1.ListEmailTemplatesRequest) (*v1.ListEmailTemplatesResponse, error) {
	result, err := query.Send(ctx, s.queryBus, query.ListTemplates{})
	if err != nil {
		return nil, err
	}

	templates := make([]*v1.EmailTemplate, len(result.Templates))

	for i, template := range result.Templates {
		templates[i] = &v1.EmailTemplate{
			Name:    resourcename.Sprint(emailTemplateResourcePattern, template.ID),
			Alias:   template.Alias,
			Subject: template.Subject,
			Html:    template.HTML,
			Text:    template.Text,
		}
	}

	return &v1.ListEmailTemplatesResponse{
		EmailTemplates: templates,
	}, nil
}

func (s *Service) CreateEmailTemplate(ctx context.Context, req *v1.CreateEmailTemplateRequest) (*v1.EmailTemplate, error) {
	result, err := command.Send(ctx, s.commandBus, command.CreateTemplate{
		Alias:   req.EmailTemplate.Alias,
		Subject: req.EmailTemplate.Subject,
		HTML:    req.EmailTemplate.Html,
		Text:    req.EmailTemplate.Text,
	})
	if err != nil {
		return nil, err
	}

	return &v1.EmailTemplate{
		Name:    resourcename.Sprint(emailTemplateResourcePattern, result.Template.ID),
		Alias:   result.Template.Alias,
		Subject: result.Template.Subject,
		Html:    result.Template.HTML,
		Text:    result.Template.Text,
	}, nil
}

func (s *Service) UpdateEmailTemplate(ctx context.Context, req *v1.UpdateEmailTemplateRequest) (*v1.EmailTemplate, error) {
	var templateID string
	if err := resourcename.Sscan(req.EmailTemplate.Name, emailTemplateResourcePattern, &templateID); err != nil {
		return nil, ErrInvalidEmailTemplateResourceName
	}

	var alias, subject, html, text *string

	paths := req.UpdateMask.GetPaths()

	if paths == nil {
		if req.EmailTemplate.Alias != "" {
			value := req.EmailTemplate.Alias
			alias = &value
		}
		if req.EmailTemplate.Subject != "" {
			value := req.EmailTemplate.Subject
			subject = &value
		}
		if req.EmailTemplate.Html != "" {
			value := req.EmailTemplate.Html
			html = &value
		}
		if req.EmailTemplate.Text != "" {
			value := req.EmailTemplate.Text
			text = &value
		}
	} else {
		for _, path := range paths {
			switch path {
			case "*":
				aliasValue := req.EmailTemplate.Alias
				alias = &aliasValue

				subjectValue := req.EmailTemplate.Subject
				subject = &subjectValue

				htmlValue := req.EmailTemplate.Html
				html = &htmlValue

				textValue := req.EmailTemplate.Text
				text = &textValue

			case "alias":
				value := req.EmailTemplate.Alias
				alias = &value

			case "subject":
				value := req.EmailTemplate.Subject
				subject = &value

			case "html":
				value := req.EmailTemplate.Html
				html = &value

			case "text":
				value := req.EmailTemplate.Text
				text = &value

			default:
				return nil, ErrInvalidUpdateMaskPath
			}
		}
	}

	result, err := command.Send(ctx, s.commandBus, command.UpdateTemplate{
		TemplateID: templateID,
		Alias:      alias,
		Subject:    subject,
		HTML:       html,
		Text:       text,
	})
	if err != nil {
		return nil, err
	}

	return &v1.EmailTemplate{
		Name:    resourcename.Sprint(emailTemplateResourcePattern, result.Template.ID),
		Alias:   result.Template.Alias,
		Subject: result.Template.Subject,
		Html:    result.Template.HTML,
		Text:    result.Template.Text,
	}, nil
}

func (s *Service) DeleteEmailTemplate(ctx context.Context, req *v1.DeleteEmailTemplateRequest) (*emptypb.Empty, error) {
	var templateID string
	if err := resourcename.Sscan(req.Name, emailTemplateResourcePattern, &templateID); err != nil {
		return nil, ErrInvalidEmailTemplateResourceName
	}

	if _, err := command.Send(ctx, s.commandBus, command.DeleteTemplate{
		TemplateID: templateID,
	}); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) SendEmail(ctx context.Context, req *v1.SendEmailRequest) (*v1.SendEmailResponse, error) {
	if content := req.Email.GetCustom(); content != nil {
		if _, err := command.Send(ctx, s.commandBus, command.SendEmail{
			From:           req.Email.From,
			To:             req.Email.To,
			CC:             req.Email.Cc,
			BCC:            req.Email.Bcc,
			Subject:        content.Subject,
			HTML:           content.Html,
			Text:           content.Text,
			IdempotencyKey: req.IdempotencyKey,
		}); err != nil {
			return nil, err
		}
	}
	if content := req.Email.GetTemplate(); content != nil {
		var templateID string
		if err := resourcename.Sscan(content.Name, emailTemplateResourcePattern, &templateID); err != nil {
			return nil, ErrInvalidEmailTemplateResourceName
		}

		if _, err := command.Send(ctx, s.commandBus, command.SendEmailFromTemplate{
			From:           req.Email.From,
			To:             req.Email.To,
			CC:             req.Email.Cc,
			BCC:            req.Email.Bcc,
			TemplateID:     templateID,
			Variables:      content.Variables,
			IdempotencyKey: req.IdempotencyKey,
		}); err != nil {
			return nil, err
		}
	}

	return &v1.SendEmailResponse{}, nil
}
