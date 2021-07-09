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
	//a := func(){}
	////b := func(){}
	//
	//if reflect.ValueOf(a) == reflect.ValueOf(a){
	//	fmt.Println("TRUE", reflect.ValueOf(a))
	//}else{
	//	fmt.Println("FALSE")
	//}

	//bus.CreateTopicIfNotExists(bus.SyncTopic)
	topic, _ := bus.CreateTopic("topic", bus.SyncTopic)

	//bus.CreateTopic("topic-2", bus.WorkerTopic)
	h := bus.EventHandler(func(evt bus.Event) {
		fmt.Println("H1", evt)
	})

	err := topic.Subscribe(h)

	//err = bus.Subscribe("topic-2", bus.EventHandler(func(evt bus.Event) {
	//	fmt.Println("H2", evt)
	//}))

	if err != nil {
		fmt.Println("sub failed", err)
	}

	for i := 0; i < 100; i++ {
		topic.Publish("test", "test1234")
		topic.Unsubscribe(h)
		//topic.Publish("test", "test1234")
		//bus.Publish("topic-2", 1, 2, 3, 4)
	}

	time.Sleep(1 * time.Second)

}
