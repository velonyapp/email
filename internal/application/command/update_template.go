package command

import (
	"context"

	"github.com/velonyapp/email/internal/application/common"
	"github.com/velonyapp/email/internal/domain/entity"
	"github.com/velonyapp/email/internal/domain/repo"
	"github.com/velonyapp/email/internal/domain/vo"
)

type UpdateTemplate struct {
	TemplateID string
	Alias      *string
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
}

func NewUpdateTemplateHandler(
	templateRepo repo.Template,
) UpdateTemplateHandler {
	return &updateTemplateHandler{
		templateRepo: templateRepo,
	}
}

func (h *updateTemplateHandler) Handle(
	ctx context.Context,
	cmd UpdateTemplate,
) (UpdateTemplateResult, error) {
	templateID, _ := vo.NewTemplateID(cmd.TemplateID)

	var template *entity.Template

	template, err := h.templateRepo.FindByID(ctx, templateID)
	if err != nil {
		return UpdateTemplateResult{}, err
	}
	if template == nil {
		return UpdateTemplateResult{}, common.ErrTemplateNotFound
	}

	if cmd.Alias != nil {
		alias, err := vo.NewAlias(*cmd.Alias)
		if err != nil {
			return UpdateTemplateResult{}, err
		}
		if !template.Alias().Equal(alias) {
			template.ChangeAlias(alias)
		}
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
			return UpdateTemplateResult{}, err
		}
	case cmd.HTML != nil:
		if err := template.ChangeHTML(vo.NewHTML(*cmd.HTML)); err != nil {
			return UpdateTemplateResult{}, err
		}
	case cmd.Text != nil:
		if err := template.ChangeText(vo.NewText(*cmd.Text)); err != nil {
			return UpdateTemplateResult{}, err
		}
	}

	if err := h.templateRepo.Save(ctx, template); err != nil {
		return UpdateTemplateResult{}, err
	}

	return UpdateTemplateResult{
		Template: common.TemplateResult{
			ID:      template.ID().String(),
			Alias:   template.Alias().String(),
			Subject: template.Subject().String(),
			HTML:    template.HTML().String(),
			Text:    template.Text().String(),
		},
	}, nil
}
