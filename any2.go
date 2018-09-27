package async

import (
	"context"
)

type Any2Future interface {
	Await() (interface{}, interface{}, error)
	AwaitContext(ctx context.Context) (interface{}, interface{}, error)
}

type any2Future struct {
	await func(ctx context.Context) (interface{}, interface{}, error)
}

func (f any2Future) Await() (interface{}, interface{}, error) {
	return f.await(context.Background())
}

func (f any2Future) AwaitContext(ctx context.Context) (interface{}, interface{}, error) {
	return f.await(ctx)
}

func Any2(f func() (interface{}, interface{}, error)) Any2Future {
	var result1 interface{}
	var result2 interface{}
	var err error
	c := make(chan struct{}, 1)
	go func() {
		defer close(c)
		result1, result2, err = f()
	}()
	return any2Future{
		await: func(ctx context.Context) (interface{}, interface{}, error) {
			select {
			case <-ctx.Done():
				return nil, nil, ctx.Err()
			case <-c:
				return result1, result2, err
			}
		},
	}
}
