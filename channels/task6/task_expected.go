package main

func FanOut(in <-chan int, count int) []<-chan int {
	out := make([]<-chan int, count)
	for i := 0; i < count; i++ {
		ch := make(chan int)
		out[i] = ch

		go func() {
			defer close(ch)

			for n := range in {
				ch <- n
			}
		}()
	}

	return out
}
