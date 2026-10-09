package service

import (
	"context"
	"errors"

	"github.com/velonyapp/email/internal/domain/repo"
	"github.com/velonyapp/email/internal/domain/vo"
)

var (
	ErrTemplateIDAlreadyExists = errors.New("template ID already exists")
)

type TemplateIDPolicy struct {
	templateRepo repo.Template
}

func NewTemplateIDPolicy(
	templateRepo repo.Template,
) *TemplateIDPolicy {
	return &TemplateIDPolicy{
		templateRepo: templateRepo,
	}
}

func (p *TemplateIDPolicy) CanUse(ctx context.Context, templateID vo.TemplateID) error {
	template, err := p.templateRepo.FindByID(ctx, templateID)
	if err != nil {
		return err
	}
	if template != nil {
		return ErrTemplateIDAlreadyExists
	}

	return nil
}
