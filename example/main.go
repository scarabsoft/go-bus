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
	//bus.CreateTopicIfNotExists(bus.SyncTopic)
	//
	//_ = bus.Subscribe("test", bus.EventHandler(func(event bus.Event) {
	//	fmt.Println(event)
	//}))
	//
	//for i := 0; i < 100; i++ {
	//	func() {
	//		if err := bus.Publish("test", "test"); err != nil {
	//			fmt.Println(err)
	//		}
	//	}()
	//}
	//
	//time.Sleep(1 * time.Second)

	//p := pool.NewPool(pool.Options{
	//	MaxQueueSize: 1,
	//	MaxWorkers:   1,
	//})

	//for i := 0; i < 100; i++ {
	//	_ = p.Submit(pool.Task{
	//		Payload: pool.TaskPayload{
	//			ID:      uint64(i),
	//			Name:    "XYZ",
	//			Payload: struct{}{},
	//		},
	//		Handler: bus.EventHandler(func(evt bus.Event) {
	//			fmt.Println(evt)
	//			time.Sleep(500 * time.Millisecond)
	//		}),
	//	})
	//}
	//time.Sleep(1 * time.Second)
	//
	//p.Close()
	//
	//time.Sleep(1 * time.Second)
	//
	//err := p.Submit(pool.Task{
	//	Payload: pool.TaskPayload{
	//		ID:      uint64(10),
	//		Name:    "XYZ",
	//		Payload: struct{}{},
	//	},
	//	Handler: bus.EventHandler(func(evt bus.Event) {
	//		fmt.Println(evt)
	//	}),
	//})
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//time.Sleep(100 * time.Second)

	////topic, _ := bus.CreateTopic("topic", bus.AsyncTopic)
	//topic, _ := bus.CreateTopic("topic", bus.SyncTopic)

	//p := pool.NewPool(pool.Options{
	//	MaxQueueSize: 1,
	//	MaxWorkers:   1,
	//})

	//topic, _ := bus.CreateTopic("topic", bus.AsyncTopic)
	topic, _ := bus.CreateTopic("topic", bus.SyncTopic)

	bus.CreateTopic("topic-2", bus.WorkerTopic)

	err := topic.Subscribe(bus.EventHandler(func(evt bus.Event) {
		fmt.Println("H1", evt)
	}))

	err = bus.Subscribe("topic-2", bus.EventHandler(func(evt bus.Event) {
		fmt.Println("H2", evt)
	}))

	if err != nil {
		fmt.Println("sub failed", err)
	}

	for i := 0; i < 100; i++ {
		topic.Publish("test", "test1234")
		bus.Publish("topic-2", 1, 2, 3, 4)
	}

	time.Sleep(1 * time.Second)

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
