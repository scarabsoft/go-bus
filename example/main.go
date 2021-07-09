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

func main() {
	topic, _ := bus.CreateTopic("topic", bus.AsyncTopic)

	//err := topic.Subscribe(func(event bus.Event) {
	err := topic.Subscribe(func(ID uint64, topic string, payload interface{}) {
		switch topic {
		case "topic":
			fmt.Println(ID, topic, payload)
		default:
			fmt.Println("ignored")

		}
		//fmt.Println(event.String())
		//fmt.Println(event)
	})

	//err = topic.Subscribe(bus.EventHandler(func(event bus.Event) {
	//	fmt.Println(event.ID, event.Topic, event.Payload)
	//	fmt.Println(event)
	//}))

	err = topic.Publish("HaHaHa")
	topic.Close()
	time.Sleep(1 * time.Second)
	fmt.Println(err)
}
