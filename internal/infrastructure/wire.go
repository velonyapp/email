package infrastructure

import (
	"github.com/velonyapp/notification/internal/infrastructure/email/resend"
	"github.com/velonyapp/notification/internal/infrastructure/observability"
	"github.com/velonyapp/notification/internal/infrastructure/transport"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	resend.NewClient,
	resend.NewSender,
	transport.NewGRPCServer,
	transport.NewHTTPServer,
	transport.NewTracesMiddleware,
	transport.NewMetricsMiddleware,
	transport.NewValidationMiddleware,
	observability.NewServerMetrics,
	observability.NewOpenTelemetry,
)
