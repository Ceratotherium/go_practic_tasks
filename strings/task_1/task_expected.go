package task

import (
	"unicode"
)

func IsPalindrome(str string) bool {
	runeString := []rune(str)

	for i, j := 0, len(runeString)-1; i < j; {
		if !unicode.IsLetter(runeString[i]) && !unicode.IsNumber(runeString[i]) {
			i++
			continue
		}
		if !unicode.IsLetter(runeString[j]) && !unicode.IsNumber(runeString[j]) {
			j--
			continue
		}

		if unicode.ToLower(runeString[i]) != unicode.ToLower(runeString[j]) {
			return false
		}
		i++
		j--
	}
	return true
}
