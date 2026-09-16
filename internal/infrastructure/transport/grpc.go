package transport

import (
	v1 "github.com/velonyapp/notification/gen/api/v1"
	"github.com/velonyapp/notification/internal/conf"
	"github.com/velonyapp/notification/internal/presentation/api"

	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/grpc"
)

func NewGRPCServer(
	c *conf.Transport,
	service *api.Service,
	tracingMiddleware TracesMiddleware,
	metricsMiddleware MetricsMiddleware,
	validationMiddleware ValidationMiddleware,
) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.Address(c.Grpc.Address),
		grpc.Middleware(
			recovery.Recovery(),
			middleware.Middleware(tracingMiddleware),
			middleware.Middleware(metricsMiddleware),
			middleware.Middleware(validationMiddleware),
		),
	}

	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	srv := grpc.NewServer(opts...)

	v1.RegisterNotificationServiceServer(srv, service)

	return srv
}
