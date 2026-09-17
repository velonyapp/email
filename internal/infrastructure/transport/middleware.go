package transport

import (
	"context"

	"github.com/velonyapp/email/internal/infrastructure/observability"

	"github.com/go-kratos/kratos/contrib/otel/v3/metrics"
	"github.com/go-kratos/kratos/contrib/otel/v3/tracing"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/middleware/validate"
	"github.com/go-kratos/kratos/v3/transport"
	"go.einride.tech/aip/fieldbehavior"
	"google.golang.org/protobuf/proto"
)

func optionalAuthentication(auth middleware.Middleware) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		authenticated := auth(next)

		return func(ctx context.Context, req any) (any, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return next(ctx, req)
			}

			if tr.RequestHeader().Get("Authorization") == "" {
				return next(ctx, req)
			}

			return authenticated(ctx, req)
		}
	}
}

type MetricsMiddleware middleware.Middleware

func NewMetricsMiddleware(serverMetrics *observability.ServerMetrics) MetricsMiddleware {
	return MetricsMiddleware(
		metrics.Server(
			metrics.WithSeconds(serverMetrics.Seconds),
			metrics.WithRequests(serverMetrics.Requests),
		),
	)
}

type TracesMiddleware middleware.Middleware

func NewTracesMiddleware() TracesMiddleware {
	return TracesMiddleware(tracing.Server())
}

type ValidationMiddleware middleware.Middleware

func NewValidationMiddleware() ValidationMiddleware {
	return ValidationMiddleware(
		validate.Validator(func(req any) error {
			message, ok := req.(proto.Message)
			if !ok {
				return nil
			}

			return fieldbehavior.ValidateRequiredFields(message)
		}),
	)
}
