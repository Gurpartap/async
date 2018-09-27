# async/await

Context cancellable futures in Go

[![GoDoc](https://godoc.org/github.com/Gurpartap/async?status.svg)](https://godoc.org/github.com/Gurpartap/async)

Provides a simple parallel execution primitive that we've all come to miss from Go: async/await futures.

Each future is a straightforward implementation of a single go routine and a cancellable result channel.

Accessing a future's result may require type casting from `interface{}`.

## Usage

```sh
go get -u github.com/Gurpartap/async
```

1. Create a `f := async.Any(...)`.
2. Do something else while it runs in the background.
3. `v, err := f.Await()` or `f.AwaitContext(ctx)`
4. Cast the resulting `value := v.(*MyType)`.

```go
import (
	"github.com/Gurpartap/async"
)

future := async.Any(func() (interface{}, error) {
	return getUser() // a slow database or api call
})

// do other stuff until we need user

v, err := future.Await()
value := v.(*User)
```

```go
// with timeout
timeout := 5*time.Second
ctx, _ := context.WithTimeout(context.Background(), timeout)
v, err := future.Await()
if err != nil {
	// handle if the err was a context timeout
}
value := v.(*User)
```

## Available Futures

The async package provides a bunch future types, each useful for a varying number of return params:

For a value and an error:
```go
func Any(func() (interface{}, error))
func Any(context.Context, func() (interface{}, error))
```

For 2 values and an error:
```go
func Any2(func() (interface{}, interface{}, error))
func Any2(context.Context, func() (interface{}, interface{}, error))
```

For a value only:
```go
func Value(func() interface{})
func Value(context.Context, func() interface{})
```

For an error only:
```go
func Err(func() error)
func Err(context.Context, func() error)
```

## Full Example

```go
package async_test

// ...
// see the complete example code in example_test.go
// ...

func ExampleSync() {
	// if the db calls take 1 second each, the usual synchronous way
	// will take 4 seconds to get results

	user, err1 := db.GetUser()
	name, age, err2 := db.GetNameAndAge()
	location := db.GetLocation()
	err := db.SetName("Another Name")

	// use results
}

func ExampleAsync() {
	// execute the db calls asynchronously
	// and get results in 1 second

	anyFuture := async.Any(func() (interface{}, error) {
		return db.GetUser()
	})
	any2Future := async.Any2(func() (interface{}, interface{}, error) {
		return db.GetNameAndAge()
	})
	valueFuture := async.Value(func() interface{} {
		return db.GetLocation()
	})
	errFuture := async.Err(func() error {
		return db.SetName("Another Name")
	})

	// wait for results
	any, err1 := anyFuture.Await()
	user := any.(User)

	any1, any2, err2 := any2Future.Await()
	name := any1.(string)
	age := any2.(int)

	value := valueFuture.Await()
	location := value.(string)

	// cancellable
	ctx, _ := context.WithCancel(context.Background())
	err := errFuture.AwaitContext(ctx)

	// use results
}

```

## License

Copyright 2018 Gurpartap Singh

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
