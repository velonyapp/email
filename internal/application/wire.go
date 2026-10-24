package application

import (
	"github.com/velonyapp/email/internal/application/command"
	"github.com/velonyapp/email/internal/application/query"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	command.NewBus,
	command.NewCreateTemplateHandler,
	command.NewUpdateTemplateHandler,
	command.NewDeleteTemplateHandler,
	command.NewRenameTemplateHandler,
	command.NewSendEmailHandler,
	query.NewBus,
	query.NewGetTemplateHandler,
	query.NewListTemplatesHandler,
)
