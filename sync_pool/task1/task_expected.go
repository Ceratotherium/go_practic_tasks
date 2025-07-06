package main

import (
	"sync"
)

type Connection interface {
	Get() string
}

type ConnectionPool interface {
	Acquire() Connection
	Release(Connection)
}

type pool struct {
	syncPool *sync.Pool
	sem      chan struct{}
}

func NewConnectionPool(poolSize int, openConnection func() Connection) ConnectionPool {
	syncPool := &sync.Pool{
		New: func() any { return openConnection() },
	}
	for range poolSize {
		syncPool.Put(openConnection())
	}

	return &pool{
		syncPool: syncPool,
		sem:      make(chan struct{}, poolSize),
	}
}

func (p *pool) Acquire() Connection {
	p.sem <- struct{}{}

	return p.syncPool.Get().(Connection)
}

func (p *pool) Release(c Connection) {
	defer func() {
		<-p.sem
	}()

	p.syncPool.Put(c)
}
