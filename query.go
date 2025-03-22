package bus

import "context"

type Asker interface {
	Ask(ctx context.Context, query, result any) error
}

type QueryHandler interface {
	Handle(ctx context.Context, query, result any) error
}
