package main

import (
	"sync"
)

func FanOut(inputCh <-chan string, workers int) []<-chan string {
	resultChannels := make([]<-chan string, workers)

	for i := 0; i < workers; i++ {
		ch := make(chan string)
		resultChannels[i] = ch

		go func(ch chan<- string) {
			defer close(ch)

			for input := range inputCh {
				ch <- input
			}
		}(ch)
	}

	return resultChannels
}

func FanIn(inputChs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	wg.Add(len(inputChs))

	for _, c := range inputChs {
		go func(c <-chan int) {
			for n := range c {
				out <- n
			}
			wg.Done()
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func ProcessFunc(in <-chan string, callback func(string) int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for n := range in {
			out <- callback(n)
		}
	}()

	return out
}

func SumNumbers(in <-chan int, count int) <-chan int {
	result := make(chan int)

	go func() {
		defer close(result)

		for {
			sum := 0
			for num := range count {
				value, ok := <-in
				if !ok {
					if num != 0 {
						result <- sum
					}

					return
				}

				sum += value
			}

			result <- sum
		}
	}()

	return result
}

func Process(in <-chan string, converter func(string) int, procNum int, sumCount int) <-chan int {
	procChannels := make([]<-chan int, 0, procNum)
	for _, ch := range FanOut(in, procNum) {
		procChannels = append(procChannels, ProcessFunc(ch, converter))
	}

	return SumNumbers(FanIn(procChannels...), sumCount)
}
