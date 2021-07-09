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
//func (b busImpl) Subscribe(topic string, handler event.Handler) error {
//	return nil
//}

//func (b busImpl) Unsubscribe(handler event.Handler) error {
//	return nil
//}

type CustomEvent bus.Event

func (c CustomEvent) String() string {
	return fmt.Sprintf("custom event %d", c.ID)
}

func (c CustomEvent) Test() string {
	return c.Payload.(string)
}

func main() {
	bus.CreateTopicIfNotExists(bus.AsyncTopic)

	_ = bus.Subscribe("test", bus.EventHandler(func(event bus.Event) {
		fmt.Println(event)
	}))

	for i := 0; i < 100; i++ {
		func() {
			if err := bus.Publish("test", "test"); err != nil {
				fmt.Println(err)
			}
		}()
	}

	time.Sleep(1 * time.Second)

	////topic, _ := bus.CreateTopic("topic", bus.AsyncTopic)
	//topic, _ := bus.CreateTopic("topic", bus.SyncTopic)
	////topic, _ := bus.CreateTopic("topic", bus.WorkerTopic)
	//
	////err := topic.Subscribe(func(event bus.Event) {
	//err := topic.Subscribe(func(ID uint64, topic string, payload interface{}) {
	//	switch topic {
	//	case "topic":
	//		fmt.Println(ID, topic, payload)
	//	default:
	//		fmt.Println("ignored")
	//
	//	}
	//	//fmt.Println(event.String())
	//	//fmt.Println(event)
	//})
	//
	//err = topic.Subscribe(bus.EventHandler(func(event bus.Event) {
	//	switch event.Topic {
	//	case "topic":
	//		var c = CustomEvent(event)
	//		fmt.Println(c)
	//		fmt.Println(c.Test())
	//	default:
	//		fmt.Println(event)
	//	}
	//}))
	//
	//for i := 0; i < 100; i++ {
	//	err = topic.Publish("HaHaHa")
	//}
	//topic.Close()
	//time.Sleep(1 * time.Second)
	//fmt.Println(err)
}
