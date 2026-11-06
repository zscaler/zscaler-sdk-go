package zscaler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
)

const editLockBody = `{"code":"EDIT_LOCK_NOT_AVAILABLE","message":"Failed during enter Org barrier"}`

// deletedOnFirstAttempt serves a rule that the first DELETE deletes although
// that DELETE is answered with status/body. Every later DELETE gets 404, as
// the API answers for a rule that no longer exists. With applied=false the
// rule never existed and every DELETE gets 404.
func deletedOnFirstAttempt(attempts *int32, status int, body string, applied bool) http.HandlerFunc {
	var deleted atomic.Bool
	return func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(attempts, 1)
		w.Header().Set("Content-Type", "application/json")
		if applied && n == 1 {
			deleted.Store(true)
			w.WriteHeader(status)
			_, _ = w.Write([]byte(body))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":"RESOURCE_NOT_FOUND","message":"Rule with id 1495551 does not exist."}`))
	}
}

// The retryable HTTP client built by getHTTPClient records every attempt of a
// DELETE: SUP-4290 / terraform-provider-zia#619, where ZIA deleted the rule
// but answered with an edit lock conflict, and the retry got 404.
func TestOneAPIGetHTTPClient_RecordsDeleteAttempts(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(deletedOnFirstAttempt(&attempts, http.StatusConflict, editLockBody, true))
	defer srv.Close()

	cfg := &Configuration{}
	cfg.Zscaler.Client.RateLimit.MaxRetries = 5
	cfg.Zscaler.Client.RateLimit.RetryWaitMin = time.Millisecond
	cfg.Zscaler.Client.RateLimit.RetryWaitMax = 2 * time.Millisecond
	httpClient := getHTTPClient(logger.GetDefaultLogger("oneapi-test: "), nil, cfg)

	ctx := errorx.TrackDeleteAttempts(context.Background(), http.MethodDelete)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, srv.URL+"/zia/api/v1/firewallFilteringRules/1495551", nil)
	require.NoError(t, err)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, int32(2), atomic.LoadInt32(&attempts), "the edit lock conflict must be retried")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.True(t, errorx.DeleteAlreadyApplied(ctx, http.MethodDelete, resp), "the 404 follows an attempt that deleted the rule")
}

// newDeleteRetryTestClient routes the OneAPI client to srv without the
// retryable HTTP client, so ExecuteRequest's own loop resends the request.
// The transport records attempts as getHTTPClient's does.
func newDeleteRetryTestClient(t *testing.T, srv *httptest.Server) *Client {
	t.Helper()
	target, err := url.Parse(srv.URL)
	require.NoError(t, err)
	c := newSessionTestClient(t, srv)
	httpClient := &http.Client{Transport: errorx.TrackAttempts(rewriteTransport{target: target})}
	cfg := c.oauth2Credentials
	cfg.HTTPClient = httpClient
	cfg.ZIAHTTPClient = httpClient
	cfg.Zscaler.Client.RateLimit.RetryWaitMin = time.Millisecond
	cfg.Zscaler.Client.RateLimit.RetryWaitMax = 2 * time.Millisecond
	return c
}

func TestExecuteRequest_DeleteAppliedBeforeRetryIsDeleted(t *testing.T) {
	for _, first := range []struct {
		name   string
		status int
		body   string
	}{
		{"edit lock conflict", http.StatusConflict, editLockBody},
		{"gateway error", http.StatusBadGateway, "<html>502</html>"},
	} {
		t.Run(first.name, func(t *testing.T) {
			var attempts int32
			srv := httptest.NewServer(deletedOnFirstAttempt(&attempts, first.status, first.body, true))
			defer srv.Close()
			c := newDeleteRetryTestClient(t, srv)

			err := c.Delete(context.Background(), "/zia/api/v1/firewallFilteringRules/1495551")

			require.NoError(t, err, "the rule was deleted by the first attempt")
			assert.Equal(t, int32(2), atomic.LoadInt32(&attempts))
		})
	}
}

func TestExecuteRequest_DeleteOfMissingObjectStillFails(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(deletedOnFirstAttempt(&attempts, 0, "", false))
	defer srv.Close()
	c := newDeleteRetryTestClient(t, srv)

	err := c.Delete(context.Background(), "/zia/api/v1/firewallFilteringRules/1495551")

	require.Error(t, err, "a 404 on the first attempt means the rule did not exist")
	respErr, ok := errorx.AsErrorResponse(err)
	require.True(t, ok)
	require.NotNil(t, respErr.Parsed)
	assert.Equal(t, "RESOURCE_NOT_FOUND", respErr.Parsed.Code)
	assert.Equal(t, int32(1), atomic.LoadInt32(&attempts))
}

// Only a DELETE is reported as done: a 404 after a resent GET is still an
// error.
func TestExecuteRequest_ResentGetKeeps404(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(deletedOnFirstAttempt(&attempts, http.StatusConflict, editLockBody, true))
	defer srv.Close()
	c := newDeleteRetryTestClient(t, srv)

	_, _, _, err := c.ExecuteRequest(context.Background(), http.MethodGet, "/zia/api/v1/firewallFilteringRules/1495551", nil, nil, contentTypeJSON)

	require.Error(t, err)
	assert.Equal(t, int32(2), atomic.LoadInt32(&attempts))
}
