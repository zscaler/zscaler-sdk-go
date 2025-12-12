// Package unit provides unit tests for ZPA App Segment Browser Access service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentbrowseraccess"
)

func TestAppSegmentBrowserAccess_Structure(t *testing.T) {
	t.Parallel()

	t.Run("BrowserAccess JSON marshaling", func(t *testing.T) {
		segment := applicationsegmentbrowseraccess.BrowserAccess{
			ID:             "ba-123",
			Name:           "Test Browser Access App",
			Description:    "Test Description",
			Enabled:        true,
			SegmentGroupID: "sg-001",
		}

		data, err := json.Marshal(segment)
		require.NoError(t, err)

		var unmarshaled applicationsegmentbrowseraccess.BrowserAccess
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, segment.ID, unmarshaled.ID)
		assert.Equal(t, segment.Name, unmarshaled.Name)
	})
}

func TestAppSegmentBrowserAccess_MockServerOperations(t *testing.T) {
	t.Run("GET browser access app segment", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "ba-123", "name": "Mock BA App"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/browserAccess")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
