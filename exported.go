package go_bus

var (
	std = new()
)

func Get(topic string) Topic {
	return std.Get(topic)
}

func CreateTopic(name string, options ...Option) error {
	return std.CreateTopic(name, options...)
}
