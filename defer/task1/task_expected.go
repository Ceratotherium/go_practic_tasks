package main

import (
	"errors"
	"fmt"
)

func SafeExecute(fn func() error) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()

	return fn()
}
