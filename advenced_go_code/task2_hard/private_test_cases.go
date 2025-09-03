package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

var privateTestCases = []TestCase{
	{
		name: "Большое количество адресов и запросов",
		check: func() bool {
			address := make([]string, 0, 100)
			for i := range 100 {
				address = append(address, fmt.Sprintf("addr%d", i))
			}

			requestCount := 1000
			responseByAddress := make(map[string]int, len(address))
			b := makeBalancerForTest(address...)

			responses := runRequests(b, requestCount)

			hasErr := false
			for resp := range responses {
				hasErr = hasErr || resp.err != nil
				_, actualAddr := flatResp(resp.response)
				responseByAddress[actualAddr]++
			}

			if hasErr {
				return false
			}

			if len(responseByAddress) != len(address) {
				return false
			}

			for _, count := range responseByAddress {
				if count != requestCount/len(address) {
					return false
				}
			}

			return true
		},
	},
	{
		name: "Один из бекендов тормозит",
		check: func() bool {
			addr1 := "addr1"
			address := []string{addr1, "addr2", "addr3"}
			requestCount := 9
			responseByAddress := make(map[string]int, len(address))

			b := NewBalancer(
				address,
				suspendAfterErrors,
				suspendFor,
				func(s string) Backend {
					if s == addr1 {
						return &MockBackend{
							addr:  s,
							delay: time.Millisecond * 200,
						}
					}
					return NewMockBackend(s)
				},
			)

			started := time.Now()
			ch := make(chan Response, 1)

			errGroup := errgroup.Group{}
			errGroup.SetLimit(3)

			go func() {
				for range requestCount {
					errGroup.Go(func() error {
						resp, err := b.Invoke(context.Background(), "request")
						ch <- resp
						return err
					})
				}

				_ = errGroup.Wait()
				close(ch)
			}()

			for resp := range ch {
				_, actualAddr := flatResp(resp)
				responseByAddress[actualAddr]++
			}

			elapsed := time.Since(started)
			if elapsed >= time.Millisecond*200*3 { // 3 запроса должно было попасть на медленный бекенд
				return false
			}

			if len(responseByAddress) != len(address) {
				return false
			}

			if responseByAddress[addr1] != 1 {
				return false
			}

			return true
		},
	},
}
