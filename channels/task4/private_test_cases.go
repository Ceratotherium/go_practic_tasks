package main

import "time"

type privateTestCase struct {
	name  string
	check func() bool
}

var privateTestCases = []privateTestCase{
	{
		name: "Ожидание всех данных из канала",
		check: func() bool {
			var started time.Time
			ch := make(chan int)

			go func() {
				defer close(ch)
				ch <- 1

				started = time.Now()
				time.Sleep(time.Millisecond * 100)
			}()

			_ = ChannelToSlice(ch)
			elapsed := time.Since(started)

			return elapsed >= time.Millisecond*20
		},
	},
	{
		name: "Большое количество данных",
		check: func() bool {
			ch := make(chan int)

			go func() {
				defer close(ch)
				for i := range 1000 {
					ch <- i
				}
			}()

			values := ChannelToSlice(ch)
			return len(values) == 1000
		},
	},
}
