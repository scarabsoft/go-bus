package pool

import (
	"errors"
	"sync"
)

var (
	ErrPoolIsDone = errors.New("pool is done")
)

type TaskPayload struct {
	ID      uint64
	Name    string
	Payload interface{}
}

type Task struct {
	Payload TaskPayload
	Handler func(ID uint64, name string, payload interface{})
}

type Pool interface {
	Submit(Task) error

	Start()

	Stop()
}

type Options struct {
	MaxWorkers   int
	MaxQueueSize int
}

type poolImpl struct {
	Options Options
	Queue   chan Task
	Done    bool

	wg   sync.WaitGroup
	stop chan struct{}
}

func (p *poolImpl) Submit(job Task) error {
	if p.Done {
		return ErrPoolIsDone
	}
	p.Queue <- job
	return nil
}

func (p *poolImpl) Start() {
	for i := 0; i < p.Options.MaxWorkers; i++ {
		p.wg.Add(1)
		go func() {
			for {
				select {
				case <-p.stop:
					p.wg.Done()
				case j := <-p.Queue:
					p := j.Payload
					j.Handler(p.ID, p.Name, p.Payload)
				}
			}
		}()
	}
}

func (p *poolImpl) Stop() {
	for i := 0; i < p.Options.MaxWorkers; i++ {
		p.stop <- struct{}{}
	}
	p.wg.Wait()
	p.Done = true
}

func NewPool(options Options) Pool {
	return &poolImpl{
		Options: options,
		Queue:   make(chan Task, options.MaxQueueSize),
		Done:    false,
		wg:      sync.WaitGroup{},
		stop:    make(chan struct{}, options.MaxQueueSize),
	}
}
