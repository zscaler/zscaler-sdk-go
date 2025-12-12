// Package unit provides unit tests for ZPA Isolation Profile service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// IsolationProfile represents the isolation profile for testing
type IsolationProfile struct {
	ID                 string `json:"id,omitempty"`
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	Enabled            bool   `json:"enabled"`
	IsolationProfileID string `json:"isolationProfileId,omitempty"`
	IsolationTenantID  string `json:"isolationTenantId,omitempty"`
	IsolationURL       string `json:"isolationUrl"`
	CreationTime       string `json:"creationTime,omitempty"`
	ModifiedBy         string `json:"modifiedBy,omitempty"`
	ModifiedTime       string `json:"modifiedTime,omitempty"`
}

// TestIsolationProfile_Structure tests the struct definitions
func TestIsolationProfile_Structure(t *testing.T) {
	t.Parallel()

	t.Run("IsolationProfile JSON marshaling", func(t *testing.T) {
		profile := IsolationProfile{
			ID:                 "iso-123",
			Name:               "Corporate Isolation Profile",
			Description:        "Main isolation profile for corporate users",
			Enabled:            true,
			IsolationProfileID: "iso-profile-001",
			IsolationTenantID:  "iso-tenant-001",
			IsolationURL:       "https://isolation.zscaler.com",
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		var unmarshaled IsolationProfile
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, profile.ID, unmarshaled.ID)
		assert.Equal(t, profile.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
		assert.Equal(t, profile.IsolationURL, unmarshaled.IsolationURL)
	})

	t.Run("IsolationProfile from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "iso-456",
			"name": "Guest Isolation Profile",
			"description": "Isolation profile for guest users",
			"enabled": true,
			"isolationProfileId": "iso-profile-002",
			"isolationTenantId": "iso-tenant-002",
			"isolationUrl": "https://guest-isolation.zscaler.com",
			"creationTime": "1609459200000",
			"modifiedBy": "admin@example.com",
			"modifiedTime": "1612137600000"
		}`

		var profile IsolationProfile
		err := json.Unmarshal([]byte(apiResponse), &profile)
		require.NoError(t, err)

		assert.Equal(t, "iso-456", profile.ID)
		assert.Equal(t, "Guest Isolation Profile", profile.Name)
		assert.True(t, profile.Enabled)
		assert.Equal(t, "iso-tenant-002", profile.IsolationTenantID)
	})
}

// TestIsolationProfile_ResponseParsing tests parsing of API responses
func TestIsolationProfile_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse isolation profile list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "Profile 1", "enabled": true},
				{"id": "2", "name": "Profile 2", "enabled": true},
				{"id": "3", "name": "Profile 3", "enabled": false}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []IsolationProfile `json:"list"`
			TotalPages int                `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.True(t, listResp.List[0].Enabled)
		assert.False(t, listResp.List[2].Enabled)
	})
}

// TestIsolationProfile_MockServerOperations tests operations
func TestIsolationProfile_MockServerOperations(t *testing.T) {
	t.Run("GET all isolation profiles", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/isolation/profiles")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Isolation A", "enabled": true},
					{"id": "2", "name": "Isolation B", "enabled": true}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/isolation/profiles")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET isolation profile by name (search)", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			search := r.URL.Query().Get("search")
			assert.NotEmpty(t, search)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Corporate Isolation", "enabled": true}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/isolation/profiles?search=Corporate")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestIsolationProfile_SpecialCases tests edge cases
func TestIsolationProfile_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Disabled isolation profile", func(t *testing.T) {
		profile := IsolationProfile{
			ID:      "iso-disabled",
			Name:    "Disabled Isolation",
			Enabled: false,
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		var unmarshaled IsolationProfile
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.False(t, unmarshaled.Enabled)
	})

	t.Run("Isolation URL variations", func(t *testing.T) {
		urls := []string{
			"https://isolation.zscaler.com",
			"https://isolation.zscalerone.net",
			"https://isolation.zscalertwo.net",
			"https://isolation.zscloud.net",
		}

		for _, url := range urls {
			profile := IsolationProfile{
				ID:           "iso-url",
				Name:         "URL Test",
				IsolationURL: url,
			}

			data, err := json.Marshal(profile)
			require.NoError(t, err)

			assert.Contains(t, string(data), url)
		}
	})
}

