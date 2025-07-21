package main

import (
	"sync"
)

var cache = map[int]int{}
var mu sync.RWMutex

func CachedLongCalculation(n int) int {
	mu.RLock()
	found, ok := cache[n]
	mu.RUnlock()

	if ok {
		return found
	}

	value := LongCalculation(n)
	mu.Lock()
	defer mu.Unlock()
	cache[n] = value
	return value
}
