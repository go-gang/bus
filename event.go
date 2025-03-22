package bus

import "context"

type Publisher interface {
	Publish(ctx context.Context, event any) error
}

type EventHandler interface {
	Handle(ctx context.Context, event any) error
}
