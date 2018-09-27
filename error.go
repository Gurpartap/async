package async

import (
	"context"
)

type ErrorFuture interface {
	Await() error
	AwaitContext(ctx context.Context) error
}

type errorFuture struct {
	await func(ctx context.Context) error
}

func (f errorFuture) Await() error {
	return f.await(context.Background())
}

func (f errorFuture) AwaitContext(ctx context.Context) error {
	return f.await(ctx)
}

func Err(f func() error) ErrorFuture {
	var err error
	c := make(chan struct{}, 1)
	go func() {
		defer close(c)
		err = f()
	}()
	return errorFuture{
		await: func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				return err
			}
		},
	}
}
