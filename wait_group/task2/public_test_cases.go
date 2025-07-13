package main

import "sync/atomic"

type TestCase struct {
	name  string
	check func() bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	// Тесткейсы в помощь
	{
		name: "Проверка выполнения всех функций",
		check: func() bool {
			counter := atomic.Int64{}
			callbacks := []func(){
				func() { counter.Add(1) },
				func() { counter.Add(1) },
				func() { counter.Add(1) },
			}

			err := ConcurrentlyRun(callbacks...)
			return err == nil && counter.Load() == 3
		},
	},
	{
		name: "Проверка на пустой список функций",
		check: func() bool {
			err := ConcurrentlyRun()
			return err == nil
		},
	},
	{
		name: "Проверка параллельного выполнения",
		check: func() bool {
			ch := make(chan struct{})
			callbacks := []func(){
				func() { ch <- struct{}{} },
				func() { <-ch },
			}

			err := ConcurrentlyRun(callbacks...)
			if err != nil {
				return false
			}

			return true
		},
	},
	{
		name: "Проверка обработки паники в одной из функций",
		check: func() bool {
			completed := atomic.Bool{}
			callbacks := []func(){
				func() { panic("ошибка в горутине") },
				func() { completed.Store(true) },
			}

			err := ConcurrentlyRun(callbacks...)
			return err != nil && completed.Load()
		},
	},
}
