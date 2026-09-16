package application

import (
	"github.com/velonyapp/notification/internal/application/usecase"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	usecase.NewSendEmailHandler,
)
