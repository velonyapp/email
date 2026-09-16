package presentation

import (
	"github.com/velonyapp/notification/internal/presentation/api"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	api.NewService,
)
