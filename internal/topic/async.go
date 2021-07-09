package topic

import "github.com/scarabsoft/go-bus/internal/pool"

type asyncTopicBuilderImpl struct {
	topic workerTopicImpl
}

func NewAsyncBuilder(defaultPool pool.Pool) *asyncTopicBuilderImpl {
	return &asyncTopicBuilderImpl{topic: workerTopicImpl{
		abstractTopicImpl: newAbstractTopicImpl(),
		pool:              defaultPool,
	}}
}

func (wtb *asyncTopicBuilderImpl) Name(name string) Builder {
	wtb.topic.name = name
	return wtb
}

func (wtb *asyncTopicBuilderImpl) Build() Topic {
	return &wtb.topic
}
