package main

import (
	"context"
	"errors"
	"time"
)

type TestCase struct {
	name    string
	prepare []Request
	check   func([]Request) bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Успешное выполнение всех запросов",
		prepare: []Request{
			func(ctx context.Context) error { return nil },
			func(ctx context.Context) error { return nil },
			func(ctx context.Context) error { return nil },
		},
		check: func(requests []Request) bool {
			// Все запросы должны выполниться успешно
			ctx := context.Background()
			err := DoRequests(ctx, 2, requests...)
			return err == nil
		},
	},
	{
		name: "Ошибка в одном из запросов",
		prepare: []Request{
			func(ctx context.Context) error { return nil },
			func(ctx context.Context) error { return errors.New("ошибка выполнения") },
			func(ctx context.Context) error { return nil },
		},
		check: func(requests []Request) bool {
			// Должна вернуться первая ошибка
			ctx := context.Background()
			err := DoRequests(ctx, 2, requests...)
			return err != nil && err.Error() == "ошибка выполнения"
		},
	},
	{
		name: "Нет ограничения на количество горутин",
		prepare: func() []Request {
			requests := make([]Request, 0, 1000)
			for range 1000 {
				requests = append(requests, func(_ context.Context) error {
					time.Sleep(time.Millisecond * 100)
					return nil
				})
			}
			return requests
		}(),
		check: func(requests []Request) bool {
			start := time.Now()
			ctx := context.Background()
			err := DoRequests(ctx, 0, requests...)
			duration := time.Since(start)

			// I/O bound задача, должна отработать за ~100мс
			return err == nil && duration >= 100*time.Millisecond && duration < 200*time.Millisecond
		},
	},
	// Тесткейсы в помощь
	{
		name: "Отмена контекста во время выполнения",
		prepare: []Request{
			func(ctx context.Context) error {
				time.Sleep(100 * time.Millisecond)
				return ctx.Err()
			},
			func(ctx context.Context) error {
				time.Sleep(200 * time.Millisecond)
				return ctx.Err()
			},
			func(ctx context.Context) error {
				time.Sleep(300 * time.Millisecond)
				return ctx.Err()
			},
		},
		check: func(requests []Request) bool {
			// Проверяем отмену контекста
			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				time.Sleep(50 * time.Millisecond)
				cancel()
			}()

			err := DoRequests(ctx, 2, requests...)
			return errors.Is(err, context.Canceled)
		},
	},
	{
		name: "Ограничение количества одновременных запросов",
		prepare: []Request{
			func(ctx context.Context) error {
				time.Sleep(100 * time.Millisecond)
				return nil
			},
			func(ctx context.Context) error {
				time.Sleep(100 * time.Millisecond)
				return nil
			},
			func(ctx context.Context) error {
				time.Sleep(100 * time.Millisecond)
				return nil
			},
		},
		check: func(requests []Request) bool {
			// Проверяем, что лимит concurrent соблюдается
			start := time.Now()
			ctx := context.Background()
			err := DoRequests(ctx, 2, requests...)
			duration := time.Since(start)

			// 3 запроса с лимитом 2 должны выполняться ~200ms (100+100)
			return err == nil && duration >= 200*time.Millisecond && duration < 300*time.Millisecond
		},
	},
	{
		name:    "Пустой список запросов",
		prepare: []Request{},
		check: func(requests []Request) bool {
			// Должен вернуться nil без ошибок
			ctx := context.Background()
			err := DoRequests(ctx, 2, requests...)
			return err == nil
		},
	},
	{
		name: "Долгий запрос с таймаутом",
		prepare: []Request{
			func(ctx context.Context) error {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(time.Second):
					return nil
				}
			},
			func(ctx context.Context) error {
				return ctx.Err()
			},
		},
		check: func(requests []Request) bool {
			// Проверяем обработку таймаута
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			start := time.Now()
			err := DoRequests(ctx, 2, requests...)
			duration := time.Since(start)

			return errors.Is(err, context.DeadlineExceeded) &&
				duration < 200*time.Millisecond &&
				duration >= 100*time.Millisecond
		},
	},
}
