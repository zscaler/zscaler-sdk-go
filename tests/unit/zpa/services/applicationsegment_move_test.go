// Package unit provides unit tests for ZPA App Segment Move service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment_move"
)

func TestAppSegmentMove_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AppSegmentMicrotenantMoveRequest JSON marshaling", func(t *testing.T) {
		req := applicationsegment_move.AppSegmentMicrotenantMoveRequest{
			ApplicationID:        "app-123",
			TargetMicrotenantID:  "mt-target",
			TargetSegmentGroupID: "sg-target",
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var unmarshaled applicationsegment_move.AppSegmentMicrotenantMoveRequest
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, req.ApplicationID, unmarshaled.ApplicationID)
		assert.Equal(t, req.TargetMicrotenantID, unmarshaled.TargetMicrotenantID)
	})
}

func TestAppSegmentMove_MockServerOperations(t *testing.T) {
	t.Run("POST move app segment", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/move", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}
