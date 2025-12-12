// Package unit provides unit tests for ZPA Application Segment service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
)

// TestApplicationSegment_Structure tests the struct definitions
func TestApplicationSegment_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ApplicationSegmentResource JSON marshaling", func(t *testing.T) {
		appSeg := applicationsegment.ApplicationSegmentResource{
			ID:              "app-123",
			Name:            "Test App Segment",
			Description:     "Test Description",
			Enabled:         true,
			DomainNames:     []string{"app.example.com", "www.app.example.com"},
			SegmentGroupID:  "sg-001",
			DoubleEncrypt:   false,
			IpAnchored:      true,
			HealthCheckType: "DEFAULT",
			BypassType:      "NEVER",
		}

		data, err := json.Marshal(appSeg)
		require.NoError(t, err)

		var unmarshaled applicationsegment.ApplicationSegmentResource
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, appSeg.ID, unmarshaled.ID)
		assert.Equal(t, appSeg.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
		assert.Len(t, unmarshaled.DomainNames, 2)
	})

	t.Run("ApplicationSegmentResource from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "app-456",
			"name": "Production App",
			"description": "Production application segment",
			"enabled": true,
			"domainNames": ["prod.example.com"],
			"segmentGroupId": "sg-002",
			"segmentGroupName": "Production Group",
			"doubleEncrypt": true,
			"ipAnchored": false,
			"healthCheckType": "DEFAULT",
			"bypassType": "NEVER",
			"passiveHealthEnabled": true,
			"selectConnectorCloseToApp": true,
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var appSeg applicationsegment.ApplicationSegmentResource
		err := json.Unmarshal([]byte(apiResponse), &appSeg)
		require.NoError(t, err)

		assert.Equal(t, "app-456", appSeg.ID)
		assert.True(t, appSeg.Enabled)
		assert.True(t, appSeg.DoubleEncrypt)
	})
}

// TestApplicationSegment_MockServerOperations tests CRUD operations
func TestApplicationSegment_MockServerOperations(t *testing.T) {
	t.Run("GET application segment by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "app-123", "name": "Mock App", "enabled": true}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/application/app-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create application segment", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "new-app", "name": "New App"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/application", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("DELETE application segment", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/application/app-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestApplicationSegment_SpecialCases tests edge cases
func TestApplicationSegment_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Bypass types", func(t *testing.T) {
		types := []string{"NEVER", "ALWAYS", "ON_NET"}

		for _, bypassType := range types {
			appSeg := applicationsegment.ApplicationSegmentResource{
				ID:         "app-" + bypassType,
				Name:       bypassType + " App",
				BypassType: bypassType,
			}

			data, err := json.Marshal(appSeg)
			require.NoError(t, err)
			assert.Contains(t, string(data), bypassType)
		}
	})

	t.Run("Health check types", func(t *testing.T) {
		types := []string{"DEFAULT", "NONE"}

		for _, hcType := range types {
			appSeg := applicationsegment.ApplicationSegmentResource{
				ID:              "app-" + hcType,
				Name:            hcType + " App",
				HealthCheckType: hcType,
			}

			data, err := json.Marshal(appSeg)
			require.NoError(t, err)
			assert.Contains(t, string(data), hcType)
		}
	})
}
