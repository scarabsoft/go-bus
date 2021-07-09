package bus

type syncTopic struct {
	abstractTopicImpl
}

func (s *syncTopic) Publish(data ...interface{}) error {
	if err := s.abstractTopicImpl.Publish(data); err != nil {
		return err
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	for _, d := range data {
		e := newEvent(s.idGenerator(), s.name, d)
		for _, handler := range s.handlers {
			_ = handler(e)
		}
	}
	return nil
}

func (s *syncTopic) Close() error {
	return nil
}

type SyncTopicBuilder struct {
	topic syncTopic
}

func newSyncTopicBuilder(name string) *SyncTopicBuilder {
	return &SyncTopicBuilder{topic: syncTopic{
		abstractTopicImpl: newAbstractTopicImpl(name),
	}}
}

func (stb *SyncTopicBuilder) Build() Topic {
	return &stb.topic
}
