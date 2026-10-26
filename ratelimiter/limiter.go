package ratelimiter

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/logger"
)

type RateLimiter struct {
	mu                    sync.Mutex
	getRequests           []time.Time
	postPutDeleteRequests []time.Time
	postPutRequests       []time.Time
	deleteRequests        []time.Time
	getLimit              int
	postPutDeleteLimit    int
	getFreq               int
	postPutDeleteFreq     int
	// Hourly limits
	getHourlyLimit     int
	postPutHourlyLimit int
	deleteHourlyLimit  int
	// perEndpointHourly applies the hourly limits to each endpoint separately
	// (method + path, with IDs normalized), as the ZIA API enforces them,
	// instead of to all requests of a method together.
	perEndpointHourly bool
	endpointRequests  map[string][]time.Time
}

func NewRateLimiter(getLimit, postPutDeleteLimit, getFreq, postPutDeleteFreq int) *RateLimiter {
	return &RateLimiter{
		getRequests:           []time.Time{},
		postPutDeleteRequests: []time.Time{},
		postPutRequests:       []time.Time{},
		deleteRequests:        []time.Time{},
		getLimit:              getLimit,
		postPutDeleteLimit:    postPutDeleteLimit,
		getFreq:               getFreq,
		postPutDeleteFreq:     postPutDeleteFreq,
		getHourlyLimit:        0, // disabled by default
		postPutHourlyLimit:    0, // disabled by default
		deleteHourlyLimit:     0, // disabled by default
	}
}

// NewRateLimiterWithHourly creates a rate limiter with both per-second and hourly limits
func NewRateLimiterWithHourly(getLimit, postPutDeleteLimit, getFreq, postPutDeleteFreq, getHourly, postPutHourly, deleteHourly int) *RateLimiter {
	return &RateLimiter{
		getRequests:           []time.Time{},
		postPutDeleteRequests: []time.Time{},
		postPutRequests:       []time.Time{},
		deleteRequests:        []time.Time{},
		getLimit:              getLimit,
		postPutDeleteLimit:    postPutDeleteLimit,
		getFreq:               getFreq,
		postPutDeleteFreq:     postPutDeleteFreq,
		getHourlyLimit:        getHourly,
		postPutHourlyLimit:    postPutHourly,
		deleteHourlyLimit:     deleteHourly,
	}
}

// NewRateLimiterWithEndpointHourly creates a rate limiter whose per-second
// limits apply to all requests of a method, and whose hourly limits apply to
// each endpoint separately: "PUT /zia/api/v1/firewallFilteringRules/123" and
// "PUT /zia/api/v1/firewallFilteringRules/456" share one hourly budget, while
// "POST /zia/api/v1/firewallFilteringRules" has its own.
func NewRateLimiterWithEndpointHourly(getLimit, postPutDeleteLimit, getFreq, postPutDeleteFreq, getHourly, postPutHourly, deleteHourly int) *RateLimiter {
	rl := NewRateLimiterWithHourly(getLimit, postPutDeleteLimit, getFreq, postPutDeleteFreq, getHourly, postPutHourly, deleteHourly)
	rl.perEndpointHourly = true
	rl.endpointRequests = map[string][]time.Time{}
	return rl
}

// WaitFor is Wait for a request to path. The path is used only by limiters
// created with NewRateLimiterWithEndpointHourly; other limiters behave as Wait.
func (rl *RateLimiter) WaitFor(method, path string) (bool, time.Duration) {
	if !rl.perEndpointHourly {
		return rl.Wait(method)
	}
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	rl.cleanOldRequests(now.Add(-1 * time.Hour))

	key := method + " " + normalizeEndpoint(path)
	hourly := rl.endpointRequests[key]
	if limit := rl.hourlyLimit(method); limit > 0 && len(hourly) >= limit {
		return true, time.Until(hourly[0].Add(time.Hour))
	}
	if wait, d := rl.takePerSecond(method, now); wait {
		return true, d
	}
	rl.endpointRequests[key] = append(hourly, now)
	return false, 0
}

func (rl *RateLimiter) hourlyLimit(method string) int {
	switch method {
	case http.MethodGet:
		return rl.getHourlyLimit
	case http.MethodDelete:
		return rl.deleteHourlyLimit
	default:
		return rl.postPutHourlyLimit
	}
}

