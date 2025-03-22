package direct_test

import (
	"context"
	"errors"
	"github.com/go-gang/bus"
	"github.com/go-gang/bus/direct"
	"github.com/go-gang/bus/internal/assert"
	"testing"
)

type query struct {
	value int
	err   error
}

type result struct {
	value int
}

func TestAsker(t *testing.T) {
	handler := direct.NewQueryHandler(func(ctx context.Context, q *query, r *result) error {
		r.value = q.value + 1
		return q.err
	})
	a := direct.NewAsker(handler)
	ctx := context.Background()

	q := &query{
		value: 1,
		err:   nil,
	}
	r := &result{
		value: 0,
	}
	err := errors.New("test")

	assert.NoError(t, a.Ask(ctx, q, r))
	assert.Equal(t, 2, r.value)
	assert.ErrorIs(t, a.Ask(ctx, &query{err: err}, &result{}), err)
	assert.ErrorIs(t, a.Ask(ctx, query{}, &result{}), bus.ErrNonPointer)
	assert.ErrorIs(t, a.Ask(ctx, &query{}, result{}), bus.ErrNonPointer)
	assert.ErrorIs(t, a.Ask(ctx, nil, &result{}), bus.ErrNonPointer)
	assert.ErrorIs(t, a.Ask(ctx, &query{}, nil), bus.ErrNonPointer)
	assert.ErrorIs(t, a.Ask(ctx, &result{}, &query{}), bus.ErrHandlerNotFound)
}

// goos: linux
// goarch: amd64
// cpu: AMD Ryzen 5 6600HS Creator Edition
// BenchmarkAsker-12    	98966337	        11.38 ns/op
func BenchmarkAsker(b *testing.B) {
	handler := direct.NewQueryHandler(func(ctx context.Context, q *query, r *result) error {
		r.value = q.value + 1
		return q.err
	})
	a := direct.NewAsker(handler)
	ctx := context.Background()
	q := &query{}
	r := &result{}

	for i := 0; i < b.N; i++ {
		_ = a.Ask(ctx, q, r)
	}
}

// goos: linux
// goarch: amd64
// cpu: AMD Ryzen 5 6600HS Creator Edition
// BenchmarkAskerParallel-12    	121281715	         9.693 ns/op
func BenchmarkAskerParallel(b *testing.B) {
	handler := direct.NewQueryHandler(func(ctx context.Context, q *query, r *result) error {
		r.value = q.value + 1
		return q.err
	})
	a := direct.NewAsker(handler)
	ctx := context.Background()
	q := &query{}
	r := &result{}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = a.Ask(ctx, q, r)
		}
	})
}
