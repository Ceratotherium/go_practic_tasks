package main

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"golang.org/x/sync/errgroup"
)

var backendError = errors.New("backendError")

type TestCase struct {
	name  string
	check func() bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Один адрес",
		check: func() bool {
			address := "addr1"
			b := NewBalancer([]string{
				address,
			},
				NewMockBackend,
			)

			for i := range 3 {
				request := fmt.Sprintf("request_%d", i)
				resp, err := b.Invoke(context.Background(), request)
				if err != nil {
					return false
				}

				actualReq, actualAddr := flatResp(resp)

				if actualReq != request || actualAddr != address {
					return false
				}
			}
			return true
		},
	},
	{
		name: "Несколько адресов",
		check: func() bool {
			address := []string{"addr1", "addr2", "addr3"}
			requestCount := 9
			responseByAddress := make(map[string][]string, len(address))

			b := NewBalancer(address, NewMockBackend)

			for i := range requestCount {
				resp, err := b.Invoke(context.Background(), fmt.Sprintf("request_%d", i))
				if err != nil {
					return false
				}

				actualReq, actualAddr := flatResp(resp)
				responseByAddress[actualAddr] = append(responseByAddress[actualAddr], actualReq)
			}

			if len(responseByAddress) != len(address) {
				return false
			}

			for _, requests := range responseByAddress {
				if len(requests) != requestCount/len(address) {
					return false
				}

				if slices.Contains(requests, "request_0") &&
					slices.Contains(requests, "request_3") &&
					slices.Contains(requests, "request_6") ||
					slices.Contains(requests, "request_1") &&
						slices.Contains(requests, "request_4") &&
						slices.Contains(requests, "request_7") ||
					slices.Contains(requests, "request_2") &&
						slices.Contains(requests, "request_5") &&
						slices.Contains(requests, "request_8") {
					continue
				}

				return false
			}

			return true
		},
	},
	// Тесткейсы в помощь
	{
		name: "Опрос несколько адресов из нескольких горутин",
		check: func() bool {
			address := []string{"addr1", "addr2", "addr3"}
			requestCount := 9
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
		name: "Один из бекендов возвращает ошибки",
		check: func() bool {
			address := []string{"addr1", "addr2", "addr3"}
			requestCount := 9
			responseByAddress := make(map[string]int, len(address))
			errorsCount := 0

			b := NewBalancer(address, func(s string) Backend {
				if s == "addr1" {
					return &MockBackend{
						addr: s,
						err:  backendError,
					}
				}
				return NewMockBackend(s)
			})

			for range requestCount {
				resp, err := b.Invoke(context.Background(), "request")
				if err != nil {
					errorsCount++
				} else {
					_, actualAddr := flatResp(resp)
					responseByAddress[actualAddr]++
				}
			}

			if len(responseByAddress) != 2 {
				return false
			}

			if errorsCount != requestCount/len(address) {
				return false
			}

			return true
		},
	},
	{
		name: "Нет бекендов",
		check: func() bool {
			b := NewBalancer(nil, NewMockBackend)
			_, err := b.Invoke(context.Background(), "request")
			return err != nil
		},
	},
}

type MockBackend struct {
	addr string
	err  error
}

func NewMockBackend(addr string) Backend {
	return &MockBackend{
		addr: addr,
	}
}

func (b *MockBackend) Invoke(_ context.Context, r Request) (Response, error) {
	return mockResp{
		request: r.(string),
		addr:    b.addr,
	}, b.err
}

type mockResp struct {
	request string
	addr    string
}

func flatResp(resp Response) (string, string) {
	typedResp, ok := resp.(mockResp)
	if !ok {
		return "", ""
	}

	return typedResp.request, typedResp.addr
}
