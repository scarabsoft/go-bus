package topic

import (
	"errors"
)

var (
	ErrDoesNotExists = errors.New("name does not exists")
	ErrAlreadyExists = errors.New("name already exists")
	ErrAlreadyClosed = errors.New("name already closed")
)

type Topic interface {
	Name() string
	Publish(payloads ...interface{}) error
	Subscribe(handlers ...func(ID uint64, topic string, payload interface{})) error
	Close() error
}

type Builder interface {
	Name(name string) Builder
	Build() Topic
}


