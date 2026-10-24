package repo

import (
	"context"
	"errors"

	"github.com/velonyapp/email/internal/domain/entity"
	"github.com/velonyapp/email/internal/domain/vo"
)

var (
	ErrAliasAlreadyExists      = errors.New("alias already exists")
	ErrTemplateVersionConflict = errors.New("template version conflict")
)

type Template interface {
	FindByID(ctx context.Context, templateID vo.TemplateID) (*entity.Template, error)
	FindByAlias(ctx context.Context, alias vo.Alias) (*entity.Template, error)
	FindAll(ctx context.Context) ([]*entity.Template, error)

	Save(ctx context.Context, template *entity.Template) error
}
