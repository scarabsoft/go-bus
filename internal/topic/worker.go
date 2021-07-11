package topic

import (
	"github.com/scarabsoft/go-bus/internal/pool"
)

type workerTopicImpl struct {
	abstractTopicImpl

	pool pool.Pool
}

func (w *workerTopicImpl) Publish(payloads ...interface{}) error {
	w.lock.RLock()
	defer w.lock.RUnlock()

	if w.closed {
		return ErrorAlreadyClosed{w.name}
	}

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

type workerTopicBuilderImpl struct {
	topic workerTopicImpl
}

func NewWorkerBuilder(defaultPool pool.Pool) *workerTopicBuilderImpl {
	return &workerTopicBuilderImpl{topic: workerTopicImpl{
		abstractTopicImpl: newAbstractTopicImpl(),
		pool:              defaultPool,
	}}
}

func (wtb *workerTopicBuilderImpl) Name(name string) Builder {
	wtb.topic.name = name
	return wtb
}

func (wtb *workerTopicBuilderImpl) Pool(p pool.Pool) *workerTopicBuilderImpl {
	wtb.topic.pool = p
	return wtb
}

func (wtb *workerTopicBuilderImpl) Build() Topic {
	return &wtb.topic
}
