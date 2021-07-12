# go-bus

[![Code Coverage](https://codecov.io/gh/scarabsoft/go-bus/branch/main/graph/badge.svg)](https://codecov.io/gh/scarabsoft/go-bus)
[![Go Report Card](https://goreportcard.com/badge/github.com/scarabsoft/go-bus)](https://goreportcard.com/report/github.com/scarabsoft/go-bus)
[![Go Reference](https://pkg.go.dev/badge/github.com/scarabsoft/go-bus.svg)](https://pkg.go.dev/github.com/scarabsoft/go-bus)
[![GitHub license](https://img.shields.io/github/license/scarabsoft/go-bus.svg)](https://github.com/scarabsoft/go-bus/blob/main/LICENSE)

**go-bus** is a tiny local pub/sub library I needed for a project to decouple different components. Example usage can be found in examples/


```go
func TestAutoCreate(t *testing.T) {
	// we want every not existing topic run on the same worker pool by that we are able to configure
	// tasks running together and find good queue sizes
	p := pool.New(pool.Options{
		MaxQueueSize: 4,
		MaxWorkers:   2,
	})

	require := hamcrest.NewRequirement(t)

	// every not existing topic should be build by the provided builder
	bus.SetDefaultTopicBuilder(bus.WorkerTopic.Pool(p))

	err := bus.Subscribe("autoCreatedTopic", bus.EventHandler(func(event bus.Event) {
		fmt.Println("PrintHandler", event)
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
```  