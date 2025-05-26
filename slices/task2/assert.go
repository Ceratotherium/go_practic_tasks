package main

import (
	"fmt"
	"os"
)

func AssertEqual[T comparable, IN any](message string, expected T, testFunc func(IN) T, input IN) {
	AssertEqualT[T, IN](message, expected, testFunc, input, compareSimpleTypes[T])
}

func AssertEqualValues[T comparable, IN any](message string, expected []T, testFunc func(IN) []T, input IN) {
	AssertEqualT[[]T, IN](message, expected, testFunc, input, compareSliceValues[T])
}

func AssertEqualT[T any, IN any](message string, expected T, testFunc func(IN) T, input IN, compare func(T, T) bool) {
	defer catchPanic(message)()

	actual := testFunc(input)

	if !compare(expected, actual) {
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

func CustomTestBody[T any](message string, prepare func() T, check func(T) bool) {
	defer catchPanic(message)()

	isSuccess := check(prepare())

	if !isSuccess {
		_, _ = fmt.Fprintf(
			os.Stderr,
			"Тест кейс %q - провал\n",
			message,
		)
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(os.Stderr, "Тест кейс %q - успех\n", message)
}

func compareSimpleTypes[T comparable](expected T, actual T) bool {
	return expected == actual
}

func compareSliceValues[T comparable](expected []T, actual []T) bool {
	if len(expected) != len(actual) {
		return false
	}

	for i := 0; i < len(expected); i++ {
		if expected[i] != actual[i] {
			return false
		}
	}

	return true
}

func catchPanic(message string) func() {
	return func() {
		if r := recover(); r != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Тест кейс %q - Паника: %s\n", message, r)
			os.Exit(1)
		}
	}
}
