package bus

type syncTopicImpl struct {
	abstractTopicImpl
}

func (s *syncTopicImpl) Publish(payloads ...interface{}) error {
	if err := s.abstractTopicImpl.Publish(payloads); err != nil {
		return err
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	for _, payload := range payloads {
		evt := newEvent(s.idGenerator, s.name, payload)
		for _, handler := range s.handlers {
			_ = handler(evt)
		}
	}
	return nil
}

func (s *syncTopicImpl) Close() error {
	return nil
}

type syncTopicBuilder struct {
	topic syncTopicImpl
}

func newSyncTopicBuilder(name string) *syncTopicBuilder {
	return &syncTopicBuilder{topic: syncTopicImpl{
		abstractTopicImpl: newAbstractTopicImpl(name),
	}}
}

func (stb *syncTopicBuilder) build() Topic {
	return &stb.topic
}
