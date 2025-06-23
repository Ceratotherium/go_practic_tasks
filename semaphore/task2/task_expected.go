package main

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type Result struct {
	Value int
	Err   error
}

type Job interface {
	Execute(ctx context.Context) Result
}

type WorkerPool interface {
	Run(ctx context.Context, jobsBulk ...Job) <-chan Result
}

type workerPool struct {
	sem          *semaphore.Weighted
	workersCount int64
	jobs         chan Job
}

func New(workersCount int64) WorkerPool {
	return &workerPool{
		sem:          semaphore.NewWeighted(workersCount),
		workersCount: workersCount,
		jobs:         make(chan Job, workersCount),
	}
}

func (wp *workerPool) Run(ctx context.Context, jobsBulk ...Job) <-chan Result {
	results := make(chan Result, wp.workersCount)

	go func() {
		defer close(results)

		for _, job := range jobsBulk {
			if err := wp.sem.Acquire(ctx, 1); err != nil {
				break
			}

			go func() {
				defer wp.sem.Release(1)
				results <- job.Execute(ctx)
			}()
		}

		_ = wp.sem.Acquire(context.WithoutCancel(ctx), wp.workersCount)
		wp.sem.Release(wp.workersCount)
	}()

	return results
}
