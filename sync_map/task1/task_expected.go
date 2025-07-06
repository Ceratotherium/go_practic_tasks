package main

import (
	"strings"
	"sync"
)

func WordCount(text string, wordCount *sync.Map) {
	words := strings.Fields(text)

	for _, word := range words {
		word = strings.ToLower(word)
		_, loaded := wordCount.LoadOrStore(word, 1)

		if loaded {
			swapped := false
			for !swapped {
				count, _ := wordCount.Load(word)
				swapped = wordCount.CompareAndSwap(word, count.(int), count.(int)+1)
			}
		}
	}
}
