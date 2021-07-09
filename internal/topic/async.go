package topic

import "github.com/scarabsoft/go-bus/internal/pool"

type AsyncTopicBuilder struct {
	topic workerTopicImpl
}

func NewAsyncBuilder(defaultPool pool.Pool) *AsyncTopicBuilder {
	return &AsyncTopicBuilder{topic: workerTopicImpl{
		abstractTopicImpl: newAbstractTopicImpl(),
		pool:              defaultPool,
	}}
}

func (wtb *AsyncTopicBuilder) Name(name string) Builder {
	wtb.topic.name = name
	return wtb
}

func (wtb *AsyncTopicBuilder) Build() Topic {
	return &wtb.topic
}
