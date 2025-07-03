package main

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Request func(ctx context.Context) error

func DoRequests(ctx context.Context, n int, requests ...Request) error {
	group, ctx := errgroup.WithContext(ctx)

	if n > 0 {
		group.SetLimit(n)
	}

	for _, request := range requests {
		select {
		case <-ctx.Done():
			break // Прерываем цикл и дожидаемся запущенные горутины
		default:
		}

		group.Go(func() error {
			return request(ctx)
		})
	}

	return group.Wait()
}
