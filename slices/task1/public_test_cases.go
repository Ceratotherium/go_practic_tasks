package main

type TestCase struct {
	input    []int
	chunkLen int
	expected [][]int
	name     string
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name:     "Размер чанка 2",
		input:    []int{1, 2, 3, 4},
		chunkLen: 2,
		expected: [][]int{
			{1, 2},
			{3, 4},
		},
	},
	{
		name:     "Размер чанка 3",
		input:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		chunkLen: 3,
		expected: [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		},
	},
	// Тесткейсы в помощь
	{
		name:     "Размер чанка 1",
		input:    []int{1, 2, 3, 4},
		chunkLen: 1,
		expected: [][]int{
			{1}, {2}, {3}, {4},
		},
	},
	{
		name:     "Размер чанка равен размеру слайса",
		input:    []int{1, 2, 3, 4},
		chunkLen: 4,
		expected: [][]int{
			{1, 2, 3, 4},
		},
	},
}
