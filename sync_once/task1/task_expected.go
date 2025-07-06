package main

import "sync"

type Database interface {
	Connect()
}

var (
	once     sync.Once
	database Database
)

func GetDatabase() Database {
	once.Do(func() {
		database = MakeDatabase()
		database.Connect()
	})
	return database
}
