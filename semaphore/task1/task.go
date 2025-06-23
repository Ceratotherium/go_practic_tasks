//go:build task_template

package main

type Semaphore interface {
	Acquire()
	Release()
}

func New(limit int) Semaphore {
	return nil
}
