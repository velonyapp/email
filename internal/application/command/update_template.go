package command

import (
	"context"

	"github.com/velonyapp/email/internal/application/common"
	"github.com/velonyapp/email/internal/application/port"
	"github.com/velonyapp/email/internal/domain/entity"
	"github.com/velonyapp/email/internal/domain/repo"
	"github.com/velonyapp/email/internal/domain/vo"
)

type UpdateTemplate struct {
	TemplateID string
	Subject    *string
	HTML       *string
	Text       *string
}

type UpdateTemplateResult struct {
	Template common.TemplateResult
}

func (UpdateTemplate) resultType() UpdateTemplateResult {
	return UpdateTemplateResult{}
}

type UpdateTemplateHandler Handler[UpdateTemplate, UpdateTemplateResult]

type updateTemplateHandler struct {
	templateRepo repo.Template
	unitOfWork   port.UnitOfWork
}

func NewUpdateTemplateHandler(
	templateRepo repo.Template,
	unitOfWork port.UnitOfWork,
) UpdateTemplateHandler {
	return &updateTemplateHandler{
		templateRepo: templateRepo,
		unitOfWork:   unitOfWork,
	}
}

func (h *updateTemplateHandler) Handle(
	ctx context.Context,
	cmd UpdateTemplate,
) (UpdateTemplateResult, error) {
	templateID, _ := vo.NewTemplateID(cmd.TemplateID)

	var template *entity.Template

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		var err error

		template, err = h.templateRepo.FindByID(ctx, templateID)
		if err != nil {
			return err
		}
		if template == nil {
			return common.ErrTemplateNotFound
		}

		if cmd.Subject != nil {
			template.ChangeSubject(vo.NewSubject(*cmd.Subject))
		}
		switch {
		case cmd.HTML != nil && cmd.Text != nil:
			if err := template.ChangeHTMLAndText(
				vo.NewHTML(*cmd.HTML),
				vo.NewText(*cmd.Text),
			); err != nil {
				return err
			}
		case cmd.HTML != nil:
			if err := template.ChangeHTML(vo.NewHTML(*cmd.HTML)); err != nil {
				return err
			}
		case cmd.Text != nil:
			if err := template.ChangeText(vo.NewText(*cmd.Text)); err != nil {
				return err
			}
		}

		return h.templateRepo.Save(ctx, template)
	}); err != nil {
		return UpdateTemplateResult{}, err
	}

	return UpdateTemplateResult{
		Template: common.TemplateResult{
			ID:      template.ID().String(),
			Subject: template.Subject().String(),
			HTML:    template.HTML().String(),
			Text:    template.Text().String(),
		},
	}, nil
}
