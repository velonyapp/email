package command

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

func NewBus(
	createTemplate CreateTemplateHandler,
	updateTemplate UpdateTemplateHandler,
	deleteTemplate DeleteTemplateHandler,
	sendEmail SendEmailHandler,
	sendEmailFromTemplate SendEmailFromTemplateHandler,
) *Bus {
	bus := &Bus{}

	Register(bus, createTemplate)
	Register(bus, updateTemplate)
	Register(bus, deleteTemplate)
	Register(bus, sendEmail)
	Register(bus, sendEmailFromTemplate)

	return bus
}

type Bus struct {
	handlers sync.Map
}

func Register[R any, C Command[R]](bus *Bus, handler Handler[C, R]) {
	t := typeOf[C]()

	bus.handlers.Store(t, erasedHandler(func(ctx context.Context, cmd any) (any, error) {
		return handler.Handle(ctx, cmd.(C))
	}))
}

func Send[R any, C Command[R]](ctx context.Context, bus *Bus, cmd C) (R, error) {
	var zero R

	raw, ok := bus.handlers.Load(typeOf[C]())
	if !ok {
		return zero, fmt.Errorf("no handler registered for %T", cmd)
	}

	result, err := raw.(erasedHandler)(ctx, cmd)
	if err != nil {
		return zero, err
	}

	return result.(R), nil
}

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

type erasedHandler func(context.Context, any) (any, error)
