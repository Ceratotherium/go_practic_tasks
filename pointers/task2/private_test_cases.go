package main

import "math"

var privateTestCases = []TestCase{
	{
		name:     "Два самых маленьких числа",
		input:    toPointersArray(math.MinInt32, math.MinInt32+1),
		expected: 1,
	},
	{
		name:     "Два самых больших числа",
		input:    toPointersArray(math.MaxInt32-1, math.MaxInt32),
		expected: 1,
	},
	{
		name:     "Массив с nil",
		input:    []*int{pointer(1), pointer(2), nil, pointer(3), nil},
		expected: 3,
	},
	{
		name:     "Массив с одними nil",
		input:    []*int{nil, nil, nil, nil, nil, nil, nil},
		expected: 0,
	},
	{
		name:     "Одно значение MinInt32, остальные nil",
		input:    []*int{nil, nil, pointer(math.MinInt32), nil, nil, nil, nil},
		expected: 2,
	},
}

func pointer(value int) *int {
	return &value
}
