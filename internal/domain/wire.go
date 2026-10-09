package domain

import (
	"github.com/velonyapp/email/internal/domain/service"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	service.NewTemplateIDPolicy,
)
