package main

import (
	"context"
	"fmt"
	"net/http"
)

var privateTestCases = []TestCase{
	{
		name:     "Server error 500",
		input:    makeServerError(http.StatusInternalServerError),
		expected: false,
	},
	{
		name:     "Database error no timeout",
		input:    &dataBaseError{isTimeout: false},
		expected: false,
	},
	{
		name: "Обернутая server error",
		input: fmt.Errorf("fail to request: %w", serverError{
			statusCode: http.StatusServiceUnavailable,
			message:    "changed text",
		}),
		expected: true,
	},
	{
		name:     "Обернутая context.DeadlineExceeded",
		input:    fmt.Errorf("fail to request: %w", context.DeadlineExceeded),
		expected: true,
	},
	{
		name:     "Неправильно обернутая context.DeadlineExceeded",
		input:    fmt.Errorf("fail to request: %s", context.DeadlineExceeded),
		expected: false,
	},
}
