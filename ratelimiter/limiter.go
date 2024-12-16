package ratelimiter

import (
	"net/http"
	"sync"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/logger"
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

// ZDX GLOBAL RATE LIMIT MANAGEMENT

// GlobalRateLimiter enforces global rate limits across all requests.
type GlobalRateLimiter struct {
	mu          sync.Mutex
	allRequests []time.Time // Tracks timestamps of all requests made
	Limit       int         // Maximum number of requests allowed
	Freq        int         // Frequency window in seconds
}

// NewGlobalRateLimiter creates a new instance of a global rate limiter.
func NewGlobalRateLimiter(limit, freq int) *GlobalRateLimiter {
	return &GlobalRateLimiter{
		allRequests: []time.Time{},
		Limit:       limit,
		Freq:        freq,
	}
}

// Wait checks if the rate limit is exceeded and calculates the wait time.
func (rl *GlobalRateLimiter) Wait() (bool, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Remove outdated requests outside the time window
	windowStart := now.Add(-time.Duration(rl.Freq) * time.Second)
	filteredRequests := []time.Time{}
	for _, t := range rl.allRequests {
		if t.After(windowStart) {
			filteredRequests = append(filteredRequests, t)
		}
	}
	rl.allRequests = filteredRequests

	// Check if the rate limit is exceeded
	if len(rl.allRequests) >= rl.Limit {
		oldestRequest := rl.allRequests[0]
		delay := time.Duration(rl.Freq)*time.Second - now.Sub(oldestRequest)
		return true, delay
	}

	// Add the current request to the list
	rl.allRequests = append(rl.allRequests, now)
	return false, 0
}

// RateLimitTransport wraps the HTTP transport to apply rate limiting.
type RateLimitTransport struct {
	Base            http.RoundTripper
	Limiter         *GlobalRateLimiter
	WaitFunc        func() (bool, time.Duration) // Wait function reference
	Logger          logger.Logger
	AdditionalDelay time.Duration // Optional constant delay
}

// RoundTrip implements the http.RoundTripper interface for rate limiting.
func (rlt *RateLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Check rate limit
	shouldWait, delay := rlt.WaitFunc()
	if shouldWait {
		rlt.Logger.Printf("[INFO] Rate limit exceeded. Waiting for %v before making request.", delay)
		time.Sleep(delay + rlt.AdditionalDelay)
	}

	// Execute the actual HTTP request
	if rlt.Base == nil {
		rlt.Base = http.DefaultTransport
	}
	return rlt.Base.RoundTrip(req)
}
