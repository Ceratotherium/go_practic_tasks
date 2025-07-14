package main

import "sync/atomic"

var privateTestCases = []TestCase{
	{
		name: "Много тасок",
		check: func() bool {
			counter := atomic.Int64{}
			tasks := make([]Task, 0, 1000)

			for i := 0; i < 1000; i++ {
				tasks = append(tasks, mockTask{
					ready: i%2 == 0,
					run: func() {
						counter.Add(1)
					},
				})
			}

			RunTasks(tasks)

			return counter.Load() == 500
		},
	},
}
