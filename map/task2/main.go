package main

func main() {
	for _, tt := range testCases {
		AssertEqualValues(tt.name, tt.expected, Unique, tt.input)
	}

	for _, tt := range privateTestCases {
		CustomTestBody(tt.name, func() []int { return Unique(tt.input) }, tt.check)
	}
}
