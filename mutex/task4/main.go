package main

func main() {
	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		ConcurrentCustomTestBody(
			tt.name,
			func() struct{} {
				cleanup()
				return struct{}{}
			},
			func(_ struct{}) bool {
				return tt.check()
			},
		)
	}
}
