package command

import (
	"context"

	"github.com/velonyapp/email/internal/application/common"
	"github.com/velonyapp/email/internal/application/port"
	"github.com/velonyapp/email/internal/domain/entity"
	"github.com/velonyapp/email/internal/domain/repo"
	"github.com/velonyapp/email/internal/domain/service"
	"github.com/velonyapp/email/internal/domain/vo"
)

type CreateTemplate struct {
	TemplateID string
	Subject    string
	HTML       string
	Text       string
}

type CreateTemplateResult struct {
	Template common.TemplateResult
}

func (CreateTemplate) resultType() CreateTemplateResult {
	return CreateTemplateResult{}
}

type CreateTemplateHandler Handler[CreateTemplate, CreateTemplateResult]

type createTemplateHandler struct {
	templateRepo     repo.Template
	templateIDPolicy *service.TemplateIDPolicy
	unitOfWork       port.UnitOfWork
}

func NewCreateTemplateHandler(
	templateRepo repo.Template,
	templateIDPolicy *service.TemplateIDPolicy,
	unitOfWork port.UnitOfWork,
) CreateTemplateHandler {
	return &createTemplateHandler{
		templateRepo:     templateRepo,
		templateIDPolicy: templateIDPolicy,
		unitOfWork:       unitOfWork,
	}
}

func (h *createTemplateHandler) Handle(
	ctx context.Context,
	cmd CreateTemplate,
) (CreateTemplateResult, error) {
	templateID, err := vo.NewTemplateID(cmd.TemplateID)
	if err != nil {
		return CreateTemplateResult{}, err
	}

	var template *entity.Template

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		if err := h.templateIDPolicy.CanUse(ctx, templateID); err != nil {
			return err
		}

		subject := vo.NewSubject(cmd.Subject)
		html := vo.NewHTML(cmd.HTML)
		text := vo.NewText(cmd.Text)

		template, err = entity.NewTemplate(
			templateID,
			subject,
			html,
			text,
		)
		if err != nil {
			return err
		}

		return h.templateRepo.Save(ctx, template)
	}); err != nil {
		return CreateTemplateResult{}, err
	}

	return CreateTemplateResult{
		Template: common.TemplateResult{
			ID:      template.ID().String(),
			Subject: template.Subject().String(),
			HTML:    template.HTML().String(),
			Text:    template.Text().String(),
		},
	}, nil
}
