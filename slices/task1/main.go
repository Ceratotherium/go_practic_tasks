package main

func main() {
	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		AssertEqualT(tt.name, tt.expected, func(values []int) [][]int {
			return Chunk(values, tt.chunkLen)
		}, tt.input, func(expected [][]int, actual [][]int) bool {
			if len(actual) != len(expected) {
				return false
			}

			for i := 0; i < len(actual); i++ {
				if !compareSliceValues(actual[i], expected[i]) {
					return false
				}
			}
			return true
		})
	}
}
