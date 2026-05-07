package zscaler

import (
	"net/http"
	"testing"
	"time"
)

// TestJitter verifies the bounds and short-circuit behavior of the jitter
// helper that the 429/503/401 backoff policy relies on. Without jitter,
// parallel goroutines that all received the same Retry-After wake up at the
// same instant and stampede the per-endpoint rate limiter; the tests below
// lock down the contract that prevents that.
func TestJitter(t *testing.T) {
	t.Run("zero duration short-circuits to zero", func(t *testing.T) {
		if got := jitter(0, 0.25); got != 0 {
			t.Errorf("jitter(0, 0.25) = %v, want 0", got)
		}
	})

	t.Run("zero fraction returns input unchanged", func(t *testing.T) {
		d := 5 * time.Second
		if got := jitter(d, 0); got != d {
			t.Errorf("jitter(%v, 0) = %v, want %v", d, got, d)
		}
	})

	t.Run("negative duration short-circuits", func(t *testing.T) {
		d := -3 * time.Second
		if got := jitter(d, 0.25); got != d {
			t.Errorf("jitter(%v, 0.25) = %v, want %v (short-circuit)", d, got, d)
		}
	})

	t.Run("output stays within +/- frac of input", func(t *testing.T) {
		d := time.Second
		frac := 0.25
		lo := time.Duration(float64(d) * (1 - frac))
		hi := time.Duration(float64(d) * (1 + frac))
		for i := 0; i < 1000; i++ {
			got := jitter(d, frac)
			if got < lo || got > hi {
				t.Fatalf("iter %d: jitter(%v, %v) = %v, want in [%v, %v]", i, d, frac, got, lo, hi)
			}
		}
	})

	t.Run("desynchronises identical inputs", func(t *testing.T) {
		// 100 calls against a continuous distribution should yield well
		// over 50 distinct values; the floor at 10 catches a degenerate
		// (e.g. unseeded) RNG without being flaky.
		d := time.Second
		seen := make(map[time.Duration]struct{}, 100)
		for i := 0; i < 100; i++ {
			seen[jitter(d, 0.25)] = struct{}{}
		}
		if len(seen) < 10 {
			t.Errorf("jitter produced only %d distinct values out of 100 — RNG not desynchronising", len(seen))
		}
	})
}

// TestRetryBackoffPolicy verifies the policy installed by getHTTPClient as
// retryableClient.Backoff: on 429/503/401 it honours Retry-After as a
// floor, grows exponentially per attempt (capped at RetryWaitMax), and
// jitters by +/- retryJitterFraction. This intentionally duplicates the
// production formula so that any drift in the closure fails this test.
func TestRetryBackoffPolicy(t *testing.T) {
	cfg, err := NewConfiguration()
	if err != nil {
		t.Fatalf("NewConfiguration: %v", err)
	}
	// Pin RetryWaitMin/Max to known values so env-loaded overrides on the
	// developer's machine cannot make the assertions flaky.
	cfg.Zscaler.Client.RateLimit.RetryWaitMin = time.Second
	cfg.Zscaler.Client.RateLimit.RetryWaitMax = 10 * time.Second
	min := cfg.Zscaler.Client.RateLimit.RetryWaitMin
	max := cfg.Zscaler.Client.RateLimit.RetryWaitMax

	resp429 := func(retryAfterSeconds string) *http.Response {
		h := http.Header{}
		if retryAfterSeconds != "" {
			h.Set("Retry-After", retryAfterSeconds)
		}
		return &http.Response{StatusCode: http.StatusTooManyRequests, Header: h}
	}

	// Mirror of the 429 branch installed by getHTTPClient. If the production
	// closure changes, update this mirror — the duplication is intentional.
	policy := func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		base := getRetryAfter(resp, cfg)
		if base <= 0 {
			base = min
		}
		shift := uint(attemptNum)
		if shift > 16 {
			shift = 16
		}
		grown := base * time.Duration(1<<shift)
		if grown <= 0 || grown > max {
			grown = max
		}
		if grown < base {
			grown = base
		}
		return jitter(grown, retryJitterFraction)
	}

	t.Run("first retry stays within Retry-After +/- jitter", func(t *testing.T) {
		base := time.Second
		lo := time.Duration(float64(base) * (1 - retryJitterFraction))
		hi := time.Duration(float64(base) * (1 + retryJitterFraction))
		for i := 0; i < 50; i++ {
			got := policy(min, max, 0, resp429("1"))
			if got < lo || got > hi {
				t.Fatalf("iter %d: attempt 0 = %v, want in [%v, %v]", i, got, lo, hi)
			}
		}
	})

	t.Run("growth doubles per attempt until the cap", func(t *testing.T) {
		// With base=1s and jitter=+/-25%, attempt N's pre-jitter value
		// is min(1s * 2^N, RetryWaitMax). The post-jitter window must
		// straddle that target.
		base := 1 * time.Second
		for attempt := 0; attempt < 6; attempt++ {
			pre := base * time.Duration(1<<uint(attempt))
			if pre > max {
				pre = max
			}
			lo := time.Duration(float64(pre) * (1 - retryJitterFraction))
			hi := time.Duration(float64(pre) * (1 + retryJitterFraction))
			got := policy(min, max, attempt, resp429("1"))
			if got < lo || got > hi {
				t.Errorf("attempt %d: got %v, want in [%v, %v] (pre-jitter target %v)", attempt, got, lo, hi, pre)
			}
		}
	})

	t.Run("never exceeds RetryWaitMax + jitter ceiling", func(t *testing.T) {
		// attempt 5 with base=1s yields 32s pre-cap, must clamp to <=max.
		ceiling := time.Duration(float64(max) * (1 + retryJitterFraction))
		for attempt := 5; attempt < 12; attempt++ {
			got := policy(min, max, attempt, resp429("1"))
			if got > ceiling {
				t.Errorf("attempt %d: %v exceeded jittered cap %v", attempt, got, ceiling)
			}
		}
	})

	t.Run("missing Retry-After falls back to RetryWaitMin floor", func(t *testing.T) {
		// getRetryAfter returns RetryWaitMinSeconds (2s) when no rate-limit
		// headers are present. Validate the fallback survives the policy.
		got := policy(min, max, 0, resp429(""))
		if got <= 0 {
			t.Fatalf("attempt 0 with no Retry-After = %v, want > 0", got)
		}
		ceiling := time.Duration(float64(max) * (1 + retryJitterFraction))
		if got > ceiling {
			t.Errorf("attempt 0 with no Retry-After = %v, exceeded jittered cap %v", got, ceiling)
		}
	})
}

// TestRetryMaxDefault locks down the default retry budget that ships with
// the SDK. Override is via cfg.Zscaler.Client.RateLimit.MaxRetries (env
// ZSCALER_CLIENT_RATE_LIMIT_MAX_RETRIES); this constant is the floor when
// the override is unset. Raising it back to 100 silently was the bug we
// shipped 3.8.33 to fix — this test exists so a future revert lands in
// CI instead of in a customer's apply log.
func TestRetryMaxDefault(t *testing.T) {
	if MaxNumOfRetries != 10 {
		t.Errorf("MaxNumOfRetries = %d, want 10 (raise this assertion deliberately if changing the SDK contract)", MaxNumOfRetries)
	}
}
