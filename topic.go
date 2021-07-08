package bus

import "sync/atomic"

type Topic interface {
	Publish(data interface{}) error       // FIXME better would be ...interface{}
	Subscribe(handler EventHandler) error // FIXME better would be ...event.EventHandler

	//FIXME add Close() 				 // FIXME proper closing topic in aysn case use wait groups
}

type abstractTopicImpl struct {
	ID          uint
	name        string
	handlers    []EventHandler
	idGenerator func() uint64
}

//generates topic id which guarantees to be thread safe and monotonous
func topicIdGenerator() func() uint64 {
	var idx uint64 = 0
	return func() uint64 {
		return atomic.AddUint64(&idx, 1)
	}
}
