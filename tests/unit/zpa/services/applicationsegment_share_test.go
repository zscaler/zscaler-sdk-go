// Package unit provides unit tests for ZPA App Segment Share service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment_share"
)

func TestAppSegmentShare_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AppSegmentSharedToMicrotenant JSON marshaling", func(t *testing.T) {
		req := applicationsegment_share.AppSegmentSharedToMicrotenant{
			ApplicationID:       "app-123",
			ShareToMicrotenants: []string{"mt-001", "mt-002"},
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var unmarshaled applicationsegment_share.AppSegmentSharedToMicrotenant
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, req.ApplicationID, unmarshaled.ApplicationID)
		assert.Len(t, unmarshaled.ShareToMicrotenants, 2)
	})
}

func TestAppSegmentShare_MockServerOperations(t *testing.T) {
	t.Run("POST share app segment", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/share", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}
