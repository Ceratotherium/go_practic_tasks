package main

import (
	"sync"
)

type TestCase struct {
	name    string
	prepare []string
	check   func(*sync.Map) bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name:    "По одному слову в строке",
		prepare: getStrings("one", "two", "three", "four", "five", "six"),
		check: func(words *sync.Map) bool {
			counter := counter{words}
			return counter.Count("one") == 1 && counter.Count("three") == 1 && counter.Count("six") == 1
		},
	},
	{
		name:    "Несколько слов в строке",
		prepare: getStrings("one two", "three four", "five six"),
		check: func(words *sync.Map) bool {
			counter := counter{words}
			return counter.Count("one") == 1 && counter.Count("three") == 1 && counter.Count("six") == 1
		},
	},
	{
		name:    "Одинаковые слова в одной строке",
		prepare: getStrings("one one", "three three six", "six six"),
		check: func(words *sync.Map) bool {
			counter := counter{words}
			return counter.Count("one") == 2 && counter.Count("three") == 2 && counter.Count("six") == 3
		},
	},
	// Тесткейсы в помощь
	{
		name:    "Слова в разном регистре",
		prepare: getStrings("one two", "One tWo", "onE TWO"),
		check: func(words *sync.Map) bool {
			counter := counter{words}
			return counter.Count("one") == 3 && counter.Count("two") == 3
		},
	},
	{
		name:    "Пустые строки",
		prepare: getStrings("", "", ""),
		check: func(words *sync.Map) bool {
			isEmpty := true
			words.Range(func(_, _ any) bool {
				isEmpty = false
				return true
			})
			return isEmpty
		},
	},
	{
		name:    "Есть все слова",
		prepare: getStrings("one", "two", "three", "four", "five", "six"),
		check: func(words *sync.Map) bool {
			existWords := map[string]bool{
				"one":   false,
				"two":   false,
				"three": false,
				"four":  false,
				"five":  false,
				"six":   false,
			}

			words.Range(func(word, _ any) bool {
				existWords[word.(string)] = true
				return true
			})

			for _, ok := range existWords {
				if !ok {
					return false
				}
			}
			return true
		},
	},
	{
		name:    "Нет лишних слов",
		prepare: getStrings("one", "two", "three", "four", "five", "six"),
		check: func(words *sync.Map) bool {
			existWords := map[string]struct{}{
				"one":   {},
				"two":   {},
				"three": {},
				"four":  {},
				"five":  {},
				"six":   {},
			}

			hasAnotherWords := false

			words.Range(func(word, _ any) bool {
				if _, ok := existWords[word.(string)]; !ok {
					hasAnotherWords = true
				}
				return true
			})

			return !hasAnotherWords
		},
	},
	{
		name: "Большой текст одной строкой",
		prepare: getStrings("Linda wants to buy a new car",
			"She has an old car",
			"Her old car is a white Honda",
			"Linda wants to buy a new Honda",
			"She wants to buy a new red Honda",
			"She has saved $1,000",
			"She will use $1,000 to help buy the new car",
			"She will give $1,000 to the Honda dealer",
			"The Honda dealer will give her a contract to sign",
			"The contract will require her to pay $400 a month for seven years",
			"Her new red Honda will cost Linda a lot of money",
			"But that's okay, because Linda makes a lot of money"),
		check: func(words *sync.Map) bool {
			counter := counter{words}

			return counter.Count("car") == 4 && counter.Count("honda") == 6 && counter.Count("her") == 4
		},
	},
	{
		name: "Большой текст несколькими строками",
		prepare: getStrings("Linda wants to buy a new car " +
			"She has an old car " +
			"Her old car is a white Honda " +
			"Linda wants to buy a new Honda " +
			"She wants to buy a new red Honda " +
			"She has saved $1,000 " +
			"She will use $1,000 to help buy the new car " +
			"She will give $1,000 to the Honda dealer " +
			"The Honda dealer will give her a contract to sign " +
			"The contract will require her to pay $400 a month for seven years " +
			"Her new red Honda will cost Linda a lot of money " +
			"But that's okay, because Linda makes a lot of money"),
		check: func(words *sync.Map) bool {
			counter := counter{words}

			return counter.Count("car") == 4 && counter.Count("honda") == 6 && counter.Count("her") == 4
		},
	},
}

func testBody(testName string, lines []string) *sync.Map {
	words := &sync.Map{}
	wg := sync.WaitGroup{}
	wg.Add(len(lines))

	for _, line := range lines {
		go func() {
			defer catchPanic(testName)()
			defer wg.Done()

			WordCount(line, words)
		}()
	}

	wg.Wait()
	return words
}

func getStrings(str ...string) []string {
	return str
}

type counter struct {
	*sync.Map
}

func (s *counter) Count(key string) int {
	val, ok := s.Load(key)
	if !ok {
		return 0
	}
	return val.(int)
}
