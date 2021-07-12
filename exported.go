package bus

import (
	"fmt"
	"github.com/scarabsoft/go-bus/internal/bus"
	"github.com/scarabsoft/go-bus/internal/pool"
	"github.com/scarabsoft/go-bus/internal/topic"
	"runtime"
	"sync"
)

type Event struct {
	ID      uint64
	Topic   string
	Payload interface{}
}

func (e Event) String() string {
	return fmt.Sprintf("[%s-%d]: %[3]T(%[3]v)", e.Topic, e.ID, e.Payload)
}
func EventHandler(handler func(event Event)) func(ID uint64, name string, payload interface{}) {
	return func(ID uint64, name string, payload interface{}) {
		handler(Event{ID: ID, Topic: name, Payload: payload})
	}
}

var (
	std = bus.New()

	defaultAsyncPoolOnce           = sync.Once{}
	defaultAsyncPool     pool.Pool = nil

	defaultWorkerPoolOnce           = sync.Once{}
	defaultWorkerPool     pool.Pool = nil

	SyncTopic = &SyncTopicBuilder{}

	AsyncTopic = NewAsyncTopicBuilder()

	WorkerTopic = NewWorkerTopicBuilder()
)

func DefaultBus() bus.Bus {
	return std
}

func DefaultAsyncPool() pool.Pool {
	defaultAsyncPoolOnce.Do(func() {
		defaultAsyncPool = pool.New(pool.Options{
			MaxQueueSize: 1,
			MaxWorkers:   1,
		})
	})
	return defaultAsyncPool
}

func DefaultWorkerPool() pool.Pool {
	defaultWorkerPoolOnce.Do(func() {
		defaultWorkerPool = pool.New(pool.Options{
			MaxQueueSize: 100,
			MaxWorkers:   runtime.NumCPU(),
		})
	})
	return defaultWorkerPool
}

func Publish(topic string, payloads ...interface{}) error {
	return std.Publish(topic, payloads...)
}

func Subscribe(name string, handlers ...func(ID uint64, name string, payload interface{})) error {
	return std.Subscribe(name, handlers...)
}

func CreateTopic(name string, tb TopicBuilder) (topic.Topic, error) {
	return std.CreateTopic(name, tb.Build())
}

func Get(name string) (topic.Topic, error) {
	return std.Get(name)
}

func DeleteTopic(name string) error {
	return std.DeleteTopic(name)
}

func SetDefaultTopicBuilder(tb TopicBuilder) {
	std.SetDefaultTopicBuilder(tb.Build())
}

type TopicBuilder interface {
	Build() topic.Builder
}

type SyncTopicBuilder struct {
}

func (s SyncTopicBuilder) Build() topic.Builder {
	return topic.NewSyncBuilder()
}

type AsyncTopicBuilder struct {
	p pool.Pool
}

func NewAsyncTopicBuilder() *AsyncTopicBuilder {
	return &AsyncTopicBuilder{p: DefaultAsyncPool()}
}

func (a AsyncTopicBuilder) Build() topic.Builder {
	return topic.NewAsyncBuilder(a.p)
}

type WorkerTopicBuilder struct {
	p pool.Pool
}

func NewWorkerTopicBuilder() *WorkerTopicBuilder {
	return &WorkerTopicBuilder{p: DefaultWorkerPool()}
}

func (w *WorkerTopicBuilder) Pool(p pool.Pool) *WorkerTopicBuilder {
	w.p = p
	return w
}

func (w *WorkerTopicBuilder) Build() topic.Builder {
	return topic.NewWorkerBuilder(w.p)
}
