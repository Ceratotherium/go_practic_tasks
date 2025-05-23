package main

import (
	"fmt"
	"os"
)

func AssertEqual[T comparable, IN any](message string, expected T, testFunc func(IN) T, input IN) {
	defer func() {
		if r := recover(); r != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Тест кейс %q - Паника: %s\n", message, r)
			os.Exit(1)
		}
	}()

	actual := testFunc(input)

	if expected != actual {
		_, _ = fmt.Fprintf(
			os.Stderr,
			"Тест кейс %q - провал\n\tОжидаемый результат - %v\n\tТекущий результат - %v\n\tВходные данные - %v\n",
			message,
			expected,
			actual,
			input,
		)
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(os.Stderr, "Тест кейс %q - успех\n", message)
}
