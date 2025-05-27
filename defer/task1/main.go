package main

func main() {
	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		CustomTestBody(
			tt.name,
			func() error {
				return SafeExecute(tt.execute)
			},
			tt.check,
		)
	}
}
