package main

import (
	"sync"
	"time"
)

var privateTestCases = []TestCase{
	{
		name: "Проверка долгой блокировки одной горутины",
		prepare: func() Semaphore {
			return New(2)
		},
		check: func(s Semaphore) bool {
			lockChan := make(chan struct{})

			go func() {
				defer s.Release()
				s.Acquire()
				<-lockChan
			}()

			var wg sync.WaitGroup
			wg.Add(10)

			for range 10 {
				go func() {
					defer wg.Done()
					defer s.Release()
					s.Acquire()
					time.Sleep(time.Millisecond * 10)
				}()
			}

			wg.Wait()

			lockChan <- struct{}{}

			return true
		},
	},
}
