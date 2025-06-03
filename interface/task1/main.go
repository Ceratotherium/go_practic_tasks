package main

import (
	"encoding/json"
	"reflect"
)

func main() {
	tests := append(testCases, privateTestCases...)

	for _, tt := range tests {
		CustomTestBody(
			tt.name,
			func() string {
				return tt.input
			},
			func(input string) bool {
				actual := Message{}
				err := json.Unmarshal([]byte(input), &actual)

				return err == nil && compare(tt.expected, actual)
			},
		)
	}
}

func compare(expected, actual Message) bool {
	for key, expectedValue := range expected {
		actualValue, ok := actual[key]
		if !ok {
			return false
		}

		expectedStruct := map[string]interface{}{}
		err := json.Unmarshal([]byte(expectedValue), &expectedStruct)

		if err != nil && expectedValue != actualValue { // Если значение не парсится как json (для обычных типов)
			return false
		} else { // Иначе считаем, что должна быть структура, и сравниваем структуры
			actualStruct := map[string]interface{}{}
			err = json.Unmarshal([]byte(actualValue), &actualStruct)
			if !reflect.DeepEqual(actualStruct, expectedStruct) {
				return false
			}
		}
	}
	return true
}
