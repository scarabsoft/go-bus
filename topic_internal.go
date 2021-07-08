package bus

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrTopicClosed = errors.New("topic already closed")

type Topic interface {
	Publish(data interface{}) error       // FIXME better would be ...interface{}
	Subscribe(handler EventHandler) error // FIXME better would be ...event.EventHandler
	Close() error
}

type abstractTopicImpl struct {
	name        string
	handlers    []EventHandler
	idGenerator func() uint64
	lock        sync.RWMutex

	closed bool
}

func newAbstractTopicImpl(name string) abstractTopicImpl {
	return abstractTopicImpl{
		name:        name,
		handlers:    []EventHandler{},
		idGenerator: topicIdGenerator(),
		closed:      false,
	}
}

func (a *abstractTopicImpl) Subscribe(handler EventHandler) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.closed {
		return ErrTopicClosed
	}
	a.handlers = append(a.handlers, handler)
	return nil
}

func (a *abstractTopicImpl) Close() error {
	return nil
}

//generates topic id which guarantees to be thread safe and monotonous
func topicIdGenerator() func() uint64 {
	var idx uint64 = 0
	return func() uint64 {
		return atomic.AddUint64(&idx, 1)
	}
}
