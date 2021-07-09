package topic

import (
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/has"
	"github.com/scarabsoft/go-hamcrest/is"
	"sync"
	"testing"
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
