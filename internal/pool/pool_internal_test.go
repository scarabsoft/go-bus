package pool

import (
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/has"
	"github.com/scarabsoft/go-hamcrest/is"
	"sync/atomic"
	"testing"
	"time"
)

const (
	maxWorkers   = 1
	maxQueueSize = 1
)

var givenOptions = Options{
	MaxWorkers:   maxWorkers,
	MaxQueueSize: maxQueueSize,
}

func TestNew(t *testing.T) {
	assert := hamcrest.NewAssertion(t)

	testInstance := New(givenOptions)

	assert.That(testInstance, is.NotNil())
	assert.That(testInstance.closed, is.False())
	assert.That(testInstance.queue, is.Empty())
	assert.That(testInstance.options, is.EqualTo(givenOptions))

	assert.That(testInstance.Close(), is.Nil())
}

func TestSubmit(t *testing.T) {
	t.Run("pool is already closed", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := New(givenOptions)
		_ = testInstance.Close()

		err := testInstance.Submit(Task{})
		assert.That(err, is.EqualTo(ErrPoolIsClosed))
	})

	t.Run("submit and execute task", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := New(givenOptions)

		handlerCalled := false
		givenId := uint64(42)
		givenName := "givenName"
		givenPayload := "givenPayload"

		givenHandler := func(ID uint64, name string, payload interface{}) {
			assert.That(ID, is.EqualTo(givenId))
			assert.That(name, is.EqualTo(givenName))
			assert.That(payload, is.EqualTo(givenPayload))
			handlerCalled = true
		}

		err := testInstance.Submit(Task{Payload: TaskPayload{ID: givenId, Name: givenName, Payload: givenPayload}, Handler: givenHandler})
		assert.That(err, is.Nil())
		_ = testInstance.Close()

		assert.That(handlerCalled, is.True())
	})

	t.Run("submit more jobs than queue size and execute all", func(t *testing.T) {
		assert := hamcrest.NewAssertion(t)

		testInstance := New(Options{MaxWorkers: 2, MaxQueueSize: 2})
		defer testInstance.Close()

		maxIterations := 100
		idChan := make(chan uint64, maxIterations)
		handlerCounter := uint64(0)
		givenName := "givenName"
		givenPayload := "givenPayload"

		givenHandler := func(ID uint64, name string, payload interface{}) {
			assert.That(name, is.EqualTo(givenName))
			assert.That(payload, is.EqualTo(givenPayload))
			atomic.AddUint64(&handlerCounter, 1)
			idChan <- ID
		}

		for i := 0; i < maxIterations; i++ {
			err := testInstance.Submit(Task{Payload: TaskPayload{ID: uint64(i), Name: givenName, Payload: givenPayload}, Handler: givenHandler})
			assert.That(err, is.Nil())
		}

		handledIds := []uint64{}

		for {
			select {
			case id := <-idChan:
				handledIds = append(handledIds, id)
				break
			case <-time.After(250 * time.Millisecond):
				assert.That(handlerCounter, is.EqualTo(uint64(maxIterations)))
				assert.That(handledIds, has.Length(maxIterations))

				for i := 1; i < maxIterations-1; i++ {
					if handledIds[i-1] > handledIds[i] {
						return
					}
				}
				t.Fatal("expected unordered id slice as execution was concurrent")
				return
			}
		}
	})
}
