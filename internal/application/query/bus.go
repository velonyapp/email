package query

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

func NewBus(
	getTemplate GetTemplateHandler,
	listTemplates ListTemplatesHandler,
) *Bus {
	bus := &Bus{}

	Register(bus, getTemplate)
	Register(bus, listTemplates)

	return bus
}

type Bus struct {
	handlers sync.Map
}

func Register[R any, Q Query[R]](bus *Bus, handler Handler[Q, R]) {
	t := typeOf[Q]()

	bus.handlers.Store(t, erasedHandler(func(ctx context.Context, query any) (any, error) {
		return handler.Handle(ctx, query.(Q))
	}))
}

func Send[R any, Q Query[R]](ctx context.Context, bus *Bus, query Q) (R, error) {
	var zero R

	raw, ok := bus.handlers.Load(typeOf[Q]())
	if !ok {
		return zero, fmt.Errorf("no handler registered for %T", query)
	}

	result, err := raw.(erasedHandler)(ctx, query)
	if err != nil {
		return zero, err
	}

	return result.(R), nil
}

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

type erasedHandler func(context.Context, any) (any, error)
