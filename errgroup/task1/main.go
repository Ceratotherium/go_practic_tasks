package main

func main() {
	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		ConcurrentCustomTestBody(tt.name, func() []Request { return tt.prepare }, tt.check)
	}
}
