package presentation

import (
	"github.com/velonyapp/email/internal/presentation/api"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	api.NewService,
)
