package main

import (
	"context"
	"sync"
	"time"
)

type Config struct {
	PodsCount    int
	StartTimeout time.Duration
}

type Loader interface {
	Load() Config
}

type ConfigUpdater struct {
	config Config
	mutex  sync.RWMutex
	loader Loader
}

func New(loader Loader, updateInterval time.Duration) (*ConfigUpdater, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	updater := &ConfigUpdater{
		loader: loader,
	}

	updater.update()

	go func() {
		t := time.NewTicker(updateInterval)

		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				updater.update()
			}
		}
	}()

	return updater, cancel
}

func (c *ConfigUpdater) update() {
	config := c.loader.Load()

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.config = config
}

func (c *ConfigUpdater) GetConfig() Config {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.config
}
