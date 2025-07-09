package main

import "sync"

func FanIn(channels ...<-chan int) <-chan int {
	results := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	for _, channel := range channels {
		go func() {
			defer wg.Done()
			for v := range channel {
				results <- v
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
