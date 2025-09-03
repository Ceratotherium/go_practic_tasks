package main

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"
)

type Request interface{}

type Response interface{}

type Backend interface {
	Invoke(ctx context.Context, req Request) (Response, error)
}

type BalancerBackend struct {
	Backend

	lastError   time.Time
	errorsCount int

	load int
}

var _ Backend = &BalancerBackend{}

type Balancer struct {
	mx       sync.Mutex
	backends []*BalancerBackend

	suspendAfterErrors int
	suspendFor         time.Duration
}

var _ Backend = &Balancer{}

func NewBalancer(addrs []string, suspendAfterErrors int, suspendFor time.Duration, makeBackend func(string) Backend) *Balancer {
	var bs []*BalancerBackend

	for _, a := range addrs {
		bs = append(bs, &BalancerBackend{Backend: makeBackend(a)})
	}

	return &Balancer{
		backends:           bs,
		suspendFor:         suspendFor,
		suspendAfterErrors: suspendAfterErrors,
	}
}

func (b *Balancer) Invoke(ctx context.Context, req Request) (Response, error) {
	var backend *BalancerBackend

	b.mx.Lock()
	sort.Sort(balancerBackendsByLoadDesc(b.backends))

	for _, bb := range b.backends {
		if bb.errorsCount >= b.suspendAfterErrors {
			if time.Since(bb.lastError) < b.suspendFor {
				continue
			}

			bb.errorsCount = 0
		}

		backend = bb
		break
	}

	if backend == nil {
		b.mx.Unlock()

		return nil, errors.New("all backends suspended")
	}

	backend.load++
	b.mx.Unlock()

	res, err := backend.Invoke(ctx, req)

	b.mx.Lock()
	backend.load--

	if err != nil {
		backend.errorsCount++
		backend.lastError = time.Now()
	} else {
		backend.errorsCount = 0
	}
	b.mx.Unlock()

	return res, err
}

type balancerBackendsByLoadDesc []*BalancerBackend

func (bs balancerBackendsByLoadDesc) Len() int           { return len(bs) }
func (bs balancerBackendsByLoadDesc) Less(i, j int) bool { return bs[i].load < bs[j].load }
func (bs balancerBackendsByLoadDesc) Swap(i, j int)      { bs[i], bs[j] = bs[j], bs[i] }
