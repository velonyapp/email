package query

import (
	"context"

	"github.com/velonyapp/email/internal/application/common"
	"github.com/velonyapp/email/internal/domain/repo"
	"github.com/velonyapp/email/internal/domain/vo"
)

type GetTemplate struct {
	Alias string
}

type GetTemplateResult struct {
	Template common.TemplateResult
}

func (GetTemplate) resultType() GetTemplateResult {
	return GetTemplateResult{}
}

type GetTemplateHandler Handler[GetTemplate, GetTemplateResult]

type getTemplateHandler struct {
	templateRepo repo.Template
}

func NewGetTemplateHandler(
	templateRepo repo.Template,
) GetTemplateHandler {
	return &getTemplateHandler{
		templateRepo: templateRepo,
	}
}

func (h *getTemplateHandler) Handle(
	ctx context.Context,
	qry GetTemplate,
) (GetTemplateResult, error) {
	alias, _ := vo.NewAlias(qry.Alias)

	template, err := h.templateRepo.FindByAlias(ctx, alias)
	if err != nil {
		return GetTemplateResult{}, err
	}
	if template == nil {
		return GetTemplateResult{}, common.ErrTemplateNotFound
	}

	return GetTemplateResult{
		Template: common.TemplateResult{
			Alias:   template.Alias().String(),
			Subject: template.Subject().String(),
			HTML:    template.HTML().String(),
			Text:    template.Text().String(),
		},
	}, nil
}
