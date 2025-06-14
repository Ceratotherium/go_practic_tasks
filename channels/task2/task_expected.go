package main

import "fmt"

func PrintFive() {
	ch := make(chan int)

	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
		}
		close(ch)
	}()

	for n := range ch {
		fmt.Println(n)
	}
}
