package main

import "strings"

func ReverseString(str string) string {
	builder := strings.Builder{}
	builder.Grow(len(str))
	for i := len(str) - 1; i >= 0; i-- {
		builder.WriteByte(str[i])
	}
	return builder.String()
}
