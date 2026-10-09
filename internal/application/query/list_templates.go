package query

import (
	"context"

	"github.com/velonyapp/email/internal/application/common"
	"github.com/velonyapp/email/internal/domain/repo"
)

type ListTemplates struct{}

type ListTemplatesResult struct {
	Templates []common.TemplateResult
}

func (ListTemplates) resultType() ListTemplatesResult {
	return ListTemplatesResult{}
}

type ListTemplatesHandler Handler[ListTemplates, ListTemplatesResult]

type listTemplatesHandler struct {
	templateRepo repo.Template
}

func NewListTemplatesHandler(
	templateRepo repo.Template,
) ListTemplatesHandler {
	return &listTemplatesHandler{
		templateRepo: templateRepo,
	}
}

func (h *listTemplatesHandler) Handle(
	ctx context.Context,
	qry ListTemplates,
) (ListTemplatesResult, error) {
	templates, err := h.templateRepo.FindAll(ctx)
	if err != nil {
		return ListTemplatesResult{}, err
	}

	results := make([]common.TemplateResult, 0, len(templates))

	for _, template := range templates {
		results = append(results, common.TemplateResult{
			ID:      template.ID().String(),
			Alias:   template.Alias().String(),
			Subject: template.Subject().String(),
			HTML:    template.HTML().String(),
			Text:    template.Text().String(),
		})
	}

	return ListTemplatesResult{
		Templates: results,
	}, nil
}
