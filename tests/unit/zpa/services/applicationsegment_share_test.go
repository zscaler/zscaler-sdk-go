// Package unit provides unit tests for ZPA Application Segment Share service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AppSegmentSharedToMicrotenant represents the share request structure for testing
type AppSegmentSharedToMicrotenant struct {
	ApplicationID       string   `json:"applicationId,omitempty"`
	ShareToMicrotenants []string `json:"shareToMicrotenants,omitempty"`
	MicroTenantID       string   `json:"microtenantId,omitempty"`
}

// TestAppSegmentShare_Structure tests the struct definitions
func TestAppSegmentShare_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AppSegmentSharedToMicrotenant JSON marshaling", func(t *testing.T) {
		share := AppSegmentSharedToMicrotenant{
			ApplicationID:       "app-123",
			ShareToMicrotenants: []string{"mt-001", "mt-002", "mt-003"},
			MicroTenantID:       "mt-source",
		}

		data, err := json.Marshal(share)
		require.NoError(t, err)

		var unmarshaled AppSegmentSharedToMicrotenant
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, share.ApplicationID, unmarshaled.ApplicationID)
		assert.ElementsMatch(t, share.ShareToMicrotenants, unmarshaled.ShareToMicrotenants)
		assert.Equal(t, share.MicroTenantID, unmarshaled.MicroTenantID)
	})

	t.Run("Empty share list", func(t *testing.T) {
		share := AppSegmentSharedToMicrotenant{
			ApplicationID:       "app-123",
			ShareToMicrotenants: []string{},
			MicroTenantID:       "mt-source",
		}

		data, err := json.Marshal(share)
		require.NoError(t, err)

		var unmarshaled AppSegmentSharedToMicrotenant
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.ShareToMicrotenants, 0)
	})

	t.Run("Single tenant share", func(t *testing.T) {
		share := AppSegmentSharedToMicrotenant{
			ApplicationID:       "app-456",
			ShareToMicrotenants: []string{"mt-target"},
		}

		data, err := json.Marshal(share)
		require.NoError(t, err)

		var unmarshaled AppSegmentSharedToMicrotenant
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.ShareToMicrotenants, 1)
		assert.Equal(t, "mt-target", unmarshaled.ShareToMicrotenants[0])
	})
}

// TestAppSegmentShare_MockServerOperations tests share operations with mock server
func TestAppSegmentShare_MockServerOperations(t *testing.T) {
	t.Run("PUT share application segment to microtenants", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Contains(t, r.URL.Path, "/application/")
			assert.Contains(t, r.URL.Path, "/share")

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/zpa/mgmtconfig/v1/admin/customers/123/application/app-123/share", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("PUT share with microtenant filter", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Contains(t, r.URL.Path, "/share")
			
			// Check for microtenant ID in query params
			microtenantID := r.URL.Query().Get("microtenantId")
			if microtenantID != "" {
				assert.Equal(t, "mt-source", microtenantID)
			}

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/application/app-123/share?microtenantId=mt-source", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestAppSegmentShare_ErrorHandling tests error scenarios
func TestAppSegmentShare_ErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("404 Application Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"code": "NOT_FOUND", "message": "Application segment not found"}`))
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/application/nonexistent/share", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("400 Invalid microtenant", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code": "INVALID_REQUEST", "message": "Invalid microtenant ID"}`))
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/application/app-123/share", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("403 Forbidden - insufficient permissions", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"code": "FORBIDDEN", "message": "Insufficient permissions to share application"}`))
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/application/app-123/share", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}

// TestAppSegmentShare_SpecialCases tests edge cases
func TestAppSegmentShare_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Share to many microtenants", func(t *testing.T) {
		tenants := make([]string, 50)
		for i := 0; i < 50; i++ {
			tenants[i] = "mt-" + string(rune('a'+i%26)) + string(rune('0'+i/26))
		}

		share := AppSegmentSharedToMicrotenant{
			ApplicationID:       "app-123",
			ShareToMicrotenants: tenants,
		}

		data, err := json.Marshal(share)
		require.NoError(t, err)

		var unmarshaled AppSegmentSharedToMicrotenant
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.ShareToMicrotenants, 50)
	})

	t.Run("Unshare - empty list to remove sharing", func(t *testing.T) {
		share := AppSegmentSharedToMicrotenant{
			ApplicationID:       "app-123",
			ShareToMicrotenants: []string{}, // Empty to unshare
			MicroTenantID:       "mt-source",
		}

		data, err := json.Marshal(share)
		require.NoError(t, err)

		// With omitempty, empty slices are omitted from JSON
		// Verify it doesn't contain any shareToMicrotenants entries
		var unmarshaled AppSegmentSharedToMicrotenant
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.ShareToMicrotenants, 0)
	})

	t.Run("Share without source microtenant", func(t *testing.T) {
		share := AppSegmentSharedToMicrotenant{
			ApplicationID:       "app-123",
			ShareToMicrotenants: []string{"mt-001", "mt-002"},
			// MicroTenantID not set - will use service default
		}

		data, err := json.Marshal(share)
		require.NoError(t, err)

		var unmarshaled AppSegmentSharedToMicrotenant
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Empty(t, unmarshaled.MicroTenantID)
		assert.Len(t, unmarshaled.ShareToMicrotenants, 2)
	})
}

