package main

import (
	"context"
	"errors"
	"net/http"
)

type TestCase struct {
	name     string
	input    error
	expected bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name:     "Обычная ошибка",
		input:    errors.New("test error"),
		expected: false,
	},
	{
		name:     "serverNotReady",
		input:    serverNotReady,
		expected: true,
	},
	{
		name:     "context.DeadlineExceeded",
		input:    context.DeadlineExceeded,
		expected: true,
	},
	{
		name:     "Server error 504",
		input:    makeServerError(http.StatusGatewayTimeout),
		expected: true,
	},
	{
		name: "Database error timeout",
		input: &dataBaseError{
			isTimeout: true,
			message:   "query timeout",
		},
		expected: true,
	},
	// Тесткейсы в помощь
	{
		name:     "context.Canceled",
		input:    context.Canceled,
		expected: false,
	},
	{
		name:     "Server error 503",
		input:    makeServerError(http.StatusServiceUnavailable),
		expected: true,
	},
	{
		name: "Server error with changed message",
		input: serverError{
			statusCode: http.StatusServiceUnavailable,
			message:    "changed text",
		},
		expected: true,
	},
	{
		name:     "Nil ошибка",
		input:    nil,
		expected: false,
	},
}

func makeServerError(statusCode int) serverError {
	return serverError{
		statusCode: statusCode,
		message:    http.StatusText(statusCode),
	}
}
