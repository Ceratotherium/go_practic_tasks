package main

func RotateLeft(values [5]int, rotateFactor int) [5]int {
	rotateFactor *= -1

	if rotateFactor > 5 {
		rotateFactor = rotateFactor % 5
	}

	if rotateFactor < 0 {
		rotateFactor = (rotateFactor % 5) + 5
	}

	result := [5]int{}
	for index, value := range values {
		result[(index+rotateFactor)%5] = value
	}

	return result
}
