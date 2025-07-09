package main

func ChannelToSlice(ch <-chan int) []int {
	values := make([]int, 0)
	for value := range ch {
		values = append(values, value)
	}

	return values
}
