package main

type Set struct {
	values map[int]struct{}
}

func NewSet(length int) *Set {
	return &Set{
		values: make(map[int]struct{}, length),
	}
}

func (s *Set) Add(value ...int) {
	for _, v := range value {
		s.values[v] = struct{}{}
	}
}

func (s *Set) Remove(value int) {
	delete(s.values, value)
}

func (s *Set) Contains(value int) bool {
	_, ok := s.values[value]
	return ok
}

func (s *Set) Len() int {
	return len(s.values)
}

func (s *Set) ToSlice() []int {
	values := make([]int, 0, s.Len())
	for v := range s.values {
		values = append(values, v)
	}
	return values
}
