package main

type TestCase struct {
	name     string
	input    tuple
	expected tuple
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Два разных числа",
		input: tuple{
			value1: 1,
			value2: 2,
		},
		expected: tuple{
			value1: 2,
			value2: 1,
		},
	},
	{
		name: "Два одинаковых числа",
		input: tuple{
			value1: 2,
			value2: 2,
		},
		expected: tuple{
			value1: 2,
			value2: 2,
		},
	},
	// Тесткейсы в помощь
	{
		name: "Числа разных знаков",
		input: tuple{
			value1: -1,
			value2: 2,
		},
		expected: tuple{
			value1: 2,
			value2: -1,
		},
	},
	{
		name: "Отрицательные числа",
		input: tuple{
			value1: -1,
			value2: -2,
		},
		expected: tuple{
			value1: -2,
			value2: -1,
		},
	},
}

type tuple struct {
	value1 int
	value2 int
}
