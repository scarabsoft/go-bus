package bus

import (
	"fmt"
	"github.com/scarabsoft/go-bus/internal/bus"
	"github.com/scarabsoft/go-bus/internal/topic"
)

type Event struct {
	ID      uint64
	Topic   string
	Payload interface{}
}

func (e Event) String() string {
	return fmt.Sprintf("[%s-%d]: %s", e.Topic, e.ID, e.Payload)
}
func EventHandler(handler func(event Event)) func(ID uint64, name string, payload interface{}) {
	return func(ID uint64, name string, payload interface{}) {
		handler(Event{ID: ID, Topic: name, Payload: payload})
	}
}

var (
	std = bus.NewBus()

	SyncTopic = func(t topic.RootBuilder) topic.Builder {
		return topic.NewSyncBuilder(t.Name)
	}

	AsyncTopic = func(t topic.RootBuilder) topic.Builder {
		return topic.NewAsyncBuilder(t.Name)
	}

	WorkerTopic = func(t topic.RootBuilder) topic.Builder {
		return topic.NewWorkerBuilder(t.Name)
	}
)

func Get(topic string) topic.Topic {
	return std.Get(topic)
}

func Publish(topic string, payloads ...interface{}) (topic.Topic, error) {
	return std.Publish(topic, payloads...)
}

func CreateTopic(name string, fn func(t topic.RootBuilder) topic.Builder) (topic.Topic, error) {
	return std.CreateTopic(name, fn)
}
