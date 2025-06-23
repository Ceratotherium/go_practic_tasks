package main

import (
	"context"
	"fmt"
	"time"
)

type TestCase struct {
	name    string
	prepare func() []Job
	check   func([]Job) bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Проверка корректного выполнения задач",
		prepare: func() []Job {
			return []Job{
				&mockJob{value: 1, err: nil},
				&mockJob{value: 2, err: nil},
				&mockJob{value: 3, err: nil},
			}
		},
		check: func(jobs []Job) bool {
			results := New(1).Run(context.Background(), jobs...)

			return (<-results).Value == 1 &&
				(<-results).Value == 2 &&
				(<-results).Value == 3
		},
	},
	{
		name: "Проверка обработки ошибок в задачах",
		prepare: func() []Job {
			return []Job{
				&mockJob{value: 1, err: nil},
				&mockJob{value: 0, err: fmt.Errorf("ошибка выполнения")},
				&mockJob{value: 3, err: nil},
			}
		},
		check: func(jobs []Job) bool {
			results := New(1).Run(context.Background(), jobs...)

			return (<-results).Value == 1 &&
				(<-results).Err != nil &&
				(<-results).Value == 3
		},
	},
	// Тесткейсы в помощь
	{
		name: "Проверка отмены контекста",
		prepare: func() []Job {
			jobs := make([]Job, 100)
			for i := 0; i < 100; i++ {
				jobs[i] = &mockTimeoutJob{time.Millisecond * 500}
			}
			return jobs
		},
		check: func(jobs []Job) bool {
			ctx, cancel := context.WithCancel(context.Background())
			results := New(2).Run(ctx, jobs...)

			cancel()
			count := 0

			for range results {
				count++
			}

			return count != 100
		},
	},
	{
		name: "Проверка пула с нулевым количеством задач",
		prepare: func() []Job {
			return []Job{}
		},
		check: func(jobs []Job) bool {
			results := New(1).Run(context.Background(), jobs...)

			_, ok := <-results

			return !ok
		},
	},
	{
		name: "Проверка пула с одним воркером и множеством задач",
		prepare: func() []Job {
			jobs := make([]Job, 100)
			for i := 0; i < 100; i++ {
				jobs[i] = &mockJob{value: i, err: nil}
			}
			return jobs
		},
		check: func(jobs []Job) bool {
			results := New(1).Run(context.Background(), jobs...)

			for i := range 100 {
				result := <-results
				if result.Value != i {
					return false
				}
			}

			_, ok := <-results

			return !ok
		},
	},
	{
		name: "Проверка пула с несколькими воркерами и множеством задач",
		prepare: func() []Job {
			jobs := make([]Job, 100)
			for i := 0; i < 100; i++ {
				jobs[i] = &mockJob{value: i, err: nil}
			}
			return jobs
		},
		check: func(jobs []Job) bool {
			results := New(10).Run(context.Background(), jobs...)

			values := make(map[int]struct{}, 100)

			// Собираем все значения
			for result := range results {
				if result.Err != nil {
					return false
				}

				values[result.Value] = struct{}{}
			}

			// Проверяем, что прочитали 100 значений
			if len(values) != len(jobs) {
				return false
			}

			// Проверяем, что в канале больше ничего нет
			if _, ok := <-results; ok {
				return false
			}

			// Проверяем, что прочитали нужные значения
			for i := range 100 {
				if _, ok := values[i]; !ok {
					return false
				}
			}
			return true
		},
	},
}

// mockJob - реализация интерфейса Job для тестов
type mockJob struct {
	value    int
	err      error
	timeount *time.Duration
}

func (m *mockJob) Execute(_ context.Context) Result {
	if m.timeount != nil {
		time.Sleep(*m.timeount)
	}

	return Result{
		Value: m.value,
		Err:   m.err,
	}
}

// mockTimeoutJob - реализация интерфейса Job для тестов с тайматом
type mockTimeoutJob struct {
	timeount time.Duration
}

func (m *mockTimeoutJob) Execute(_ context.Context) Result {
	time.Sleep(m.timeount)

	return Result{
		Value: 0,
		Err:   nil,
	}
}
