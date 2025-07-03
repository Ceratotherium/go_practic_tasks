//go:build task_template

package main

type Request func(ctx context.Context) error

func DoRequests(ctx context.Context, concurrent int, requests ...Request) error {
	return nil
}
