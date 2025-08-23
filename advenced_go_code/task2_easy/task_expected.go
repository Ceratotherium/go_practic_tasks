package main

import (
	"context"
	"errors"
	"sync/atomic"
)

type Request interface{}

type Response interface{}

type Backend interface {
	Invoke(ctx context.Context, req Request) (Response, error)
}

type Balancer struct {
	backends []Backend
	current  atomic.Int64
}

var _ Backend = &Balancer{}

// addrs содержат адреса всех балансируемых экземпляров
func NewBalancer(addrs []string, makeBackend func(string) Backend) *Balancer {
	backends := make([]Backend, 0, len(addrs))

	for _, addr := range addrs {
		backends = append(backends, makeBackend(addr))
	}

	return &Balancer{backends: backends}
}

func (b *Balancer) Invoke(ctx context.Context, req Request) (Response, error) {
	if len(b.backends) == 0 {
		return nil, errors.New("no backends")
	}

	index := b.current.Add(1) % int64(len(b.backends))
	return b.backends[index].Invoke(ctx, req)
}
