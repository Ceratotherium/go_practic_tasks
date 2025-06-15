//go:build task_template

package main

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
}

func New(cleanupInterval time.Duration) (Cache, func()) {
	return nil, func() {}
}
