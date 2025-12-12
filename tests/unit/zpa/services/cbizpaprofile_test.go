// Package unit provides unit tests for ZPA CBI ZPA Profile service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ZPAProfiles represents the ZPA profiles for CBI for testing
type ZPAProfiles struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	Enabled      bool   `json:"enabled"`
	CreationTime string `json:"creationTime,omitempty"`
	ModifiedBy   string `json:"modifiedBy,omitempty"`
	ModifiedTime string `json:"modifiedTime,omitempty"`
	CBITenantID  string `json:"cbiTenantId,omitempty"`
	CBIProfileID string `json:"cbiProfileId,omitempty"`
	CBIURL       string `json:"cbiUrl"`
}

// TestZPAProfiles_Structure tests the struct definitions
func TestZPAProfiles_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ZPAProfiles JSON marshaling", func(t *testing.T) {
		profile := ZPAProfiles{
			ID:           "zpa-profile-123",
			Name:         "ZPA CBI Integration Profile",
			Description:  "Profile for ZPA and CBI integration",
			Enabled:      true,
			CBITenantID:  "tenant-001",
			CBIProfileID: "cbi-001",
			CBIURL:       "https://cbi.zscaler.com/isolation",
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		var unmarshaled ZPAProfiles
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, profile.ID, unmarshaled.ID)
		assert.Equal(t, profile.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
		assert.Equal(t, profile.CBIURL, unmarshaled.CBIURL)
	})

	t.Run("ZPAProfiles from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "zpa-profile-456",
			"name": "Production CBI Profile",
			"description": "Production integration profile",
			"enabled": true,
			"creationTime": "1609459200000",
			"modifiedBy": "admin@example.com",
			"modifiedTime": "1612137600000",
			"cbiTenantId": "tenant-002",
			"cbiProfileId": "cbi-002",
			"cbiUrl": "https://production-cbi.zscaler.com"
		}`

		var profile ZPAProfiles
		err := json.Unmarshal([]byte(apiResponse), &profile)
		require.NoError(t, err)

		assert.Equal(t, "zpa-profile-456", profile.ID)
		assert.Equal(t, "Production CBI Profile", profile.Name)
		assert.True(t, profile.Enabled)
		assert.Equal(t, "tenant-002", profile.CBITenantID)
	})
}

// TestZPAProfiles_ResponseParsing tests parsing of API responses
func TestZPAProfiles_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse ZPA profiles list response", func(t *testing.T) {
		response := `[
			{"id": "1", "name": "Profile 1", "enabled": true, "cbiUrl": "https://cbi1.zscaler.com"},
			{"id": "2", "name": "Profile 2", "enabled": true, "cbiUrl": "https://cbi2.zscaler.com"},
			{"id": "3", "name": "Profile 3", "enabled": false, "cbiUrl": "https://cbi3.zscaler.com"}
		]`

		var profiles []ZPAProfiles
		err := json.Unmarshal([]byte(response), &profiles)
		require.NoError(t, err)

		assert.Len(t, profiles, 3)
		assert.True(t, profiles[0].Enabled)
		assert.False(t, profiles[2].Enabled)
	})
}

// TestZPAProfiles_MockServerOperations tests operations
func TestZPAProfiles_MockServerOperations(t *testing.T) {
	t.Run("GET ZPA profile by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/zpaprofiles")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `[
				{"id": "profile-123", "name": "Mock ZPA Profile", "enabled": true}
			]`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cbi/api/customers/123/zpaprofiles")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all ZPA profiles", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `[
				{"id": "1", "name": "ZPA Profile A", "enabled": true},
				{"id": "2", "name": "ZPA Profile B", "enabled": false}
			]`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cbi/api/customers/123/zpaprofiles")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestZPAProfiles_SpecialCases tests edge cases
func TestZPAProfiles_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Disabled ZPA profile", func(t *testing.T) {
		profile := ZPAProfiles{
			ID:      "profile-disabled",
			Name:    "Disabled Profile",
			Enabled: false,
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		var unmarshaled ZPAProfiles
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.False(t, unmarshaled.Enabled)
	})

	t.Run("CBI URL formats", func(t *testing.T) {
		urls := []string{
			"https://cbi.zscaler.com",
			"https://cbi.zscalerone.net",
			"https://cbi.zscalertwo.net",
			"https://cbi.zscloud.net",
		}

		for _, url := range urls {
			profile := ZPAProfiles{
				ID:     "profile-url",
				Name:   "URL Test",
				CBIURL: url,
			}

			data, err := json.Marshal(profile)
			require.NoError(t, err)

			assert.Contains(t, string(data), url)
		}
	})
}

