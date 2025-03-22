# Golang command/event/query bus

### This repository contains contracts for working with commands, events, requests

A bus for Go. Zero allocations during pub/sub.

Please note:
- This repository contains contracts for working with commands, events, requests
- A working example of implementation with direct messaging can be found in [direct](./direct)

## Install
```
go get github.com/go-gang/bus
```

## Usage

### Command dispatcher

```go
package main

import (
	"context"
	"github.com/go-gang/bus/direct"
)

func main() {
	type command struct {
		name string
	}
	
	handler := direct.NewCommandHandler(func(ctx context.Context, cmd *command) error {
		println(cmd.name)

		return nil
	})

	bus := direct.NewDispatcher(handler)

	_ = bus.Dispatch(context.Background(), &command{
		name: "test",
	})
}

```

If no handler is specified for the command, an error will be returned `bus.ErrHandlerNotFound`.



### Event pusblisher
```go
package main

import (
	"context"
	"github.com/go-gang/bus/direct"
)

func main() {
	type event struct {
		name string
	}

	handler := direct.NewEventGroupHandler(
		func(ctx context.Context, e *event) error {
			println("group handler 1", e.name)

			return nil
		},

		func(ctx context.Context, e *event) error {
			println("group handler 2", e.name)

			return nil
		},
	)

	 bus := direct.NewPublisher(
		 handler,
		 direct.NewEventHandler(
			 func(ctx context.Context, e *event) error {
				 println("single handler", e.name)

				 return nil
			 }, 
		 ),
	 )

	_ = bus.Publish(context.Background(), &event{
		name: "test",
	})
}

```

If no handler is specified for the command, no error is returned.

### Query asker
```go
package main

import (
	"context"
	"github.com/go-gang/bus/direct"
)

func main() {
	type query struct {
		value string
	}
	
	type result1 struct {
		value string
	}

	type result2 struct {
		value string
	}

	bus := direct.NewAsker(
		direct.NewQueryHandler(
			func(ctx context.Context, query *query, result *result1) error {
				result.value = query.value + "1"
				return nil
			},
		),
		direct.NewQueryHandler(
			func(ctx context.Context, query *query, result *result1) error {
				result.value = query.value + "2"
				return nil
			},
		),
	)
	
	r1 := &result1{}
	r2 := &result2{}

	_ = bus.Ask(context.Background(), &query{
		value: "test",
	}, r1)


	_ = bus.Ask(context.Background(), &query{
		value: "test",
	}, r2)
}
```

If no handler is specified for the command, an error will be returned `bus.ErrHandlerNotFound`.