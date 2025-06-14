package main

import "strings"

type TestCase struct {
	name  string
	check func() bool
}

var testCases = []TestCase{
	{
		name: "Проверка вывода",
		check: func() bool {
			output := catchPrint(MakeDrink)

			return strings.Contains(output, "Coke") && strings.Contains(output, "Fanta")
		},
	},
}
