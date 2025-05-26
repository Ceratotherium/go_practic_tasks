package main

var privateTestCases = []TestCase{
	{
		name: "Пустой стек - пустой",
		testBody: func() *Stack {
			return &Stack{}
		},
		check: func(stack *Stack) bool {
			return stack.Len() == 0 && stack.Max() == 0 && compareSliceValues(stack.Values(), []int{})
		},
	},
	{
		name: "Pop на пустом стеке",
		testBody: func() *Stack {
			return &Stack{}
		},
		check: func(stack *Stack) bool {
			return stack.Pop() == 0
		},
	},
	{
		name: "Изменение слайса Values не меняет внутреннее состояние стека",
		testBody: func() *Stack {
			stack := Stack{}
			stack.Push(2)
			stack.Push(42)

			return &stack
		},
		check: func(stack *Stack) bool {
			values := stack.Values()
			values[0] = 0
			return compareSliceValues(stack.Values(), []int{2, 42})
		},
	},
}
