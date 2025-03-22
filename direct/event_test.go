package direct_test

import (
	"context"
	"errors"
	"github.com/go-gang/bus"
	"github.com/go-gang/bus/direct"
	"github.com/go-gang/bus/internal/assert"
	"sync/atomic"
	"testing"
)

type event struct {
	value int
	err   error
}

func TestPublisher(t *testing.T) {
	handler := direct.NewEventGroupHandler(
		func(ctx context.Context, event *event) error {
			event.value += 1

			return event.err
		},

		func(ctx context.Context, event *event) error {
			event.value += 2

			return event.err
		},
	)

	p := direct.NewPublisher(handler)
	ctx := context.Background()

	e := &event{
		value: 1,
		err:   nil,
	}
	err := errors.New("test error")

	assert.NoError(t, p.Publish(ctx, e))
	assert.Equal(t, e.value, 4)
	assert.ErrorIs(t, p.Publish(ctx, &event{err: err}), err)
	assert.ErrorIs(t, p.Publish(ctx, event{}), bus.ErrNonPointer)
	assert.ErrorIs(t, p.Publish(ctx, nil), bus.ErrNonPointer)
	assert.NoError(t, p.Publish(ctx, &err))
}

// goos: linux
// goarch: amd64
// cpu: AMD Ryzen 5 6600HS Creator Edition
// BenchmarkPublisher-12    	152014070	         7.671 ns/op
func BenchmarkPublisher(b *testing.B) {
	handler := direct.NewEventGroupHandler(
		func(ctx context.Context, event *event) error {
			event.value += 1

			return event.err
		},

		func(ctx context.Context, event *event) error {
			event.value += 2

			return event.err
		},
	)

	p := direct.NewPublisher(handler)
	ctx := context.Background()
	e := &event{}

	for i := 0; i < b.N; i++ {
		_ = p.Publish(ctx, e)
	}
}

// goos: linux
// goarch: amd64
// cpu: AMD Ryzen 5 6600HS Creator Edition
// BenchmarkPublisherParallel-12    	44976279	        27.67 ns/op
func BenchmarkPublisherParallel(b *testing.B) {
	type event struct {
		value uint32
	}

	handler := direct.NewEventGroupHandler(
		func(ctx context.Context, event *event) error {
			atomic.AddUint32(&event.value, 1)

			return nil
		},

		func(ctx context.Context, event *event) error {
			atomic.AddUint32(&event.value, 2)

			return nil
		},
	)

	p := direct.NewPublisher(handler)
	ctx := context.Background()
	e := &event{}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = p.Publish(ctx, e)
		}
	})
}
