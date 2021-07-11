package topic

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/scarabsoft/go-bus/internal/pool"
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/is"
	"testing"
)

var testErrPoolSubmission = errors.New("some test submission error")

func TestNewAsyncBuilder(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	givenPool := pool.New(pool.Options{})

	testInstance := NewAsyncBuilder(givenPool)
	assert.That(testInstance.topic.name, is.Empty())
	assert.That(testInstance.topic.handlers, is.Empty())
	assert.That(testInstance.topic.closed, is.False())
	assert.That(testInstance.topic.generateID, is.NotNil())
	assert.That(testInstance.topic.pool, is.EqualTo(givenPool))
}

func TestAsyncTopicBuilder_Name(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	givenPool := pool.New(pool.Options{})

	testInstance := NewAsyncBuilder(givenPool)
	testInstance.Name(givenName)

	assert.That(testInstance.topic.name, is.EqualTo(givenName))
}

func TestAsyncTopicBuilder_Pool(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	givenPool := pool.New(pool.Options{})

	testInstance := NewAsyncBuilder(nil)
	testInstance.Pool(givenPool)

	assert.That(testInstance.topic.pool, is.EqualTo(givenPool))
}

func TestNewWorkerBuilder(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	givenPool := pool.New(pool.Options{})

	testInstance := NewWorkerBuilder(givenPool)
	assert.That(testInstance.topic.name, is.Empty())
	assert.That(testInstance.topic.handlers, is.Empty())
	assert.That(testInstance.topic.closed, is.False())
	assert.That(testInstance.topic.generateID, is.NotNil())
	assert.That(testInstance.topic.pool, is.EqualTo(givenPool))
}

func TestWorkerTopicBuilder_Name(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	givenPool := pool.New(pool.Options{})

	testInstance := NewWorkerBuilder(givenPool)
	testInstance.Name(givenName)

	assert.That(testInstance.topic.name, is.EqualTo(givenName))
}

func TestWorkerTopicBuilder_Pool(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	givenPool := pool.New(pool.Options{})

	testInstance := NewWorkerBuilder(nil)
	testInstance.Pool(givenPool)

	assert.That(testInstance.topic.pool, is.EqualTo(givenPool))
}

func TestWorkerTopicImpl_Publish(t *testing.T) {
	t.Run("already closed", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		builder := NewWorkerBuilder(nil)
		builder.topic.name = givenName
		builder.topic.closed = true
		testInstance := builder.Build()

		err := testInstance.Publish(givenPayload)
		assert.That(err, is.NotNil())
		assert.That(err, is.EqualTo(ErrorAlreadyClosed{Name: givenName}))
	})

	t.Run("publish", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		givenHandler := func(ID uint64, name string, payload interface{}) {}

		mockedPool := pool.NewMockPool(ctrl)
		mockedPool.EXPECT().Submit(gomock.Any()).Times(1).Return(nil).Do(func(arg interface{}) {
			if task, ok := arg.(pool.Task); !ok {
				t.Fatal("arg is not a pool.Task")
			} else {
				assert.That(task.Payload.ID, is.EqualTo(uint64(1)))
				assert.That(task.Payload.Name, is.EqualTo(givenName))
				assert.That(task.Payload.Payload, is.EqualTo(givenPayload))

				assert.That(task.Handler, is.EqualTo(givenHandler))
			}
		})

		builder := NewWorkerBuilder(mockedPool)
		builder.topic.name = givenName
		builder.topic.handlers = append(builder.topic.handlers, givenHandler)
		testInstance := builder.Build()

		err := testInstance.Publish(givenPayload)
		assert.That(err, is.Nil())
	})

	t.Run("publish but no handler registered", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockedPool := pool.NewMockPool(ctrl)

		builder := NewWorkerBuilder(mockedPool)
		builder.topic.name = givenName
		testInstance := builder.Build()

		err := testInstance.Publish(givenPayload)
		assert.That(err, is.Nil())

	})

	t.Run("publish interrupt in case of error", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		givenHandler := func(ID uint64, name string, payload interface{}) {}

		mockedPool := pool.NewMockPool(ctrl)
		mockedPool.EXPECT().Submit(gomock.Any()).Times(1).Return(testErrPoolSubmission).Do(func(arg interface{}) {
			if task, ok := arg.(pool.Task); !ok {
				t.Fatal("arg is not a pool.Task")
			} else {
				assert.That(task.Payload.ID, is.EqualTo(uint64(1)))
				assert.That(task.Payload.Name, is.EqualTo(givenName))
				assert.That(task.Payload.Payload, is.EqualTo(givenPayload))

				assert.That(task.Handler, is.EqualTo(givenHandler))
			}
		})

		builder := NewWorkerBuilder(mockedPool)
		builder.topic.name = givenName
		builder.topic.handlers = append(builder.topic.handlers, givenHandler)
		testInstance := builder.Build()

		err := testInstance.Publish(givenPayload)
		assert.That(err, is.EqualTo(testErrPoolSubmission))
	})
}
