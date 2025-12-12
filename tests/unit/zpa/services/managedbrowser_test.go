// Package unit provides unit tests for ZPA Managed Browser service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ManagedBrowserProfile represents the managed browser profile for testing
type ManagedBrowserProfile struct {
	ID                   string                    `json:"id,omitempty"`
	Name                 string                    `json:"name,omitempty"`
	Description          string                    `json:"description,omitempty"`
	BrowserType          string                    `json:"browserType,omitempty"`
	CustomerID           string                    `json:"customerId,omitempty"`
	MicrotenantID        string                    `json:"microtenantId,omitempty"`
	MicrotenantName      string                    `json:"microtenantName,omitempty"`
	ChromePostureProfile ManagedBrowserPosture     `json:"chromePostureProfile,omitempty"`
	CreationTime         string                    `json:"creationTime,omitempty"`
	ModifiedBy           string                    `json:"modifiedBy,omitempty"`
	ModifiedTime         string                    `json:"modifiedTime,omitempty"`
}

// ManagedBrowserPosture represents the posture profile for testing
type ManagedBrowserPosture struct {
	ID               string `json:"id,omitempty"`
	BrowserType      string `json:"browserType,omitempty"`
	CrowdStrikeAgent bool   `json:"crowdStrikeAgent,omitempty"`
	CreationTime     string `json:"creationTime,omitempty"`
	ModifiedBy       string `json:"modifiedBy,omitempty"`
	ModifiedTime     string `json:"modifiedTime,omitempty"`
}

// TestManagedBrowser_Structure tests the struct definitions
func TestManagedBrowser_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ManagedBrowserProfile JSON marshaling", func(t *testing.T) {
		profile := ManagedBrowserProfile{
			ID:          "mb-123",
			Name:        "Chrome Enterprise Profile",
			Description: "Managed Chrome browser profile",
			BrowserType: "CHROME",
			CustomerID:  "cust-001",
			ChromePostureProfile: ManagedBrowserPosture{
				ID:               "cp-001",
				BrowserType:      "CHROME",
				CrowdStrikeAgent: true,
			},
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		var unmarshaled ManagedBrowserProfile
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, profile.ID, unmarshaled.ID)
		assert.Equal(t, profile.Name, unmarshaled.Name)
		assert.Equal(t, "CHROME", unmarshaled.BrowserType)
		assert.True(t, unmarshaled.ChromePostureProfile.CrowdStrikeAgent)
	})

	t.Run("ManagedBrowserProfile from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "mb-456",
			"name": "Edge Enterprise Profile",
			"description": "Managed Edge browser profile",
			"browserType": "EDGE",
			"customerId": "cust-002",
			"microtenantId": "mt-001",
			"microtenantName": "Production",
			"chromePostureProfile": {
				"id": "cp-002",
				"browserType": "EDGE",
				"crowdStrikeAgent": false,
				"creationTime": "1609459200000",
				"modifiedTime": "1612137600000"
			},
			"creationTime": "1609459200000",
			"modifiedBy": "admin@example.com",
			"modifiedTime": "1612137600000"
		}`

		var profile ManagedBrowserProfile
		err := json.Unmarshal([]byte(apiResponse), &profile)
		require.NoError(t, err)

		assert.Equal(t, "mb-456", profile.ID)
		assert.Equal(t, "EDGE", profile.BrowserType)
		assert.Equal(t, "Edge Enterprise Profile", profile.Name)
		assert.False(t, profile.ChromePostureProfile.CrowdStrikeAgent)
	})
}

// TestManagedBrowser_ResponseParsing tests parsing of API responses
func TestManagedBrowser_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse managed browser list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "Chrome Profile", "browserType": "CHROME"},
				{"id": "2", "name": "Edge Profile", "browserType": "EDGE"},
				{"id": "3", "name": "Firefox Profile", "browserType": "FIREFOX"}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []ManagedBrowserProfile `json:"list"`
			TotalPages int                     `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.Equal(t, "CHROME", listResp.List[0].BrowserType)
	})
}

// TestManagedBrowser_MockServerOperations tests operations
func TestManagedBrowser_MockServerOperations(t *testing.T) {
	t.Run("GET all managed browser profiles", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/managedBrowserProfile/search")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Profile A", "browserType": "CHROME"},
					{"id": "2", "name": "Profile B", "browserType": "EDGE"}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/managedBrowserProfile/search")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestManagedBrowser_SpecialCases tests edge cases
func TestManagedBrowser_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Browser types", func(t *testing.T) {
		types := []string{"CHROME", "EDGE", "FIREFOX", "SAFARI"}

		for _, browserType := range types {
			profile := ManagedBrowserProfile{
				ID:          "mb-" + browserType,
				Name:        browserType + " Profile",
				BrowserType: browserType,
			}

			data, err := json.Marshal(profile)
			require.NoError(t, err)

			assert.Contains(t, string(data), browserType)
		}
	})

	t.Run("CrowdStrike agent enabled", func(t *testing.T) {
		posture := ManagedBrowserPosture{
			ID:               "cp-123",
			CrowdStrikeAgent: true,
		}

		data, err := json.Marshal(posture)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"crowdStrikeAgent":true`)
	})
}

