package main

import "sync"

type Task interface {
	IsReady() bool
	Run()
}

func RunTasks(tasks []Task) {
	wg := sync.WaitGroup{}

	for _, task := range tasks {
		if !task.IsReady() {
			continue
		}

		wg.Add(1)
		go func() {
			defer func() {
				_ = recover()
			}()
			defer wg.Done()
			task.Run()
		}()
	}

	wg.Wait()
}
