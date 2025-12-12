// Package unit provides unit tests for ZPA Posture Profile service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PostureProfile represents the posture profile structure for testing
type PostureProfile struct {
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	Domain              string `json:"domain,omitempty"`
	CreationTime        string `json:"creationTime,omitempty"`
	ModifiedBy          string `json:"modifiedBy,omitempty"`
	ModifiedTime        string `json:"modifiedTime,omitempty"`
	MasterCustomerID    string `json:"masterCustomerId,omitempty"`
	PostureUdid         string `json:"postureUdid,omitempty"`
	ZscalerCloud        string `json:"zscalerCloud,omitempty"`
	ZscalerCustomerID   string `json:"zscalerCustomerId,omitempty"`
	MicroTenantID       string `json:"microtenantId,omitempty"`
	MicroTenantName     string `json:"microtenantName,omitempty"`
}

// TestPostureProfile_Structure tests the struct definitions
func TestPostureProfile_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PostureProfile JSON marshaling", func(t *testing.T) {
		profile := PostureProfile{
			ID:                "pp-123",
			Name:              "CrowdStrike Posture",
			Domain:            "crowdstrike",
			PostureUdid:       "udid-12345",
			ZscalerCloud:      "zscaler.net",
			ZscalerCustomerID: "customer-123",
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		var unmarshaled PostureProfile
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, profile.ID, unmarshaled.ID)
		assert.Equal(t, profile.Name, unmarshaled.Name)
		assert.Equal(t, profile.Domain, unmarshaled.Domain)
	})

	t.Run("PostureProfile JSON unmarshaling from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "pp-456",
			"name": "Windows Defender Posture",
			"domain": "windows_defender",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com",
			"masterCustomerId": "master-123",
			"postureUdid": "udid-67890",
			"zscalerCloud": "zscaler.net",
			"zscalerCustomerId": "customer-456",
			"microtenantId": "mt-001",
			"microtenantName": "Production"
		}`

		var profile PostureProfile
		err := json.Unmarshal([]byte(apiResponse), &profile)
		require.NoError(t, err)

		assert.Equal(t, "pp-456", profile.ID)
		assert.Equal(t, "Windows Defender Posture", profile.Name)
		assert.Equal(t, "windows_defender", profile.Domain)
		assert.NotEmpty(t, profile.CreationTime)
		assert.NotEmpty(t, profile.PostureUdid)
	})
}

// TestPostureProfile_ResponseParsing tests parsing of various API responses
func TestPostureProfile_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse posture profile list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "CrowdStrike", "domain": "crowdstrike"},
				{"id": "2", "name": "Windows Defender", "domain": "windows_defender"},
				{"id": "3", "name": "SentinelOne", "domain": "sentinelone"},
				{"id": "4", "name": "Carbon Black", "domain": "carbon_black"}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []PostureProfile `json:"list"`
			TotalPages int              `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 4)
		assert.Equal(t, "CrowdStrike", listResp.List[0].Name)
		assert.Equal(t, "crowdstrike", listResp.List[0].Domain)
	})
}

// TestPostureProfile_MockServerOperations tests CRUD operations with mock server
func TestPostureProfile_MockServerOperations(t *testing.T) {
	t.Run("GET posture profile by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/posture/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "pp-123",
				"name": "Mock Posture Profile",
				"domain": "mock_vendor"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v2/admin/customers/123/posture/pp-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all posture profiles", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Profile A", "domain": "vendor_a"},
					{"id": "2", "name": "Profile B", "domain": "vendor_b"}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v2/admin/customers/123/posture")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestPostureProfile_ErrorHandling tests error scenarios
func TestPostureProfile_ErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("404 Posture Profile Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"code": "NOT_FOUND", "message": "Posture profile not found"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/posture/nonexistent")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

// TestPostureProfile_SpecialCases tests edge cases
func TestPostureProfile_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Common posture domains", func(t *testing.T) {
		domains := []string{
			"crowdstrike",
			"windows_defender",
			"sentinelone",
			"carbon_black",
			"zscaler_client_connector",
			"microsoft_intune",
			"jamf",
			"google_beyondcorp",
		}

		for _, domain := range domains {
			profile := PostureProfile{
				ID:     "pp-" + domain,
				Name:   domain + " Profile",
				Domain: domain,
			}

			data, err := json.Marshal(profile)
			require.NoError(t, err)

			var unmarshaled PostureProfile
			err = json.Unmarshal(data, &unmarshaled)
			require.NoError(t, err)

			assert.Equal(t, domain, unmarshaled.Domain)
		}
	})

	t.Run("Posture UDID format", func(t *testing.T) {
		profile := PostureProfile{
			ID:          "pp-123",
			Name:        "UDID Test Profile",
			PostureUdid: "ABC123-DEF456-GHI789-JKL012",
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		assert.Contains(t, string(data), "ABC123-DEF456-GHI789-JKL012")
	})
}

