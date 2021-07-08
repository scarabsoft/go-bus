package bus

import (
	"fmt"
	"sync"
)

type Bus interface {
	Get(topic string) Topic

	Publish(topic string, data interface{}) (Topic, error)
	//
	//Subscribe(topic string, handler event.EventHandler) error

	//Unsubscribe(handler event.EventHandler) error

	CreateTopic(name string) (Topic, error)
}

type busImpl struct {
	topics map[string]Topic
	lock   sync.RWMutex

	// FIXME optional topic creation if missing
	// FIXME set default topic options
}

func (b *busImpl) Get(topic string) Topic {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.topics[topic]
}

func (b *busImpl) Publish(topic string, data interface{}) (Topic, error) {
	//FIXME if topic does not exist create default one and publish
	t := b.Get(topic)
	if err := t.Publish(data); err != nil {
		return nil, err
	}
	return t, nil
}

//FIXME should be method which accepts builder instead of options
func (b *busImpl) CreateTopic(name string) (Topic, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	fmt.Println("Create Topic")
	// FIXME make it thread safe
	// FIXME check whether topic already exists, if so --> error
	// FIXME add topic

	//result := syncTopic{name: name, handlers: []event.EventHandler{}}
	//result := topic.asyncTopic{name: name, handlers: []event.EventHandler{}}
	result := newAsyncTopic(name)
	fmt.Println(result)

	b.topics[name] = result
	return result, nil
}

func newBus() Bus {
	return &busImpl{
		topics: make(map[string]Topic),
		lock:   sync.RWMutex{},
	}
}
