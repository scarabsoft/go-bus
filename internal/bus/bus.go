package bus

import (
	"github.com/scarabsoft/go-bus/internal/topic"
	"sync"
)

type Bus interface {
	Publish(name string, payloads ...interface{}) error

	Subscribe(name string, handlers ...func(ID uint64, name string, payload interface{})) error

	//Unsubscribe(name string, handler event.EventHandler) error

	CreateTopic(name string, fn func(topic topic.RootBuilder) topic.Builder) (topic.Topic, error)

	Get(name string) (topic.Topic, error)

	//DeleteTopic(name string)
}

type busImpl struct {
	topics map[string]topic.Topic
	lock   sync.RWMutex

	defaultTopicBuilder func(t topic.RootBuilder) topic.Builder
}

func (b *busImpl) Get(name string) (topic.Topic, error) {
	return b.getOrCreateEventually(name, b.defaultTopicBuilder)
}

func (b *busImpl) Publish(name string, payloads ...interface{}) error {
	var t topic.Topic
	var err error

	if t, err = b.getOrCreateEventually(name, b.defaultTopicBuilder); err != nil {
		return err
	}

	if err := t.Publish(payloads...); err != nil {
		return err
	}
	return nil
}

func (b *busImpl) Subscribe(name string, handlers ...func(ID uint64, name string, payload interface{})) error {
	var t topic.Topic
	var err error

	if t, err = b.getOrCreateEventually(name, b.defaultTopicBuilder); err != nil {
		return err
	}

	if err := t.Subscribe(handlers...); err != nil {
		return err
	}
	return nil
}

func (b *busImpl) getOrCreateEventually(name string, fn func(topic topic.RootBuilder) topic.Builder) (topic.Topic, error) {
	b.lock.RLock()
	if t, ok := b.topics[name]; !ok {
		if fn != nil {
			b.lock.RUnlock()
			b.lock.Lock()
			defer b.lock.Unlock()

			// this makes sure that we dont overwrite an existing topic
			if _, ok := b.topics[name]; ok {
				return nil, topic.ErrAlreadyExists
			}

			t = fn(topic.NewTopicInit(name)).Build()
			b.topics[name] = t
			return t, nil
		}
		b.lock.RUnlock()
		return nil, topic.ErrDoesNotExists
	} else {
		b.lock.RUnlock()
		return t, nil
	}
}

func (b *busImpl) CreateTopic(name string, fn func(topic topic.RootBuilder) topic.Builder) (topic.Topic, error) {
	return b.getOrCreateEventually(name, fn)
}

func (b *busImpl) CreateTopicIfNotExists(fn func(t topic.RootBuilder) topic.Builder) {
	b.defaultTopicBuilder = fn
}

func NewBus() *busImpl {
	return &busImpl{
		topics:              make(map[string]topic.Topic),
		lock:                sync.RWMutex{},
		defaultTopicBuilder: nil,
	}
}
