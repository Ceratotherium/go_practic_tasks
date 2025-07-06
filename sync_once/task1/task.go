//go:build task_template

package main

type Database interface {
	Connect()
}

func GetDatabase() Database {
	return MakeDatabase()
}

type MockDatabase struct{}

func (m *MockDatabase) Connect() {}

func MakeDatabase() Database {
	return MockDatabase{}
}
