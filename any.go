package async

import (
	"context"
)

type AnyFuture interface {
	Await() (interface{}, error)
	AwaitContext(ctx context.Context) (interface{}, error)
}

type anyFuture struct {
	await func(ctx context.Context) (interface{}, error)
}

func (f anyFuture) Await() (interface{}, error) {
	return f.await(context.Background())
}

func (f anyFuture) AwaitContext(ctx context.Context) (interface{}, error) {
	return f.await(ctx)
}

func Any(f func() (interface{}, error)) AnyFuture {
	var result interface{}
	var err error
	c := make(chan struct{}, 1)
	go func() {
		defer close(c)
		result, err = f()
	}()
	return anyFuture{
		await: func(ctx context.Context) (interface{}, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-c:
				return result, err
			}
		},
	}
}
