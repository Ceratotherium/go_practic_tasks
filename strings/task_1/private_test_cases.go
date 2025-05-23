package main

var privateTestCases = []TestCase{
	{
		name:     "Палиндром с заглавной буквой",
		input:    "Eve",
		expected: true,
	},
	{
		name:     "Фраза палиндром",
		input:    "A man, a plan, a canal: Panama",
		expected: true,
	},
	{
		name:     "Число палиндром",
		input:    "2002",
		expected: true,
	},
	{
		name:     "Один символ",
		input:    "a",
		expected: true,
	},
	{
		name:     "Пустая строка",
		input:    "",
		expected: true,
	},
	{
		name:     "Палиндром на руссом языке",
		input:    "Искать такси",
		expected: true,
	},
}
