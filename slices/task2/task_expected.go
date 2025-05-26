package main

type Stack struct {
	values    []int
	maxValues []int
}

func (s *Stack) Push(v int) {
	s.values = append(s.values, v)

	if len(s.maxValues) == 0 {
		s.maxValues = append(s.maxValues, v)
	} else {
		s.maxValues = append(s.maxValues, max(s.maxValues[len(s.maxValues)-1], v))
	}
}

func (s *Stack) Pop() int {
	if len(s.values) == 0 {
		return 0
	}

	v := s.values[len(s.values)-1]
	s.values = s.values[:len(s.values)-1]
	s.maxValues = s.maxValues[:len(s.maxValues)-1]
	return v
}

func (s *Stack) Len() int {
	return len(s.values)
}

func (s *Stack) Values() []int {
	res := make([]int, len(s.values))
	copy(res, s.values)
	return res
}

func (s *Stack) Max() int {
	if len(s.maxValues) == 0 {
		return 0
	}

	return s.maxValues[len(s.maxValues)-1]
}
