package main

import (
	"errors"
	"io"
	"reflect"
)

func WriteHelloWorld(out io.Writer) error {
	if out == nil || reflect.ValueOf(out).IsZero() {
		return errors.New("nil interface")
	}

	_, err := out.Write([]byte("Hello world"))
	return err
}
