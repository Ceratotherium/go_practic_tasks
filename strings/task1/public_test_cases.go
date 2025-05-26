package main

type TestCase struct {
	input    string
	expected string
	name     string
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name:     "Просто строка",
		input:    "hello",
		expected: "olleh",
	},
	{
		name:     "Строка с числом",
		input:    "number 42",
		expected: "24 rebmun",
	},
	// Тесткейсы в помощь
	{
		name:     "Четное количество символов",
		input:    "name",
		expected: "eman",
	},
	{
		name:     "Два одинаковых символа",
		input:    "kk",
		expected: "kk",
	},
}
