package main

import (
	"sync"
)

type TestCase struct {
	name    string
	prepare func() Counter
	check   func(counter Counter) bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Одиночное добавление положительного числа",
		prepare: func() Counter {
			c := NewCounter()
			c.Add(42)
			return c
		},
		check: func(c Counter) bool {
			return c.Get() == 42
		},
	},
	{
		name: "Одиночное добавление отрицательного числа",
		prepare: func() Counter {
			c := NewCounter()
			c.Add(-15)
			return c
		},
		check: func(c Counter) bool {
			return c.Get() == -15
		},
	},
	// Тесткейсы в помощь
	{
		name: "Инициализация с нулевым значением",
		prepare: func() Counter {
			return NewCounter()
		},
		check: func(c Counter) bool {
			return c.Get() == 0
		},
	},
	{
		name: "Множественные последовательные добавления",
		prepare: func() Counter {
			c := NewCounter()
			c.Add(10)
			c.Add(20)
			c.Add(-5)
			return c
		},
		check: func(c Counter) bool {
			return c.Get() == 25
		},
	},
	{
		name: "Параллельные добавления из горутин",
		prepare: func() Counter {
			c := NewCounter()
			var wg sync.WaitGroup
			wg.Add(1000)
			for i := 0; i < 1000; i++ {
				go func() {
					defer wg.Done()
					c.Add(1)
				}()
			}
			wg.Wait()
			return c
		},
		check: func(c Counter) bool {
			return c.Get() == 1000
		},
	},
}
