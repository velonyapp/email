package repo

import (
	"context"

	"github.com/velonyapp/email/internal/domain/entity"
	"github.com/velonyapp/email/internal/domain/vo"
)

type Template interface {
	FindByID(ctx context.Context, templateID vo.TemplateID) (*entity.Template, error)
	FindAll(ctx context.Context) ([]*entity.Template, error)

	Save(ctx context.Context, template *entity.Template) error
}
