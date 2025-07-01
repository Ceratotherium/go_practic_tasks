package main

type TestCase struct {
	name         string
	rotateFactor int
	input        [5]int
	expected     [5]int
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name:         "Сдвиг на 1",
		rotateFactor: 1,
		input:        [...]int{1, 2, 3, 4, 5},
		expected:     [...]int{2, 3, 4, 5, 1},
	},
	{
		name:         "Сдвиг на 3",
		rotateFactor: 3,
		input:        [...]int{1, 2, 3, 4, 5},
		expected:     [...]int{4, 5, 1, 2, 3},
	},
	// Тесткейсы в помощь
	{
		name:         "Сдвиг на 0 не меняет массив",
		rotateFactor: 0,
		input:        [...]int{1, 2, 3, 4, 5},
		expected:     [...]int{1, 2, 3, 4, 5},
	},
	{
		name:         "Сдвиг на 5 не меняет массив",
		rotateFactor: 5,
		input:        [...]int{1, 2, 3, 4, 5},
		expected:     [...]int{1, 2, 3, 4, 5},
	},
	{
		name:         "Сдвиг на -2 (эквивалентно сдвигу на 3)",
		rotateFactor: -2,
		input:        [...]int{1, 2, 3, 4, 5},
		expected:     [...]int{4, 5, 1, 2, 3},
	},
	{
		name:         "Сдвиг на 8 (эквивалентно сдвигу на 3)",
		rotateFactor: 8,
		input:        [...]int{1, 2, 3, 4, 5},
		expected:     [...]int{4, 5, 1, 2, 3},
	},
}
