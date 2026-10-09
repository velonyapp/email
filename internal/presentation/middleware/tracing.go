package middleware

import (
	"github.com/go-kratos/kratos/contrib/otel/v3/tracing"
	"github.com/go-kratos/kratos/v3/middleware"
)

type Tracing middleware.Middleware

func NewTracing() Tracing {
	return Tracing(tracing.Server())
}

func (m Tracing) Middleware() middleware.Middleware {
	return middleware.Middleware(m)
}
