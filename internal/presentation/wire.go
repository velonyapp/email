package presentation

import (
	"github.com/velonyapp/email/internal/presentation/middleware"
	"github.com/velonyapp/email/internal/presentation/service"
	"github.com/velonyapp/email/internal/presentation/transport"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	service.NewAPI,
	service.NewJob,
	transport.NewGRPCServer,
	transport.NewHTTPServer,
	transport.NewRabbitMQConsumer,
	middleware.NewTracing,
	middleware.NewMetrics,
	middleware.NewErrorMapper,
	middleware.NewValidation,
)
