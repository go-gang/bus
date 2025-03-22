package direct_test

import (
	"context"
	"errors"
	"github.com/go-gang/bus"
	"github.com/go-gang/bus/direct"
	"github.com/go-gang/bus/internal/assert"
	"testing"
)

type command struct {
	value int
	err   error
}

func TestDispatcher(t *testing.T) {
	handler := direct.NewCommandHandler(func(ctx context.Context, c *command) error {
		c.value++
		return c.err
	})
	d := direct.NewDispatcher(handler)
	ctx := context.Background()

	cmd := &command{
		value: 0,
		err:   nil,
	}
	err := errors.New("test")

	assert.NoError(t, d.Dispatch(ctx, cmd))
	assert.Equal(t, 1, cmd.value)
	assert.ErrorIs(t, d.Dispatch(ctx, &command{err: err}), err)
	assert.ErrorIs(t, d.Dispatch(ctx, command{}), bus.ErrNonPointer)
	assert.ErrorIs(t, d.Dispatch(ctx, nil), bus.ErrNonPointer)
	assert.ErrorIs(t, d.Dispatch(ctx, &err), bus.ErrHandlerNotFound)
}

// goos: linux
// goarch: amd64
// cpu: AMD Ryzen 5 6600HS Creator Edition
// BenchmarkDispatcher-12    	148222830	         7.973 ns/op
func BenchmarkDispatcher(b *testing.B) {
	handler := direct.NewCommandHandler(func(ctx context.Context, c *command) error {
		c.value += 1
		return nil
	})
	d := direct.NewDispatcher(handler)
	ctx := context.Background()
	cmd := &command{}

	for i := 0; i < b.N; i++ {
		_ = d.Dispatch(ctx, cmd)
	}
}

// goos: linux
// goarch: amd64
// cpu: AMD Ryzen 5 6600HS Creator Edition
// BenchmarkDispatcherParallel-12    	749080788	         1.577 ns/op
func BenchmarkDispatcherParallel(b *testing.B) {
	handler := direct.NewCommandHandler(func(ctx context.Context, c *command) error {
		c.value += 0
		return nil
	})
	d := direct.NewDispatcher(handler)
	ctx := context.Background()
	cmd := &command{}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = d.Dispatch(ctx, cmd)
		}
	})
}
