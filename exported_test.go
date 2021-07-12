package bus_test

import (
	"github.com/scarabsoft/go-bus"
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/is"
	"testing"
)

const (
	givenID      uint64 = 42
	givenTopic          = "TestTopic"
	givenPayload        = "Payload"
)

func TestEvent_String(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	givenEvent := bus.Event{
		ID:      givenID,
		Topic:   givenTopic,
		Payload: givenPayload,
	}

	assert.That(givenEvent.String(), is.EqualTo("[TestTopic-42]: string(Payload)"))
}

func TestEventHandler(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	bus.EventHandler(func(event bus.Event) {
		assert.That(event.ID, is.EqualTo(givenID))
		assert.That(event.Topic, is.EqualTo(givenTopic))
		assert.That(event.Payload, is.EqualTo(givenPayload))
	})(givenID, givenTopic, givenPayload)
}
