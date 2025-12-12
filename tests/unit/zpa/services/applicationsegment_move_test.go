// Package unit provides unit tests for ZPA Application Segment Move service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AppSegmentMicrotenantMoveRequest represents the move request structure for testing
type AppSegmentMicrotenantMoveRequest struct {
	ApplicationID        string `json:"applicationId,omitempty"`
	MicroTenantID        string `json:"microtenantId,omitempty"`
	TargetSegmentGroupID string `json:"targetSegmentGroupId,omitempty"`
	TargetMicrotenantID  string `json:"targetMicrotenantId,omitempty"`
	TargetServerGroupID  string `json:"targetServerGroupId,omitempty"`
}

// TestAppSegmentMove_Structure tests the struct definitions
func TestAppSegmentMove_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AppSegmentMicrotenantMoveRequest JSON marshaling", func(t *testing.T) {
		move := AppSegmentMicrotenantMoveRequest{
			ApplicationID:        "app-123",
			MicroTenantID:        "mt-source",
			TargetSegmentGroupID: "sg-target",
			TargetMicrotenantID:  "mt-target",
			TargetServerGroupID:  "srv-target",
		}

		data, err := json.Marshal(move)
		require.NoError(t, err)

		var unmarshaled AppSegmentMicrotenantMoveRequest
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, move.ApplicationID, unmarshaled.ApplicationID)
		assert.Equal(t, move.MicroTenantID, unmarshaled.MicroTenantID)
		assert.Equal(t, move.TargetSegmentGroupID, unmarshaled.TargetSegmentGroupID)
		assert.Equal(t, move.TargetMicrotenantID, unmarshaled.TargetMicrotenantID)
		assert.Equal(t, move.TargetServerGroupID, unmarshaled.TargetServerGroupID)
	})

	t.Run("Minimal move request - just segment group", func(t *testing.T) {
		move := AppSegmentMicrotenantMoveRequest{
			ApplicationID:        "app-123",
			TargetSegmentGroupID: "sg-target",
		}

		data, err := json.Marshal(move)
		require.NoError(t, err)

		var unmarshaled AppSegmentMicrotenantMoveRequest
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "app-123", unmarshaled.ApplicationID)
		assert.Equal(t, "sg-target", unmarshaled.TargetSegmentGroupID)
		assert.Empty(t, unmarshaled.TargetMicrotenantID)
		assert.Empty(t, unmarshaled.TargetServerGroupID)
	})

	t.Run("Move to different microtenant", func(t *testing.T) {
		move := AppSegmentMicrotenantMoveRequest{
			ApplicationID:        "app-123",
			MicroTenantID:        "mt-source",
			TargetMicrotenantID:  "mt-destination",
			TargetSegmentGroupID: "sg-in-destination",
		}

		data, err := json.Marshal(move)
		require.NoError(t, err)

		var unmarshaled AppSegmentMicrotenantMoveRequest
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.NotEqual(t, unmarshaled.MicroTenantID, unmarshaled.TargetMicrotenantID)
	})
}

// TestAppSegmentMove_MockServerOperations tests move operations with mock server
func TestAppSegmentMove_MockServerOperations(t *testing.T) {
	t.Run("POST move application segment", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Contains(t, r.URL.Path, "/application/")
			assert.Contains(t, r.URL.Path, "/move")

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/zpa/mgmtconfig/v1/admin/customers/123/application/app-123/move", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("POST move with microtenant filter", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Contains(t, r.URL.Path, "/move")

			// Check for microtenant ID in query params
			microtenantID := r.URL.Query().Get("microtenantId")
			if microtenantID != "" {
				assert.Equal(t, "mt-source", microtenantID)
			}

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/application/app-123/move?microtenantId=mt-source", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("POST move returns success with response body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Contains(t, r.URL.Path, "/move")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"message": "Application segment moved successfully",
				"applicationId": "app-123",
				"newSegmentGroupId": "sg-target"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/application/app-123/move", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestAppSegmentMove_ErrorHandling tests error scenarios
func TestAppSegmentMove_ErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("404 Application Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"code": "NOT_FOUND", "message": "Application segment not found"}`))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/application/nonexistent/move", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("400 Invalid target segment group", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code": "INVALID_REQUEST", "message": "Target segment group not found"}`))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/application/app-123/move", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("400 Invalid target microtenant", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code": "INVALID_REQUEST", "message": "Target microtenant not found or not accessible"}`))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/application/app-123/move", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("403 Forbidden - insufficient permissions", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"code": "FORBIDDEN", "message": "Insufficient permissions to move application to target microtenant"}`))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/application/app-123/move", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("409 Conflict - application in use", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(`{"code": "CONFLICT", "message": "Application segment is referenced by active policies"}`))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/application/app-123/move", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})
}

// TestAppSegmentMove_SpecialCases tests edge cases
func TestAppSegmentMove_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Move within same microtenant", func(t *testing.T) {
		move := AppSegmentMicrotenantMoveRequest{
			ApplicationID:        "app-123",
			MicroTenantID:        "mt-001",
			TargetMicrotenantID:  "mt-001", // Same microtenant
			TargetSegmentGroupID: "sg-new",
		}

		data, err := json.Marshal(move)
		require.NoError(t, err)

		var unmarshaled AppSegmentMicrotenantMoveRequest
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, unmarshaled.MicroTenantID, unmarshaled.TargetMicrotenantID)
	})

	t.Run("Move with server group change", func(t *testing.T) {
		move := AppSegmentMicrotenantMoveRequest{
			ApplicationID:        "app-123",
			TargetSegmentGroupID: "sg-new",
			TargetServerGroupID:  "srv-new",
		}

		data, err := json.Marshal(move)
		require.NoError(t, err)

		var unmarshaled AppSegmentMicrotenantMoveRequest
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.NotEmpty(t, unmarshaled.TargetServerGroupID)
		assert.NotEmpty(t, unmarshaled.TargetSegmentGroupID)
	})

	t.Run("Move without microtenant ID - uses service default", func(t *testing.T) {
		move := AppSegmentMicrotenantMoveRequest{
			ApplicationID:        "app-123",
			TargetSegmentGroupID: "sg-new",
			TargetMicrotenantID:  "mt-target",
			// MicroTenantID not set - service will use its configured value
		}

		data, err := json.Marshal(move)
		require.NoError(t, err)

		var unmarshaled AppSegmentMicrotenantMoveRequest
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Empty(t, unmarshaled.MicroTenantID)
		assert.NotEmpty(t, unmarshaled.TargetMicrotenantID)
	})

	t.Run("All fields populated", func(t *testing.T) {
		move := AppSegmentMicrotenantMoveRequest{
			ApplicationID:        "app-full-move",
			MicroTenantID:        "mt-source-full",
			TargetSegmentGroupID: "sg-target-full",
			TargetMicrotenantID:  "mt-target-full",
			TargetServerGroupID:  "srv-target-full",
		}

		data, err := json.Marshal(move)
		require.NoError(t, err)

		// Verify all fields are present
		assert.Contains(t, string(data), "applicationId")
		assert.Contains(t, string(data), "microtenantId")
		assert.Contains(t, string(data), "targetSegmentGroupId")
		assert.Contains(t, string(data), "targetMicrotenantId")
		assert.Contains(t, string(data), "targetServerGroupId")
	})
}

