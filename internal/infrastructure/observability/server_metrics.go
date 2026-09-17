package observability

import (
	"github.com/go-kratos/kratos/contrib/otel/v3/metrics"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const instrumentationName = "github.com/velonyapp/email/internal/infrastructure/observability"

type ServerMetrics struct {
	Requests metric.Int64Counter
	Seconds  metric.Float64Histogram
}

func NewServerMetrics() (*ServerMetrics, error) {
	meter := otel.Meter(instrumentationName)

	Requests, err := metrics.DefaultRequestsCounter(
		meter,
		metrics.DefaultServerRequestsCounterName,
	)
	if err != nil {
		return nil, err
	}

	Seconds, err := metrics.DefaultSecondsHistogram(
		meter,
		metrics.DefaultServerSecondsHistogramName,
	)
	if err != nil {
		return nil, err
	}

	return &ServerMetrics{
		Requests: Requests,
		Seconds:  Seconds,
	}, nil
}
