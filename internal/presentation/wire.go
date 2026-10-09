package presentation

import (
	"github.com/velonyapp/email/internal/presentation/api"
	"github.com/velonyapp/email/internal/presentation/middleware"
	"github.com/velonyapp/email/internal/presentation/transport"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	api.NewService,
	transport.NewGRPCServer,
	transport.NewHTTPServer,
	transport.NewRabbitMQConsumer,
	middleware.NewTracing,
	middleware.NewMetrics,
	middleware.NewErrorMapper,
	middleware.NewValidation,
)
