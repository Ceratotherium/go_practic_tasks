package main

import (
	"strings"
)

var privateTestCases = []TestCase{
	{
		name: "Паника - это строка",
		execute: func() error {
			panic("panic is string")
			return nil
		},
		check: func(err error) bool {
			return strings.Contains(err.Error(), "panic is string")
		},
	},
}
