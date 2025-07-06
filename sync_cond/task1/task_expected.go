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
	connections  map[Connection]bool
	cond         *sync.Cond
	mutex        *sync.Mutex
	acquireCount int
}

func NewConnectionPool(poolSize int, openConnection func() Connection) ConnectionPool {
	connections := make(map[Connection]bool, poolSize)
	for range poolSize {
		conn := openConnection()
		connections[conn] = false
	}

	mutex := &sync.Mutex{}

	return &pool{
		connections:  connections,
		cond:         sync.NewCond(mutex),
		mutex:        mutex,
		acquireCount: 0,
	}
}

func (p *pool) Acquire() Connection {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for p.acquireCount >= len(p.connections) {
		p.cond.Wait()
	}

	return p.getUnusedConnection()
}

func (p *pool) Release(c Connection) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, ok := p.connections[c]; !ok || p.acquireCount <= 0 {
		panic("Connection not acquired")
	}

	p.releaseConnection(c)
	p.cond.Signal()
}

func (p *pool) getUnusedConnection() Connection {
	p.acquireCount++

	for conn, used := range p.connections {
		if !used {
			p.connections[conn] = true
			return conn
		}
	}
	return nil
}

func (p *pool) releaseConnection(c Connection) {
	p.connections[c] = false
	p.acquireCount--
}
