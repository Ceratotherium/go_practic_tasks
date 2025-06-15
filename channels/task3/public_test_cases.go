package main

import "time"

type TestCase struct {
	name    string
	prepare func() []Backend
	check   func([]string) bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Один бэкенд с успешным ответом",
		prepare: func() []Backend {
			return []Backend{
				&mockBackend{response: "ok", delay: 0},
			}
		},
		check: func(results []string) bool {
			return len(results) == 1 && results[0] == "ok"
		},
	},
	{
		name: "Бэкенд с задержкой больше таймаута",
		prepare: func() []Backend {
			return []Backend{
				&mockBackend{response: "slow", delay: 3 * time.Second},
			}
		},
		check: func(results []string) bool {
			return len(results) == 0
		},
	},
	{
		name: "Несколько бэкендов с успешными ответами",
		prepare: func() []Backend {
			return []Backend{
				&mockBackend{response: "1", delay: 0},
				&mockBackend{response: "2", delay: 0},
				&mockBackend{response: "3", delay: 0},
			}
		},
		check: func(results []string) bool {
			return len(results) == 3 &&
				СontainsAll(results, "1", "2", "3")
		},
	},
	// Тесткейсы в помощь
	{
		name: "Пустой список бэкендов",
		prepare: func() []Backend {
			return []Backend{}
		},
		check: func(results []string) bool {
			return len(results) == 0
		},
	},
	{
		name: "Бэкенд с задержкой меньше таймаута",
		prepare: func() []Backend {
			return []Backend{
				&mockBackend{response: "fast", delay: 100 * time.Millisecond},
			}
		},
		check: func(results []string) bool {
			return len(results) == 1 && results[0] == "fast"
		},
	},
	{
		name: "Смесь быстрых и медленных бэкендов",
		prepare: func() []Backend {
			return []Backend{
				&mockBackend{response: "fast1", delay: 100 * time.Millisecond},
				&mockBackend{response: "slow", delay: 3 * time.Second},
				&mockBackend{response: "fast2", delay: 150 * time.Millisecond},
			}
		},
		check: func(results []string) bool {
			return len(results) == 2 &&
				СontainsAll(results, "fast1", "fast2")
		},
	},
}

// Mock реализация Backend для тестирования
type mockBackend struct {
	response string
	delay    time.Duration
}

func (m *mockBackend) DoRequest() string {
	time.Sleep(m.delay)
	return m.response
}
