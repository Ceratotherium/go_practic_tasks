package main

import "sync"

var privateTestCases = []TestCase{
	{
		name: "Конкурентная инициализация (много горутин)",
		check: func() bool {
			dbs := make([]*mockDatabase, 1000)

			wg := sync.WaitGroup{}
			wg.Add(len(dbs))
			for i := range dbs {
				go func() {
					defer wg.Done()
					dbs[i] = castToMock(GetDatabase())
				}()
			}

			wg.Wait()

			for _, db := range dbs[1:] {
				if dbs[0] != db {
					return false
				}
			}
			return dbs[0].ConnectCount() == 1
		},
	}}
