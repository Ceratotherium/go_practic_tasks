package main

type TestCase struct {
	name     string
	expected string
}

var testCases = []TestCase{
	{
		name:     "Проверка вывода",
		expected: "0\n1\n2\n3\n4\n",
	},
}
