package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func ConcurrentlyRun(callback ...func()) error {
	wg := sync.WaitGroup{}
	wg.Add(len(callback))

	panicErr := atomic.Pointer[error]{}

	for _, fn := range callback {
		fn := fn
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						panicErr.Store(&err)
					} else {
						err := fmt.Errorf("%v", err)
						panicErr.Store(&err)
					}
				}
			}()

			fn()
		}()
	}

	wg.Wait()

	if err := panicErr.Load(); err != nil {
		return *err
	}
	return nil
}
