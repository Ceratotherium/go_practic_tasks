package main

func Unique(values []int) []int {
	results := make([]int, 0, len(values))
	existValues := make(map[int]struct{}, len(values))

	for _, v := range values {
		if _, ok := existValues[v]; !ok {
			existValues[v] = struct{}{}
			results = append(results, v)
		}
	}

	return results
}
