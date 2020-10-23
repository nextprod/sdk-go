package runtime

import (
	"context"
	"errors"
	"log"
)

type startFunc struct {
	f func(ctx context.Context, hander Handler) error
}

var (
	rpcStartFunc = &startFunc{
		f: func(c context.Context, h Handler) error {
			return errors.New("function was compiled without rpc support")
		},
	}
	startFunctions = []*startFunc{rpcStartFunc}
)

// Start starts extension runtime with context.
func Start(handler interface{}) {
	start(context.Background(), NewHandler(handler))
}

// StartWithContext starts extension runtime with context.
func StartWithContext(ctx context.Context, handler Handler) {
	start(ctx, NewHandler(handler))
}

func start(ctx context.Context, handler Handler) {
	for _, start := range startFunctions {
		if err := start.f(ctx, handler); err != nil {
			log.Fatalf("%v", err)
		}
	}
}
