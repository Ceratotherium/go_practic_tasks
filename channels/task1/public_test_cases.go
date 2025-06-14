package main

type TestCase struct {
	name  string
	input []int
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name:  "Пустой слайс",
		input: []int{},
	},
	{
		name:  "Слайс с одним элементом",
		input: []int{42},
	},
	{
		name:  "Слайс с несколькими элементами",
		input: []int{1, 2, 3, 4, 5},
	},
	// Тесткейсы в помощь
	{
		name:  "Слайс с отрицательными числами",
		input: []int{-1, -2, -3},
	},
	{
		name:  "Слайс с нулевыми значениями",
		input: []int{0, 0, 0},
	},
}
