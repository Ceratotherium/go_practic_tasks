package main

import (
	"fmt"
	"time"
)

var privateTestCases = []TestCase{
	{
		name: "Большое количество быстрых бэкендов (1000)",
		prepare: func() []Backend {
			backends := make([]Backend, 1000)
			for i := 0; i < 1000; i++ {
				backends[i] = &mockBackend{
					response: fmt.Sprintf("res%d", i),
					delay:    0,
				}
			}
			return backends
		},
		check: func(results []string) bool {
			if len(results) != 1000 {
				return false
			}
			// Проверяем что все уникальные ответы присутствуют
			resSet := make(map[string]struct{})
			for _, r := range results {
				resSet[r] = struct{}{}
			}
			for i := 0; i < 1000; i++ {
				if _, exists := resSet[fmt.Sprintf("res%d", i)]; !exists {
					return false
				}
			}
			return true
		},
	},
	{
		name: "Большое количество бэкендов с разными задержками",
		prepare: func() []Backend {
			backends := make([]Backend, 500)
			for i := 0; i < 500; i++ {
				delay := time.Duration(i%10) * 100 * time.Millisecond
				backends[i] = &mockBackend{
					response: fmt.Sprintf("res%d", i),
					delay:    delay,
				}
			}
			return backends
		},
		check: func(results []string) bool {
			// Должны вернуться все, так как максимальная задержка 900ms < 2s timeout
			return len(results) == 500
		},
	},
	{
		name: "Очень большое количество бэкендов (10 000) с частичным таймаутом",
		prepare: func() []Backend {
			backends := make([]Backend, 10000)
			for i := 0; i < 10000; i++ {
				var delay time.Duration
				if i < 5000 {
					delay = 100 * time.Millisecond // быстрые
				} else {
					delay = 3 * time.Second // медленные (не успеют)
				}
				backends[i] = &mockBackend{
					response: fmt.Sprintf("res%d", i),
					delay:    delay,
				}
			}
			return backends
		},
		check: func(results []string) bool {
			// Должны вернуться только быстрые (первые 5000)
			return len(results) == 5000
		},
	},
}
