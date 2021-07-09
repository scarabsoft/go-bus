package bus

type Event struct {
	ID      uint64
	Topic   string
	Payload interface{}
}

func newEvent(generateId func() uint64, topic string, payload interface{}) Event {
	return Event{ID: generateId(), Topic: topic, Payload: payload}
}

type EventHandler func(event Event) error
