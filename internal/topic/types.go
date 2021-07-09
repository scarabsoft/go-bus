package topic

import (
	"errors"
)

var (
	ErrDoesNotExists     = errors.New("topic does not exists")
	ErrAlreadyExists     = errors.New("topic already exists")
	ErrAlreadyClosed     = errors.New("topic already closed")
	ErrAlreadySubscribed = errors.New("handler already subscribed to topic")
)

type Topic interface {
	Name() string
	Publish(payloads ...interface{}) error
	Subscribe(handlers ...func(ID uint64, topic string, payload interface{})) error
	Unsubscribe(handlers ...func(ID uint64, topic string, payload interface{})) error
	Close() error
}

type Builder interface {
	Name(name string) Builder
	Build() Topic
}
