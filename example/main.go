package main

import (
	"fmt"
	"sync"
	"time"
)

type Event struct {
	Topic   string
	Payload interface{}
}

type EventHandler func(event Event) error

type Topic interface {
	Publish(data interface{}) error       // FIXME better would be ...interface{}
	Subscribe(handler EventHandler) error // FIXME better would be ...EventHandler

	//FIXME add Close() 				 // FIXME proper closing topic in aysn case use wait groups
}

type SyncTopic struct {
	name     string
	handlers []EventHandler
	// FIXME FIELD stop on error
}

func (s SyncTopic) Publish(data interface{}) error {
	e := Event{Topic: s.name, Payload: data}
	for _, handler := range s.handlers {
		_ = handler(e)
	}
	return nil
}

func (s *SyncTopic) Subscribe(handler EventHandler) error {
	fmt.Println("subscribe", s.name)
	s.handlers = append(s.handlers, handler)
	return nil
}

type AsyncTopic struct {
	name     string
	handlers []EventHandler
}

func (a AsyncTopic) Publish(data interface{}) error {
	//FIXME error if finished
	go func() {
		e := Event{Topic: a.name, Payload: data} // FIXME this should be a simple event generator to have auto increment ids
		for _, handler := range a.handlers {
			_ = handler(e)
		}
	}()
	return nil
}

func (a *AsyncTopic) Subscribe(handler EventHandler) error {
	//FIXME error if finished
	fmt.Println("subscribe async", a.name)
	a.handlers = append(a.handlers, handler)
	return nil
}

type Option func(topic Topic) error

type Bus interface {
	Get(topic string) Topic

	//Publish(topic string, data interface{}) error
	//
	//Subscribe(topic string, handler EventHandler) error

	//Unsubscribe(handler EventHandler) error

	CreateTopic(name string, options ...Option) error
}

type busImpl struct {
	topics map[string]Topic
	lock   sync.RWMutex
}

//func (b busImpl) Publish(topic string, data interface{}) error {
//	return nil
//}
//
//func (b busImpl) Subscribe(topic string, handler EventHandler) error {
//	return nil
//}

//func (b busImpl) Unsubscribe(handler EventHandler) error {
//	return nil
//}

// sync topics block async topics from firing
// FIXME should async topics fire async?

// topicBuilder().Sync()... only sync options	// sync fire and forget
// topicBuilder().Async()... only async options // async fire and forget - publish async & run all event handler async in single go routine

// topicBuilder().SyncBuffer()... only worker options  // keep n elements for history for late joiner
// topicBuilder().AsyncPool()... only worker options // use default pool or individual pool for a topic
// ? topicBuilder().AsyncBufferedPool()... only worker options // use default pool or individual pool for a topic

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

	//s := SyncTopic{name: name, handlers: []EventHandler{}}
	s := AsyncTopic{name: name, handlers: []EventHandler{}}
	fmt.Println(s)

	b.topics[name] = &s
	return nil
}

func main() {
	bus := &busImpl{
		topics: make(map[string]Topic),
		lock:   sync.RWMutex{},
	}
	fmt.Println(bus)

	_ = bus.CreateTopic("someTopic")

	t := bus.Get("someTopic")

	_ = t.Subscribe(func(event Event) error {
		fmt.Println("Handle:", event.Topic, event.Payload)
		return nil
	})

	_ = t.Publish("Hello")
	_ = t.Publish("World")

	//FIXME close / wait
	time.Sleep(5 * time.Second)

	//fmt.Println(t)
}
