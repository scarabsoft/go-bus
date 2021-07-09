package bus

var (
	std = newBus()

	SyncTopic = func(topic TopicInit) TopicBuilder {
		return newSyncTopicBuilder(topic.name)
	}

	AsyncTopic = func(topic TopicInit) TopicBuilder {
		return newAsyncTopicBuilder(topic.name)
	}
)

func Get(topic string) Topic {
	return std.Get(topic)
}

func Publish(topic string, payloads ...interface{}) (Topic, error) {
	return std.Publish(topic, payloads...)
}

func CreateTopic(name string, fn func(topic TopicInit) TopicBuilder) (Topic, error) {
	return std.CreateTopic(name, fn)
}
