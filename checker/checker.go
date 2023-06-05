package checker

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type SiteCheck struct {
	sites    map[string]float64
	requests map[string]int
	rwm      sync.RWMutex
}

type Checker interface {
	CheckWebsite(contxt context.Context, webs []string)
	GetAccessTime(url string) (map[string]float64, bool)
	GetMaxAccessTime() map[string]float64
	GetMinAccessTime() map[string]float64
	IncCounter(endpoint string)
	GetCounts() map[string]int
}

func New() Checker {
	req := map[string]int{
		"/access": 0,
		"/min":    0,
		"/max":    0,
	}
	scheck := SiteCheck{
		sites:    make(map[string]float64),
		requests: req,
	}

	return &scheck
}

func (sc *SiteCheck) CheckWebsite(contxt context.Context, urls []string) {

	for i := 0; i < len(urls); i++ {
		ii := i
		go func() {

			start := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			req, _ := http.NewRequestWithContext(ctx, "GET", urls[ii], nil)

			client := http.DefaultClient
			resp, err := client.Do(req)

			if err == nil && resp.StatusCode == 200 {
				resp.Body.Close()
			} else {
				return
			}

			elapsed := time.Since(start).Seconds()

			sc.rwm.Lock()
			sc.sites[urls[ii]] = elapsed
			sc.rwm.Unlock()
		}()
	}
}

func (sc *SiteCheck) GetAccessTime(url string) (map[string]float64, bool) {

	sc.rwm.RLock()
	defer sc.rwm.RUnlock()
	v, ok := sc.sites[url]

	res := map[string]float64{
		url: v,
	}

	return res, ok
}

func (sc *SiteCheck) GetMinAccessTime() map[string]float64 {

	var minURL string
	var minTime float64 = 24 * 7
	res := make(map[string]float64, 0)

	sc.rwm.RLock()
	defer sc.rwm.RUnlock()
	for url, val := range sc.sites {
		if val < minTime {
			minTime = val
			minURL = url
		}
	}

	res[minURL] = minTime

	return res
}

func (sc *SiteCheck) GetMaxAccessTime() map[string]float64 {

	var maxURL string
	var maxTime float64 = 0
	res := make(map[string]float64)

	sc.rwm.Lock()
	defer sc.rwm.Unlock()
	for url, val := range sc.sites {
		if val > maxTime {
			maxTime = val
			maxURL = url
		}
	}

	res[maxURL] = maxTime

	return res
}

func (sc *SiteCheck) IncCounter(endpoint string) {
	sc.rwm.Lock()
	defer sc.rwm.Unlock()
	sc.requests[endpoint]++
}

func (sc *SiteCheck) GetCounts() map[string]int {
	sc.rwm.RLock()
	defer sc.rwm.RUnlock()

	counts := make(map[string]int)
	for endpoint, count := range sc.requests {
		counts[endpoint] = count
	}

	return counts
}
