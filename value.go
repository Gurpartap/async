package async

import (
	"context"
)

type ValueFuture interface {
	Await() interface{}
	AwaitContext(ctx context.Context) interface{}
}

type valueFuture struct {
	await func(ctx context.Context) interface{}
}

func (f valueFuture) Await() interface{} {
	return f.await(context.Background())
}

func (f valueFuture) AwaitContext(ctx context.Context) interface{} {
	return f.await(ctx)
}

func Value(f func() interface{}) ValueFuture {
	var result interface{}
	c := make(chan struct{}, 1)
	go func() {
		defer close(c)
		result = f()
	}()
	return valueFuture{
		await: func(ctx context.Context) interface{} {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				return result
			}
		},
	}
}
