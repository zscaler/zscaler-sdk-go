package zscaler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A write through one ZPA API version must invalidate cached reads of the same
// object under the other version: segment groups are updated with
// PUT /mgmtconfig/v2/... and read with GET /mgmtconfig/v1/..., so without this
// the read right after the update was served stale from cache.
func TestExecuteRequest_V2WriteInvalidatesV1CachedRead(t *testing.T) {
	const v1 = "/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup/9"
	const v2 = "/zpa/mgmtconfig/v2/admin/customers/123/segmentGroup/9"
	var gets int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oauth2/v1/token":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"token_type":"Bearer","access_token":"tok","expires_in":3599}`))
		case r.Method == http.MethodGet && r.URL.Path == v1:
			n := atomic.AddInt32(&gets, 1)
			w.Header().Set("Content-Type", "application/json")
			if n == 1 {
				_, _ = w.Write([]byte(`{"id":"9","description":"before"}`))
			} else {
				_, _ = w.Write([]byte(`{"id":"9","description":"after"}`))
			}
		case r.Method == http.MethodPut && r.URL.Path == v2:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	c := newSessionTestClient(t, srv)
	cfg := c.oauth2Credentials
	cfg.ZPAHTTPClient = cfg.HTTPClient
	cfg.Zscaler.Client.Cache.Enabled = true
	cfg.Zscaler.Client.Cache.DefaultTtl = 10 * time.Minute
	cfg.Zscaler.Client.Cache.DefaultTti = 8 * time.Minute
	cfg.Zscaler.Client.Cache.DefaultCacheMaxSizeMB = 1
	cfg.CacheManager = newCache(cfg)
	ctx := context.Background()

	body, _, _, err := c.ExecuteRequest(ctx, http.MethodGet, v1, nil, nil, "")
	require.NoError(t, err)
	assert.Contains(t, string(body), "before")

	_, _, _, err = c.ExecuteRequest(ctx, http.MethodPut, v2, nil, nil, "")
	require.NoError(t, err)

	body, _, _, err = c.ExecuteRequest(ctx, http.MethodGet, v1, nil, nil, "")
	require.NoError(t, err)
	assert.Contains(t, string(body), "after", "read after a v2 write must not be served from the v1 cache")
	assert.Equal(t, int32(2), atomic.LoadInt32(&gets))
}

// newCachingTestClient returns a OneAPI client for srv with the response cache
// enabled.
func newCachingTestClient(t *testing.T, srv *httptest.Server) *Client {
	t.Helper()
	c := newSessionTestClient(t, srv)
	cfg := c.oauth2Credentials
	cfg.Zscaler.Client.Cache.Enabled = true
	cfg.Zscaler.Client.Cache.DefaultTtl = 10 * time.Minute
	cfg.Zscaler.Client.Cache.DefaultTti = 8 * time.Minute
	cfg.Zscaler.Client.Cache.DefaultCacheMaxSizeMB = 1
	cfg.CacheManager = newCache(cfg)
	return c
}

// Replays the interleaving from a terraform-provider-zia apply that updated two
// URL categories in parallel: the list read for one category's refresh reached
// the API while the other category's PUT was still being processed, received
// the list from before the PUT, and was cached after the PUT had invalidated
// the cache. The refresh of the updated category was then served that stale
// list, and its state kept the URLs the PUT had removed.
func TestExecuteRequest_ReadOverlappingWriteIsNotCached(t *testing.T) {
	const list = "/zia/api/v1/urlCategories"
	const item = "/zia/api/v1/urlCategories/CUSTOM_14"
	var (
		gets       int32
		updated    atomic.Bool
		putArrived = make(chan struct{})
		listServed = make(chan struct{})
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oauth2/v1/token":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"token_type":"Bearer","access_token":"tok","expires_in":3599}`))
		case r.Method == http.MethodGet && r.URL.Path == list:
			atomic.AddInt32(&gets, 1)
			w.Header().Set("Content-Type", "application/json")
			if updated.Load() {
				_, _ = w.Write([]byte(`[{"id":"CUSTOM_14","urls":[]}]`))
			} else {
				_, _ = w.Write([]byte(`[{"id":"CUSTOM_14","urls":["before.example.com"]}]`))
			}
		case r.Method == http.MethodPut && r.URL.Path == item:
			close(putArrived)
			<-listServed // the API is still processing the PUT while the list is read
			updated.Store(true)
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	c := newCachingTestClient(t, srv)
	ctx := context.Background()

	putDone := make(chan error, 1)
	go func() {
		_, _, _, err := c.ExecuteRequest(ctx, http.MethodPut, item, nil, nil, "")
		putDone <- err
	}()
	<-putArrived

	body, _, _, err := c.ExecuteRequest(ctx, http.MethodGet, list+"?customOnly=true", nil, nil, "")
	require.NoError(t, err)
	assert.Contains(t, string(body), "before.example.com", "the overlapping read receives the list from before the PUT")
	close(listServed)
	require.NoError(t, <-putDone)

	body, _, _, err = c.ExecuteRequest(ctx, http.MethodGet, list+"?customOnly=true", nil, nil, "")
	require.NoError(t, err)
	assert.NotContains(t, string(body), "before.example.com", "a read after the PUT must not be served the overlapping read's response")
	assert.Equal(t, int32(2), atomic.LoadInt32(&gets), "the overlapping read must not have been cached")
}

// Reads that do not overlap a write are still cached.
func TestExecuteRequest_ReadWithoutOverlappingWriteIsCached(t *testing.T) {
	const list = "/zia/api/v1/urlCategories"
	var gets int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oauth2/v1/token":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"token_type":"Bearer","access_token":"tok","expires_in":3599}`))
		case r.Method == http.MethodGet && r.URL.Path == list:
			atomic.AddInt32(&gets, 1)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"id":"CUSTOM_14"}]`))
		case r.Method == http.MethodPut:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	c := newCachingTestClient(t, srv)
	ctx := context.Background()

	// A finished write to another collection does not affect caching.
	_, _, _, err := c.ExecuteRequest(ctx, http.MethodPut, "/zia/api/v1/firewallFilteringRules/1", nil, nil, "")
	require.NoError(t, err)
	for i := 0; i < 3; i++ {
		_, _, _, err := c.ExecuteRequest(ctx, http.MethodGet, list, nil, nil, "")
		require.NoError(t, err)
	}
	assert.Equal(t, int32(1), atomic.LoadInt32(&gets), "repeated reads without writes must be served from cache")
}
