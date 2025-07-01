package main

func main() {
	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		AssertEqualValues(tt.name, tt.expected[:], func(values [5]int) []int {
			result := RotateLeft(values, tt.rotateFactor)
			return result[:]
		}, tt.input)
	}
}
