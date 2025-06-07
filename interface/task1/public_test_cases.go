package main

import (
	"bytes"
	"errors"
	"io"
)

type TestCase struct {
	name  string
	input io.Writer
	check func(writer io.Writer, err error) bool
}

var testCases = []TestCase{
	// Публичные тесты
	{
		name:  "Успешная запись в буфер",
		input: &bytes.Buffer{},
		check: func(writer io.Writer, err error) bool {
			if err != nil {
				return false
			}
			buf := writer.(*bytes.Buffer)
			return buf.String() == "Hello world"
		},
	},
	{
		name:  "Ошибка записи в writer",
		input: &errorWriter{},
		check: func(writer io.Writer, err error) bool {
			return errors.Is(err, testError)
		},
	},
	// Тесты в помощь
	{
		name:  "Ошибка при nil интерфейсе",
		input: nil,
		check: func(writer io.Writer, err error) bool {
			return err != nil
		},
	},
}

var testError = errors.New("write error")

// Вспомогательная структура для эмуляции ошибки записи
type errorWriter struct{}

func (w *errorWriter) Write(p []byte) (n int, err error) {
	return 0, testError
}
