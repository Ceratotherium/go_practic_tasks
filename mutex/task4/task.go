//go:build task_template

package main

import (
	"math/rand"
	"sync"
	"time"
)

// LongCalculation is idempotent, goroutine-safe
func LongCalculation(n int) int {
	secondsToSleep := rand.Float64() * float64(n)
	time.Sleep(time.Duration(secondsToSleep))
	return n + 1
}

var cache = map[int]int{}

func CachedLongCalculation(n int) int {
	var mu sync.Mutex
	mu.Lock()
	found, ok := cache[n]
	mu.Unlock()

	if !ok {
		value := LongCalculation(n)
		mu.Lock()
		cache[n] = value
		mu.Unlock()
		return value
	}

	mu.Unlock()
	return found
}
