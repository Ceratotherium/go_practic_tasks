package main

import (
	"sync"
	"time"
)

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

type ttlCache struct {
	items map[string]cacheItem
	mu    sync.RWMutex
	stop  chan struct{}
}

func New(cleanupInterval time.Duration) (Cache, func()) {
	c := &ttlCache{
		items: make(map[string]cacheItem),
		stop:  make(chan struct{}),
	}

	if cleanupInterval > 0 {
		go c.cleanupLoop(cleanupInterval)
	}

	return c, func() { c.StopCleanup() }
}

func (c *ttlCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

func (c *ttlCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

func (c *ttlCache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if now.After(item.expiration) {
			delete(c.items, key)
		}
	}
}

func (c *ttlCache) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Cleanup()
		case <-c.stop:
			return
		}
	}
}

func (c *ttlCache) StopCleanup() {
	close(c.stop)
}
