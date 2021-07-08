package bus

type Event struct {
	ID      uint64
	Topic   string
	Payload interface{}
}

func newEvent(ID uint64, topic string, payload interface{}) Event {
	return Event{ID: ID, Topic: topic, Payload: payload}
}

type EventHandler func(event Event) error
