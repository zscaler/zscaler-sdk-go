package zpa

import (
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

// deleteRetryServer serves an object that the first DELETE deletes although a
// gateway answers it with 502; every later DELETE gets 404. With
// applied=false the object never existed and every DELETE gets 404.
func deleteRetryServer(attempts *int32, applied bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(attempts, 1)
		w.Header().Set("Content-Type", "application/json")
		if applied && n == 1 {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("<html>502</html>"))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":"RESOURCE_NOT_FOUND","message":"Object does not exist."}`))
	}))
}

func newDeleteRetryTestClient() *Client {
	cfg := &Configuration{}
	cfg.Logger = logger.GetDefaultLogger("zpa-test: ")
	cfg.ZPA.Client.RateLimit.MaxRetries = 5
	cfg.ZPA.Client.RateLimit.RetryWaitMin = time.Millisecond
	cfg.ZPA.Client.RateLimit.RetryWaitMax = 2 * time.Millisecond
	cfg.HTTPClient = getHTTPClient(cfg.Logger, nil, cfg)
	return &Client{Config: cfg}
}

func deleteRequest(t *testing.T, srv *httptest.Server) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodDelete, srv.URL+"/mgmtconfig/v1/admin/customers/1/segmentGroup/1", nil)
	require.NoError(t, err)
	return req
}

// A DELETE that the API applied although the client saw a retryable error is
// retried, and the retry gets 404. The object is deleted (SUP-4290).
func TestZPADelete_AppliedBeforeRetryIsDeleted(t *testing.T) {
	var attempts int32
	srv := deleteRetryServer(&attempts, true)
	defer srv.Close()

	_, err := newDeleteRetryTestClient().do(deleteRequest(t, srv), nil, time.Now(), "test")

	require.NoError(t, err)
	assert.Equal(t, int32(2), atomic.LoadInt32(&attempts))
}

func TestZPADelete_MissingObjectStillFails(t *testing.T) {
	var attempts int32
	srv := deleteRetryServer(&attempts, false)
	defer srv.Close()

	_, err := newDeleteRetryTestClient().do(deleteRequest(t, srv), nil, time.Now(), "test")

	require.Error(t, err, "a 404 on the first attempt means the object did not exist")
	respErr, ok := errorx.AsErrorResponse(err)
	require.True(t, ok)
	require.NotNil(t, respErr.Parsed)
	assert.Equal(t, "RESOURCE_NOT_FOUND", respErr.Parsed.Code)
	assert.Equal(t, int32(1), atomic.LoadInt32(&attempts))
}
