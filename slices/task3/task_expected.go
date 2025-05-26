package main

func removeElement(values []int, pos int) []int {
	if pos >= len(values) || pos < 0 || len(values) == 0 {
		return values
	}

	result := make([]int, 0, len(values)-1)
	return append(append(result, values[:pos]...), values[pos+1:]...)
}
