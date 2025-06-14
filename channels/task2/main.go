package main

func main() {
	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		AssertPrint(tt.name, tt.expected, PrintFive)
	}
}
