package bus

import (
	"errors"
	"github.com/scarabsoft/go-bus/internal/topic"
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/has"
	"github.com/scarabsoft/go-hamcrest/is"
	"testing"
)

var (
	givenTopicName = "givenTopicName"
	givenTopic     = topic.NewSyncBuilder().Build()
	givenPayload   = "givenPayload"
	givenHandler   = func(ID uint64, name string, payload interface{}) {}
)

type testTopicBuilder struct {
	topic topic.Topic
}

func (t testTopicBuilder) Name(name string) topic.Builder {
	return t
}

func (t testTopicBuilder) Build() topic.Topic {
	return t.topic
}

type spyTopicImpl struct {
	NameCount        int
	PublishCount     int
	SubscribeCount   int
	UnsubscribeCount int
	CloseCount       int
}

func (s *spyTopicImpl) Name() string {
	s.NameCount++
	return givenTopicName
}

func (s *spyTopicImpl) Publish(payloads ...interface{}) error {
	s.PublishCount++
	return nil
}

func (s *spyTopicImpl) Subscribe(handlers ...func(ID uint64, topic string, payload interface{})) error {
	s.SubscribeCount++
	return nil
}

func (s *spyTopicImpl) Unsubscribe(handlers ...func(ID uint64, topic string, payload interface{})) error {
	s.UnsubscribeCount++
	return nil
}

func (s *spyTopicImpl) Close() error {
	s.CloseCount++
	return nil
}

type spyTopicBuilder struct {
	topic *spyTopicImpl
}

func (s spyTopicBuilder) Name(name string) topic.Builder {
	return s
}

func (s spyTopicBuilder) Build() topic.Topic {
	return s.topic
}

var ErrTestError = errors.New("some test error")

type failingTopic struct{}

func (f failingTopic) Name() string {
	return ""
}

func (f failingTopic) Publish(payloads ...interface{}) error {
	return ErrTestError
}

func (f failingTopic) Subscribe(handlers ...func(ID uint64, topic string, payload interface{})) error {
	return ErrTestError
}

func (f failingTopic) Unsubscribe(handlers ...func(ID uint64, topic string, payload interface{})) error {
	return ErrTestError
}

func (f failingTopic) Close() error {
	return ErrTestError
}

func TestCreateTopic(t *testing.T) {
	t.Run("topic found", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := New()
		testInstance.topics[givenTopicName] = givenTopic

		res, err := testInstance.CreateTopic(givenTopicName, nil)

		assert.That(err, is.Nil())
		assert.That(res, is.EqualTo(givenTopic))
	})

	t.Run("topic not found and no builder provided", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := New()

		res, err := testInstance.CreateTopic(givenTopicName, nil)

		assert.That(err, is.EqualTo(topic.ErrDoesNotExists))
		assert.That(res, is.Nil())

		assert.That(testInstance.topics, is.Empty())
	})

	t.Run("topic not found but topic builder provided", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := New()

		res, err := testInstance.CreateTopic(givenTopicName, testTopicBuilder{topic: givenTopic})
		assert.That(err, is.Nil())
		assert.That(res, is.EqualTo(givenTopic))

		assert.That(testInstance.topics, has.Length(1))
		assert.That(testInstance.topics[givenTopicName], is.EqualTo(givenTopic))
	})
}