// takePerSecond applies the per-second limits and records the request when it
// may proceed. The caller holds rl.mu.
func (rl *RateLimiter) takePerSecond(method string, now time.Time) (bool, time.Duration) {
	if method == http.MethodGet {
		if len(rl.getRequests) >= rl.getLimit {
			oldest := rl.getRequests[0]
			if window := time.Duration(rl.getFreq) * time.Second; now.Sub(oldest) < window {
				return true, window - now.Sub(oldest)
			}
			rl.getRequests = rl.getRequests[1:]
		}
		rl.getRequests = append(rl.getRequests, now)
		return false, 0
	}
	if len(rl.postPutDeleteRequests) >= rl.postPutDeleteLimit {
		oldest := rl.postPutDeleteRequests[0]
		if window := time.Duration(rl.postPutDeleteFreq) * time.Second; now.Sub(oldest) < window {
			return true, window - now.Sub(oldest)
		}
		rl.postPutDeleteRequests = rl.postPutDeleteRequests[1:]
	}
	rl.postPutDeleteRequests = append(rl.postPutDeleteRequests, now)
	return false, 0
}

// normalizeEndpoint replaces numeric path segments (resource IDs) with "{id}",
// so every request to the same API resource shares one hourly budget.
func normalizeEndpoint(path string) string {
	segments := strings.Split(path, "/")
	for i, seg := range segments {
		if seg == "" {
			continue
		}
		numeric := true
		for _, c := range seg {
			if c < '0' || c > '9' {
				numeric = false
				break
			}
		}
		if numeric {
			segments[i] = "{id}"
		}
	}
	return strings.Join(segments, "/")
}

func (rl *RateLimiter) Wait(method string) (bool, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)

	// Clean up old requests older than 1 hour for hourly tracking
	rl.cleanOldRequests(oneHourAgo)

	switch method {
	case http.MethodGet:
		// Check hourly limit first (if enabled)
		if rl.getHourlyLimit > 0 {
			if len(rl.getRequests) >= rl.getHourlyLimit {
				oldestRequest := rl.getRequests[0]
				if oldestRequest.After(oneHourAgo) {
					// We've hit the hourly limit
					d := time.Until(oldestRequest.Add(time.Hour))
					return true, d
				}
			}
		}

		// Check per-second limit
		if len(rl.getRequests) >= rl.getLimit {
			oldestRequest := rl.getRequests[0]
			if now.Sub(oldestRequest) < time.Duration(rl.getFreq)*time.Second {
				d := time.Duration(rl.getFreq)*time.Second - now.Sub(oldestRequest)
				return true, d
			}
			rl.getRequests = rl.getRequests[1:]
		}
		rl.getRequests = append(rl.getRequests, now)

	case http.MethodPost, http.MethodPut:
		// Check hourly limit first (if enabled)
		if rl.postPutHourlyLimit > 0 {
			if len(rl.postPutRequests) >= rl.postPutHourlyLimit {
				oldestRequest := rl.postPutRequests[0]
				if oldestRequest.After(oneHourAgo) {
					// We've hit the hourly limit
					d := time.Until(oldestRequest.Add(time.Hour))
					return true, d
				}
			}
		}

		// Check per-second limit
		if len(rl.postPutDeleteRequests) >= rl.postPutDeleteLimit {
			oldestRequest := rl.postPutDeleteRequests[0]
			if now.Sub(oldestRequest) < time.Duration(rl.postPutDeleteFreq)*time.Second {
				d := time.Duration(rl.postPutDeleteFreq)*time.Second - now.Sub(oldestRequest)
				return true, d
			}
			rl.postPutDeleteRequests = rl.postPutDeleteRequests[1:]
		}
		rl.postPutDeleteRequests = append(rl.postPutDeleteRequests, now)
		rl.postPutRequests = append(rl.postPutRequests, now)

	case http.MethodDelete:
		// Check hourly limit first (if enabled)
		if rl.deleteHourlyLimit > 0 {
			if len(rl.deleteRequests) >= rl.deleteHourlyLimit {
				oldestRequest := rl.deleteRequests[0]
				if oldestRequest.After(oneHourAgo) {
					// We've hit the hourly limit
					d := time.Until(oldestRequest.Add(time.Hour))
					return true, d
				}
			}
		}

		// Check per-second limit
		if len(rl.postPutDeleteRequests) >= rl.postPutDeleteLimit {
			oldestRequest := rl.postPutDeleteRequests[0]
			if now.Sub(oldestRequest) < time.Duration(rl.postPutDeleteFreq)*time.Second {
				d := time.Duration(rl.postPutDeleteFreq)*time.Second - now.Sub(oldestRequest)
				return true, d
			}
			rl.postPutDeleteRequests = rl.postPutDeleteRequests[1:]
		}
		rl.postPutDeleteRequests = append(rl.postPutDeleteRequests, now)
		rl.deleteRequests = append(rl.deleteRequests, now)
	}

	return false, 0
}

