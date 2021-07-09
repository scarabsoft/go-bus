package bus

var (
	std = newBus()
)

func Get(topic string) Topic {
	return std.Get(topic)
}

func Publish(topic string, data ...interface{}) (Topic, error) {
	return std.Publish(topic, data...)
}

func CreateTopic(name string, fn func(topic *TopicBuilder) Topic) (Topic, error) {
	return std.CreateTopic(name, fn)
}
