package main

func GetMaxIndex(values []*int) int {
	var maxValue *int
	maxIndex := 0

	for index, value := range values {
		if value == nil {
			continue
		}

		if maxValue == nil || *maxValue < *value {
			maxValue = value
			maxIndex = index
		}
	}

	return maxIndex
}
