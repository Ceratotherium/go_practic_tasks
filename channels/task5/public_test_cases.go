package main

type TestCase struct {
	name  string
	check func() bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "1 канал",
		check: func() bool {
			ch := generateChan(0, 1, 2, 3, 4)

			results := FanIn(ch)

			expected := 0
			for value := range results {
				if value != expected {
					return false
				}
				expected++
			}

			return true
		},
	},
	{
		name: "Несколько каналов",
		check: func() bool {
			ch1 := generateChan(1, 2, 3, 4, 5)
			ch2 := generateChan(1, 2, 3, 4, 5)
			ch3 := generateChan(1, 2, 3, 4, 5)

			results := FanIn(ch1, ch2, ch3)

			values := map[int]int{}
			for value := range results {
				values[value]++
			}

			if len(values) != 5 {
				return false
			}

			for _, count := range values {
				if count != 3 {
					return false
				}
			}

			return true
		},
	},
	// Тесткейсы в помощь
	{
		name: "Нет каналов на вход",
		check: func() bool {
			results := FanIn()
			return isClosed(results)
		},
	},
	{
		name: "Данные только в одном канале",
		check: func() bool {
			ch1 := generateChan(1, 2, 3, 4, 5)
			ch2 := make(chan int)

			results := FanIn(ch1, ch2)

			for range 5 {
				<-results
			}

			close(ch2)

			return isClosed(results)
		},
	},
	{
		name: "Буфферизированные каналы",
		check: func() bool {
			ch1 := make(chan int, 5)
			ch2 := make(chan int, 5)
			for val := range 5 {
				ch1 <- val
				ch2 <- val
			}

			results := FanIn(ch1, ch2)

			for range 10 {
				<-results
			}

			close(ch1)
			close(ch2)

			return isClosed(results)
		},
	},
}

func isClosed(ch <-chan int) bool {
	_, opened := <-ch
	return !opened
}

func generateChan(values ...int) <-chan int {
	ch := make(chan int)

	go func() {
		for _, v := range values {
			ch <- v
		}
		close(ch)
	}()
	return ch
}
