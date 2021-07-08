package go_bus

import (
	"fmt"
)

type syncTopic struct {
	topicImpl
}

func (s syncTopic) Publish(data interface{}) error {
	e := Event{Topic: s.name, Payload: data}
	for _, handler := range s.handlers {
		_ = handler(e)
	}
	return nil
}

func (s *syncTopic) Subscribe(handler EventHandler) error {
	fmt.Println("subscribe", s.name)
	s.handlers = append(s.handlers, handler)
	return nil
}
