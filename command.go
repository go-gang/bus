package bus

import "context"

type Dispatcher interface {
	Dispatch(ctx context.Context, command any) error
}

type CommandHandler interface {
	Handle(ctx context.Context, command any) error
}
