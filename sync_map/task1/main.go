package main

func main() {
	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		ConcurrentCustomTestBody(
			tt.name,
			func() []string {
				return tt.prepare
			},
			func(values []string) bool {
				return tt.check(testBody(tt.name, values))
			},
		)
	}
}
