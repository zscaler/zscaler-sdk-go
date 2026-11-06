package ztw

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
)

// deleteRetryServer serves a rule that the first DELETE deletes although the
// API answers it with an edit lock conflict; every later DELETE gets 404. With
// applied=false the rule never existed and every DELETE gets 404.
func deleteRetryServer(attempts *int32, applied bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(attempts, 1)
		w.Header().Set("Content-Type", "application/json")
		if applied && n == 1 {
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(`{"code":"EDIT_LOCK_NOT_AVAILABLE","message":"Failed during enter Org barrier"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":"RESOURCE_NOT_FOUND","message":"Rule with id 1495551 does not exist."}`))
	}))
}

func newDeleteRetryTestClient(srv *httptest.Server) *Client {
	cfg := &Configuration{}
	cfg.ZTW.Client.RateLimit.MaxRetries = 5
	cfg.ZTW.Client.RateLimit.RetryWaitMin = time.Millisecond
	cfg.ZTW.Client.RateLimit.RetryWaitMax = 2 * time.Millisecond
	l := logger.GetDefaultLogger("ztw-test: ")
	return &Client{
		session:          &Session{JSessionID: "test"},
		sessionRefreshed: time.Now(),
		sessionTimeout:   time.Hour,
		URL:              srv.URL + "/api/v1",
		HTTPClient:       getHTTPClient(l, nil, cfg),
		Logger:           l,
	}
}

// Same case as the ZIA client test (SUP-4290): the API deleted the rule but
// answered with an edit lock conflict, and the retry got 404. The rule is
// deleted.
func TestLegacyDelete_AppliedBeforeRetryIsDeleted(t *testing.T) {
	var attempts int32
	srv := deleteRetryServer(&attempts, true)
	defer srv.Close()

	err := newDeleteRetryTestClient(srv).Delete(context.Background(), "/firewallFilteringRules/1495551")

	require.NoError(t, err)
	assert.Equal(t, int32(2), atomic.LoadInt32(&attempts))
}

func TestLegacyDelete_MissingObjectStillFails(t *testing.T) {
	var attempts int32
	srv := deleteRetryServer(&attempts, false)
	defer srv.Close()

	err := newDeleteRetryTestClient(srv).Delete(context.Background(), "/firewallFilteringRules/1495551")

	require.Error(t, err, "a 404 on the first attempt means the rule did not exist")
	respErr, ok := errorx.AsErrorResponse(err)
	require.True(t, ok)
	require.NotNil(t, respErr.Parsed)
	assert.Equal(t, "RESOURCE_NOT_FOUND", respErr.Parsed.Code)
	assert.Equal(t, int32(1), atomic.LoadInt32(&attempts))
}
