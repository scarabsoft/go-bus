package bus

import "github.com/scarabsoft/go-bus/internal/topic"

type Bus interface {
	Publish(name string, payloads ...interface{}) error

	Subscribe(name string, handlers ...func(ID uint64, name string, payload interface{})) error

	Unsubscribe(name string, handlers ...func(ID uint64, name string, payload interface{})) error

	CreateTopic(name string, tb topic.Builder) (topic.Topic, error)

	Get(name string) (topic.Topic, error)

	DeleteTopic(name string) error
}
