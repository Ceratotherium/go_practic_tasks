package main

func main() {
	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		CustomTestBody(tt.name, tt.testBody, tt.check)
	}
}
