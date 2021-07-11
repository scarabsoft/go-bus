package pool

import (
	"sync"
)

type poolImpl struct {
	options Options
	queue   chan Task
	closed  bool

	wg   sync.WaitGroup
	stop chan struct{}
}

func (p *poolImpl) Submit(job Task) error {
	if p.closed {
		return ErrPoolIsClosed
	}
	p.queue <- job
	return nil
}

func (p *poolImpl) start() {
	for i := 0; i < p.options.MaxWorkers; i++ {
		p.wg.Add(1)
		go func() {
			for {
				select {
				case <-p.stop:
					p.wg.Done()
				case j := <-p.queue:
					p := j.Payload
					j.Handler(p.ID, p.Name, p.Payload)
				}
			}
		}()
	}
}

func (p *poolImpl) Close() error {
	for i := 0; i < p.options.MaxWorkers; i++ {
		p.stop <- struct{}{}
	}
	p.wg.Wait()
	p.closed = true
	return nil
}

func New(options Options) *poolImpl {
	result := &poolImpl{
		options: options,
		queue:   make(chan Task, options.MaxQueueSize),
		closed:  false,
		wg:      sync.WaitGroup{},
		stop:    make(chan struct{}, options.MaxQueueSize),
	}
	result.start()
	return result
}
