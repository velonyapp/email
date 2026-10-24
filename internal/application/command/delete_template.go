package command

import (
	"context"

	"github.com/velonyapp/email/internal/application/common"
	"github.com/velonyapp/email/internal/application/port"
	"github.com/velonyapp/email/internal/domain/repo"
	"github.com/velonyapp/email/internal/domain/vo"
)

type DeleteTemplate struct {
	Alias string
}

type DeleteTemplateResult struct{}

func (DeleteTemplate) resultType() DeleteTemplateResult {
	return DeleteTemplateResult{}
}

type DeleteTemplateHandler Handler[DeleteTemplate, DeleteTemplateResult]

type deleteTemplateHandler struct {
	templateRepo repo.Template
	unitOfWork   port.UnitOfWork
}

func NewDeleteTemplateHandler(
	templateRepo repo.Template,
	unitOfWork port.UnitOfWork,
) DeleteTemplateHandler {
	return &deleteTemplateHandler{
		templateRepo: templateRepo,
		unitOfWork:   unitOfWork,
	}
}

func (h *deleteTemplateHandler) Handle(
	ctx context.Context,
	cmd DeleteTemplate,
) (DeleteTemplateResult, error) {
	alias, _ := vo.NewAlias(cmd.Alias)

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		template, err := h.templateRepo.FindByAlias(ctx, alias)
		if err != nil {
			return err
		}
		if template == nil {
			return common.ErrTemplateNotFound
		}

		template.Delete()

		return h.templateRepo.Save(ctx, template)
	}); err != nil {
		return DeleteTemplateResult{}, err
	}

	return DeleteTemplateResult{}, nil
}
