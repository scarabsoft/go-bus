package bus

var (
	std = new()
)

func Get(topic string) Topic {
	return std.Get(topic)
}

func CreateTopic(name string) error {
	return std.CreateTopic(name)
}
