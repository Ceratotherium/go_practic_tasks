//go:build task_template

package main

type Result struct {
	Value int
	Err   error
}

type Job interface {
	Execute(ctx context.Context) Result
}

type WorkerPool interface {
	Run(ctx context.Context, jobs ...Job) <-chan Result
}

func New(workersCount int) WorkerPool {
	return nil
}
