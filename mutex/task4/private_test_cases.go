package main

import (
	"errors"
	"time"

	"golang.org/x/sync/errgroup"
)

var privateTestCases = []TestCase{
	{
		name: "Большое количество вызовов с разными аргументами",
		check: func() bool {
			errGroup := errgroup.Group{}

			for value := range 1000 {
				errGroup.Go(func() error {
					res, elapsed := measure(CachedLongCalculation, value)
					if elapsed < time.Millisecond*time.Duration(value) || // value мс
						elapsed > time.Millisecond*time.Duration(value+10) ||
						res != value+1 {
						return errors.New("incorrect result")
					}
					return nil
				})
			}

			return errGroup.Wait() == nil
		},
	},
	{
		name: "Большое количество вызовов с одним аргументом",
		check: func() bool {
			errGroup := errgroup.Group{}

			for range 1000 {
				errGroup.Go(func() (err error) {
					res, elapsed := measure(CachedLongCalculation, 10)
					if elapsed < time.Millisecond*10 || elapsed > time.Millisecond*20 || res != 11 { // 10мс
						return errors.New("incorrect result")
					}
					return nil
				})
			}

			return errGroup.Wait() == nil
		},
	},
}
