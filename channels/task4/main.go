package main

func main() {
	for _, tt := range testCases {
		ConcurrentCustomTestBody(
			tt.name,
			func() <-chan int {
				return tt.input
			},
			func(values <-chan int) bool {
				return compareSliceValues(tt.expected, ChannelToSlice(values))
			},
		)
	}

	for _, tt := range privateTestCases {
		ConcurrentCustomTestBody(
			tt.name,
			func() struct{} {
				return struct{}{}
			},
			func(_ struct{}) bool {
				return tt.check()
			})
	}
}
