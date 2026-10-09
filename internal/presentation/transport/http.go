package transport

import (
	v1 "github.com/velonyapp/email/gen/api/v1"
	"github.com/velonyapp/email/internal/conf"
	"github.com/velonyapp/email/internal/presentation/api"
	"github.com/velonyapp/email/internal/presentation/middleware"

	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/http"
)

func NewHTTPServer(
	c *conf.Transport,
	service *api.Service,
	tracing middleware.Tracing,
	metrics middleware.Metrics,
	errorMapper middleware.ErrorMapper,
	validation middleware.Validation,
) *http.Server {
	opts := []http.ServerOption{
		http.Address(c.Http.Address),
		http.Middleware(
			recovery.Recovery(),
			tracing.Middleware(),
			metrics.Middleware(),
			errorMapper.Middleware(),
			validation.Middleware(),
		),
	}

	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)

	v1.RegisterEmailServiceHTTPServer(srv, service)

	return srv
}
