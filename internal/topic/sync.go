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
		id := s.generateID()
		for _, handler := range s.handlers {
			handler(id, s.name, payload)
		}
	}
	return nil
}

type syncTopicBuilder struct {
	topic syncTopicImpl
}

func NewSyncBuilder() *syncTopicBuilder {
	return &syncTopicBuilder{topic: syncTopicImpl{
		abstractTopicImpl: newAbstractTopicImpl(),
	}}
}

func (stb *syncTopicBuilder) Name(name string) Builder {
	stb.topic.name = name
	return stb
}

func (stb *syncTopicBuilder) Build() Topic {
	return &stb.topic
}
