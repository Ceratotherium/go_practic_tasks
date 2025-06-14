package main

func main() {
	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		ConcurrentCustomTestBody(
			tt.name,
			func() <-chan int {
				return SliceToChannel(tt.input)
			},
			func(ch <-chan int) bool {
				for _, value := range tt.input {
					if chValue, closed := <-ch; !closed || chValue != value {
						return false
					}
				}

				_, closed := <-ch
				return !closed
			},
		)
	}

	ConcurrentCustomTestBody(
		"Проверка капасити канала",
		func() <-chan int {
			return SliceToChannel(make([]int, 1000))
		},
		func(ch <-chan int) bool {
			return cap(ch) == 0
		},
	)
}
