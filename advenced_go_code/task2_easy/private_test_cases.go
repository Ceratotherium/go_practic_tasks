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

			b := NewBalancer(address, NewMockBackend)

			errGroup := errgroup.Group{}
			ch := make(chan Response)

			for range requestCount {
				errGroup.Go(func() error {
					resp, err := b.Invoke(context.Background(), "request")
					ch <- resp
					return err
				})
			}

			for range requestCount {
				resp := <-ch
				_, actualAddr := flatResp(resp)
				responseByAddress[actualAddr]++
			}

			if err := errGroup.Wait(); err != nil {
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
			address := []string{"addr1", "addr2", "addr3"}
			requestCount := 9
			responseByAddress := make(map[string]int, len(address))

			b := NewBalancer(address, func(s string) Backend {
				b := NewMockBackend(s)
				if s == "addr1" {
					return slowBackend{
						Backend: b,
						delay:   time.Millisecond * 200,
					}
				}
				return b
			})

			started := time.Now()

			for range requestCount {
				resp, err := b.Invoke(context.Background(), "request")
				if err != nil {
					return false
				}

				_, actualAddr := flatResp(resp)
				responseByAddress[actualAddr]++
			}

			elapsed := time.Since(started)
			if elapsed < time.Millisecond*200*3 { // 3 запроса должно было попасть на медленный бекенд
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
}

type slowBackend struct {
	Backend
	delay time.Duration
}

func (b slowBackend) Invoke(ctx context.Context, r Request) (Response, error) {
	time.Sleep(b.delay)
	return b.Backend.Invoke(ctx, r)
}
