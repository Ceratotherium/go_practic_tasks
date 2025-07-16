package main

import (
	"sync/atomic"
	"time"
)

type TestCase struct {
	name  string
	check func() bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Получение конфига",
		check: func() bool {
			updated, cancel := New(returnConfig(mockConfig()), time.Minute)
			defer cancel()

			config := updated.GetConfig()

			return config == mockConfig()
		},
	},
	{
		name: "Обновление конфига",
		check: func() bool {
			config1 := Config{
				PodsCount:    1,
				StartTimeout: time.Second,
			}
			config2 := Config{
				PodsCount:    2,
				StartTimeout: time.Second,
			}
			isFirstCall := atomic.Bool{}
			isFirstCall.Store(true)

			updated, cancel := New(loadFunc(func() Config {
				if isFirstCall.Load() {
					isFirstCall.Store(false)
					return config1
				}
				return config2
			}), time.Millisecond*20)
			defer cancel()

			if updated.GetConfig() != config1 {
				return false
			}

			time.Sleep(time.Millisecond * 50)

			return updated.GetConfig() == config2
		},
	},
	// Тесткейсы в помощь
	{
		name: "Проверка первичной инициализации",
		check: func() bool {
			started := time.Now()
			_, cancel := New(loadFunc(func() Config {
				time.Sleep(time.Millisecond * 50)
				return mockConfig()
			}), time.Minute)
			defer cancel()

			elapsed := time.Since(started)

			return elapsed >= time.Millisecond*50 && elapsed < time.Millisecond*100
		},
	},
	// Тесткейсы в помощь
	{
		name: "Загрузка конфига не блокирует чтение",
		check: func() bool {
			isFirstCall := atomic.Bool{}
			isFirstCall.Store(true)
			locker := make(chan struct{})

			updater, cancel := New(
				loadFunc(func() Config {
					if isFirstCall.Load() {
						isFirstCall.Store(false)
						return mockConfig()
					}

					locker <- struct{}{}
					<-locker
					return mockConfig()
				}),
				time.Millisecond*10,
			)
			defer cancel()

			<-locker
			_ = updater.GetConfig() // Если загрузка выполняется под мьютексом, то будет дедлок
			close(locker)
			return true
		},
	},
	{
		name: "Прекращение обновлений после вызова колбека",
		check: func() bool {
			closed := atomic.Bool{}
			calledAfterCancel := atomic.Bool{}

			_, cancel := New(
				loadFunc(func() Config {
					if closed.Load() {
						calledAfterCancel.Store(true)
					}

					return mockConfig()
				}),
				time.Millisecond*5,
			)

			time.Sleep(time.Millisecond * 20) // Ждем, пока точно запустится поток загрузки конфига

			cancel()
			closed.Store(true)

			time.Sleep(time.Millisecond * 20) // Даем время на предположительный апдейт

			return !calledAfterCancel.Load()
		},
	},
}

func mockConfig() Config {
	return Config{
		PodsCount:    2,
		StartTimeout: time.Second,
	}
}

type loadFunc func() Config

func (f loadFunc) Load() Config {
	return f()
}

func returnConfig(c Config) Loader {
	return loadFunc(func() Config {
		return c
	})
}
