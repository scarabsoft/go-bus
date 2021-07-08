package go_bus

import (
	"fmt"
)

type asyncTopic struct {
	topicImpl
}

func newAsyncTopic(name string) Topic {
	return &asyncTopic{
		topicImpl: topicImpl{
			name:     name,
			handlers: []EventHandler{},
		},
	}
}

func (a asyncTopic) Publish(data interface{}) error {
	//FIXME error if finished
	go func() {
		e := Event{Topic: a.name, Payload: data} // FIXME this should be a simple event generator to have auto increment ids
		for _, handler := range a.handlers {
			_ = handler(e)
		}
	}()
	return nil
}

func (a *asyncTopic) Subscribe(handler EventHandler) error {
	//FIXME error if finished
	fmt.Println("subscribe async", a.name)
	a.handlers = append(a.handlers, handler)
	return nil
}
