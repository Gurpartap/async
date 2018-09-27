package async_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Gurpartap/async"
)

var db = &UsersStorage{
	User: User{
		Name:     "Full Name",
		Age:      100,
		Location: "Mars",
	},
}

type User struct {
	Name     string
	Age      int
	Location string
}

type UsersStorage struct {
	User User
}

func (s *UsersStorage) GetUser() (User, error) {
	time.Sleep(1 * time.Second)
	return s.User, nil
}

func (s *UsersStorage) GetLocation() string {
	time.Sleep(1 * time.Second)
	return s.User.Location
}

func (s *UsersStorage) GetNameAndAge() (string, int, error) {
	time.Sleep(1 * time.Second)
	return s.User.Name, s.User.Age, nil
}

func (s *UsersStorage) SetName(name string) error {
	time.Sleep(1 * time.Second)
	return errors.New("could not set name")
}

func ExampleSync() {
	// if the db calls take 1 second each, the usual synchronous way
	// will take 4 seconds to get results

	user, err1 := db.GetUser()
	name, age, err2 := db.GetNameAndAge()
	location := db.GetLocation()
	err := db.SetName("Another Name")

	// use results
	fmt.Printf("Any: user = %v, err = %v\n", user, err1)
	fmt.Printf("Any2: name = %s, age = %d, err = %v\n", name, age, err2)
	fmt.Printf("Value: location = %s\n", location)
	fmt.Printf("Err: err = %v\n\n", err)
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
	fmt.Printf("Any: user = %v, err = %v\n", user, err1)
	fmt.Printf("Any2: name = %s, age = %d, err = %v\n", name, age, err2)
	fmt.Printf("Value: location = %s\n", location)
	fmt.Printf("Err: err = %v\n\n", err)
}
