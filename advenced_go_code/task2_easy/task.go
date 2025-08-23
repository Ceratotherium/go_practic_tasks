//go:build task_template

package main

type Request interface{}

type Response interface{}

type Backend interface {
	Invoke(ctx context.Context, req Request) (Response, error)
}

type Balancer struct {
	//TODO
}

var _ Backend = &Balancer{}

// addrs содержат адреса всех балансируемых экземпляров
func NewBalancer(addrs []string, makeBackend func(string) Backend) *Balancer {
	//TODO
}
