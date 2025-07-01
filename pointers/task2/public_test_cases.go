package main

type TestCase struct {
	name     string
	input    []*int
	expected int
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name:     "Поиск максимального",
		input:    toPointersArray(1, 2, 3, 4, 5),
		expected: 4,
	},
	{
		name:     "Все значения одинаковые",
		input:    toPointersArray(1, 1, 1, 1, 1),
		expected: 0,
	},
	{
		name:     "Слайс с нулем",
		input:    toPointersArray(0, 1, 2, 3, 4),
		expected: 4,
	},
	// Тесткейсы в помощь
	{
		name:     "Максимальное значение первое",
		input:    toPointersArray(5, 4, 3, 2, 1),
		expected: 0,
	},
	{
		name:     "Максимальное значение в середине",
		input:    toPointersArray(1, 2, 5, 4, 3),
		expected: 2,
	},
	{
		name:     "Есть 2 локальных максимума",
		input:    toPointersArray(1, 2, 5, 4, 10, 3),
		expected: 4,
	},
	{
		name:     "Только отрицательные числа",
		input:    toPointersArray(-5, -4, -3, -2, -1),
		expected: 4,
	},
	{
		name:     "Максимальное отрицательное число в середине",
		input:    toPointersArray(-5, -4, -1, -2, -3),
		expected: 2,
	},
	{
		name:     "Пустой слайс",
		input:    []*int{},
		expected: 0,
	},
	{
		name:     "Nil слайс",
		input:    nil,
		expected: 0,
	},
}

func toPointersArray(input ...int) []*int {
	output := make([]*int, len(input))
	for i, value := range input {
		output[i] = &value
	}

	return output
}
