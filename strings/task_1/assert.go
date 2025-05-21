package main

import (
	"fmt"
	"os"
)

func AssertEqual[T comparable, IN any](expected T, actual T, message string, input ...IN) {
	if expected != actual {
		_, _ = fmt.Fprintf(
			os.Stderr,
			"Тест кейс %q - провал\n\tОжидаемый результат - %v\n\tТекущий результат - %v\n\tВходные данные - %+v\n",
			message,
			expected,
			actual,
			input,
		)
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(os.Stderr, "Тест кейс %q - успех\n", message)
}
