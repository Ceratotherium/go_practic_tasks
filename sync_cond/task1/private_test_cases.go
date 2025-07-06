package main

import (
	"sync"
	"time"
)

var privateTestCases = []TestCase{
	{
		name:    "Использование многими горутинами",
		prepare: makePool(100),
		check: func(pool ConnectionPool) bool {
			wg := sync.WaitGroup{}
			wg.Add(1000)

			start := time.Now()

			for range 1000 {
				go func() {
					defer wg.Done()
					defer catchPanic("Использование многими горутинами")()

					conn := pool.Acquire()
					defer pool.Release(conn)

					time.Sleep(time.Millisecond * 20)
				}()
			}

			wg.Wait()
			elapsed := time.Since(start)

			// Примерное время выполнения ~200мс
			return elapsed >= time.Millisecond*200 && elapsed < time.Millisecond*300
		},
	},
	{
		name:    "Долгое ожидание доступного соединения",
		prepare: makePool(10),
		check: func(pool ConnectionPool) bool {
			var elapsed time.Duration
			done := make(chan struct{})

			started := time.Now()

			firstConn := pool.Acquire()

			for range 9 {
				conn := pool.Acquire()
				defer pool.Release(conn)
			}

			go func() {
				defer catchPanic("Использование многими горутинами")()

				conn := pool.Acquire()
				elapsed = time.Since(started)
				pool.Release(conn)

				done <- struct{}{}
			}()

			time.Sleep(time.Millisecond * 100)
			pool.Release(firstConn)

			<-done

			// Примерное время выполнения ~100мс
			return elapsed >= time.Millisecond*100 && elapsed < time.Millisecond*200
		},
	},
	{
		name:    "Переиспользование соединений (много горутин)",
		prepare: makePool(100),
		check: func(pool ConnectionPool) bool {
			wg := sync.WaitGroup{}
			wg.Add(1000)

			usedConnections := make(map[Connection]struct{})
			mutex := sync.Mutex{}

			for range 1000 {
				go func() {
					defer wg.Done()
					defer catchPanic("Переиспользование соединений (много горутин)")()

					conn := pool.Acquire()
					defer pool.Release(conn)

					time.Sleep(time.Millisecond * 20)

					mutex.Lock()
					defer mutex.Unlock()
					usedConnections[conn] = struct{}{}
				}()
			}

			wg.Wait()

			return len(usedConnections) == 100
		},
	},
}
