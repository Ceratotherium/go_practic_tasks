package main

var privateTestCases = []TestCase{
	{
		name:     "Обработка пустого объекта",
		input:    `{}`,
		expected: Message{},
	},
	{
		name:     "Обработка null",
		input:    `null`,
		expected: Message{},
	},
	{
		name: "Парсинг специальных символов в строках",
		input: `{
	       "text": "Line1\nLine2\tTab",
	       "quotes": "\"test\"",
	       "backslash": "\\path\\to\\file"
	   }`,
		expected: Message{
			"text":      "Line1\nLine2\tTab",
			"quotes":    `"test"`,
			"backslash": `\path\to\file`,
		},
	},
	{
		name: "Парсинг чисел в разных форматах",
		input: `{
            "int": 123,
            "negative": -456,
            "scientific": 1.23e4,
            "big": 12345678901234567890
        }`,
		expected: Message{
			"int":        "123",
			"negative":   "-456",
			"scientific": "12300",
			"big":        "12345678901234567890",
		},
	},
}
