package main

import "sync/atomic"

type TestCase struct {
	name  string
	check func() bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Все таски готовы",
		check: func() bool {
			var task1Run, task2Run atomic.Bool
			tasks := []Task{
				mockTask{
					ready: true,
					run: func() {
						task1Run.Store(true)
					},
				},
				mockTask{
					ready: true,
					run: func() {
						task2Run.Store(true)
					},
				},
			}

			RunTasks(tasks)

			return task1Run.Load() && task2Run.Load()
		},
	},
	{
		name: "Все таски не готовы",
		check: func() bool {
			var task1Run, task2Run bool
			tasks := []Task{
				mockTask{
					ready: false,
					run: func() {
						task1Run = true
					},
				},
				mockTask{
					ready: false,
					run: func() {
						task2Run = true
					},
				},
			}

			RunTasks(tasks)

			return !task1Run && !task2Run
		},
	},
	{
		name: "Запуск только готовых",
		check: func() bool {
			var task1Run, task2Run bool
			tasks := []Task{
				mockTask{
					ready: true,
					run: func() {
						task1Run = true
					},
				},
				mockTask{
					ready: false,
					run: func() {
						task2Run = true
					},
				},
			}

			RunTasks(tasks)

			return task1Run && !task2Run
		},
	},
	// Тесткейсы в помощь
	{
		name: "Проверка конкурентности",
		check: func() bool {
			ch := make(chan struct{})

			tasks := []Task{
				mockTask{
					ready: true,
					run: func() {
						ch <- struct{}{}
					},
				},
				mockTask{
					ready: true,
					run: func() {
						<-ch
					},
				},
			}

			RunTasks(tasks) // Если запуск не конкурентный, то зависнет

			return true
		},
	},
	{
		name: "Паника в таске",
		check: func() bool {
			tasks := []Task{
				mockTask{
					ready: true,
					run:   func() {},
				},
				mockTask{
					ready: true,
					run: func() {
						panic("panic")
					},
				},
			}

			RunTasks(tasks) // Должно завершиться

			return true
		},
	},
}

type mockTask struct {
	ready bool
	run   func()
}

func (m mockTask) IsReady() bool {
	return m.ready
}

func (m mockTask) Run() {
	m.run()
}
