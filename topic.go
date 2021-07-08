package go_bus

type Topic interface {
	Publish(data interface{}) error       // FIXME better would be ...interface{}
	Subscribe(handler EventHandler) error // FIXME better would be ...event.EventHandler

	//FIXME add Close() 				 // FIXME proper closing topic in aysn case use wait groups
}

type topicImpl struct {
	name     string
	handlers []EventHandler
}

type Option func(topic Topic) error
