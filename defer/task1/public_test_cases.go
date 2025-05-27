package main

import (
	"errors"
	"strings"
)

var testErr = errors.New("test error")

type TestCase struct {
	name    string
	execute func() error
	check   func(error) bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name:    "Паники нет",
		execute: func() error { return nil },
		check: func(err error) bool {
			return err == nil
		},
	},
	{
		name: "Паника есть",
		execute: func() error {
			panic(errors.New("test error"))
			return nil
		},
		check: func(err error) bool {
			return err != nil && strings.Contains(err.Error(), "test error")
		},
	},
	// Тесткейсы в помощь
	{
		name:    "Функция возвращает ошибку",
		execute: func() error { return testErr },
		check: func(err error) bool {
			return errors.Is(err, testErr)
		},
	},
}

func init() {

}
