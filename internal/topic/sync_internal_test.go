package topic

import (
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/is"
	"testing"
)

var (
	givenID      = uint64(42)
	givenName    = "givenName"
	givenPayload = "givenPayload"
)

func TestNewSyncBuilder(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	testInstance := NewSyncBuilder()
	assert.That(testInstance.topic.name, is.Empty())
	assert.That(testInstance.topic.handlers, is.Empty())
	assert.That(testInstance.topic.closed, is.False())
	assert.That(testInstance.topic.generateID, is.NotNil())
}

func TestSyncTopicBuilder_Name(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	testInstance := NewSyncBuilder()
	testInstance.Name(givenName)

	assert.That(testInstance.topic.name, is.EqualTo(givenName))
}

func TestSyncTopicImpl_Publish(t *testing.T) {

	t.Run("publish to closed topic", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := NewSyncBuilder().Build()
		_ = testInstance.Close()

		err := testInstance.Publish(givenPayload)

		assert.That(err, is.NotNil())
		assert.That(err, is.EqualTo(ErrAlreadyClosed))
	})

	t.Run("publish once", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		fakeGeneratorCounter := 0
		fakeGenerator := func() uint64 {
			fakeGeneratorCounter++
			return givenID
		}

		spyHandlerCounter := 0
		spyHandler := func(ID uint64, name string, payload interface{}) {
			assert.That(ID, is.EqualTo(givenID))
			assert.That(name, is.EqualTo(givenName))
			assert.That(payload, is.EqualTo(givenPayload))
			spyHandlerCounter++
		}

		testInstance := &syncTopicImpl{
			abstractTopicImpl: newAbstractTopicImpl(),
		}

		testInstance.name = givenName
		testInstance.generateID = fakeGenerator
		testInstance.handlers = append(testInstance.handlers, spyHandler)

		err := testInstance.Publish(givenPayload)
		assert.That(err, is.Nil())

		assert.That(fakeGeneratorCounter, is.EqualTo(1))
		assert.That(spyHandlerCounter, is.EqualTo(1))
	})

	t.Run("publish multiple", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		fakeGeneratorCounter := 0
		fakeGenerator := func() uint64 {
			fakeGeneratorCounter++
			return givenID
		}

		spyHandlerCounter := 0
		spyHandler := func(ID uint64, name string, payload interface{}) {
			assert.That(ID, is.EqualTo(givenID))
			assert.That(name, is.EqualTo(givenName))
			assert.That(payload, is.EqualTo(givenPayload))
			spyHandlerCounter++
		}

		testInstance := &syncTopicImpl{
			abstractTopicImpl: newAbstractTopicImpl(),
		}

		testInstance.name = givenName
		testInstance.generateID = fakeGenerator
		testInstance.handlers = append(testInstance.handlers, spyHandler)

		err := testInstance.Publish(givenPayload, givenPayload, givenPayload)
		assert.That(err, is.Nil())

		assert.That(fakeGeneratorCounter, is.EqualTo(3))
		assert.That(spyHandlerCounter, is.EqualTo(3))
	})
}
