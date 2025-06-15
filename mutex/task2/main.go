package main

func main() {
	type testArgs struct {
		c      Cache
		cancel func()
	}

	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		ConcurrentCustomTestBody(
			tt.name,
			func() testArgs {
				c, cancel := tt.prepare()
				return testArgs{
					c:      c,
					cancel: cancel,
				}
			},
			func(args testArgs) bool {
				defer args.cancel()
				return tt.check(args.c)
			},
		)
	}
}
