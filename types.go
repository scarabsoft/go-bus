package bus


type Event interface {
	Id() uint64
	Topic() string
	Payload() interface{}
}

type Dispatcher interface {

}

