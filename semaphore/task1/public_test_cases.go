package main

import (
	"fmt"
	"sync"
	"time"
)

type TestCase struct {
	name    string
	prepare func() Semaphore
	check   func(Semaphore) bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Проверка базового Acquire/Release",
		prepare: func() Semaphore {
			return New(1)
		},
		check: func(s Semaphore) bool {
			s.Acquire()
			s.Release()

			return true
		},
	},
	{
		name: "Проверка ограничения количества Acquire",
		prepare: func() Semaphore {
			return New(2)
		},
		check: func(s Semaphore) bool {
			s.Acquire()
			s.Acquire()

			// Пытаемся занять больше лимита (должно блокироваться)
			acquired := make(chan bool)
			go func() {
				s.Acquire() // Это должно заблокироваться
				acquired <- true
			}()

			select {
			case <-acquired:
				return false // Не должно сработать
			case <-time.After(100 * time.Millisecond):
				s.Release()
				s.Release()
				return true // Ожидаемое поведение - блокировка
			}
		},
	},
	// Тесткейсы в помощь
	{
		name: "Проверка Release без Acquire",
		prepare: func() Semaphore {
			return New(1)
		},
		check: func(s Semaphore) (success bool) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered from panic:", r)
					success = true
				}
			}()

			s.Release()  // Должно вызвать панику или deadlock
			return false // Не должно сюда дойти
		},
	},
	{
		name: "Проверка работы с несколькими горутинами",
		prepare: func() Semaphore {
			return New(3)
		},
		check: func(s Semaphore) bool {
			var wg sync.WaitGroup
			counter := 0
			mu := sync.Mutex{}

			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					s.Acquire()
					defer s.Release()

					mu.Lock()
					counter++
					mu.Unlock()

					time.Sleep(10 * time.Millisecond)
				}()
			}

			wg.Wait()
			return counter == 10
		},
	},
	{
		name: "Проверка нулевого лимита семафора",
		prepare: func() Semaphore {
			return New(0)
		},
		check: func(s Semaphore) bool {
			acquired := make(chan bool)
			go func() {
				s.Acquire() // Должно заблокироваться навсегда
				acquired <- true
			}()

			select {
			case <-acquired:
				return false
			case <-time.After(100 * time.Millisecond):
				return true // Ожидаемое поведение - блокировка
			}
		},
	},
}
