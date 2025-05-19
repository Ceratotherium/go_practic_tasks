package task_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	task "go_practic_tasks/strings/task_1"
)

func TestIsPalindrome(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected assert.BoolAssertionFunc
	}{
		{
			name:     "Палиндром с четным количеством символов",
			input:    "noon",
			expected: assert.True,
		},
		{
			name:     "Палиндром с нечетным количеством символов",
			input:    "level",
			expected: assert.True,
		},
		{
			name:     "Палиндром с заглавной буквой",
			input:    "Eve",
			expected: assert.True,
		},
		{
			name:     "Фраза палиндром",
			input:    "A man, a plan, a canal: Panama",
			expected: assert.True,
		},
		{
			name:     "Число палиндром",
			input:    "2002",
			expected: assert.True,
		},
		{
			name:     "Один символ",
			input:    "a",
			expected: assert.True,
		},
		{
			name:     "Пустая строка",
			input:    "",
			expected: assert.True,
		},
		{
			name:     "Палиндром на руссом языке",
			input:    "Искать такси",
			expected: assert.True,
		},
		{
			name:     "Палиндром с двумя одинаковыми символами",
			input:    "kk",
			expected: assert.True,
		},
		{
			name:     "Обычная строка с четным количеством символов",
			input:    "love",
			expected: assert.False,
		},
		{
			name:     "Обычная строка с нечетным количеством символов",
			input:    "day",
			expected: assert.False,
		},
		{
			name:     "Обычная фраза",
			input:    "day",
			expected: assert.False,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expected(t, task.IsPalindrome(tt.input), "Входная строка - %q", tt.input)
		})
	}
}
