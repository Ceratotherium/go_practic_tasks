package main

import (
	"context"
	"sync"
)

type Result struct {
	Value int
	Err   error
}

type Job interface {
	Execute(ctx context.Context) Result
}

type Pool struct {
	workersCount int
	jobs         chan Job
}

type WorkerPool interface {
	Run(ctx context.Context, jobs ...Job) <-chan Result
}

func New(workersCount int) WorkerPool {
	return &Pool{
		workersCount: workersCount,
		jobs:         make(chan Job, workersCount),
	}
}

func (wp Pool) Run(ctx context.Context, jobs ...Job) <-chan Result {
	results := make(chan Result, wp.workersCount)

	go wp.add(jobs)

	var wg sync.WaitGroup
	wg.Add(wp.workersCount)

	for i := 0; i < wp.workersCount; i++ {
		go worker(ctx, &wg, wp.jobs, results)
	}

	go func() {
		defer close(results)
		wg.Wait()
	}()

	return results
}

func (wp Pool) add(jobsBulk []Job) {
	for i, _ := range jobsBulk {
		wp.jobs <- jobsBulk[i]
	}
	close(wp.jobs)
}

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		job, ok := <-jobs
		if !ok {
			return
		}

		// Выполняем джобу и отправляем результат
		results <- job.Execute(ctx)
	}
}
