package main

type TestCase struct {
	name     string
	input    string
	expected bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name:     "Палиндром",
		input:    "noon",
		expected: true,
	},
	{
		name:     "Не палиндром",
		input:    "love",
		expected: false,
	},
	// Тесткейсы в помощь
	{
		name:     "Палиндром с нечетным количеством символов",
		input:    "level",
		expected: true,
	},
	{
		name:     "Палиндром с двумя одинаковыми символами",
		input:    "kk",
		expected: true,
	},
	{
		name:     "Не палиндром с нечетным количеством символов",
		input:    "day",
		expected: false,
	},
	{
		name:     "Обычная фраза",
		input:    "a long day",
		expected: false,
	},
}
