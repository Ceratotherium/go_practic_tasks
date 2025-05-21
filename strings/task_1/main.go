package main

func main() {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Палиндром с четным количеством символов",
			input:    "noon",
			expected: true,
		},
		{
			name:     "Палиндром с нечетным количеством символов",
			input:    "level",
			expected: true,
		},
		{
			name:     "Палиндром с заглавной буквой",
			input:    "Eve",
			expected: true,
		},
		{
			name:     "Фраза палиндром",
			input:    "A man, a plan, a canal: Panama",
			expected: true,
		},
		{
			name:     "Число палиндром",
			input:    "2002",
			expected: true,
		},
		{
			name:     "Один символ",
			input:    "a",
			expected: true,
		},
		{
			name:     "Пустая строка",
			input:    "",
			expected: true,
		},
		{
			name:     "Палиндром на руссом языке",
			input:    "Искать такси",
			expected: true,
		},
		{
			name:     "Палиндром с двумя одинаковыми символами",
			input:    "kk",
			expected: true,
		},
		{
			name:     "Обычная строка с четным количеством символов",
			input:    "love",
			expected: false,
		},
		{
			name:     "Обычная строка с нечетным количеством символов",
			input:    "day",
			expected: false,
		},
		{
			name:     "Обычная фраза",
			input:    "day",
			expected: false,
		},
	}

	for _, tt := range tests {
		AssertEqual(tt.expected, IsPalindrome(tt.input), tt.name, tt.input)
	}
}
