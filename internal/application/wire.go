package application

import (
	"github.com/velonyapp/email/internal/application/usecase"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	usecase.NewSendEmailHandler,
)
