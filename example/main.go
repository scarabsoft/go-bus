package main

import (
	"fmt"
	"github.com/scarabsoft/go-bus"
	"time"
)

//func (b busImpl) Publish(topic string, data interface{}) error {
//	return nil
//}
//
//func (b busImpl) Subscribe(topic string, handler event.EventHandler) error {
//	return nil
//}

//func (b busImpl) Unsubscribe(handler event.EventHandler) error {
//	return nil
//}

// sync topics block async topics from firing
// FIXME should async topics fire async?

// topicBuilder().syncTopic()... only sync options	// sync fire and forget
// topicBuilder().asyncTopic()... only async options // async fire and forget - publish async & run all event handler async in single go routine

// topicBuilder().SyncBuffer()... only worker options  // keep n elements for history for late joiner
// topicBuilder().AsyncPool()... only worker options // use default pool or individual pool for a topic
// ? topicBuilder().AsyncBufferedPool()... only worker options // use default pool or individual pool for a topic

func main() {

	_ = bus.CreateTopic("someTopic")
	t := bus.Get("someTopic")

	_ = t.Subscribe(func(event bus.Event) error {
		fmt.Println("Handle:", event.ID, event.Topic, event.Payload)
		return nil
	})

	go func() {
		for {
			_ = t.Publish("Hello")
			_ = t.Publish("World")
		}
	}()

	//FIXME close / wait
	time.Sleep(1 * time.Second)

	fmt.Println(t)
}
