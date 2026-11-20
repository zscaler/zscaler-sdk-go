package ratelimiter

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const ztw = "/ztw/api/v1"

func newZTWLimiterAt(clock *fakeClock) *EndpointRateLimiter {
	rl := NewZTWEndpointRateLimiter()
	rl.now = clock.now
	return rl
}

func TestZTWEndpointLimiter_DefaultsPerMethodAndResource(t *testing.T) {
	clock := &fakeClock{t: time.Unix(1_000_000, 0)}
	rl := newZTWLimiterAt(clock)

	// POST 1/sec per resource.
	admit(t, rl, http.MethodPost, ztw+"/networkServiceGroups")
	d := refused(t, rl, http.MethodPost, ztw+"/networkServiceGroups")
	require.Equal(t, 1100*time.Millisecond, d, "1/sec plus the 10% window margin")

	// Other methods and other resources have their own windows.
	admit(t, rl, http.MethodGet, ztw+"/networkServiceGroups")
	admit(t, rl, http.MethodGet, ztw+"/networkServiceGroups")
	refused(t, rl, http.MethodGet, ztw+"/networkServiceGroups")
	admit(t, rl, http.MethodPost, ztw+"/networkServices")
	admit(t, rl, http.MethodPut, ztw+"/networkServiceGroups/1")
	admit(t, rl, http.MethodDelete, ztw+"/networkServiceGroups/1")
	refused(t, rl, http.MethodDelete, ztw+"/networkServiceGroups/2")

	clock.advance(1100 * time.Millisecond)
	admit(t, rl, http.MethodPost, ztw+"/networkServiceGroups")
}

func TestZTWEndpointLimiter_HourlyBudget(t *testing.T) {
	clock := &fakeClock{t: time.Unix(1_000_000, 0)}
	rl := newZTWLimiterAt(clock)

	// 400/hr less the 5% margin, per method and resource.
	for i := 0; i < 380; i++ {
		admit(t, rl, http.MethodPost, ztw+"/ipDestinationGroups")
		clock.advance(1100 * time.Millisecond)
	}
	d := refused(t, rl, http.MethodPost, ztw+"/ipDestinationGroups")
	require.Greater(t, d, 50*time.Minute)

	// Another resource still has its own budget.
	admit(t, rl, http.MethodPost, ztw+"/ipSourceGroups")
}

func TestZTWEndpointLimiter_StricterResourcesSharedWithZIA(t *testing.T) {
	clock := &fakeClock{t: time.Unix(1_000_000, 0)}
	rl := newZTWLimiterAt(clock)

	// GET /dnsGateways: 1/sec, as on ZIA.
	admit(t, rl, http.MethodGet, ztw+"/dnsGateways")
	refused(t, rl, http.MethodGet, ztw+"/dnsGateways")

	// DELETE /workloadGroups/{id}: 10/min, as on ZIA.
	for i := 0; i < 10; i++ {
		admit(t, rl, http.MethodDelete, ztw+"/workloadGroups/7")
	}
	refused(t, rl, http.MethodDelete, ztw+"/workloadGroups/8")
}

func TestZTWEndpointLimiter_EveryZIARuleIsReachable(t *testing.T) {
	rl := NewZTWEndpointRateLimiter()
	for _, r := range ZIAEndpointRules() {
		var segs []string
		for _, s := range strings.Split(r.Path, "/") {
			if strings.HasPrefix(s, "{") {
				s = "123"
			}
			segs = append(segs, s)
		}
		key, limit := rl.resolve(r.Method, ztw+strings.Join(segs, "/"))
		require.Equal(t, r.Method+" "+r.Path, key, "rule %s %s is shadowed", r.Method, r.Path)
		require.Equal(t, r.Limit, limit)
	}
}

func TestZTWEndpointLimiter_SeparateFromZIA(t *testing.T) {
	clock := &fakeClock{t: time.Unix(1_000_000, 0)}
	zia := newZIALimiterAt(clock)
	ztwRL := newZTWLimiterAt(clock)

	// The ZIA and ZTW clients each have their own limiter, so the same
	// resource on both services is paced separately.
	admit(t, zia, http.MethodPost, "/zia/api/v1/networkServiceGroups")
	admit(t, ztwRL, http.MethodPost, ztw+"/networkServiceGroups")
}
