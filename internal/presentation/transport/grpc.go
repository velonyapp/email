package transport

import (
	v1 "github.com/velonyapp/email/gen/api/v1"
	"github.com/velonyapp/email/internal/conf"
	"github.com/velonyapp/email/internal/presentation/api"
	"github.com/velonyapp/email/internal/presentation/middleware"

	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/grpc"
)

func NewGRPCServer(
	c *conf.Transport,
	service *api.Service,
	tracing middleware.Tracing,
	metrics middleware.Metrics,
	errorMapper middleware.ErrorMapper,
	validation middleware.Validation,
) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.Address(c.Grpc.Address),
		grpc.Middleware(
			recovery.Recovery(),
			tracing.Middleware(),
			metrics.Middleware(),
			errorMapper.Middleware(),
			validation.Middleware(),
		),
	}

	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	srv := grpc.NewServer(opts...)

	v1.RegisterEmailServiceServer(srv, service)

	return srv
}
