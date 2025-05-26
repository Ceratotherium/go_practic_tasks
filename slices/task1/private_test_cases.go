package main

var privateTestCases = []TestCase{
	{
		name:     "Не бьется на равные части",
		input:    []int{1, 2, 3, 4},
		chunkLen: 3,
		expected: [][]int{
			{1, 2, 3},
			{4},
		},
	},
	{
		name:     "Размер чанка больше входного слайса",
		input:    []int{1, 2, 3, 4},
		chunkLen: 5,
		expected: [][]int{
			{1, 2, 3, 4},
		},
	},
	{
		name:     "Перед нулевой размер чанка",
		input:    []int{1, 2, 3, 4},
		chunkLen: 0,
		expected: [][]int{},
	},
	{
		name:     "Пустой слайс",
		input:    []int{},
		chunkLen: 3,
		expected: [][]int{},
	},
	{
		name:     "Nil слайс",
		input:    nil,
		chunkLen: 3,
		expected: [][]int{},
	},
	{
		name:     "Одинаковые числа в слайсе",
		input:    []int{1, 1, 1, 1, 1, 1},
		chunkLen: 4,
		expected: [][]int{
			{1, 1, 1, 1},
			{1, 1},
		},
	},
}
