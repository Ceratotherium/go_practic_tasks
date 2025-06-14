package main

func SliceToChannel(values []int) <-chan int {
	ch := make(chan int)

	go func() {
		for _, value := range values {
			ch <- value
		}
		close(ch)
	}()

	return ch
}
