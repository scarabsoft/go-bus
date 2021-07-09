package bus

type asyncTopicImpl struct {
	abstractTopicImpl
}

func (a *asyncTopicImpl) Publish(payloads ...interface{}) error {
	a.lock.RLock()
	defer a.lock.RUnlock()

	if a.closed {
		return ErrTopicClosed
	}

	go func() {
		for _, payload := range payloads {
			e := newEvent(a.idGenerator, a.name, payload)
			for _, handler := range a.handlers {
				_ = handler(e)
			}
		}
	}()
	return nil
}
func (a *asyncTopicImpl) Close() error {
	return nil
}

type asyncTopicBuilder struct {
	topic asyncTopicImpl
}

func newAsyncTopicBuilder(name string) *asyncTopicBuilder {
	return &asyncTopicBuilder{topic: asyncTopicImpl{
		abstractTopicImpl: newAbstractTopicImpl(name),
	}}
}

func (atb *asyncTopicBuilder) build() Topic {
	return &atb.topic
}
