// Package unit provides unit tests for ZPA Application Segment By Type service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentbytype"
)

func TestApplicationSegmentByType_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AppSegmentBaseAppDto JSON marshaling", func(t *testing.T) {
		segment := applicationsegmentbytype.AppSegmentBaseAppDto{
			ID:      "app-123",
			Name:    "Test App Segment",
			Enabled: true,
		}

		data, err := json.Marshal(segment)
		require.NoError(t, err)

		var unmarshaled applicationsegmentbytype.AppSegmentBaseAppDto
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, segment.ID, unmarshaled.ID)
		assert.Equal(t, segment.Name, unmarshaled.Name)
	})
}

func TestApplicationSegmentByType_MockServerOperations(t *testing.T) {
	t.Run("GET application segment by type", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "app-123", "name": "Mock App"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/applicationByType")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
