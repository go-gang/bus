package direct

import (
	"context"
	"fmt"
	"github.com/go-gang/bus"
)

type queries map[[2]uint32]bus.QueryHandler

type QueryHandler interface {
	bus.QueryHandler

	key() [2]uint32
}

type genericQueryHandlerFunc[Query, Result any] func(ctx context.Context, query *Query, result *Result) error

type genericQueryHandler[Query, Result any] struct {
	typeHash [2]uint32
	handler  genericQueryHandlerFunc[Query, Result]
}

func NewAsker(handlers ...QueryHandler) bus.Asker {
	asker := make(queries)

	for _, handler := range handlers {
		asker[handler.key()] = handler
	}

	return asker
}

func (q queries) Ask(ctx context.Context, query, result any) error {
	if !isPointer(query) {
		return fmt.Errorf("query %T %w", query, bus.ErrNonPointer)
	}

	if !isPointer(result) {
		return fmt.Errorf("result %T %w", result, bus.ErrNonPointer)
	}

	key := [2]uint32{
		typeHash(query),
		typeHash(result),
	}

	if handler, ok := q[key]; ok {
		return handler.Handle(ctx, query, result)
	}

	return fmt.Errorf("query (%T, %T) %w", query, result, bus.ErrHandlerNotFound)
}

func NewQueryHandler[Query, Result any](handler genericQueryHandlerFunc[Query, Result]) QueryHandler {
	zeroValueQuery := new(Query)
	zeroValueResult := new(Result)

	return &genericQueryHandler[Query, Result]{
		typeHash: [2]uint32{
			typeHash(zeroValueQuery),
			typeHash(zeroValueResult),
		},
		handler: handler,
	}
}

func (h *genericQueryHandler[Query, Result]) key() [2]uint32 {
	return h.typeHash
}

func (h *genericQueryHandler[Query, Result]) Handle(ctx context.Context, query, result any) error {
	return h.handler(ctx, query.(*Query), result.(*Result))
}
