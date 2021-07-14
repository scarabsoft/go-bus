package examples

import (
	"fmt"
	"github.com/scarabsoft/go-bus"
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/is"
	"strconv"
	"strings"
	"testing"
)

func TestAutoCreate(t *testing.T) {
	// we want every not existing topic run on the same worker pool by that we are able to configure
	// tasks running together and find good queue sizes
	p := bus.NewPool(
		bus.WithMaxQueueSize(4),
		bus.WithMaxWorkers(2),
	)

	require := hamcrest.NewRequirement(t)

	// every not existing topic should be build by the provided builder
	bus.SetDefaultTopicBuilder(bus.WorkerTopic.Pool(p))

	err := bus.Subscribe("autoCreatedTopic", bus.EventHandler(func(evt bus.Event) {
		fmt.Println("PrintHandler", evt)
	}))
	require.That(err, is.Nil())

	for i := 0; i < 100; i++ {
		err = bus.Publish("autoCreatedTopic", strings.Repeat("_", i)+strconv.Itoa(i))
		require.That(err, is.Nil())
	}

	// deleting the pool does not delete/stop the underlying pool as this can be shared among different topics
	err = bus.DeleteTopic("autoCreatedTopic")
	require.That(err, is.Nil())

	// we dont need the pool anymore so we can close it ourselves
	err = p.Close()
	require.That(err, is.Nil())
}
