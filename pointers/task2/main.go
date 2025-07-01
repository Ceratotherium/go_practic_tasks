package main

func main() {
	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		AssertEqual(tt.name, tt.expected, GetMaxIndex, tt.input)
	}
}
