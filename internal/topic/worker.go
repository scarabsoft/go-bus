package topic

import (
	"github.com/scarabsoft/go-bus/internal/pool"
)

type workerTopicImpl struct {
	abstractTopicImpl

	pool pool.Pool
}

func (w *workerTopicImpl) Publish(payloads ...interface{}) error {
	for _, payload := range payloads {
		id := w.generateID()
		for _, handler := range w.handlers {
			if err := w.pool.Submit(
				pool.Task{
					Payload: pool.TaskPayload{
						ID:      id,
						Name:    w.name,
						Payload: payload,
					},
					Handler: handler,
				}); err != nil {
				return err
			}
		}
	}
	return nil
}



type WorkerTopicBuilder struct {
	topic workerTopicImpl
}

func NewWorkerBuilder(defaultPool pool.Pool) *WorkerTopicBuilder {
	return &WorkerTopicBuilder{topic: workerTopicImpl{
		abstractTopicImpl: newAbstractTopicImpl(),
		pool:              defaultPool,
	}}
}

func (wtb *WorkerTopicBuilder) Name(name string) Builder {
	wtb.topic.name = name
	return wtb
}

func (wtb *WorkerTopicBuilder) Pool(p pool.Pool) *WorkerTopicBuilder {
	wtb.topic.pool = p
	return wtb
}

func (wtb *WorkerTopicBuilder) Build() Topic {
	return &wtb.topic
}