func TestGetTopic(t *testing.T) {
	t.Run("topic found", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := New()
		testInstance.topics[givenTopicName] = givenTopic

		res, err := testInstance.Get(givenTopicName)

		assert.That(err, is.Nil())
		assert.That(res, is.EqualTo(givenTopic))
	})

	t.Run("topic not found and no default builder provided", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := New()

		res, err := testInstance.Get(givenTopicName)

		assert.That(err, is.EqualTo(topic.ErrDoesNotExists))
		assert.That(res, is.Nil())

		assert.That(testInstance.topics, is.Empty())
	})

	t.Run("topic not found but default topic builder provided", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := New()
		testInstance.SetDefaultTopicBuilder(testTopicBuilder{topic: givenTopic})

		res, err := testInstance.Get(givenTopicName)
		assert.That(err, is.Nil())
		assert.That(res, is.EqualTo(givenTopic))

		assert.That(testInstance.topics, has.Length(1))
		assert.That(testInstance.topics[givenTopicName], is.EqualTo(givenTopic))
	})
}
func TestPublish(t *testing.T) {
	t.Run("publish with topic does not exists and no default topic builder provided", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := New()

		err := testInstance.Publish(givenTopicName, givenPayload)
		assert.That(err, is.EqualTo(topic.ErrDoesNotExists))

		assert.That(testInstance.topics, is.Empty())
	})

	t.Run("publish with topic not found but default topic builder provided", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		spyTopic := &spyTopicImpl{}

		testInstance := New()
		testInstance.SetDefaultTopicBuilder(spyTopicBuilder{topic: spyTopic})

		err := testInstance.Publish(givenTopicName, givenPayload, givenPayload, givenPayload)
		assert.That(err, is.Nil())

		assert.That(testInstance.topics, has.Length(1))
		assert.That(testInstance.topics[givenTopicName], is.EqualTo(spyTopic))

		assert.That(spyTopic.PublishCount, is.EqualTo(1))
		assert.That(spyTopic.CloseCount, is.EqualTo(0))
		assert.That(spyTopic.NameCount, is.EqualTo(0))
		assert.That(spyTopic.SubscribeCount, is.EqualTo(0))
		assert.That(spyTopic.UnsubscribeCount, is.EqualTo(0))
	})

	t.Run("publish with topic exists", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		spyTopic := &spyTopicImpl{}

		testInstance := New()
		testInstance.topics[givenTopicName] = spyTopic

		err := testInstance.Publish(givenTopicName, givenPayload, givenPayload, givenPayload)
		assert.That(err, is.Nil())

		assert.That(spyTopic.PublishCount, is.EqualTo(1))
		assert.That(spyTopic.CloseCount, is.EqualTo(0))
		assert.That(spyTopic.NameCount, is.EqualTo(0))
		assert.That(spyTopic.SubscribeCount, is.EqualTo(0))
		assert.That(spyTopic.UnsubscribeCount, is.EqualTo(0))
	})

	t.Run("publishing failed", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := New()
		testInstance.topics[givenTopicName] = &failingTopic{}

		err := testInstance.Publish(givenTopicName, givenPayload, givenPayload, givenPayload)
		assert.That(err, is.NotNil())
		assert.That(err, is.EqualTo(ErrTestError))

	})
}

func TestSubscribe(t *testing.T) {
	t.Run("subscribe with topic does not exists and no default topic builder provided", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := New()

		err := testInstance.Subscribe(givenTopicName, givenHandler)
		assert.That(err, is.EqualTo(topic.ErrDoesNotExists))

		assert.That(testInstance.topics, is.Empty())
	})

	t.Run("subscribe with topic not found but default topic builder provided", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		spyTopic := &spyTopicImpl{}

		testInstance := New()
		testInstance.SetDefaultTopicBuilder(spyTopicBuilder{topic: spyTopic})

		err := testInstance.Subscribe(givenTopicName, givenHandler, givenHandler, givenHandler)
		assert.That(err, is.Nil())

		assert.That(testInstance.topics, has.Length(1))
		assert.That(testInstance.topics[givenTopicName], is.EqualTo(spyTopic))

		assert.That(spyTopic.PublishCount, is.EqualTo(0))
		assert.That(spyTopic.CloseCount, is.EqualTo(0))
		assert.That(spyTopic.NameCount, is.EqualTo(0))
		assert.That(spyTopic.SubscribeCount, is.EqualTo(1))
		assert.That(spyTopic.UnsubscribeCount, is.EqualTo(0))
	})

	t.Run("subscribe with topic exists", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		spyTopic := &spyTopicImpl{}

		testInstance := New()
		testInstance.topics[givenTopicName] = spyTopic

		err := testInstance.Subscribe(givenTopicName, givenHandler, givenHandler, givenHandler)
		assert.That(err, is.Nil())

		assert.That(spyTopic.PublishCount, is.EqualTo(0))
		assert.That(spyTopic.CloseCount, is.EqualTo(0))
		assert.That(spyTopic.NameCount, is.EqualTo(0))
		assert.That(spyTopic.SubscribeCount, is.EqualTo(1))
		assert.That(spyTopic.UnsubscribeCount, is.EqualTo(0))
	})

	t.Run("publishing failed", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := New()
		testInstance.topics[givenTopicName] = &failingTopic{}

		err := testInstance.Subscribe(givenTopicName, givenHandler, givenHandler, givenHandler)
		assert.That(err, is.NotNil())
		assert.That(err, is.EqualTo(ErrTestError))

	})
}