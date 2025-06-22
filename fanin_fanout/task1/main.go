package main

func main() {
	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		ConcurrentCustomTestBody(
			tt.name,
			func() chan string {
				result := make(chan string)

				go func() {
					defer close(result)
					for _, str := range tt.in {
						result <- str
					}
				}()

				return result
			},
			func(in chan string) bool {
				return tt.check(Process(in, tt.convert, tt.procNum, tt.sumCount))
			},
		)
	}
}
