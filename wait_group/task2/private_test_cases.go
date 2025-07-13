package main

import (
	"sync/atomic"
	"time"
)

var privateTestCases = []TestCase{
	{
		name: "Проверка ожидания всех горутин",
		check: func() bool {
			var funcFinish time.Time
			callbacks := []func(){
				func() {
					time.Sleep(100 * time.Millisecond)
					funcFinish = time.Now()
				},
			}

			err := ConcurrentlyRun(callbacks...)
			return err == nil && time.Now().After(funcFinish)
		},
	},
	{
		name: "Много горутин",
		check: func() bool {
			counter := atomic.Int64{}
			callbacks := make([]func(), 0, 1000)

			for range 1000 {
				callbacks = append(callbacks, func() {
					counter.Add(1)
				})
			}

			err := ConcurrentlyRun(callbacks...)
			return err == nil && counter.Load() == 1000
		},
	},
	{
		name: "Проверка конкурентной работы",
		check: func() bool {
			callbacks := make([]func(), 0, 1000)

			for range 1000 {
				callbacks = append(callbacks, func() {
					time.Sleep(100 * time.Millisecond)
				})
			}

			start := time.Now()
			err := ConcurrentlyRun(callbacks...)
			return err == nil && time.Since(start) < 200*time.Millisecond // Должно отработать ~100мс
		},
	},
}
