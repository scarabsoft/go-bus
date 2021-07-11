package topic

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrAlreadySubscribed = errors.New("handler already subscribed")
)

type ErrorNotExists struct {
	Name string
}

func (e ErrorNotExists) Error() string {
	return fmt.Sprintf("%s does not exists", e.Name)
}

type ErrorAlreadyExists struct {
	Name string
}

func (e ErrorAlreadyExists) Error() string {
	return fmt.Sprintf("%s already exists", e.Name)
}

type ErrorAlreadyClosed struct {
	Name string
}

func (e ErrorAlreadyClosed) Error() string {
	return fmt.Sprintf("%s already closed", e.Name)
}

type Publisher interface {
	Publish(payloads ...interface{}) error
}

type Subscription interface {
	Subscribe(handlers ...func(ID uint64, topic string, payload interface{})) error
	Unsubscribe(handlers ...func(ID uint64, topic string, payload interface{})) error
}

type Topic interface {
	io.Closer

	Publisher
	Subscription

	Name() string
}

type Builder interface {
	Name(name string) Builder
	Build() Topic
}
