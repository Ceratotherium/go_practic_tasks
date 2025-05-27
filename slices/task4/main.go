package main

func main() {
	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		AssertEqualValues(tt.name, tt.expected, func(values []int) []int {
			return removeElement(values, tt.pos)
		}, tt.input)
	}

	CustomTestBody("Проверка, что создается копия массива",
		func() []int {
			return []int{1, 2, 3}
		},
		func(values []int) bool {
			result := removeElement(values, 1)
			result[0] = 0

			return values[0] == 0
		},
	)
}
