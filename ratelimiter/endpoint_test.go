package ratelimiter

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
)

// fakeClock drives an EndpointRateLimiter without sleeping.
type fakeClock struct{ t time.Time }

func (c *fakeClock) now() time.Time          { return c.t }
func (c *fakeClock) advance(d time.Duration) { c.t = c.t.Add(d) }

func newZIALimiterAt(clock *fakeClock) *EndpointRateLimiter {
	rl := NewZIAEndpointRateLimiter()
	rl.now = clock.now
	return rl
}

func admit(t *testing.T, rl *EndpointRateLimiter, method, path string) Ticket {
	t.Helper()
	wait, delay, ticket := rl.Acquire(method, path)
	require.False(t, wait, "%s %s should be admitted, got wait %v", method, path, delay)
	return ticket
}

func refused(t *testing.T, rl *EndpointRateLimiter, method, path string) time.Duration {
	t.Helper()
	wait, delay, _ := rl.Acquire(method, path)
	require.True(t, wait, "%s %s should wait", method, path)
	require.Greater(t, delay, time.Duration(0))
	return delay
}

const zia = "/zia/api/v1"

func TestEndpointLimiter_DefaultWritesOnePerSecondPerEndpoint(t *testing.T) {
	clock := &fakeClock{t: time.Unix(1_000_000, 0)}
	rl := newZIALimiterAt(clock)

	admit(t, rl, http.MethodPost, zia+"/networkServiceGroups")
	d := refused(t, rl, http.MethodPost, zia+"/networkServiceGroups")
	require.Equal(t, 1100*time.Millisecond, d, "1/sec plus the 10% window margin")

	// Other methods and other resources have their own windows.
	admit(t, rl, http.MethodGet, zia+"/networkServiceGroups")
	admit(t, rl, http.MethodPost, zia+"/networkServices")
	admit(t, rl, http.MethodPut, zia+"/networkServiceGroups/1")
	admit(t, rl, http.MethodDelete, zia+"/networkServiceGroups/1")

	clock.advance(1100 * time.Millisecond)
	admit(t, rl, http.MethodPost, zia+"/networkServiceGroups")
}

func TestEndpointLimiter_DefaultReadsTwoPerSecond(t *testing.T) {
	clock := &fakeClock{t: time.Unix(1_000_000, 0)}
	rl := newZIALimiterAt(clock)

	// Paged reads share the budget of the resource: the query is not part of the path.
	admit(t, rl, http.MethodGet, zia+"/networkServiceGroups")
	admit(t, rl, http.MethodGet, zia+"/networkServiceGroups")
	refused(t, rl, http.MethodGet, zia+"/networkServiceGroups")
}

func TestEndpointLimiter_IDsShareTheResourceBudget(t *testing.T) {
	clock := &fakeClock{t: time.Unix(1_000_000, 0)}
	rl := newZIALimiterAt(clock)

	admit(t, rl, http.MethodPut, zia+"/firewallFilteringRules/101")
	refused(t, rl, http.MethodPut, zia+"/firewallFilteringRules/202")
	// A literal segment is its own resource.
	admit(t, rl, http.MethodGet, zia+"/networkServiceGroups/lite")
	admit(t, rl, http.MethodGet, zia+"/networkServiceGroups/lite")
	refused(t, rl, http.MethodGet, zia+"/networkServiceGroups/lite")
	admit(t, rl, http.MethodGet, zia+"/networkServiceGroups/77")
}

func TestEndpointLimiter_DefaultHourlyBudget(t *testing.T) {
	clock := &fakeClock{t: time.Unix(1_000_000, 0)}
	rl := newZIALimiterAt(clock)

	for i := 0; i < 380; i++ {
		admit(t, rl, http.MethodPost, zia+"/networkServiceGroups")
		clock.advance(1100 * time.Millisecond)
	}
	d := refused(t, rl, http.MethodPost, zia+"/networkServiceGroups")
	require.InDelta(t, (time.Hour - 380*1100*time.Millisecond).Seconds(), d.Seconds(), 0.001,
		"waits until the first request of the hour ages out")

	// Another resource is not affected.
	admit(t, rl, http.MethodPost, zia+"/networkServices")
}

func TestEndpointLimiter_StricterResourcesFromTheSummary(t *testing.T) {
	clock := &fakeClock{t: time.Unix(1_000_000, 0)}
	rl := newZIALimiterAt(clock)

	// GET 1/sec on rule types and time intervals.
	admit(t, rl, http.MethodGet, zia+"/timeIntervals/5")
	refused(t, rl, http.MethodGet, zia+"/timeIntervals/6")
	admit(t, rl, http.MethodGet, zia+"/bandwidthControlRules")
	refused(t, rl, http.MethodGet, zia+"/bandwidthControlRules")

	// Activation: 10/min and 40/hr.
	for i := 0; i < 10; i++ {
		admit(t, rl, http.MethodPost, zia+"/status/activate")
	}
	d := refused(t, rl, http.MethodPost, zia+"/status/activate")
	require.Equal(t, 66*time.Second, d)

	// GET /users: 10/min and 40/hr.
	for i := 0; i < 10; i++ {
		admit(t, rl, http.MethodGet, zia+"/users")
	}
	refused(t, rl, http.MethodGet, zia+"/users")
	// GET /users/{id} is not in the summary's stricter list.
	admit(t, rl, http.MethodGet, zia+"/users/42")
	admit(t, rl, http.MethodGet, zia+"/users/43")
}

