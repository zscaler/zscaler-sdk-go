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
