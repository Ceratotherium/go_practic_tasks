package main

type privateTestCase struct {
	name  string
	input []int
	check func([]int) bool
}

var privateTestCases = []privateTestCase{
	{
		name:  "Память под слайс выделена заранее",
		input: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		check: func(values []int) bool {
			return len(values) == cap(values)
		},
	},
}
