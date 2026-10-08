package zscaler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
)

// rewriteTransport sends every request to the test server, keeping the path,
// so the fixed OneAPI auth and API hosts resolve to httptest.
type rewriteTransport struct{ target *url.URL }

func (t rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	r.URL.Scheme = t.target.Scheme
	r.URL.Host = t.target.Host
	r.Host = t.target.Host
	return http.DefaultTransport.RoundTrip(r)
}

// newSessionTestClient returns a OneAPI client holding a still-unexpired token
// ("tok-0") whose HTTP clients are routed to srv.
func newSessionTestClient(t *testing.T, srv *httptest.Server) *Client {
	t.Helper()
	target, err := url.Parse(srv.URL)
	require.NoError(t, err)
	httpClient := &http.Client{Transport: rewriteTransport{target: target}}

	cfg := &Configuration{}
	cfg.Logger = logger.GetDefaultLogger("session-test: ")
	cfg.Context = context.Background()
	cfg.HTTPClient = httpClient
	cfg.ZIAHTTPClient = httpClient
	cfg.Zscaler.Client.ClientID = "client-id"
	cfg.Zscaler.Client.ClientSecret = "client-secret"
	cfg.Zscaler.Client.VanityDomain = "example"
	cfg.Zscaler.Client.RateLimit.MaxRetries = 10
	cfg.Zscaler.Client.RateLimit.MaxSessionNotValidRetries = 3
	cfg.Zscaler.Client.AuthToken = &AuthToken{AccessToken: "tok-0", Expiry: time.Now().Add(time.Hour)}
	return &Client{oauth2Credentials: cfg}
}

// tokenHandler mints tok-1, tok-2, ... and counts how many tokens were issued.
func tokenHandler(issued *int32, w http.ResponseWriter) {
	n := atomic.AddInt32(issued, 1)
	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprintf(w, `{"token_type":"Bearer","access_token":"tok-%d","expires_in":3599}`, n)
}

// A 403 "Resource Access Blocked" (API session ended while the OAuth token is
// still unexpired) must trigger a token refresh and a successful retry.
func TestExecuteRequest_403SessionEnded_RefreshesTokenAndRetries(t *testing.T) {
	var issued, apiCalls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth2/v1/token" {
			tokenHandler(&issued, w)
			return
		}
		atomic.AddInt32(&apiCalls, 1)
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("Authorization") == "Bearer tok-0" {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"message":"Resource Access Blocked"}`))
			return
		}
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newSessionTestClient(t, srv)
	body, resp, _, err := c.ExecuteRequest(context.Background(), http.MethodGet, "/zia/api/v1/urlFilteringRules", nil, nil, "")
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "[]", string(body))
	assert.Equal(t, int32(1), atomic.LoadInt32(&issued), "exactly one token refresh expected")
	assert.Equal(t, int32(2), atomic.LoadInt32(&apiCalls), "original call plus one retry expected")
}

// A 403 that is not a session error (e.g. permissions) must fail immediately,
// without refreshing the token or retrying.
func TestExecuteRequest_403OtherReason_FailsWithoutRetry(t *testing.T) {
	var issued, apiCalls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth2/v1/token" {
			tokenHandler(&issued, w)
			return
		}
		atomic.AddInt32(&apiCalls, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"code":"ACCESS_DENIED","message":"Insufficient permissions"}`))
	}))
	defer srv.Close()

	c := newSessionTestClient(t, srv)
	_, resp, _, err := c.ExecuteRequest(context.Background(), http.MethodGet, "/zia/api/v1/urlFilteringRules", nil, nil, "")
	require.Error(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	assert.Equal(t, int32(0), atomic.LoadInt32(&issued), "no token refresh expected")
	assert.Equal(t, int32(1), atomic.LoadInt32(&apiCalls), "no retry expected")
}

// If the session error persists after refreshing, retries stop at
// MaxSessionNotValidRetries.
func TestExecuteRequest_403SessionEnded_BoundedRetries(t *testing.T) {
	if testing.Short() {
		t.Skip("waits ~6s for the post-refresh delays")
	}
	var issued, apiCalls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth2/v1/token" {
			tokenHandler(&issued, w)
			return
		}
		atomic.AddInt32(&apiCalls, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"message":"Resource Access Blocked"}`))
	}))
	defer srv.Close()

	c := newSessionTestClient(t, srv)
	_, _, _, err := c.ExecuteRequest(context.Background(), http.MethodGet, "/zia/api/v1/urlFilteringRules", nil, nil, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "max SESSION_NOT_VALID retries exceeded")
	assert.Equal(t, int32(3), atomic.LoadInt32(&issued), "one refresh per allowed session retry")
	assert.Equal(t, int32(4), atomic.LoadInt32(&apiCalls), "original call plus three retries")
}
