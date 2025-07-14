//go:build task_template

package main

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

		go func() {
			wg.Add(1)
			task.Run()
			wg.Done()
		}()
	}

	wg.Wait()
}