// cleanOldRequests removes requests older than the specified cutoff time to prevent memory leaks
func (rl *RateLimiter) cleanOldRequests(cutoff time.Time) {
	// Clean GET requests
	cleaned := []time.Time{}
	for _, t := range rl.getRequests {
		if t.After(cutoff) {
			cleaned = append(cleaned, t)
		}
	}
	rl.getRequests = cleaned

	// Clean POST/PUT requests
	cleaned = []time.Time{}
	for _, t := range rl.postPutRequests {
		if t.After(cutoff) {
			cleaned = append(cleaned, t)
		}
	}
	rl.postPutRequests = cleaned

	// Clean DELETE requests
	cleaned = []time.Time{}
	for _, t := range rl.deleteRequests {
		if t.After(cutoff) {
			cleaned = append(cleaned, t)
		}
	}
	rl.deleteRequests = cleaned

	// Clean POST/PUT/DELETE combined requests
	cleaned = []time.Time{}
	for _, t := range rl.postPutDeleteRequests {
		if t.After(cutoff) {
			cleaned = append(cleaned, t)
		}
	}
	rl.postPutDeleteRequests = cleaned

	// Clean per-endpoint hourly requests
	for key, times := range rl.endpointRequests {
		kept := times[:0]
		for _, t := range times {
			if t.After(cutoff) {
				kept = append(kept, t)
			}
		}
		if len(kept) == 0 {
			delete(rl.endpointRequests, key)
		} else {
			rl.endpointRequests[key] = kept
		}
	}
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

// sleepCtx waits for d, or until the request's context is done.
func sleepCtx(req *http.Request, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-req.Context().Done():
		return req.Context().Err()
	}
}

// RateLimitTransport wraps the HTTP transport to apply rate limiting.
type RateLimitTransport struct {
	Base            http.RoundTripper
	Limiter         *RateLimiter                 // Standard rate limiter
	GlobalLimiter   *GlobalRateLimiter           // For ZDX global limiting
	WaitFunc        func() (bool, time.Duration) // Wait function reference (optional, overrides Limiter)
	Logger          logger.Logger
	AdditionalDelay time.Duration // Optional constant delay
}

// RoundTrip implements the http.RoundTripper interface for rate limiting.
func (rlt *RateLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Determine which rate limiter to use
	if rlt.WaitFunc != nil {
		// Use custom wait function if provided
		if shouldWait, delay := rlt.WaitFunc(); shouldWait {
			rlt.Logger.Printf("[INFO] Rate limit exceeded for %s request. Waiting for %v before proceeding.", req.Method, delay)
			time.Sleep(delay + rlt.AdditionalDelay)
		}
	} else if rlt.Limiter != nil || rlt.GlobalLimiter != nil {
		// Wait until the limiter admits (and records) the request: a request
		// sent after a wait without being recorded would not count against
		// the limits.
		for {
			var shouldWait bool
			var delay time.Duration
			if rlt.Limiter != nil {
				shouldWait, delay = rlt.Limiter.WaitFor(req.Method, req.URL.Path)
			} else {
				shouldWait, delay = rlt.GlobalLimiter.Wait()
			}
			if !shouldWait {
				break
			}
			rlt.Logger.Printf("[INFO] Rate limit exceeded for %s %s. Waiting for %v before proceeding.", req.Method, req.URL.Path, delay)
			if err := sleepCtx(req, delay+rlt.AdditionalDelay); err != nil {
				return nil, err
			}
		}
	}

	// Execute the actual HTTP request
	if rlt.Base == nil {
		rlt.Base = http.DefaultTransport
	}
	return rlt.Base.RoundTrip(req)
}
