package service

import (
	"context"
	"errors"

	"github.com/velonyapp/email/internal/domain/repo"
	"github.com/velonyapp/email/internal/domain/vo"
)

var (
	ErrAliasAlreadyExists = errors.New("alias already exists")
)

type AliasPolicy struct {
	templateRepo repo.Template
}

func NewAliasPolicy(
	templateRepo repo.Template,
) *AliasPolicy {
	return &AliasPolicy{
		templateRepo: templateRepo,
	}
}

func (p *AliasPolicy) CanUse(ctx context.Context, alias vo.Alias) error {
	template, err := p.templateRepo.FindByAlias(ctx, alias)
	if err != nil {
		return err
	}
	if template != nil {
		return ErrAliasAlreadyExists
	}

	return nil
}
