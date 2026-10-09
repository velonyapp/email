package command

import "context"

type Command[R any] interface {
	resultType() R
}

type Handler[C, R any] interface {
	Handle(context.Context, C) (R, error)
}
