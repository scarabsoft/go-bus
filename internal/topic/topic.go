package topic

import (
	"sync"
	"sync/atomic"
)

func NewTopicInit(name string) *RootBuilder {
	return &RootBuilder{Name: name}
}

func (r *RootBuilder) Sync() Builder {
	return NewSyncBuilder(r.Name)
}

func (r *RootBuilder) Async() Builder {
	return NewAsyncBuilder(r.Name)
}

func (r *RootBuilder) AsyncWorker() Builder {
	return NewWorkerBuilder(r.Name)
}

type abstractTopicImpl struct {
	topic      string
	handlers   []func(ID uint64, name string, payload interface{})
	generateID func() uint64
	lock       sync.RWMutex

	closed bool
}

func newAbstractTopicImpl(name string) abstractTopicImpl {
	return abstractTopicImpl{
		topic:      name,
		handlers:   []func(ID uint64, name string, payload interface{}){},
		generateID: topicIdGenerator(),
		closed:     false,
	}
}

func (a *abstractTopicImpl) Name() string {
	return a.topic
}

func (a *abstractTopicImpl) Publish(data ...interface{}) error {
	a.lock.RLock()
	defer a.lock.RUnlock()

	if a.closed {
		return ErrAlreadyClosed
	}

	return nil
}

func (a *abstractTopicImpl) Subscribe(handlers ...func(ID uint64, name string, payload interface{})) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.closed {
		return ErrAlreadyClosed
	}
	a.handlers = append(a.handlers, handlers...)
	return nil
}

func (a *abstractTopicImpl) Close() error {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.closed = true
	return nil
}

//generates topic id which guarantees to be thread safe and monotonous
func topicIdGenerator() func() uint64 {
	var idx uint64 = 0
	return func() uint64 {
		return atomic.AddUint64(&idx, 1)
	}
}
