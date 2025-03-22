package direct

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-gang/bus"
)

type events map[uint32]bus.EventHandler

type EventHandler interface {
	bus.EventHandler

	key() uint32
}

type genericEventHandlerFunc[Event any] func(ctx context.Context, event *Event) error

type genericEventHandler[Event any] struct {
	typeHash uint32
	handler  genericEventHandlerFunc[Event]
}

func NewPublisher(listeners ...EventHandler) bus.Publisher {
	publisher := make(events)

	for _, listener := range listeners {
		publisher[listener.key()] = listener
	}

	return publisher
}

func (e events) Publish(ctx context.Context, event any) error {
	if !isPointer(event) {
		return fmt.Errorf("event %T %w", event, bus.ErrNonPointer)
	}

	key := typeHash(event)

	if handler, ok := e[key]; ok {
		return handler.Handle(ctx, event)
	}

	return nil
}

func NewEventHandler[Event any](handler genericEventHandlerFunc[Event]) EventHandler {
	zeroValue := new(Event)

	return &genericEventHandler[Event]{
		typeHash: typeHash(zeroValue),
		handler:  handler,
	}
}

func NewEventGroupHandler[Event any](handlers ...genericEventHandlerFunc[Event]) EventHandler {
	return NewEventHandler(func(ctx context.Context, event *Event) error {
		errs := make([]error, 0, len(handlers))

		for _, handler := range handlers {
			errs = append(errs, handler(ctx, event))
		}

		return errors.Join(errs...)
	})
}

func (h *genericEventHandler[Event]) key() uint32 {
	return h.typeHash
}

func (h *genericEventHandler[Event]) Handle(ctx context.Context, event any) error {
	return h.handler(ctx, event.(*Event))
}
