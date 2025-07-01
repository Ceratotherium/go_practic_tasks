package main

func Swap(val1 *int, val2 *int) {
	*val1, *val2 = *val2, *val1
}
