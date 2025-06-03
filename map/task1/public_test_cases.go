package main

import "sort"

type TestCase struct {
	name    string
	prepare func() *Set
	check   func(*Set) bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Получение длины множества",
		prepare: func() *Set {
			s := NewSet(5)
			s.Add(1)
			s.Add(2)
			s.Add(3)
			s.Add(4)
			s.Add(5)
			s.Add(5)
			s.Add(5)
			return s
		},
		check: func(s *Set) bool {
			return s.Len() == 5
		},
	},
	{
		name: "Добавление одного элемента в множество",
		prepare: func() *Set {
			s := NewSet(5)
			s.Add(42)
			return s
		},
		check: func(s *Set) bool {
			return s.Len() == 1 && s.Contains(42)
		},
	},
	{
		name: "Удаление существующего элемента из множества",
		prepare: func() *Set {
			s := NewSet(5)
			s.Add(1)
			s.Add(2)
			s.Add(3)
			s.Remove(2)
			return s
		},
		check: func(s *Set) bool {
			return s.Len() == 2 && !s.Contains(2) && s.Contains(1) && s.Contains(3)
		},
	},
	{
		name: "Проверка наличия элемента в множестве",
		prepare: func() *Set {
			s := NewSet(5)
			s.Add(10)
			s.Add(20)
			return s
		},
		check: func(s *Set) bool {
			return s.Contains(10) && !s.Contains(30)
		},
	},
	{
		name: "Преобразование множества в слайс",
		prepare: func() *Set {
			s := NewSet(5)
			s.Add(1)
			s.Add(2)
			s.Add(3)
			return s
		},
		check: func(s *Set) bool {
			slice := s.ToSlice()
			sort.Ints(slice)
			return len(slice) == 3 && slice[0] == 1 && slice[1] == 2 && slice[2] == 3
		},
	},
	// Тесткейсы в помощь
	{
		name: "Создание нового множества с заданной емкостью",
		prepare: func() *Set {
			return NewSet(10)
		},
		check: func(s *Set) bool {
			return s.Len() == 0
		},
	},
	{
		name: "Добавление нескольких элементов в множество",
		prepare: func() *Set {
			s := NewSet(5)
			s.Add(1)
			s.Add(2)
			s.Add(3)
			s.Add(2)
			s.Add(1)
			return s
		},
		check: func(s *Set) bool {
			return s.Len() == 3 && s.Contains(1) && s.Contains(2) && s.Contains(3)
		},
	},
	{
		name: "Добавление нескольких элементов в множество (одной операцией)",
		prepare: func() *Set {
			s := NewSet(5)
			s.Add(1, 2, 3, 2, 1)
			return s
		},
		check: func(s *Set) bool {
			return s.Len() == 3 && s.Contains(1) && s.Contains(2) && s.Contains(3)
		},
	},
}
