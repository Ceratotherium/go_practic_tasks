package main

type TestCase struct {
	name  string
	check func() bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Один выходной канал",
		check: func() bool {
			in := makeChan(1, 2, 3, 4, 5)

			out := FanOut(in, 1)
			for range 5 {
				<-out[0]
			}

			return isClosed(out[0])
		},
	},
	{
		name: "Несколько выходных каналов",
		check: func() bool {
			in := makeChan(1, 2, 3, 4, 5)

			out := FanOut(in, 2)

			totalCount := make(chan int, 2)
			for _, ch := range out {
				go func() {
					totalCount <- countOfValuesInChannel(ch)
				}()
			}

			count1 := <-totalCount
			count2 := <-totalCount

			return count1+count2 == 5
		},
	},
	// Тесткейсы в помощь
	{
		name: "Чтение только из одного канала",
		check: func() bool {
			in := makeChan(1, 2, 3, 4, 5)

			out := FanOut(in, 2)

			count1 := countOfValuesInChannel(out[0])
			count2 := countOfValuesInChannel(out[1])

			return count1 == 4 && count2 == 1
		},
	},
	{
		name: "Буфферизированный входной канал",
		check: func() bool {
			in := make(chan int, 5)
			for range 5 {
				in <- 1
			}
			close(in)

			out := FanOut(in, 2)

			totalCount := make(chan int, 2)
			for _, ch := range out {
				go func() {
					totalCount <- countOfValuesInChannel(ch)
				}()
			}

			count1 := <-totalCount
			count2 := <-totalCount

			return count1+count2 == 5
		},
	},
}

func makeChan(values ...int) chan int {
	ch := make(chan int)

	go func() {
		for _, v := range values {
			ch <- v
		}
		close(ch)
	}()
	return ch
}

func countOfValuesInChannel(ch <-chan int) int {
	count := 0
	for range ch {
		count++
	}
	return count
}

func isClosed(ch <-chan int) bool {
	_, opened := <-ch
	return !opened
}
