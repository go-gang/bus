package direct

import (
	"context"
	"fmt"
	"github.com/go-gang/bus"
)

type commands map[uint32]bus.CommandHandler

type CommandHandler interface {
	bus.CommandHandler

	key() uint32
}

type genericCommandHandler[Command any] struct {
	typeHash uint32
	handler  func(ctx context.Context, command *Command) error
}

func NewDispatcher(handlers ...CommandHandler) bus.Dispatcher {
	dispatcher := make(commands)

	for _, h := range handlers {
		dispatcher[h.key()] = h
	}

	return dispatcher
}

func (c commands) Dispatch(ctx context.Context, command any) error {
	if !isPointer(command) {
		return fmt.Errorf("command %T %w", command, bus.ErrNonPointer)
	}

	key := typeHash(command)

	if handler, ok := c[key]; ok {
		return handler.Handle(ctx, command)
	}

	return fmt.Errorf("command %T %w", command, bus.ErrHandlerNotFound)
}

func NewCommandHandler[Command any](handler func(ctx context.Context, command *Command) error) CommandHandler {
	zeroValue := new(Command)

	return &genericCommandHandler[Command]{
		typeHash: typeHash(zeroValue),
		handler:  handler,
	}
}

func (h *genericCommandHandler[Command]) key() uint32 {
	return h.typeHash
}

func (h *genericCommandHandler[Command]) Handle(ctx context.Context, command any) error {
	return h.handler(ctx, command.(*Command))
}
