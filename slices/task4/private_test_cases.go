package main

var privateTestCases = []TestCase{
	{
		name:     "Передали пустой слайс",
		input:    []int{},
		pos:      3,
		expected: []int{},
	},
	{
		name:     "Передали nil слайс",
		input:    nil,
		pos:      3,
		expected: nil,
	},
	{
		name:     "Удаляем первый элемент из пустого слайса",
		input:    []int{},
		pos:      0,
		expected: []int{},
	},
	{
		name:     "Отрицательная позиция",
		input:    []int{1, 2, 3},
		pos:      -1,
		expected: []int{1, 2, 3},
	},
}
