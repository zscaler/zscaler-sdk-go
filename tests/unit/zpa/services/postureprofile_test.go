// Package unit provides unit tests for ZPA Posture Profile service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/postureprofile"
)

// TestPostureProfile_Structure tests the struct definitions
func TestPostureProfile_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PostureProfile JSON marshaling", func(t *testing.T) {
		profile := postureprofile.PostureProfile{
			ID:                          "pp-123",
			Name:                        "Test Posture Profile",
			Platform:                    "windows",
			PostureType:                 "DOMAIN_JOINED",
			ApplyToMachineTunnelEnabled: true,
			CRLCheckEnabled:             true,
			Domain:                      "example.com",
			ZscalerCloud:                "zscaler.net",
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		var unmarshaled postureprofile.PostureProfile
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, profile.ID, unmarshaled.ID)
		assert.Equal(t, profile.Name, unmarshaled.Name)
		assert.Equal(t, "windows", unmarshaled.Platform)
	})

	t.Run("PostureProfile from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "pp-456",
			"name": "Production Posture",
			"platform": "mac",
			"postureType": "CERTIFICATE",
			"postureUdid": "posture-udid-123",
			"applyToMachineTunnelEnabled": false,
			"crlCheckEnabled": true,
			"nonExportablePrivateKeyEnabled": true,
			"domain": "corp.example.com",
			"rootCert": "cert-data",
			"zscalerCloud": "zscaler.net",
			"zscalerCustomerId": "cust-123",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var profile postureprofile.PostureProfile
		err := json.Unmarshal([]byte(apiResponse), &profile)
		require.NoError(t, err)

		assert.Equal(t, "pp-456", profile.ID)
		assert.Equal(t, "mac", profile.Platform)
		assert.Equal(t, "CERTIFICATE", profile.PostureType)
	})
}

// TestPostureProfile_MockServerOperations tests operations
func TestPostureProfile_MockServerOperations(t *testing.T) {
	t.Run("GET posture profile by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "pp-123", "name": "Mock Posture", "platform": "windows"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/posture/pp-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all posture profiles", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"list": [{"id": "1"}, {"id": "2"}], "totalPages": 1}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/posture")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestPostureProfile_SpecialCases tests edge cases
func TestPostureProfile_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Platform types", func(t *testing.T) {
		platforms := []string{"windows", "mac", "linux", "ios", "android"}

		for _, platform := range platforms {
			profile := postureprofile.PostureProfile{
				ID:       "pp-" + platform,
				Name:     platform + " Posture",
				Platform: platform,
			}

			data, err := json.Marshal(profile)
			require.NoError(t, err)
			assert.Contains(t, string(data), platform)
		}
	})

	t.Run("Posture types", func(t *testing.T) {
		types := []string{"DOMAIN_JOINED", "CERTIFICATE", "CROWDSTRIKE", "SENTINEL_ONE"}

		for _, pType := range types {
			profile := postureprofile.PostureProfile{
				ID:          "pp-" + pType,
				Name:        pType + " Posture",
				PostureType: pType,
			}

			data, err := json.Marshal(profile)
			require.NoError(t, err)
			assert.Contains(t, string(data), pType)
		}
	})
}
