package main

import (
	"context"
	"time"
)

var privateTestCases = []TestCase{
	{
		name: "Проверка, что внутри используется несколько потоков",
		prepare: func() []Job {
			jobs := make([]Job, 100)
			for i := 0; i < 100; i++ {
				jobs[i] = &mockJob{value: i, err: nil}
			}
			return jobs
		},
		check: func(jobs []Job) bool {
			lockCh := make(chan struct{})

			allJobs := []Job{&mockLockJob{lock: lockCh}}
			allJobs = append(allJobs, jobs...)

			results := New(10).Run(context.Background(), allJobs...)

			// Вычитываем 100 значений, обработанных не залоченными воркерами
			for range 100 {
				<-results
			}

			// Проверяем, что оставшая джоба еще "выполняется"
			select {
			case <-results:
				return false
			default:
			}

			// Разблокировываем джобу
			lockCh <- struct{}{}

			// Получаем значение
			<-results

			// Проверяем, что канал закрыт
			_, ok := <-results
			return !ok
		},
	},
	{
		name: "Проверка, что после отмены контекста дожидаемся завершения воркеров",
		prepare: func() []Job {
			jobs := make([]Job, 100)
			for i := 0; i < 100; i++ {
				jobs[i] = &mockJob{value: i, err: nil}
			}
			return jobs
		},
		check: func(jobs []Job) bool {
			lock1 := make(chan struct{})
			lock2 := make(chan struct{})
			jobs[10] = &mockLockJob{lock: lock1}
			jobs[11] = &mockLockJob{lock: lock2}

			ctx, cancel := context.WithCancel(context.Background())
			results := New(2).Run(ctx, jobs...)

			go func() {
				time.Sleep(time.Second)
				cancel()

				lock1 <- struct{}{}
				lock2 <- struct{}{}
			}()

			count := 0

			for range results {
				count++
			}

			return count == 12
		},
	},
}

// mockLockJob - реализация интерфейса Job для тестов c блокировкой
type mockLockJob struct {
	lock chan struct{}
}

func (m *mockLockJob) Execute(_ context.Context) Result {
	<-m.lock

	return Result{
		Value: -1,
		Err:   nil,
	}
}
