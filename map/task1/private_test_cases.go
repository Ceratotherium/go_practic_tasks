package main

var privateTestCases = []TestCase{
	{
		name: "Удаление несуществующего элемента из множества",
		prepare: func() *Set {
			s := NewSet(5)
			s.Add(1, 2)
			s.Remove(3)
			return s
		},
		check: func(s *Set) bool {
			return s.Len() == 2 && s.Contains(1) && s.Contains(2)
		},
	},
	{
		name: "Работа с пустым множеством",
		prepare: func() *Set {
			return NewSet(0)
		},
		check: func(s *Set) bool {
			return s.Len() == 0 && !s.Contains(0) && len(s.ToSlice()) == 0
		},
	},
	{
		name: "Добавление отрицательных чисел",
		prepare: func() *Set {
			s := NewSet(5)
			s.Add(1, 2, -1, -42)
			s.Remove(-1)
			return s
		},
		check: func(s *Set) bool {
			return s.Len() == 3 && s.Contains(-42)
		},
	},
	{
		name: "Добавление 10_000 уникальных элементов",
		prepare: func() *Set {
			s := NewSet(10_000)
			for i := 0; i < 10_000; i++ {
				s.Add(i)
			}
			return s
		},
		check: func(s *Set) bool {
			if s.Len() != 10_000 {
				return false
			}
			for i := 0; i < 10_000; i++ {
				if !s.Contains(i) {
					return false
				}
			}
			return true
		},
	},
	{
		name: "Добавление 100_000 элементов с дубликатами",
		prepare: func() *Set {
			s := NewSet(100)
			for i := 0; i < 100_000; i++ {
				s.Add(i % 500) // Добавляем только 500 уникальных значений
			}
			return s
		},
		check: func(s *Set) bool {
			return s.Len() == 500
		},
	},
	{
		name: "Последовательное добавление и удаление 50_000 элементов",
		prepare: func() *Set {
			s := NewSet(50_000)
			for i := 0; i < 50_000; i++ {
				s.Add(i)
			}
			for i := 0; i < 50_000; i++ {
				s.Remove(i)
			}
			return s
		},
		check: func(s *Set) bool {
			return s.Len() == 0
		},
	},
	{
		name: "Смешанные операции с 1_000_000 элементов",
		prepare: func() *Set {
			s := NewSet(1_000_000)
			// Добавляем 1M элементов
			for i := 0; i < 1_000_000; i++ {
				s.Add(i)
			}
			// Удаляем каждый второй
			for i := 0; i < 1_000_000; i += 2 {
				s.Remove(i)
			}
			// Добавляем новые
			for i := 1_000_000; i < 1_500_000; i++ {
				s.Add(i)
			}
			return s
		},
		check: func(s *Set) bool {
			// Проверяем что остались только нечетные из первых 1M и новые 500k
			if s.Len() != 500_000+500_000 {
				return false
			}
			for i := 1; i < 1_000_000; i += 2 {
				if !s.Contains(i) {
					return false
				}
			}
			for i := 1_000_000; i < 1_500_000; i++ {
				if !s.Contains(i) {
					return false
				}
			}
			return true
		},
	},
	{
		name: "Преобразование большого множества в слайс (100_000 элементов)",
		prepare: func() *Set {
			s := NewSet(100_000)
			for i := 0; i < 100_000; i++ {
				s.Add(i)
			}
			return s
		},
		check: func(s *Set) bool {
			slice := s.ToSlice()
			if len(slice) != 100_000 {
				return false
			}

			// Проверяем что все элементы уникальны
			unique := make(map[int]bool)
			for _, v := range slice {
				if unique[v] {
					return false
				}
				unique[v] = true
			}

			return true
		},
	},
	{
		name: "Проверка Contains на 1_000_000 элементов",
		prepare: func() *Set {
			s := NewSet(1_000_000)
			for i := 0; i < 1_000_000; i++ {
				s.Add(i)
			}
			return s
		},
		check: func(s *Set) bool {
			// Проверяем существующие
			for i := 0; i < 1_000_000; i += 10_000 {
				if !s.Contains(i) {
					return false
				}
			}
			// Проверяем несуществующие
			for i := -100; i < 0; i++ {
				if s.Contains(i) {
					return false
				}
			}
			for i := 1_000_000; i < 1_000_100; i++ {
				if s.Contains(i) {
					return false
				}
			}
			return true
		},
	},
}
