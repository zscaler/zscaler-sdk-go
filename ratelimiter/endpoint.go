package ratelimiter

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

// Rate is a number of requests allowed per window, such as 1 per second or
// 10 per minute.
type Rate struct {
	Requests int
	Per      time.Duration
}

// EndpointLimit is the limit of one method on one API resource: a short-window
// rate (per second or per minute) and a number of requests per hour. Hourly 0
// means no hourly limit.
type EndpointLimit struct {
	Short  Rate
	Hourly int
}

// EndpointRule sets the limit of one method on a resource URI. Path is the
// resource URI as the API rate limit summary lists it, without the API prefix
// (for example "/timeIntervals/{id}"). A segment in braces matches any single
// path segment; other segments match case-insensitively.
type EndpointRule struct {
	Method string
	Path   string
	Limit  EndpointLimit
}

// EndpointRateLimiter applies rate limits per method and resource URI, as the
// ZIA API enforces them: every method on every resource has its own short
// window and its own hourly budget. Requests that match a rule use the rule's
// limit; all others use the default limit of their method.
//
// The short window is stretched by ShortWindowMargin, and the hourly budget is
// reduced by about 5% (at least 1 request, never below 1), so requests stay
// clear of the limits the API enforces.
type EndpointRateLimiter struct {
	mu       sync.Mutex
	prefix   string
	defaults map[string]EndpointLimit
	rules    []compiledRule
	buckets  map[string]*endpointBucket
	now      func() time.Time
}

// ShortWindowMargin stretches each short window, so requests paced at exactly
// the limit do not land in the same server-side window.
const ShortWindowMargin = 1.1

type compiledRule struct {
	method   string
	pattern  string
	segments []string // lower-cased; "" for a placeholder
	literals int
	limit    EndpointLimit
}

type endpointBucket struct {
	limit  EndpointLimit
	short  []time.Time
	hourly []time.Time
}

// Ticket identifies a request admitted by Acquire, so Release can take it back
// out of the hourly budget.
type Ticket struct {
	key string
	at  time.Time
}

// NewEndpointRateLimiter returns a limiter for request paths that start with
// prefix (for example "/zia/api/v1"). defaults holds the limit of each method
// that no rule matches; a method missing from defaults uses the PUT default.
func NewEndpointRateLimiter(prefix string, defaults map[string]EndpointLimit, rules []EndpointRule) *EndpointRateLimiter {
	rl := &EndpointRateLimiter{
		prefix:   strings.TrimSuffix(prefix, "/"),
		defaults: defaults,
		buckets:  map[string]*endpointBucket{},
		now:      time.Now,
	}
	for _, r := range rules {
		cr := compiledRule{method: strings.ToUpper(r.Method), pattern: r.Path, limit: r.Limit}
		for _, seg := range splitPath(r.Path) {
			if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
				cr.segments = append(cr.segments, "")
				continue
			}
			cr.segments = append(cr.segments, strings.ToLower(seg))
			cr.literals++
		}
		rl.rules = append(rl.rules, cr)
	}
	return rl
}

func splitPath(p string) []string {
	var out []string
	for _, s := range strings.Split(p, "/") {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

// resolve returns the bucket key and limit of a request.
func (rl *EndpointRateLimiter) resolve(method, path string) (string, EndpointLimit) {
	method = strings.ToUpper(method)
	rel := path
	if rl.prefix != "" && strings.HasPrefix(strings.ToLower(path), strings.ToLower(rl.prefix)) {
		rel = path[len(rl.prefix):]
	}
	segs := splitPath(rel)

	var best *compiledRule
	for i := range rl.rules {
		r := &rl.rules[i]
		if r.method != method || len(r.segments) != len(segs) {
			continue
		}
		ok := true
		for j, want := range r.segments {
			if want != "" && want != strings.ToLower(segs[j]) {
				ok = false
				break
			}
		}
		// The most specific rule wins: "/extranet/lite" over "/extranet/{id}".
		if ok && (best == nil || r.literals > best.literals) {
			best = r
		}
	}
	if best != nil {
		return method + " " + best.pattern, best.limit
	}

	limit, ok := rl.defaults[method]
	if !ok {
		limit = rl.defaults[http.MethodPut]
	}
	return method + " " + normalizeEndpoint("/"+strings.Join(segs, "/")), limit
}

// hourlyBudget is the hourly limit less a safety margin of about 5%.
func hourlyBudget(limit int) int {
	if limit <= 0 {
		return 0
	}
	margin := limit / 20
	if margin < 1 {
		margin = 1
	}
	if b := limit - margin; b >= 1 {
		return b
	}
	return 1
}

func shortWindow(r Rate) time.Duration {
	return time.Duration(float64(r.Per) * ShortWindowMargin)
}

// Acquire admits a request when both its short window and its hourly budget
// have room, and records it in both. Otherwise it returns how long to wait
// before trying again.
func (rl *EndpointRateLimiter) Acquire(method, path string) (bool, time.Duration, Ticket) {
	key, limit := rl.resolve(method, path)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	b := rl.buckets[key]
	if b == nil {
		b = &endpointBucket{limit: limit}
		rl.buckets[key] = b
	}
	now := rl.now()
	window := shortWindow(limit.Short)
	b.short = dropBefore(b.short, now.Add(-window))
	b.hourly = dropBefore(b.hourly, now.Add(-time.Hour))

	if budget := hourlyBudget(limit.Hourly); budget > 0 && len(b.hourly) >= budget {
		return true, b.hourly[0].Add(time.Hour).Sub(now), Ticket{}
	}
	if limit.Short.Requests > 0 && len(b.short) >= limit.Short.Requests {
		return true, b.short[0].Add(window).Sub(now), Ticket{}
	}

	b.short = append(b.short, now)
	b.hourly = append(b.hourly, now)
	return false, 0, Ticket{key: key, at: now}
}

// Release takes an admitted request back out of its hourly budget, for a
// request the API rejected without processing it (429). The request still
// counts in its short window, so pacing is unchanged.
func (rl *EndpointRateLimiter) Release(t Ticket) {
	if t.key == "" {
		return
	}
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b := rl.buckets[t.key]
	if b == nil {
		return
	}
	for i, at := range b.hourly {
		if at.Equal(t.at) {
			b.hourly = append(b.hourly[:i], b.hourly[i+1:]...)
			return
		}
	}
}

// dropBefore removes the times before cutoff from a slice sorted by time.
func dropBefore(times []time.Time, cutoff time.Time) []time.Time {
	i := 0
	for i < len(times) && !times[i].After(cutoff) {
		i++
	}
	if i == 0 {
		return times
	}
	return append(times[:0:0], times[i:]...)
}
