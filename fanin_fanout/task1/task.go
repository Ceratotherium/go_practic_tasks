//go:build task_template

package main

func fanOut(ctx context.Context, inputCh <-chan string, workers int, callback func(string) string) []<-chan string {
	return nil
}
