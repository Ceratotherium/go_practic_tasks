package main

import (
	"sync"
	"sync/atomic"
)

type TestCase struct {
	name  string
	check func() bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Функция инициализация выполняет подключение к базе",
		check: func() bool {
			db := castToMock(GetDatabase())

			return db.ConnectCount() == 1
		},
	},
	{
		name: "Функция инициализации возвращает один и тот же объект",
		check: func() bool {
			db1 := GetDatabase()
			db2 := GetDatabase()

			return db1 == db2
		},
	},
	// Тесткейсы в помощь
	{
		name: "Коннект выполняется 1 раз",
		check: func() bool {
			_ = GetDatabase()
			db := castToMock(GetDatabase())

			return db.ConnectCount() == 1
		},
	},
	{
		name: "Конкурентная инициализация",
		check: func() bool {
			var db1, db2 *mockDatabase

			wg := sync.WaitGroup{}
			wg.Add(2)

			go func() {
				defer wg.Done()
				db1 = castToMock(GetDatabase())
			}()
			go func() {
				defer wg.Done()
				db2 = castToMock(GetDatabase())
			}()

			wg.Wait()

			return db1 == db2 && db1.ConnectCount() == 1
		},
	},
}

func MakeDatabase() Database {
	return &mockDatabase{}
}

type mockDatabase struct {
	connectCount atomic.Int64
}

func (m *mockDatabase) Connect() {
	m.connectCount.Add(1)
}

func (m *mockDatabase) ConnectCount() int64 {
	return m.connectCount.Load()
}

func castToMock(db Database) *mockDatabase {
	return db.(*mockDatabase)
}
