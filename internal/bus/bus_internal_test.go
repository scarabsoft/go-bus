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

func (s *spyTopicImpl) Publish(...interface{}) error {
	s.PublishCount++
	return nil
}

func (s *spyTopicImpl) Subscribe(...func(ID uint64, topic string, payload interface{})) error {
	s.SubscribeCount++
	return nil
}

func (s *spyTopicImpl) Unsubscribe(...func(ID uint64, topic string, payload interface{})) error {
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

func (s spyTopicBuilder) Name(string) topic.Builder {
	return s
}

func (s spyTopicBuilder) Build() topic.Topic {
	return s.topic
}

var ErrTestPublishError = errors.New("ErrTestPublishError")
var ErrTestSubscribeError = errors.New("ErrTestSubscribeError")
var ErrTestUnsubscribeError = errors.New("ErrTestUnsubscribeError")
var ErrTestCloseError = errors.New("ErrTestCloseError")

type failingTopic struct {
	spyTopicImpl
}

func (f *failingTopic) Name() string {
	_ = f.spyTopicImpl.Name()
	return ""
}

func (f *failingTopic) Publish(...interface{}) error {
	_ = f.spyTopicImpl.Publish()
	return ErrTestPublishError
}

func (f *failingTopic) Subscribe(...func(ID uint64, topic string, payload interface{})) error {
	_ = f.spyTopicImpl.Subscribe()
	return ErrTestSubscribeError
}

func (f *failingTopic) Unsubscribe(...func(ID uint64, topic string, payload interface{})) error {
	_ = f.spyTopicImpl.Unsubscribe()
	return ErrTestUnsubscribeError
}

func (f *failingTopic) Close() error {
	_ = f.spyTopicImpl.Close()
	return ErrTestCloseError
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

		assert.That(err, is.EqualTo(topic.ErrorNotExists{Name: givenTopicName}))
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

		assert.That(err, is.EqualTo(topic.ErrorNotExists{Name: givenTopicName}))
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
		assert.That(err, is.EqualTo(topic.ErrorNotExists{Name: givenTopicName}))

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
		assert.That(err, is.EqualTo(ErrTestPublishError))

	})
}

func TestSubscribe(t *testing.T) {
	t.Run("subscribe with topic does not exists and no default topic builder provided", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := New()

		err := testInstance.Subscribe(givenTopicName, givenHandler)
		assert.That(err, is.EqualTo(topic.ErrorNotExists{Name: givenTopicName}))

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

	t.Run("subscribing failed", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		failTopic := &failingTopic{}

		testInstance := New()
		testInstance.topics[givenTopicName] = failTopic

		err := testInstance.Subscribe(givenTopicName, givenHandler, givenHandler, givenHandler)
		assert.That(err, is.NotNil())
		assert.That(err, is.EqualTo(ErrTestSubscribeError))

		assert.That(failTopic.PublishCount, is.EqualTo(0))
		assert.That(failTopic.CloseCount, is.EqualTo(0))
		assert.That(failTopic.NameCount, is.EqualTo(0))
		assert.That(failTopic.SubscribeCount, is.EqualTo(1))
		assert.That(failTopic.UnsubscribeCount, is.EqualTo(0))

	})
}

func TestUnsubscribe(t *testing.T) {
	t.Run("unsubscribe from topic which does not exists", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := New()

		err := testInstance.Unsubscribe(givenTopicName, givenHandler)
		assert.That(err, is.EqualTo(topic.ErrorNotExists{Name: givenTopicName}))

		assert.That(testInstance.topics, is.Empty())
	})

	t.Run("unsubscribe from topic", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		spyTopic := &spyTopicImpl{}

		testInstance := New()
		testInstance.topics[givenTopicName] = spyTopic

		err := testInstance.Unsubscribe(givenTopicName, givenHandler)
		assert.That(err, is.Nil())

		assert.That(spyTopic.PublishCount, is.EqualTo(0))
		assert.That(spyTopic.CloseCount, is.EqualTo(0))
		assert.That(spyTopic.NameCount, is.EqualTo(0))
		assert.That(spyTopic.SubscribeCount, is.EqualTo(0))
		assert.That(spyTopic.UnsubscribeCount, is.EqualTo(1))
	})

	t.Run("unsubscribing failed", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		failTopic := &failingTopic{}

		testInstance := New()
		testInstance.topics[givenTopicName] = failTopic

		err := testInstance.Unsubscribe(givenTopicName, givenHandler, givenHandler, givenHandler)
		assert.That(err, is.NotNil())
		assert.That(err, is.EqualTo(ErrTestUnsubscribeError))

		assert.That(failTopic.PublishCount, is.EqualTo(0))
		assert.That(failTopic.CloseCount, is.EqualTo(0))
		assert.That(failTopic.NameCount, is.EqualTo(0))
		assert.That(failTopic.SubscribeCount, is.EqualTo(0))
		assert.That(failTopic.UnsubscribeCount, is.EqualTo(1))
	})
}

func TestBusImpl_DeleteTopic(t *testing.T) {
	t.Run("delete topic which does not exists", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := New()

		err := testInstance.DeleteTopic(givenTopicName)
		assert.That(err, is.EqualTo(topic.ErrorNotExists{Name: givenTopicName}))

		assert.That(testInstance.topics, is.Empty())
	})

	t.Run("delete topic", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		spyTopic := &spyTopicImpl{}

		testInstance := New()
		testInstance.topics[givenTopicName] = spyTopic

		err := testInstance.DeleteTopic(givenTopicName)
		assert.That(err, is.Nil())

		assert.That(spyTopic.PublishCount, is.EqualTo(0))
		assert.That(spyTopic.CloseCount, is.EqualTo(1))
		assert.That(spyTopic.NameCount, is.EqualTo(0))
		assert.That(spyTopic.SubscribeCount, is.EqualTo(0))
		assert.That(spyTopic.UnsubscribeCount, is.EqualTo(0))

		assert.That(testInstance.topics, is.Empty())
	})

	t.Run("deleting failed", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		failTopic := &failingTopic{}

		testInstance := New()
		testInstance.topics[givenTopicName] = failTopic

		err := testInstance.DeleteTopic(givenTopicName)
		assert.That(err, is.NotNil())
		assert.That(err, is.EqualTo(ErrTestCloseError))

		assert.That(failTopic.PublishCount, is.EqualTo(0))
		assert.That(failTopic.CloseCount, is.EqualTo(1))
		assert.That(failTopic.NameCount, is.EqualTo(0))
		assert.That(failTopic.SubscribeCount, is.EqualTo(0))
		assert.That(failTopic.UnsubscribeCount, is.EqualTo(0))
	})
}
