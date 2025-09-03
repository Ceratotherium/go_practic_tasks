package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	backendError = errors.New("backendError")
	testDelay    = time.Millisecond * 5

	suspendAfterErrors = 1
	suspendFor         = time.Second
)

type TestCase struct {
	name  string
	check func() bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Один адрес один поток",
		check: func() bool {
			address := "addr1"
			b := makeBalancerForTest(address)

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
		name: "Опрос несколько адресов из нескольких горутин",
		check: func() bool {
			address := []string{"addr1", "addr2", "addr3"}
			requestCount := 9
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

			return len(responseByAddress) == len(address)
		},
	},
	// Тесткейсы в помощь
	{
		name: "Один адрес несколько потоков",
		check: func() bool {
			address := "addr1"
			requestCount := 9
			b := makeBalancerForTest(address)

			responses := runRequests(b, 3)

			responseCount := 0
			hasErr := false
			for resp := range responses {
				if resp.err != nil {
					hasErr = true
				} else {
					responseCount++
				}
			}

			if hasErr {
				return false
			}

			return responseCount != requestCount
		},
	},
	{
		name: "Один из бекендов возвращает ошибки",
		check: func() bool {
			address := []string{"addr1", "addr2", "addr3"}
			requestCount := 9
			responseByAddress := make(map[string]int, len(address))
			errorsCount := 0

			b := NewBalancer(
				address,
				suspendAfterErrors,
				suspendFor,
				func(s string) Backend {
					if s == "addr1" {
						return &MockBackend{
							addr:  s,
							err:   backendError,
							delay: testDelay,
						}
					}
					return NewMockBackend(s)
				},
			)

			// Прогоняем запросы по всем бекендам, чтобы выбить сбоящий
			responses := runRequests(b, requestCount)
			for range responses {
			}

			// Теперь сбоящего бекенда не должно быть в раздаче
			responses = runRequests(b, requestCount)

			for resp := range responses {
				if resp.err != nil {
					errorsCount++
				} else {
					_, actualAddr := flatResp(resp.response)
					responseByAddress[actualAddr]++
				}
			}

			if len(responseByAddress) != 2 {
				return false
			}

			if errorsCount != 0 {
				return false
			}

			return true
		},
	},
	{
		name: "Возвращение сбоящего бекенда через время",
		check: func() bool {
			address := []string{"addr1", "addr2", "addr3"}
			requestCount := 9
			responseByAddress := make(map[string]int, len(address))
			errorsCount := 0

			b := NewBalancer(
				address,
				suspendAfterErrors,
				time.Millisecond*50,
				func(s string) Backend {
					if s == "addr1" {
						return &MockBackend{
							addr:  s,
							err:   backendError,
							delay: testDelay,
						}
					}
					return NewMockBackend(s)
				},
			)

			// Прогоняем запросы по всем бекендам, чтобы выбить сбоящий
			responses := runRequests(b, requestCount)
			for range responses {
			}

			time.Sleep(time.Millisecond * 100)

			// Ждем таймаут, чтобы он вернулся в раздачу
			responses = runRequests(b, requestCount)

			for resp := range responses {
				if resp.err != nil {
					errorsCount++
				} else {
					_, actualAddr := flatResp(resp.response)
					responseByAddress[actualAddr]++
				}
			}

			if len(responseByAddress) != 2 {
				return false
			}

			if errorsCount == 0 {
				return false
			}

			return true
		},
	},
	{
		name: "Нет бекендов",
		check: func() bool {
			b := makeBalancerForTest()
			_, err := b.Invoke(context.Background(), "request")
			return err != nil
		},
	},
}

func makeBalancerForTest(addresses ...string) *Balancer {
	return NewBalancer(addresses, suspendAfterErrors, suspendFor, NewMockBackend)
}

type MockBackend struct {
	delay time.Duration
	addr  string
	err   error
}

func NewMockBackend(addr string) Backend {
	return &MockBackend{
		addr:  addr,
		delay: testDelay,
	}
}

func (b *MockBackend) Invoke(_ context.Context, r Request) (Response, error) {
	time.Sleep(b.delay)
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

type responseWithError struct {
	response Response
	err      error
}

func runRequests(balancer *Balancer, count int) chan responseWithError {
	ch := make(chan responseWithError, 1)

	wg := sync.WaitGroup{}

	for range count {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := balancer.Invoke(context.Background(), "request")
			ch <- responseWithError{
				response: resp,
				err:      err,
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}
