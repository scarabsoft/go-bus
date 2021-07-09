package topic

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
		//evt := event.New(s.generateID, s.topic, payload)
		id := s.generateID()
		for _, handler := range s.handlers {
			handler(id, s.topic, payload)
		}
	}
	return nil
}

type syncTopicBuilder struct {
	topic syncTopicImpl
}

func NewSyncTopicBuilder(name string) *syncTopicBuilder {
	return &syncTopicBuilder{topic: syncTopicImpl{
		abstractTopicImpl: newAbstractTopicImpl(name),
	}}
}

func (stb *syncTopicBuilder) Build() Topic {
	return &stb.topic
}
