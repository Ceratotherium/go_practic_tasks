package main

type TestCase struct {
	name     string
	input    []int
	expected []int
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name:     "Все элементы уникальны",
		input:    []int{1, 2, 3, 4},
		expected: []int{1, 2, 3, 4},
	},
	{
		name:     "Все элементы одинаковые",
		input:    []int{5, 5, 5, 5},
		expected: []int{5},
	},
	{
		name:     "Есть несколько дубликатов",
		input:    []int{1, 2, 2, 3, 4, 4, 4, 5},
		expected: []int{1, 2, 3, 4, 5},
	},
	// Тесткейсы в помощь
	{
		name:     "Пустой слайс на входе",
		input:    []int{},
		expected: []int{},
	},
	{
		name:     "Отрицательные числа и нуль",
		input:    []int{-1, 0, 1, -1, 0, 2},
		expected: []int{-1, 0, 1, 2},
	},
	{
		name: "Большой слайс с повторяющимися значениями",
		input: func() []int {
			values := make([]int, 0, 1000)
			for i := 0; i < 200; i++ {
				values = append(values, 1, 2, 3, 4, 5)
			}
			return values
		}(),
		expected: []int{1, 2, 3, 4, 5},
	},
	{
		name:     "Большой слайс с уникальными значениями",
		input:    make1000Values(),
		expected: make1000Values(),
	},
}

func make1000Values() []int {
	values := make([]int, 0, 1000)
	for i := 0; i < 1000; i++ {
		values = append(values, i)
	}
	return values
}
