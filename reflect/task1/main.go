package main

func main() {
	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		CustomTestBody(
			tt.name,
			func() struct{} {
				return struct{}{}
			},
			func(_ struct{}) bool {
				results, err := CallMethod(tt.object, tt.method, tt.input...)
				if err != nil {
					return tt.expectedErr
				}

				if tt.expectedErr {
					return false
				}

				if len(results) != len(tt.expected) {
					return false
				}

				return compareSliceValues(tt.expected, results)
			},
		)
	}
}
