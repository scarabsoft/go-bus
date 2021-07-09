package topic

import (
	"reflect"
	"sync"
	"sync/atomic"
)

type abstractTopicImpl struct {
	name       string
	handlers   []func(ID uint64, name string, payload interface{})
	generateID func() uint64
	lock       sync.RWMutex

	closed bool
}

func newAbstractTopicImpl() abstractTopicImpl {
	return abstractTopicImpl{
		handlers:   []func(ID uint64, name string, payload interface{}){},
		generateID: topicIdGenerator(),
		closed:     false,
	}
}

func (a *abstractTopicImpl) Name() string {
	return a.name
}

func (a *abstractTopicImpl) Subscribe(handlers ...func(ID uint64, name string, payload interface{})) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.closed {
		return ErrAlreadyClosed
	}

	for _, handler := range handlers {
		for _, currentHandler := range a.handlers {
			if reflect.ValueOf(currentHandler) == reflect.ValueOf(handler) {
				return ErrAlreadySubscribed
			}
		}
		a.handlers = append(a.handlers, handler)
	}

	return nil
}

func (a *abstractTopicImpl) Unsubscribe(handlers ...func(ID uint64, topic string, payload interface{})) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.closed {
		return ErrAlreadyClosed
	}
	keep := make([]func(ID uint64, topic string, payload interface{}), 0)
	for _, currentHandler := range a.handlers {
		found := false
		for _, handler := range handlers {
			if reflect.ValueOf(currentHandler) == reflect.ValueOf(handler) {
				found = true
				break
			}
		}
		if !found {
			keep = append(keep, currentHandler)
		}
	}
	a.handlers = make([]func(ID uint64, topic string, payload interface{}), len(keep))
	copy(a.handlers, keep)
	return nil
}

func (a *abstractTopicImpl) Close() error {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.closed = true
	return nil
}

//generates name id which guarantees to be thread safe and monotonous
func topicIdGenerator() func() uint64 {
	var idx uint64 = 0
	return func() uint64 {
		return atomic.AddUint64(&idx, 1)
	}
}
