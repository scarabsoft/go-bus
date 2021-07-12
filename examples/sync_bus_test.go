package examples

import (
	"fmt"
	"github.com/scarabsoft/go-bus"
	"github.com/scarabsoft/go-hamcrest"
	"github.com/scarabsoft/go-hamcrest/is"
	"testing"
	"time"
)

func TestSyncBus(t *testing.T) {
	require := hamcrest.NewRequirement(t)

	//create a new 'coolTopic' to publish cool content to
	coolTopic, err := bus.CreateTopic("coolTopic", bus.SyncTopic)
	require.That(err, is.Nil())

	//create a new 'analytics' to publish cool content to
	_, err = bus.CreateTopic("metricsTopic", bus.SyncTopic)
	require.That(err, is.Nil())

	//register a simple event handler which just prints data received
	printHandler := bus.EventHandler(func(event bus.Event) {
		fmt.Println("PrintHandler - ", event)
	})

	// simulates doing some metrics calculation and publish to metricsTopic
	metricsHandler := func(ID uint64, name string, payload interface{}) {
		err := bus.Publish("metricsTopic", ID*ID)
		require.That(err, is.Nil())
		time.Sleep(10 * time.Millisecond)
	}

	// print cool topic to the terminal and run the metrics handler
	err = coolTopic.Subscribe(printHandler, metricsHandler)
	require.That(err, is.Nil())

	// we just print the metrics stream to the terminal as well
	// equivalent to err = metricsTopic.Subscribe(printHandler)
	err = bus.Subscribe("metricsTopic", printHandler)
	require.That(err, is.Nil())

	for i := 0; i < 10; i++ {
		err = coolTopic.Publish(i)
		require.That(err, is.Nil())
	}

	// lets clean up
	err = bus.DeleteTopic("coolTopic")
	require.That(err, is.Nil())

	err = bus.DeleteTopic("metricsTopic")
	require.That(err, is.Nil())

}
