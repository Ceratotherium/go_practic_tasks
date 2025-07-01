package main

var privateTestCases = []TestCase{
	{
		name: "Проверка нуля",
		input: tuple{
			value1: -1,
			value2: 0,
		},
		expected: tuple{
			value1: 0,
			value2: -1,
		},
	},
}
