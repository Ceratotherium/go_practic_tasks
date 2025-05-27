package main

func removeElement(values []int, pos int) []int {
	if pos >= len(values) || pos < 0 || len(values) == 0 {
		return values
	}

	return append(values[:pos], values[pos+1:]...)
}
