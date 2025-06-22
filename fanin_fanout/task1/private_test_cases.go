package main

import (
	"sort"
	"time"
)

var privateTestCases = []TestCase{
	{
		name: "Проверка диапазона данных",
		in: func() []string {
			result := make([]string, 0, 30)
			for range 3 {
				for _, val := range []string{"4", "5", "6", "4", "5", "6", "4", "5", "6", "4"} {
					result = append(result, val)
				}
			}
			return result
		}(),
		convert:  convert,
		procNum:  3,
		sumCount: 3,
		check: func(in <-chan int) bool {
			for val := range in {
				if val < 12 || val > 18 {
					return false
				}
			}

			return true
		},
	},
	{
		name: "Проверка многопоточноcти",
		in: func() []string {
			result := make([]string, 0, 300)
			result = append(result, "0")
			for i := range 299 {
				if i%2 == 0 {
					result = append(result, "1")
				} else {
					result = append(result, "2")
				}
			}

			return result
		}(),
		convert: func(s string) int {
			if s == "0" {
				time.Sleep(time.Second * 3)
			}

			return convert(s)
		},
		procNum:  2,
		sumCount: 2,
		check: func(in <-chan int) bool {
			values := []int{}

			for v := range in {
				values = append(values, v)
			}

			if len(values) != 150 {
				return false
			}

			stableZone := values[10:20] // Кусок, где 100% вторая горутина была залочена
			sort.Ints(stableZone)

			return stableZone[0] == stableZone[len(stableZone)-1]
		},
	},
}
