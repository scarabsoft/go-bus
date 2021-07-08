package go_bus

import (
	"fmt"
	"sync"
)

type Bus interface {
	Get(topic string) Topic

	//Publish(topic string, data interface{}) error
	//
	//Subscribe(topic string, handler event.EventHandler) error

	//Unsubscribe(handler event.EventHandler) error

	CreateTopic(name string, options ...Option) error
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

//FIXME should be method which accepts builder instead of options
func (b *busImpl) CreateTopic(name string, options ...Option) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	fmt.Println("Create Topic")
	// FIXME make it thread safe
	// FIXME check whether topic already exists, if so --> error
	// FIXME add topic

	//s := syncTopic{name: name, handlers: []event.EventHandler{}}
	//s := topic.asyncTopic{name: name, handlers: []event.EventHandler{}}
	s := newAsyncTopic(name)
	fmt.Println(s)

	b.topics[name] = s
	return nil
}

func new() Bus {
	return &busImpl{
		topics: make(map[string]Topic),
		lock:   sync.RWMutex{},
	}
}
