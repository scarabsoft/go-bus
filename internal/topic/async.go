package topic

type asyncTopicImpl struct {
	abstractTopicImpl
}

func (a *asyncTopicImpl) Publish(payloads ...interface{}) error {
	a.lock.RLock()
	if err := a.abstractTopicImpl.Publish(payloads); err != nil {
		return err
	}

	go func() {
		defer a.lock.RUnlock()

		for _, payload := range payloads {
			id := a.generateID()
			for _, handler := range a.handlers {
				handler(id, a.topic, payload)
			}
		}
	}()
	return nil
}

type asyncTopicBuilder struct {
	topic asyncTopicImpl
}

func NewAsyncBuilder(name string) *asyncTopicBuilder {
	return &asyncTopicBuilder{topic: asyncTopicImpl{
		abstractTopicImpl: newAbstractTopicImpl(name),
	}}
}

func (atb *asyncTopicBuilder) Build() Topic {
	return &atb.topic
}
