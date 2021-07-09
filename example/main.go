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

	//tb := bus.topicBuilderImpl{}
	//
	//tb.Sync()

	//_, _ = bus.CreateTopic("syncTopic", func(topic *bus.topicBuilderImpl) bus.Topic {
	//	return topic.SyncTopic()
	//})

	//_, _ = bus.CreateTopic("syncTopic", func(topic *bus.topicBuilderImpl) bus.Topic {
	//	return topic.SyncTopic().Build()
	//})

	_, _ = bus.CreateTopic("syncTopic", bus.SyncTopic)

	//_, _ = bus.CreateTopic("asyncTopic", func(topic bus.TopicInit) bus.TopicBuilder {
	//	return topic.Async()
	//})

	_, _ = bus.CreateTopic("asyncTopic", bus.AsyncTopic)

	syncTopic := bus.Get("syncTopic")
	fmt.Println(syncTopic.Name())
	
	asyncTopic := bus.Get("asyncTopic")
	fmt.Println(asyncTopic.Name())

	handler := func(event bus.Event) error {
		fmt.Println("Handle:", event.ID, event.Topic, event.Payload)
		return nil
	}

	_ = syncTopic.Subscribe(handler)
	_ = asyncTopic.Subscribe(handler)

	//stop := make(chan struct{})

	//go func() {
	//	for {
	//		select {
	//		case <-stop:
	//		default:
	//topic, _ := bus.Publish("syncTopic", "Hello", "World")
	//_ = topic.Publish("Another Message")
	//}
	//
	//}
	//}()

	//for i := 0; i < 1000; i++ {
	//	go func() {
	//		_, _ = bus.Publish("syncTopic", "Async 1", "Async 2")
	//	}()
	//}

	for i := 0; i < 10; i++ {
		_, _ = bus.Publish("asyncTopic", "Async 1", "Async 2")
		//_, _ = bus.Publish("syncTopic", "Sync 1", "Sync 2")
		_, _ = bus.Publish("asyncTopic", "Async 3", "Async 4")
	}

	//FIXME close / wait
	time.Sleep(100 * time.Millisecond)

	//stop <- struct{}{}

	fmt.Println(syncTopic)

	syncTopic.Close()
}
