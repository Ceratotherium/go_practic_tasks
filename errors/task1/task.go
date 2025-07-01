//go:build task_template

package main

var serverNotReady = errors.New("server not ready")

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
	return false
}
