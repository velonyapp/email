package command

import (
	"context"

	"github.com/velonyapp/email/internal/application/common"
	"github.com/velonyapp/email/internal/domain/repo"
	"github.com/velonyapp/email/internal/domain/vo"
)

type RenameTemplate struct {
	Alias    string
	NewAlias string
}

type RenameTemplateResult struct {
	Template common.TemplateResult
}

func (RenameTemplate) resultType() RenameTemplateResult {
	return RenameTemplateResult{}
}

type RenameTemplateHandler Handler[RenameTemplate, RenameTemplateResult]

type renameTemplateHandler struct {
	templateRepo repo.Template
}

func NewRenameTemplateHandler(
	templateRepo repo.Template,
) RenameTemplateHandler {
	return &renameTemplateHandler{
		templateRepo: templateRepo,
	}
}

func (h *renameTemplateHandler) Handle(
	ctx context.Context,
	cmd RenameTemplate,
) (RenameTemplateResult, error) {
	alias, err := vo.NewAlias(cmd.Alias)
	if err != nil {
		return RenameTemplateResult{}, err
	}

	template, err := h.templateRepo.FindByAlias(ctx, alias)
	if err != nil {
		return RenameTemplateResult{}, err
	}
	if template == nil {
		return RenameTemplateResult{}, common.ErrTemplateNotFound
	}

	newAlias, err := vo.NewAlias(cmd.NewAlias)
	if err != nil {
		return RenameTemplateResult{}, err
	}
	if !template.Alias().Equal(newAlias) {
		template.ChangeAlias(newAlias)
	}

	if err := h.templateRepo.Save(ctx, template); err != nil {
		return RenameTemplateResult{}, err
	}

	return RenameTemplateResult{
		Template: common.TemplateResult{
			Alias:   template.Alias().String(),
			Subject: template.Subject().String(),
			HTML:    template.HTML().String(),
			Text:    template.Text().String(),
		},
	}, nil
}
