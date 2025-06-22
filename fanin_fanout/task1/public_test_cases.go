package main

import (
	"strconv"
)

type TestCase struct {
	name     string
	in       []string
	convert  func(string) int
	procNum  int
	sumCount int
	check    func(<-chan int) bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name:     "Один поток перебор чисел",
		in:       []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
		convert:  convert,
		procNum:  1,
		sumCount: 1,
		check:    expectedNumbers(0, 1, 2, 3, 4, 5, 6, 7, 8, 9),
	},
	{
		name:     "Один поток, сумма 2х",
		in:       []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
		convert:  convert,
		procNum:  1,
		sumCount: 2,
		check:    expectedNumbers(1, 5, 9, 13, 17),
	},
	// Тесткейсы в помощь
	{
		name:     "Количество данных не делится на sumCount",
		in:       []string{"0", "1", "2"},
		convert:  convert,
		procNum:  1,
		sumCount: 2,
		check:    expectedNumbers(1, 2),
	},
	{
		name:     "Отрицательные числа",
		in:       []string{"0", "1", "2", "-1", "-2", "1"},
		convert:  convert,
		procNum:  1,
		sumCount: 2,
		check:    expectedNumbers(1, 1, -1),
	},
	{
		name:     "Проверка общей суммы",
		in:       []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
		convert:  convert,
		procNum:  2,
		sumCount: 2,
		check: func(in <-chan int) bool {
			sum := 0
			for val := range in {
				sum += val
			}

			return sum == 45
		},
	},
}

func convert(str string) int {
	val, _ := strconv.Atoi(str)
	return val
}

func expectedNumbers(nums ...int) func(<-chan int) bool {
	return func(in <-chan int) bool {
		for _, expected := range nums {
			val, ok := <-in
			if !ok || val != expected {
				return false
			}
		}

		// Проверяю, что в канале больше ничего нет
		select {
		case _, ok := <-in:
			return !ok
		default:
			return true
		}
	}
}
