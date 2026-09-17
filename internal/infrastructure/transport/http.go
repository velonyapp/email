package transport

import (
	v1 "github.com/velonyapp/email/gen/api/v1"
	"github.com/velonyapp/email/internal/conf"
	"github.com/velonyapp/email/internal/presentation/api"

	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/http"
)

func NewHTTPServer(
	c *conf.Transport,
	service *api.Service,
	tracingMiddleware TracesMiddleware,
	metricsMiddleware MetricsMiddleware,
	validationMiddleware ValidationMiddleware,
) *http.Server {
	opts := []http.ServerOption{
		http.Address(c.Http.Address),
		http.Middleware(
			recovery.Recovery(),
			middleware.Middleware(tracingMiddleware),
			middleware.Middleware(metricsMiddleware),
			middleware.Middleware(validationMiddleware),
		),
	}

	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)

	v1.RegisterEmailServiceHTTPServer(srv, service)

	return srv
}
