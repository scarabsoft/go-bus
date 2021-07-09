package bus

import (
	"fmt"
	"github.com/scarabsoft/go-bus/internal/topic"
	"sync"
)

type Bus interface {
	Get(topic string) topic.Topic

	Publish(topic string, payloads ...interface{}) (topic.Topic, error)
	//
	//Subscribe(topic string, handler ...event.EventHandler) error

	//Unsubscribe(handler event.EventHandler) error

	CreateTopic(name string, fn func(topic topic.TopicInit) topic.TopicBuilder) (topic.Topic, error)
}

type busImpl struct {
	topics map[string]topic.Topic
	lock   sync.RWMutex

	// FIXME optional topic creation if missing
	// FIXME set default topic options
}

func (b *busImpl) Get(topic string) topic.Topic {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.topics[topic]
}

func (b *busImpl) Publish(topic string, payloads ...interface{}) (topic.Topic, error) {
	//FIXME if topic does not exist create default one and publish
	t := b.Get(topic)
	if err := t.Publish(payloads...); err != nil {
		return nil, err
	}
	return t, nil
}

//FIXME should be method which accepts builder instead of options
func (b *busImpl) CreateTopic(name string, fn func(topic topic.TopicInit) topic.TopicBuilder) (topic.Topic, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	fmt.Println("Create Topic")
	// FIXME make it thread safe
	// FIXME check whether topic already exists, if so --> error
	// FIXME add topic

	//result := syncTopicImpl{name: name, handlers: []event.EventHandler{}}
	//result := topic.asyncTopicImpl{name: name, handlers: []event.EventHandler{}}
	//result := newAsyncTopic(name)

	result := fn(*topic.NewTopicInit(name)).Build()
	fmt.Println(result)

	b.topics[name] = result
	return result, nil
}

func NewBus() Bus {
	return &busImpl{
		topics: make(map[string]topic.Topic),
		lock:   sync.RWMutex{},
	}
}
