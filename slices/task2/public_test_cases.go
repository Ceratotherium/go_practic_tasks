package main

type TestCase struct {
	name     string
	testBody func() *Stack
	check    func(stack *Stack) bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Добавление элемента",
		testBody: func() *Stack {
			stack := Stack{}
			stack.Push(1)
			stack.Push(2)

			return &stack
		},
		check: func(stack *Stack) bool {
			return stack.Len() == 2
		},
	},
	{
		name: "Удаление элемента",
		testBody: func() *Stack {
			stack := Stack{}
			stack.Push(1)
			stack.Push(2)

			return &stack
		},
		check: func(stack *Stack) bool {
			return stack.Pop() == 2
		},
	},
	{
		name: "Максимальный элемент",
		testBody: func() *Stack {
			stack := Stack{}
			stack.Push(1)
			stack.Push(2)

			return &stack
		},
		check: func(stack *Stack) bool {
			return stack.Max() == 2
		},
	},
	{
		name: "Получение значений стека",
		testBody: func() *Stack {
			stack := Stack{}
			stack.Push(1)
			stack.Push(2)

			return &stack
		},
		check: func(stack *Stack) bool {
			values := stack.Values()
			return values[0] == 1 && values[1] == 2
		},
	},
	// Тесткейсы в помощь
	{
		name: "Проверяем, что элементы добавляются в конец",
		testBody: func() *Stack {
			stack := Stack{}
			stack.Push(42)
			stack.Push(2)
			stack.Push(100)

			return &stack
		},
		check: func(stack *Stack) bool {
			popValue := stack.Pop()

			return popValue == 100 && compareSliceValues([]int{42, 2}, stack.Values())
		},
	},
	{
		name: "Максимальный элемент самый первый",
		testBody: func() *Stack {
			stack := Stack{}
			stack.Push(42)
			stack.Push(2)

			return &stack
		},
		check: func(stack *Stack) bool {
			return stack.Max() == 42
		},
	},
	{
		name: "Проверка максимального элемента после удаления максимального",
		testBody: func() *Stack {
			stack := Stack{}
			stack.Push(2)
			stack.Push(42)

			return &stack
		},
		check: func(stack *Stack) bool {
			stack.Pop()
			return stack.Max() == 2
		},
	},
}
