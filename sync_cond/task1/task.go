//go:build task_template

package main

type Connection interface {
	Get() string
}

type ConnectionPool interface {
	Acquire() Connection
	Release(Connection)
}

func NewConnectionPool(poolSize int, openConnection func() Connection) ConnectionPool {
	return nil
}

type mockConnection struct{}

func (m *mockConnection) Get() string {
	return "mock"
}

func makeMockConnection() Connection {
	return &mockConnection{}
}
