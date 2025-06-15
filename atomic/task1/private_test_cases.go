package main

import (
	"sync"
	"time"
)

var privateTestCases = []TestCase{
	{
		name: "Смешанные параллельные операции",
		prepare: func() Counter {
			c := NewCounter()
			var wg sync.WaitGroup
			for i := 0; i < 500; i++ {
				wg.Add(2)
				go func() {
					defer wg.Done()
					c.Add(2)
				}()
				go func() {
					defer wg.Done()
					c.Add(-1)
				}()
			}
			wg.Wait()
			return c
		},
		check: func(c Counter) bool {
			return c.Get() == 500 // (2 - 1) * 500
		},
	},
	{
		name: "Чтение при конкурентной записи",
		prepare: func() Counter {
			c := NewCounter()
			go func() {
				for i := 0; i < 100; i++ {
					c.Add(1)
					time.Sleep(time.Microsecond)
				}
			}()
			return c
		},
		check: func(c Counter) bool {
			val := c.Get()
			return val >= 0 && val <= 100
		},
	},
}
