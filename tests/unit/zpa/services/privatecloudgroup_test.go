// Package unit provides unit tests for ZPA Private Cloud Group service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PrivateCloudGroup represents the private cloud group for testing
type PrivateCloudGroup struct {
	ID                     string `json:"id,omitempty"`
	Name                   string `json:"name,omitempty"`
	Description            string `json:"description,omitempty"`
	Enabled                bool   `json:"enabled,omitempty"`
	City                   string `json:"city,omitempty"`
	CityCountry            string `json:"cityCountry,omitempty"`
	CountryCode            string `json:"countryCode,omitempty"`
	Latitude               string `json:"latitude,omitempty"`
	Longitude              string `json:"longitude,omitempty"`
	Location               string `json:"location,omitempty"`
	OverrideVersionProfile bool   `json:"overrideVersionProfile,omitempty"`
	VersionProfileID       string `json:"versionProfileId,omitempty"`
	VersionProfileName     string `json:"versionProfileName,omitempty"`
	UpgradeDay             string `json:"upgradeDay,omitempty"`
	UpgradeTimeInSecs      string `json:"upgradeTimeInSecs,omitempty"`
	MicrotenantID          string `json:"microtenantId,omitempty"`
	MicrotenantName        string `json:"microtenantName,omitempty"`
	CreationTime           string `json:"creationTime,omitempty"`
	ModifiedTime           string `json:"modifiedTime,omitempty"`
}

// TestPrivateCloudGroup_Structure tests the struct definitions
func TestPrivateCloudGroup_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PrivateCloudGroup JSON marshaling", func(t *testing.T) {
		group := PrivateCloudGroup{
			ID:                     "pcg-123",
			Name:                   "US West Group",
			Description:            "Private cloud group for US West region",
			Enabled:                true,
			City:                   "San Jose",
			CountryCode:            "US",
			Latitude:               "37.3382",
			Longitude:              "-121.8863",
			Location:               "San Jose, CA, USA",
			OverrideVersionProfile: true,
			UpgradeDay:             "SUNDAY",
			UpgradeTimeInSecs:      "7200",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled PrivateCloudGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, group.ID, unmarshaled.ID)
		assert.True(t, unmarshaled.Enabled)
		assert.Equal(t, "SUNDAY", unmarshaled.UpgradeDay)
	})

	t.Run("PrivateCloudGroup from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "pcg-456",
			"name": "EU Central Group",
			"description": "Private cloud group for EU",
			"enabled": true,
			"city": "Frankfurt",
			"cityCountry": "Frankfurt, Germany",
			"countryCode": "DE",
			"latitude": "50.1109",
			"longitude": "8.6821",
			"location": "Frankfurt, Germany",
			"overrideVersionProfile": false,
			"versionProfileId": "vp-001",
			"versionProfileName": "Default Profile",
			"upgradeDay": "SATURDAY",
			"upgradeTimeInSecs": "10800",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var group PrivateCloudGroup
		err := json.Unmarshal([]byte(apiResponse), &group)
		require.NoError(t, err)

		assert.Equal(t, "pcg-456", group.ID)
		assert.Equal(t, "DE", group.CountryCode)
		assert.Equal(t, "SATURDAY", group.UpgradeDay)
	})
}

// TestPrivateCloudGroup_MockServerOperations tests CRUD operations
func TestPrivateCloudGroup_MockServerOperations(t *testing.T) {
	t.Run("GET group by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "pcg-123", "name": "Mock Group"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/privateCloudControllerGroup/pcg-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET group summary", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/summary")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"list": [{"id": "1"}, {"id": "2"}], "totalPages": 1}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/privateCloudControllerGroup/summary")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create group", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "new-pcg", "name": "New Group"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/privateCloudControllerGroup", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("DELETE group", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/privateCloudControllerGroup/pcg-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

