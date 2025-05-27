package main

import "math/rand"

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
	{
		name: "Большой стек с рандомными данными",
		testBody: func() *Stack {
			s := &Stack{}
			rand.Seed(42) // фиксируем seed для воспроизводимости
			for i := 0; i < 10000; i++ {
				s.Push(rand.Intn(1000))
			}
			return s
		},
		check: func(s *Stack) bool {
			if s.Len() != 10000 {
				return false
			}

			// Проверяем, что Max() действительно возвращает максимальное значение
			values := s.Values()
			maxVal := 0
			for _, v := range values {
				if v > maxVal {
					maxVal = v
				}
			}
			return s.Max() == maxVal
		},
	},
	{
		name: "Большой стек с возрастающими данными",
		testBody: func() *Stack {
			s := &Stack{}
			for i := 0; i < 10000; i++ {
				s.Push(i)
			}
			return s
		},
		check: func(s *Stack) bool {
			if s.Len() != 10000 {
				return false
			}

			values := s.Values()
			for i := 0; i < 10000; i++ {
				if values[i] != i {
					return false
				}
			}
			return s.Max() == 9999
		},
	},
	{
		name: "Большой стек с уменьшающимися данными",
		testBody: func() *Stack {
			s := &Stack{}
			for i := 9999; i >= 0; i-- {
				s.Push(i)
			}
			return s
		},
		check: func(s *Stack) bool {
			if s.Len() != 10000 {
				return false
			}

			values := s.Values()
			for i := 0; i < 10000; i++ {
				if values[i] != 9999-i {
					return false
				}
			}
			return s.Max() == 9999
		},
	},
	{
		name: "Много операций удаления данных",
		testBody: func() *Stack {
			s := &Stack{}
			for i := 0; i < 100; i++ {
				s.Push(i)
			}
			for i := 0; i < 50; i++ {
				s.Pop()
			}
			return s
		},
		check: func(s *Stack) bool {
			if s.Len() != 50 {
				return false
			}

			values := s.Values()
			for i := 0; i < 50; i++ {
				if values[i] != i {
					return false
				}
			}
			return s.Max() == 49
		},
	},
	{
		name: "Провеврка максимального значения в стеке с отрицательными числами",
		testBody: func() *Stack {
			stack := Stack{}
			stack.Push(-1)
			stack.Push(-2)
			stack.Push(-3)

			return &stack
		},
		check: func(stack *Stack) bool {
			return stack.Max() == -1
		},
	},
}
