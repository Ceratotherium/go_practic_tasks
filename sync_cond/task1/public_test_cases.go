package main

import (
	"time"
)

type TestCase struct {
	name    string
	prepare func() ConnectionPool
	check   func(ConnectionPool) bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name:    "Пул с одним соединением",
		prepare: makePool(1),
		check: func(pool ConnectionPool) bool {
			conn := pool.Acquire()
			defer pool.Release(conn)

			return conn.Get() == mockData
		},
	},
	{
		name:    "Пул с несколькими соединениями",
		prepare: makePool(2),
		check: func(pool ConnectionPool) bool {
			conn1 := pool.Acquire()
			defer pool.Release(conn1)

			conn2 := pool.Acquire()
			defer pool.Release(conn2)

			return conn1 != conn2
		},
	},
	{
		name:    "Переиспользование соединений",
		prepare: makePool(1),
		check: func(pool ConnectionPool) bool {
			conn1 := pool.Acquire()
			pool.Release(conn1)

			conn2 := pool.Acquire()
			defer pool.Release(conn2)

			return conn1 == conn2
		},
	},
	// Тесткейсы в помощь
	{
		name:    "Ожидание доступного соединения",
		prepare: makePool(1),
		check: func(pool ConnectionPool) bool {
			var secondAcquireTime time.Time

			conn := pool.Acquire()

			done := make(chan struct{})

			go func() {
				defer catchPanic("Ожидание доступного соединения")()

				conn := pool.Acquire()
				secondAcquireTime = time.Now()
				pool.Release(conn)

				done <- struct{}{}
			}()

			firstReleaseTime := time.Now()
			pool.Release(conn)

			<-done

			return secondAcquireTime.After(firstReleaseTime)
		},
	},
	{
		name:    "Повторное освобождение соединения",
		prepare: makePool(1),
		check: func(pool ConnectionPool) bool {
			conn := pool.Acquire()
			pool.Release(conn)
			return AssertPanic(func() {
				pool.Release(conn)
			})
		},
	},
	{
		name:    "Освобождение неизвестного соединения",
		prepare: makePool(1),
		check: func(pool ConnectionPool) bool {
			return AssertPanic(func() {
				pool.Release(makeMockConn())
			})
		},
	},
}

func makePool(size int) func() ConnectionPool {
	return func() ConnectionPool {
		return NewConnectionPool(size, makeMockConn)
	}
}

const mockData = "mock"

type mockConn struct {
	_ int
}

func (m *mockConn) Get() string {
	return mockData
}

func makeMockConn() Connection {
	return &mockConn{}
}
