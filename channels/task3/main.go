package main

func main() {
	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		ConcurrentCustomTestBody(tt.name, tt.prepare, func(backends []Backend) bool {
			return tt.check(DoRequests(backends))
		})
	}
}
