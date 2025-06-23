package main

type Semaphore interface {
	Acquire()
	Release()
}

type semaphore struct {
	ch chan struct{}
}

func New(limit int) Semaphore {
	return &semaphore{make(chan struct{}, limit)}
}

func (s *semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *semaphore) Release() {
	select {
	case <-s.ch:
		return
	default:
		panic("Release called before Acquire")
	}
}
