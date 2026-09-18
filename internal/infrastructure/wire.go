package infrastructure

import (
	"github.com/velonyapp/email/internal/infrastructure/email/resend"
	"github.com/velonyapp/email/internal/infrastructure/observability"
	"github.com/velonyapp/email/internal/infrastructure/transport"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	resend.NewClient,
	resend.NewSender,
	transport.NewGRPCServer,
	transport.NewHTTPServer,
	transport.NewRabbitMQServer,
	transport.NewTracesMiddleware,
	transport.NewMetricsMiddleware,
	transport.NewValidationMiddleware,
	observability.NewServerMetrics,
	observability.NewOpenTelemetry,
)
