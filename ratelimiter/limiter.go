package ratelimiter

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu                    sync.Mutex
	getRequests           []time.Time
	postPutDeleteRequests []time.Time
	getLimit              int
	postPutDeleteLimit    int
	getFreq               int
	postPutDeleteFreq     int
}

func NewRateLimiter(getLimit, postPutDeleteLimit, getFreq, postPutDeleteFreq int) *RateLimiter {
	return &RateLimiter{
		getRequests:           []time.Time{},
		postPutDeleteRequests: []time.Time{},
		getLimit:              getLimit,
		postPutDeleteLimit:    postPutDeleteLimit,
		getFreq:               getFreq,
		postPutDeleteFreq:     postPutDeleteFreq,
	}
}

func (rl *RateLimiter) Wait(method string) (bool, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	switch method {
	case http.MethodGet:
		if len(rl.getRequests) >= rl.getLimit {
			oldestRequest := rl.getRequests[0]
			if now.Sub(oldestRequest) < time.Duration(rl.getFreq)*time.Second {
				d := time.Duration(rl.getFreq)*time.Second - now.Sub(oldestRequest)
				return true, d
			}
			rl.getRequests = rl.getRequests[1:]
		}
		rl.getRequests = append(rl.getRequests, now)

	case http.MethodPost, http.MethodPut, http.MethodDelete:
		if len(rl.postPutDeleteRequests) >= rl.postPutDeleteLimit {
			oldestRequest := rl.postPutDeleteRequests[0]
			if now.Sub(oldestRequest) < time.Duration(rl.postPutDeleteFreq)*time.Second {
				d := time.Duration(rl.postPutDeleteFreq)*time.Second - now.Sub(oldestRequest)
				return true, d
			}
			rl.postPutDeleteRequests = rl.postPutDeleteRequests[1:]
		}
		rl.postPutDeleteRequests = append(rl.postPutDeleteRequests, now)
	}

	return false, 0
}
