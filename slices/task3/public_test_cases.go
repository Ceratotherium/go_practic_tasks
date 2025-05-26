package main

type TestCase struct {
	name     string
	input    []int
	pos      int
	expected []int
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name:     "Удаление первого элемента",
		input:    []int{1, 2, 3},
		pos:      0,
		expected: []int{2, 3},
	},
	{
		name:     "Удаление элемента в середине",
		input:    []int{1, 2, 3},
		pos:      1,
		expected: []int{1, 3},
	},
	// Тесткейсы в помощь
	{
		name:     "Удаление последнего элемента",
		input:    []int{1, 2, 3},
		pos:      2,
		expected: []int{1, 2},
	},
	{
		name:     "Удаляемый элемент выходит за границы слайса",
		input:    []int{1, 2, 3},
		pos:      10,
		expected: []int{1, 2, 3},
	},
}
