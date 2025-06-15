package main

import (
	"sync/atomic"
)

type Counter interface {
	Add(int64)
	Get() int64
}

type counter struct {
	value int64
}

func NewCounter() Counter {
	return &counter{
		value: 0,
	}
}

func (c *counter) Add(value int64) {
	atomic.AddInt64(&c.value, value)
}

func (c *counter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}
