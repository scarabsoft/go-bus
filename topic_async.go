package bus

type asyncTopic struct {
	abstractTopicImpl
}

func newAsyncTopic(name string) Topic {
	return &asyncTopic{
		abstractTopicImpl: newAbstractTopicImpl(name),
	}
}

func (a *asyncTopic) Publish(data interface{}) error {
	a.lock.RLock()
	defer a.lock.RUnlock()

	if a.closed {
		return ErrTopicClosed
	}

	go func() {
		//e := Event{Topic: a.name, Payload: data} // FIXME this should be a simple event generator to have auto increment ids
		e := newEvent(a.idGenerator(), a.name, data) // FIXME this should be a simple event generator to have auto increment ids
		for _, handler := range a.handlers {
			_ = handler(e)
		}
	}()
	return nil
}

//type AsyncTopicBuilder struct{}
//
//func (atb* AsyncTopicBuilder) Build() Topic{
//	return newAsyncTopic()
//}