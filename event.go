package go_bus

type Event struct {
	Topic   string
	Payload interface{}
}

type EventHandler func(event Event) error
