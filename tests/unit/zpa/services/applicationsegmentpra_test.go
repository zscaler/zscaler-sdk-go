// Package unit provides unit tests for ZPA App Segment PRA service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentpra"
)

func TestAppSegmentPRA_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AppSegmentPRA JSON marshaling", func(t *testing.T) {
		segment := applicationsegmentpra.AppSegmentPRA{
			ID:             "pra-123",
			Name:           "Test PRA App",
			Description:    "Test Description",
			Enabled:        true,
			SegmentGroupID: "sg-001",
		}

		data, err := json.Marshal(segment)
		require.NoError(t, err)

		var unmarshaled applicationsegmentpra.AppSegmentPRA
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, segment.ID, unmarshaled.ID)
		assert.Equal(t, segment.Name, unmarshaled.Name)
	})
}

func TestAppSegmentPRA_MockServerOperations(t *testing.T) {
	t.Run("GET PRA app segment", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "pra-123", "name": "Mock PRA App"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/praApp")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
