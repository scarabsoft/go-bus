package topic

import (
	"errors"
)

var ErrAlreadyClosed = errors.New("topic already closed")

type Topic interface {
	Name() string
	//FIXME Options map[string]interface{} ?
	Publish(payloads ...interface{}) error
	Subscribe(handlers ...func(ID uint64, topic string, payload interface{})) error
	Close() error
}

type RootBuilder struct {
	Name string
}

type Builder interface {
	Build() Topic
}
