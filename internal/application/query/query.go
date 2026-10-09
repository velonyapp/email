package query

import "context"

type Query[R any] interface {
	resultType() R
}

type Handler[Q, R any] interface {
	Handle(context.Context, Q) (R, error)
}
