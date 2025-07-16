package main

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var privateTestCases = []TestCase{
	{
		name: "Большое количество чтений",
		check: func() bool {
			updater, cancel := New(returnConfig(mockConfig()), time.Nanosecond)
			defer cancel()

			hasError := atomic.Bool{}
			wg := sync.WaitGroup{}
			wg.Add(1000)

			for range 1000 {
				go func() {
					defer wg.Done()

					if updater.GetConfig() != mockConfig() {
						hasError.Store(true)
					}
				}()
			}

			wg.Wait()

			return !hasError.Load()
		},
	},
	{
		name: "Проверка блокировки при обновлении конфига",
		check: func() bool {
			updater, cancel := New(loadFunc(func() Config {
				randInt := rand.Intn(1000)
				return Config{
					PodsCount:    randInt,
					StartTimeout: time.Duration(randInt),
				}
			}), time.Nanosecond)
			defer cancel()

			hasError := atomic.Bool{}
			wg := sync.WaitGroup{}
			wg.Add(1000000)

			for range 1000000 {
				go func() {
					defer wg.Done()

					config := updater.GetConfig()
					if config.PodsCount != int(config.StartTimeout) { // Если не будет синхронизации, данные могут быть частично проинициализированы
						hasError.Store(true)
					}
				}()
			}

			wg.Wait()

			return !hasError.Load()
		},
	},
}
