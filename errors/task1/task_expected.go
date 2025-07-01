package main

import (
	"context"
	"errors"
	"net/http"
)

var (
	serverNotReady = errors.New("server not ready")

	namedTemporaryErrors = []error{
		context.DeadlineExceeded,
		serverNotReady,
	}

	temporaryStatusCodes = map[int]struct{}{
		http.StatusServiceUnavailable: {},
		http.StatusGatewayTimeout:     {},
	}
)

type serverError struct {
	statusCode int
	message    string
}

func (e serverError) Error() string {
	return e.message
}

func (e serverError) StatusCode() int {
	return e.statusCode
}

type dataBaseError struct {
	isTimeout bool
	message   string
}

func (e *dataBaseError) Error() string {
	return e.message
}

func (e *dataBaseError) IsTimeout() bool {
	return e.isTimeout
}

func IsTemporaryError(err error) bool {
	serverErr := serverError{}

	if errors.As(err, &serverErr) {
		_, ok := temporaryStatusCodes[serverErr.statusCode]
		return ok
	}

	databaseErr := &dataBaseError{}
	if errors.As(err, &databaseErr) {
		return databaseErr.IsTimeout()
	}

	for _, targetErr := range namedTemporaryErrors {
		if errors.Is(err, targetErr) {
			return true
		}
	}

	return false
}
