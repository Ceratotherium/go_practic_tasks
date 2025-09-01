//go:build task_template

package main

func Process(in <-chan string, converter func(string) int, procNum int, sumCount int) <-chan int {
	return nil
}
