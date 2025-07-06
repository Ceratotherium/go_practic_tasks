package main

import "sync"

var privateTestCases = []TestCase{
	{
		name: "Большое количество строк",
		prepare: func() []string {
			res := make([]string, 1000)
			for i := 0; i < 1000; i++ {
				res[i] = "one two three four five"
			}
			return res
		}(),
		check: func(words *sync.Map) bool {
			counter := counter{words}
			return counter.Count("one") == 1000 &&
				counter.Count("two") == 1000 &&
				counter.Count("three") == 1000 &&
				counter.Count("four") == 1000 &&
				counter.Count("five") == 1000
		},
	},
}
