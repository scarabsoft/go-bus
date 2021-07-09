package topic

type workerTopicImpl struct {
	abstractTopicImpl
}

func (w *workerTopicImpl) Public(payloads ...interface{}) error {
	return nil
}

type workerTopicBuilder struct {
	topic workerTopicImpl
}

func NewWorkerBuilder(name string) *workerTopicBuilder {
	return &workerTopicBuilder{topic: workerTopicImpl{
		abstractTopicImpl: newAbstractTopicImpl(name),
	}}
}

func (wtb *workerTopicBuilder) Build() Topic {
	return &wtb.topic
}
