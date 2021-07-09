package pool

import "errors"

var (
	ErrPoolIsClosed = errors.New("pool is closed")
)

type TaskPayload struct {
	ID      uint64
	Name    string
	Payload interface{}
}

type Task struct {
	Payload TaskPayload
	Handler func(ID uint64, name string, payload interface{})
}

type Pool interface {
	Submit(Task) error

	Close() error
}

type Options struct {
	MaxWorkers   int
	MaxQueueSize int
}