func TestEndpointLimiter_ActivationHourlyBudget(t *testing.T) {
	clock := &fakeClock{t: time.Unix(1_000_000, 0)}
	rl := newZIALimiterAt(clock)

	// 40/hr less the margin is 38.
	for i := 0; i < 38; i++ {
		admit(t, rl, http.MethodPost, zia+"/status/activate")
		clock.advance(7 * time.Second)
	}
	refused(t, rl, http.MethodPost, zia+"/status/activate")
}

func TestEndpointLimiter_PlaceholdersMatchNonNumericSegments(t *testing.T) {
	clock := &fakeClock{t: time.Unix(1_000_000, 0)}
	rl := newZIALimiterAt(clock)

	// PUT /webApplicationRules/{rule_type}/{ruleId}: 10/min, all rule types together.
	for i := 0; i < 5; i++ {
		admit(t, rl, http.MethodPut, zia+"/webApplicationRules/STREAMING_MEDIA/77")
		admit(t, rl, http.MethodPut, zia+"/webApplicationRules/FILE_SHARE/88")
	}
	refused(t, rl, http.MethodPut, zia+"/webApplicationRules/WEBMAIL/99")
}

func TestEndpointLimiter_MatchIsCaseInsensitive(t *testing.T) {
	clock := &fakeClock{t: time.Unix(1_000_000, 0)}
	rl := newZIALimiterAt(clock)

	for i := 0; i < 10; i++ {
		admit(t, rl, http.MethodPost, zia+"/endpointDlpRules")
	}
	refused(t, rl, http.MethodPost, zia+"/ENDPOINTDLPRULES")
}

func TestEndpointLimiter_MostSpecificRuleWins(t *testing.T) {
	clock := &fakeClock{t: time.Unix(1_000_000, 0)}
	rl := NewEndpointRateLimiter(zia, ZIADefaultLimits(), []EndpointRule{
		{http.MethodGet, "/things/{id}", EndpointLimit{Rate{1, time.Minute}, 4}},
		{http.MethodGet, "/things/lite", EndpointLimit{Rate{5, time.Second}, 100}},
	})
	rl.now = clock.now

	for i := 0; i < 5; i++ {
		admit(t, rl, http.MethodGet, zia+"/things/lite")
	}
	refused(t, rl, http.MethodGet, zia+"/things/lite")
	admit(t, rl, http.MethodGet, zia+"/things/9")
	refused(t, rl, http.MethodGet, zia+"/things/10")
}

func TestEndpointLimiter_EveryZIARuleIsReachable(t *testing.T) {
	rl := NewZIAEndpointRateLimiter()
	for _, r := range ZIAEndpointRules() {
		var segs []string
		for _, s := range strings.Split(r.Path, "/") {
			if strings.HasPrefix(s, "{") {
				s = "123"
			}
			segs = append(segs, s)
		}
		key, limit := rl.resolve(r.Method, zia+strings.Join(segs, "/"))
		require.Equal(t, r.Method+" "+r.Path, key, "rule %s %s is shadowed", r.Method, r.Path)
		require.Equal(t, r.Limit, limit)
	}
}

func TestEndpointLimiter_UnknownMethodUsesWriteDefault(t *testing.T) {
	rl := NewEndpointRateLimiter(zia, map[string]EndpointLimit{
		http.MethodPut: {Rate{1, time.Second}, 400},
	}, nil)
	_, limit := rl.resolve("PATCH", zia+"/anything")
	require.Equal(t, EndpointLimit{Rate{1, time.Second}, 400}, limit)
}

func TestHourlyBudget(t *testing.T) {
	for limit, want := range map[int]int{0: 0, 1: 1, 2: 1, 4: 3, 40: 38, 400: 380, 1000: 950} {
		require.Equal(t, want, hourlyBudget(limit), "limit %d", limit)
	}
}

func TestEndpointLimiter_ReleaseReturnsTheHourlySlot(t *testing.T) {
	clock := &fakeClock{t: time.Unix(1_000_000, 0)}
	rl := NewEndpointRateLimiter(zia, ZIADefaultLimits(), []EndpointRule{
		{http.MethodPost, "/limited", EndpointLimit{Rate{10, time.Second}, 2}}, // budget 1
	})
	rl.now = clock.now

	ticket := admit(t, rl, http.MethodPost, zia+"/limited")
	refused(t, rl, http.MethodPost, zia+"/limited")
	rl.Release(ticket)
	admit(t, rl, http.MethodPost, zia+"/limited")
	refused(t, rl, http.MethodPost, zia+"/limited")

	rl.Release(Ticket{}) // a zero ticket is ignored
}

