package main

import (
	"bytes"
	"io"
)

var privateTestCases = []TestCase{
	{
		name:  "Ошибка при нулевом значении интерфейса",
		input: (*bytes.Buffer)(nil),
		check: func(writer io.Writer, err error) bool {
			return err != nil
		},
	},
}
