//go:build task_template

package main

type Counter interface {
	Add(int64)
	Get() int64
}

func NewCounter() Counter {
	return nil
}
