package main

import (
	"context"
	"time"
)

type Backend interface {
	DoRequest() string
}

func DoRequests(backends []Backend) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch := make(chan string, len(backends))

	for _, backend := range backends {
		go func() {
			ch <- backend.DoRequest()
		}()
	}

	results := make([]string, 0, len(backends))

	for range len(backends) {
		select {
		case result := <-ch:
			results = append(results, result)
		case <-ctx.Done():
			return results
		}
	}

	return results
}
