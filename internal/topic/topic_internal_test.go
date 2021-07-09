package topic

import (
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/has"
	"github.com/scarabsoft/go-hamcrest/is"
	"sync"
	"testing"
)

var (
	givenHandler   = func(ID uint64, name string, payload interface{}) {}
	anotherHandler = func(ID uint64, name string, payload interface{}) {}
)

func TestTopicIdGenerator(t *testing.T) {
	t.Parallel()
	assert := hamcrest.NewAssertion(t)

	testInstance := topicIdGenerator()

	iterations := 1000
	parallel := 100
	waiter := sync.WaitGroup{}
	dataMutex := sync.Mutex{}
	data := make(map[uint64]struct{})

	helpFn := func() {
		id := testInstance()
		dataMutex.Lock()
		data[id] = struct{}{}
		dataMutex.Unlock()
		waiter.Done()
	}

	for i := 0; i < iterations; i++ {
		for j := 0; j < parallel; j++ {
			waiter.Add(1)
			go helpFn()
		}
	}
	waiter.Wait()
	assert.That(data, has.Length(parallel*iterations))

	//make sure that each key exists
	for i := 0; i < parallel*iterations; i++ {
		assert.That(data, has.Key(uint64(i+1)))
	}

	nextId := testInstance()
	assert.That(nextId, is.EqualTo(uint64(parallel*iterations+1)))
}

func TestAbstractTopicImpl_Unsubscribe(t *testing.T) {

	t.Run("Unsubscribe from already closed topic", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := newAbstractTopicImpl()
		testInstance.closed = true

		err := testInstance.Unsubscribe(givenHandler)
		assert.That(err, is.NotNil())
		assert.That(err, is.EqualTo(ErrAlreadyClosed))
	})

	t.Run("Unsubscribe from empty topic", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := newAbstractTopicImpl()

		err := testInstance.Unsubscribe(givenHandler)
		assert.That(err, is.Nil())
	})

	t.Run("Unsubscribe not registered handler from topic", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := newAbstractTopicImpl()
		testInstance.handlers = append(testInstance.handlers, givenHandler)

		err := testInstance.Unsubscribe(anotherHandler)
		assert.That(err, is.Nil())

		assert.That(testInstance.handlers, has.Length(1))
		assert.That(testInstance.handlers, has.Item(givenHandler))
	})

	t.Run("Unsubscribe registered handle but keep other", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := newAbstractTopicImpl()
		testInstance.handlers = append(testInstance.handlers, givenHandler, anotherHandler)

		err := testInstance.Unsubscribe(givenHandler)
		assert.That(err, is.Nil())

		assert.That(testInstance.handlers, has.Length(1))
		assert.That(testInstance.handlers, has.Item(anotherHandler))
	})

	t.Run("Unsubscribe them all", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := newAbstractTopicImpl()
		testInstance.handlers = append(testInstance.handlers, givenHandler, anotherHandler)

		err := testInstance.Unsubscribe(givenHandler, anotherHandler)
		assert.That(err, is.Nil())

		assert.That(testInstance.handlers, is.Empty())
	})
}

func TestAbstractTopicImpl_Close(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	testInstance := newAbstractTopicImpl()

	err := testInstance.Close()
	assert.That(err, is.Nil())

	assert.That(testInstance.closed, is.True())
}

func TestAbstractTopicImpl_New(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	testInstance := newAbstractTopicImpl()

	assert.That(testInstance.closed, is.False())
	assert.That(testInstance.generateID, is.NotNil())
	assert.That(testInstance.generateID(), is.EqualTo(uint64(1)))
	assert.That(testInstance.handlers, is.Empty())
}

func TestAbstractTopicImpl_Subscribe(t *testing.T) {
	t.Run("Subscribe to closed topic", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := newAbstractTopicImpl()
		testInstance.closed = true

		err := testInstance.Subscribe(givenHandler)
		assert.That(err, is.NotNil())
		assert.That(err, is.EqualTo(ErrAlreadyClosed))

		assert.That(testInstance.handlers, is.Empty())
	})

	t.Run("Subscribe to topic", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := newAbstractTopicImpl()

		err := testInstance.Subscribe(givenHandler)
		assert.That(err, is.Nil())

		assert.That(testInstance.handlers, has.Length(1))
		assert.That(testInstance.handlers, has.Item(givenHandler))
	})

	t.Run("Subscribe them all to topic", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := newAbstractTopicImpl()

		err := testInstance.Subscribe(givenHandler, anotherHandler)
		assert.That(err, is.Nil())

		assert.That(testInstance.handlers, has.Length(2))
		assert.That(testInstance.handlers, has.Items(givenHandler, anotherHandler))
	})

	t.Run("Subscribe twice and stop execution", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := newAbstractTopicImpl()

		_ = testInstance.Subscribe(givenHandler)
		err := testInstance.Subscribe(givenHandler, anotherHandler)
		assert.That(err, is.NotNil())
		assert.That(err, is.EqualTo(ErrAlreadySubscribed))

		assert.That(testInstance.handlers, has.Length(1))
		assert.That(testInstance.handlers, has.Item(givenHandler))
	})
}
