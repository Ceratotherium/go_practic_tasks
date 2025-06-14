package main

import "sync"

type Counter interface {
	Add(int64)
	Get() int64
}

type counter struct {
	value int64
	mu    *sync.RWMutex
}

func NewCounter() Counter {
	return &counter{
		value: 0,
		mu:    &sync.RWMutex{},
	}
}

func (c *counter) Add(value int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += value
}

func (c *counter) Get() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}
