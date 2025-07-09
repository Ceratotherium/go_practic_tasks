package main

import "sync/atomic"

var privateTestCases = []TestCase{
	{
		name: "Ожидание закрытия всех каналов",
		check: func() bool {
			ch1 := make(chan int)
			ch2 := make(chan int)

			results := FanIn(ch1, ch2)

			close(ch1)

			select {
			case _, open := <-results:
				if !open {
					return false // Если канал оказался закрытым
				}
			default:
			}

			chClosed := atomic.Bool{}

			go func() {
				close(ch2)
				chClosed.Store(true)
			}()

			if !isClosed(results) {
				return false
			}

			return chClosed.Load() == true
		},
	},
	{
		name: "Много значений в каналах",
		check: func() bool {
			ch1 := make(chan int)
			ch2 := make(chan int)

			go func() {
				defer close(ch1)
				defer close(ch2)

				for value := range 1000 {
					ch1 <- value
					ch2 <- value
				}
			}()

			results := FanIn(ch1, ch2)

			for range 2000 {
				<-results
			}

			return isClosed(results)
		},
	},
	{
		name: "Много каналов",
		check: func() bool {
			channels := make([]<-chan int, 0, 1000)
			for range 1000 {
				ch := make(chan int)
				channels = append(channels, ch)

				go func() {
					defer close(ch)

					for value := range 10 {
						ch <- value
					}
				}()
			}

			results := FanIn(channels...)

			for range 10000 {
				<-results
			}

			return isClosed(results)
		},
	},
}
