package bus

import (
	"github.com/scarabsoft/go-bus/internal/topic"
	"sync"
)

type busImpl struct {
	topics map[string]topic.Topic
	lock   sync.RWMutex

	defaultTopicBuilder topic.Builder
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

func (b *busImpl) getOrCreateEventually(name string, tb topic.Builder) (topic.Topic, error) {
	b.lock.RLock()
	if t, ok := b.topics[name]; !ok {
		if tb != nil {
			b.lock.RUnlock()
			b.lock.Lock()
			defer b.lock.Unlock()

			// this makes sure that we dont overwrite an existing topic
			if _, ok := b.topics[name]; ok {
				return nil, topic.ErrAlreadyExists
			}

			t = tb.Name(name).Build()
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

func (b *busImpl) CreateTopic(name string, tb topic.Builder) (topic.Topic, error) {
	return b.getOrCreateEventually(name, tb)
}

func (b *busImpl) CreateTopicIfNotExists(tb topic.Builder) {
	b.defaultTopicBuilder = tb
}

func NewBus() *busImpl {
	return &busImpl{
		topics:              make(map[string]topic.Topic),
		lock:                sync.RWMutex{},
		defaultTopicBuilder: nil,
	}
}
