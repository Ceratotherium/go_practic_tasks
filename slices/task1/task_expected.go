package main

func Chunk(values []int, chunkSize int) [][]int {
	if chunkSize < 1 || len(values) == 0 {
		return [][]int{}
	}

	var chunks [][]int

	for chunkSize < len(values) {
		values, chunks = values[chunkSize:], append(chunks, values[0:chunkSize:chunkSize])
	}
	return append(chunks, values)
}
