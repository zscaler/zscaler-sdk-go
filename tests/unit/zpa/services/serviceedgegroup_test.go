// Package unit provides unit tests for ZPA Service Edge Group service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgegroup"
)

// TestServiceEdgeGroup_Structure tests the struct definitions
func TestServiceEdgeGroup_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ServiceEdgeGroup JSON marshaling", func(t *testing.T) {
		group := serviceedgegroup.ServiceEdgeGroup{
			ID:                     "seg-123",
			Name:                   "Test Service Edge Group",
			Description:            "Test Description",
			Enabled:                true,
			CityCountry:            "San Jose, US",
			Latitude:               "37.3382",
			Longitude:              "-121.8863",
			Location:               "San Jose, CA",
			VersionProfileID:       "0",
			OverrideVersionProfile: true,
			UpgradeDay:             "SUNDAY",
			UpgradeTimeInSecs:      "66600",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled serviceedgegroup.ServiceEdgeGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, group.ID, unmarshaled.ID)
		assert.Equal(t, group.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
	})

	t.Run("ServiceEdgeGroup from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "seg-456",
			"name": "Production Service Edge Group",
			"description": "Production environment",
			"enabled": true,
			"cityCountry": "New York, US",
			"latitude": "40.7128",
			"longitude": "-74.0060",
			"location": "New York, NY",
			"versionProfileId": "1",
			"versionProfileName": "Default",
			"isPublic": "TRUE",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var group serviceedgegroup.ServiceEdgeGroup
		err := json.Unmarshal([]byte(apiResponse), &group)
		require.NoError(t, err)

		assert.Equal(t, "seg-456", group.ID)
		assert.True(t, group.Enabled)
	})
}

// TestServiceEdgeGroup_MockServerOperations tests CRUD operations
func TestServiceEdgeGroup_MockServerOperations(t *testing.T) {
	t.Run("GET service edge group by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "seg-123", "name": "Mock Group", "enabled": true}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/serviceEdgeGroup/seg-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create service edge group", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "new-seg", "name": "New Group"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/serviceEdgeGroup", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("DELETE service edge group", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/serviceEdgeGroup/seg-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}