func TestEndpointLimiter_ReleaseKeepsThePacing(t *testing.T) {
	clock := &fakeClock{t: time.Unix(1_000_000, 0)}
	rl := newZIALimiterAt(clock)

	ticket := admit(t, rl, http.MethodPost, zia+"/networkServiceGroups")
	rl.Release(ticket)
	refused(t, rl, http.MethodPost, zia+"/networkServiceGroups")
}

// statusTransport answers with the status codes of a script, then 200.
type statusTransport struct {
	mu     sync.Mutex
	script []int
	sent   []time.Time
}

func (s *statusTransport) RoundTrip(*http.Request) (*http.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sent = append(s.sent, time.Now())
	code := http.StatusOK
	if len(s.script) > 0 {
		code, s.script = s.script[0], s.script[1:]
	}
	return &http.Response{StatusCode: code, Body: http.NoBody}, nil
}

func TestRateLimitTransport_EndpointLimiterDoesNotCount429(t *testing.T) {
	rl := NewEndpointRateLimiter(zia, ZIADefaultLimits(), []EndpointRule{
		{http.MethodPost, "/limited", EndpointLimit{Rate{10, time.Second}, 2}}, // budget 1
	})
	base := &statusTransport{script: []int{http.StatusTooManyRequests, http.StatusTooManyRequests}}
	transport := &RateLimitTransport{Base: base, EndpointLimiter: rl, Logger: logger.NewNopLogger()}

	for _, want := range []int{429, 429, 200} {
		// A request that waits for the hourly budget fails instead of blocking.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.zsapi.net"+zia+"/limited", nil)
		resp, err := transport.RoundTrip(req)
		cancel()
		require.NoError(t, err, "a request rejected with 429 used the hourly budget")
		require.Equal(t, want, resp.StatusCode)
	}
	// The 200 used the only hourly slot; the two 429s did not.
	wait, _, _ := rl.Acquire(http.MethodPost, zia+"/limited")
	require.True(t, wait)
}

func TestRateLimitTransport_EndpointLimiterPacesParallelCreates(t *testing.T) {
	// Same shape as 10 parallel Terraform creates, with a 50ms window instead
	// of a second.
	rl := NewEndpointRateLimiter(zia, map[string]EndpointLimit{
		http.MethodPost: {Rate{1, 50 * time.Millisecond}, 400},
	}, nil)
	base := &statusTransport{}
	transport := &RateLimitTransport{Base: base, EndpointLimiter: rl, Logger: logger.NewNopLogger()}

	var wg sync.WaitGroup
	var ok atomic.Int32
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest(http.MethodPost, "https://api.zsapi.net"+zia+"/networkServiceGroups", nil)
			if resp, err := transport.RoundTrip(req); err == nil && resp.StatusCode == http.StatusOK {
				ok.Add(1)
			}
		}()
	}
	wg.Wait()
	require.EqualValues(t, 10, ok.Load())

	require.Len(t, base.sent, 10)
	for i := 1; i < len(base.sent); i++ {
		gap := base.sent[i].Sub(base.sent[i-1])
		require.GreaterOrEqual(t, gap, 50*time.Millisecond, "request %d sent %v after the previous one", i+1, gap)
	}
}

func TestRateLimitTransport_EndpointLimiterRetriedCreatesDoNotExhaustTheHour(t *testing.T) {
	// The run that stalled: 127 creates on one resource, many of them answered
	// with 429 and retried. Only the accepted creates count, so the hourly
	// budget (380) is never reached.
	rl := NewEndpointRateLimiter(zia, map[string]EndpointLimit{
		http.MethodPost: {Rate{100, time.Millisecond}, 400},
	}, nil)
	script := make([]int, 0, 600)
	for i := 0; i < 300; i++ {
		script = append(script, http.StatusTooManyRequests) // 300 rejected attempts
	}
	base := &statusTransport{script: script}
	transport := &RateLimitTransport{Base: base, EndpointLimiter: rl, Logger: logger.NewNopLogger()}

	created := 0
	for created < 127 {
		// A request that waits for the hourly budget fails here instead of
		// blocking the test for an hour.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.zsapi.net"+zia+"/networkServiceGroups", nil)
		start := time.Now()
		resp, err := transport.RoundTrip(req)
		cancel()
		require.NoError(t, err, "create %d waited for the hourly budget", created+1)
		require.Less(t, time.Since(start), time.Second, "no request waits for the hourly budget")
		if resp.StatusCode == http.StatusOK {
			created++
		}
	}
	require.Len(t, base.sent, 427)
}
