package main

import "math"

var privateTestCases = []TestCase{
	{
		name:         "Сдвиг на MaxInt32",
		rotateFactor: math.MaxInt32,
		input:        [...]int{1, 2, 3, 4, 5},
		expected:     [...]int{3, 4, 5, 1, 2},
	},
	{
		name:         "Сдвиг на MinInt32",
		rotateFactor: math.MinInt32,
		input:        [...]int{1, 2, 3, 4, 5},
		expected:     [...]int{3, 4, 5, 1, 2},
	},
}
