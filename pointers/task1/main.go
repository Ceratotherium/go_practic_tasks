package main

func main() {
	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		CustomTestBody(
			tt.name,
			func() tuple {
				return tt.input
			},
			func(input tuple) bool {
				Swap(&input.value1, &input.value2)
				return input == tt.expected
			},
		)
	}
}
