package main

import "io"

func main() {
	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		CustomTestBody(
			tt.name,
			func() io.Writer {
				return tt.input
			},
			func(writer io.Writer) bool {
				return tt.check(writer, WriteHelloWorld(writer))
			},
		)
	}
}
