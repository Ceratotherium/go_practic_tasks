package main

type TestCase struct {
	name     string
	input    <-chan int
	expected []int
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name:     "Канал с несколькими элементами",
		input:    generateChan(0, 1, 2, 3, 4, 5),
		expected: []int{1, 2, 3, 4, 5},
	},
	{
		name:     "Канал с повторяющимися значениями",
		input:    generateChan(0, 2, 2, 3, 3, 3, 5),
		expected: []int{2, 2, 3, 3, 3, 5},
	},
	// Тесткейсы в помощь
	{
		name: "Пустой канал",
		input: func() <-chan int {
			c := make(chan int)
			close(c)
			return c
		}(),
		expected: []int{},
	},
	{
		name:     "Канал с одним элементом",
		input:    generateChan(0, 1),
		expected: []int{1},
	},
	{
		name:     "Канал с отрицательными числами и нулём",
		input:    generateChan(0, -1, 0, 1, -2, 0),
		expected: []int{-1, 0, 1, -2, 0},
	},
	{
		name:     "Буфферизированный канал",
		input:    generateChan(6, 2, 2, 3, 3, 3, 5),
		expected: []int{2, 2, 3, 3, 3, 5},
	},
}

func generateChan(channelLen int, values ...int) <-chan int {
	ch := make(chan int, channelLen)

	go func() {
		for _, v := range values {
			ch <- v
		}
		close(ch)
	}()
	return ch
}
