package command

import (
	"context"

	"github.com/velonyapp/email/internal/application/common"
	"github.com/velonyapp/email/internal/domain/entity"
	"github.com/velonyapp/email/internal/domain/repo"
	"github.com/velonyapp/email/internal/domain/vo"
)

type CreateTemplate struct {
	Alias   string
	Subject string
	HTML    string
	Text    string
}

type CreateTemplateResult struct {
	Template common.TemplateResult
}

func (CreateTemplate) resultType() CreateTemplateResult {
	return CreateTemplateResult{}
}

type CreateTemplateHandler Handler[CreateTemplate, CreateTemplateResult]

type createTemplateHandler struct {
	templateRepo repo.Template
}

func NewCreateTemplateHandler(
	templateRepo repo.Template,
) CreateTemplateHandler {
	return &createTemplateHandler{
		templateRepo: templateRepo,
	}
}

func (h *createTemplateHandler) Handle(
	ctx context.Context,
	cmd CreateTemplate,
) (CreateTemplateResult, error) {
	alias, err := vo.NewAlias(cmd.Alias)
	if err != nil {
		return CreateTemplateResult{}, err
	}

	subject := vo.NewSubject(cmd.Subject)
	html := vo.NewHTML(cmd.HTML)
	text := vo.NewText(cmd.Text)

	template, err := entity.NewTemplate(
		alias,
		subject,
		html,
		text,
	)
	if err != nil {
		return CreateTemplateResult{}, err
	}

	if err := h.templateRepo.Save(ctx, template); err != nil {
		return CreateTemplateResult{}, err
	}

	return CreateTemplateResult{
		Template: common.TemplateResult{
			ID:      template.ID().String(),
			Alias:   template.Alias().String(),
			Subject: template.Subject().String(),
			HTML:    template.HTML().String(),
			Text:    template.Text().String(),
		},
	}, nil
}
