package main

import (
	"context"
	"errors"
	"time"
)

var privateTestCases = []TestCase{
	{
		name:    "Ошибка одного запроса отменяет контекст (проверка, что используется именно errgroup)",
		prepare: nil,
		check: func(_ []Request) bool {
			ctx := context.Background()
			testErr := errors.New("test error")

			isCtxCanceled := false

			requests := []Request{
				func(_ context.Context) error {
					return testErr
				},
				func(ctx context.Context) error {
					time.Sleep(100 * time.Millisecond)
					isCtxCanceled = errors.Is(ctx.Err(), context.Canceled)
					return nil
				},
			}

			err := DoRequests(ctx, 2, requests...)

			return errors.Is(err, testErr) && isCtxCanceled
		},
	},
	{
		name:    "Несколько запросов возвращают ошибку",
		prepare: nil,
		check: func(_ []Request) bool {
			ctx := context.Background()
			testErr1 := errors.New("test error 1")
			testErr2 := errors.New("test error 2")

			requests := []Request{
				func(_ context.Context) error {
					time.Sleep(100 * time.Millisecond)
					return testErr1
				},
				func(_ context.Context) error {
					time.Sleep(200 * time.Millisecond)
					return testErr2
				},
			}

			err := DoRequests(ctx, 2, requests...)

			return errors.Is(err, testErr1)
		},
	},
}
