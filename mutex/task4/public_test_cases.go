package main

import (
	"errors"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type TestCase struct {
	name  string
	check func() bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Повторный вызов с такими же аргументами",
		check: func() bool {
			res, elapsed := measure(CachedLongCalculation, 10)
			if elapsed > time.Millisecond*20 || elapsed < time.Millisecond*10 { // ~10мс
				return false
			}

			if res != 11 {
				return false
			}

			res, elapsed = measure(CachedLongCalculation, 10)
			if elapsed >= time.Millisecond*10 { // Должно вернуться без ожидания
				return false
			}

			return res == 11
		},
	},
	{
		name: "Повторный вызов с другими аргументами",
		check: func() bool {
			res, elapsed := measure(CachedLongCalculation, 10)
			if elapsed > time.Millisecond*20 || elapsed <= time.Millisecond*10 { // ~10мс
				return false
			}

			if res != 11 {
				return false
			}

			res, elapsed = measure(CachedLongCalculation, 11)
			if elapsed > time.Millisecond*20 || elapsed <= time.Millisecond*11 { // ~11мс
				return false
			}

			return res == 12
		},
	},
	// Тесткейсы в помощь
	{
		name: "Конкурентная работа",
		check: func() bool {
			errGroup := errgroup.Group{}

			for value := 10; value <= 100; value += 10 {
				errGroup.Go(func() error {
					res, elapsed := measure(CachedLongCalculation, value)
					if elapsed < time.Millisecond*time.Duration(value) ||
						elapsed > time.Millisecond*time.Duration(value+10) ||
						res != value+1 {
						return errors.New("incorrect result")
					}
					return nil
				})
			}

			return errGroup.Wait() == nil
		},
	},
	{
		name: "Расчет значения не блокирует чтение",
		check: func() bool {
			CachedLongCalculation(10)

			ch := make(chan struct{})
			wg := &sync.WaitGroup{}
			wg.Add(1)

			go func() {
				defer wg.Done()

				ch <- struct{}{}
				CachedLongCalculation(200)
			}()

			<-ch
			time.Sleep(time.Millisecond * 10)
			_, elapsed := measure(CachedLongCalculation, 10)
			if elapsed > time.Millisecond*10 { // Должно вернуться без ожидания
				return false
			}

			wg.Wait()
			return true
		},
	},
	{
		name: "Расчет значения не блокирует запись",
		check: func() bool {
			wg := &sync.WaitGroup{}
			wg.Add(10)

			for value := 10; value <= 100; value += 10 {
				go func() {
					defer wg.Done()

					CachedLongCalculation(value)
				}()
			}

			started := time.Now()
			wg.Wait()
			elapsed := time.Since(started)
			return elapsed >= time.Millisecond*100 && elapsed < time.Millisecond*110 //~100мс
		},
	},
}

func measure(cb func(int) int, arg int) (int, time.Duration) {
	start := time.Now()
	res := cb(arg)
	return res, time.Since(start)
}

func cleanup() {
	cache = map[int]int{}
}

func LongCalculation(n int) int {
	time.Sleep(time.Duration(n) * time.Millisecond)
	return n + 1
}
