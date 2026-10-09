package middleware

import (
	"github.com/go-kratos/kratos/contrib/otel/v3/metrics"
	"github.com/go-kratos/kratos/v3/middleware"
	"go.opentelemetry.io/otel"
)

const instrumentationName = "github.com/velonyapp/email/internal/presentation/middleware"

type Metrics middleware.Middleware

func NewMetrics() (Metrics, error) {
	meter := otel.Meter(instrumentationName)

	requests, err := metrics.DefaultRequestsCounter(
		meter,
		metrics.DefaultServerRequestsCounterName,
	)
	if err != nil {
		return nil, err
	}

	seconds, err := metrics.DefaultSecondsHistogram(
		meter,
		metrics.DefaultServerSecondsHistogramName,
	)
	if err != nil {
		return nil, err
	}

	return Metrics(
		metrics.Server(
			metrics.WithSeconds(seconds),
			metrics.WithRequests(requests),
		),
	), nil
}

func (m Metrics) Middleware() middleware.Middleware {
	return middleware.Middleware(m)
}
