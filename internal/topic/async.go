package topic

type asyncTopicImpl struct {
	abstractTopicImpl
}

func (a *asyncTopicImpl) Publish(payloads ...interface{}) error {

	if err := a.abstractTopicImpl.Publish(payloads); err != nil {
		return err
	}

	go func() {
		a.lock.RLock()
		defer a.lock.RUnlock()

		for _, payload := range payloads {
			//e := event.New(a.generateID, a.topic, payload)
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

func NewAsyncTopicBuilder(name string) *asyncTopicBuilder {
	return &asyncTopicBuilder{topic: asyncTopicImpl{
		abstractTopicImpl: newAbstractTopicImpl(name),
	}}
}

func (atb *asyncTopicBuilder) Build() Topic {
	return &atb.topic
}
