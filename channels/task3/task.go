//go:build task_template

package main

type Backend interface {
	DoRequest() string
}

func DoRequests(backends []Backend) []string {
	results := make([]string, 0, len(backends))
	for _, backend := range backends {
		results = append(results, backend.DoRequest())
	}

	return results
}
