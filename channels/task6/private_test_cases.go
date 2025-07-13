package main

import "sync/atomic"

var privateTestCases = []TestCase{
	{
		name: "Ожидание закрытия канала",
		check: func() bool {
			in := make(chan int)

			out := FanOut(in, 2)

			for _, ch := range out {
				select {
				case _, open := <-ch:
					if !open {
						return false // Если канал оказался закрытым
					}
				default:
				}
			}

			chClosed := atomic.Bool{}

			go func() {
				close(in)
				chClosed.Store(true)
			}()

			for _, ch := range out {
				if !isClosed(ch) {
					return false
				}
			}

			return chClosed.Load() == true
		},
	},
	{
		name: "Много значений",
		check: func() bool {
			in := make(chan int)

			go func() {
				for val := range 1000 * 1000 {
					in <- val
				}
				close(in)
			}()

			out := FanOut(in, 1000)

			counts := make(chan int, 1000)
			for _, ch := range out {
				go func() {
					counts <- countOfValuesInChannel(ch)
				}()
			}

			totalCount := 0
			for range 1000 {
				totalCount += <-counts
			}

			return totalCount == 1000*1000
		},
	},
}
