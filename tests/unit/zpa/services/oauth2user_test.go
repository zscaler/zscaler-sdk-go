// Package unit provides unit tests for ZPA OAuth2 User service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// OauthUser represents the OAuth user for testing
type OauthUser struct {
	ComponentGroupID     string   `json:"componentGroupId,omitempty"`
	ConfigCloudName      string   `json:"configCloudName,omitempty"`
	EnrollmentServer     string   `json:"enrollmentServer,omitempty"`
	NonceAssociationType string   `json:"nonceAssociationType,omitempty"`
	TenantID             string   `json:"tenantId,omitempty"`
	UserCodes            []string `json:"userCodes,omitempty"`
	ZcomponentID         string   `json:"zcomponentId,omitempty"`
}

// UserCodeStatusRequest represents the status request for testing
type UserCodeStatusRequest struct {
	UserCodes []string `json:"userCodes"`
}

// TestOAuth2User_Structure tests the struct definitions
func TestOAuth2User_Structure(t *testing.T) {
	t.Parallel()

	t.Run("OauthUser JSON marshaling", func(t *testing.T) {
		user := OauthUser{
			ComponentGroupID:     "cg-123",
			ConfigCloudName:      "zscaler.net",
			EnrollmentServer:     "enroll.zscaler.net",
			NonceAssociationType: "CONNECTOR_GRP",
			TenantID:             "tenant-001",
			UserCodes:            []string{"CODE1", "CODE2", "CODE3"},
			ZcomponentID:         "zc-001",
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		var unmarshaled OauthUser
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, user.ComponentGroupID, unmarshaled.ComponentGroupID)
		assert.Equal(t, user.TenantID, unmarshaled.TenantID)
		assert.Len(t, unmarshaled.UserCodes, 3)
	})

	t.Run("OauthUser from API response", func(t *testing.T) {
		apiResponse := `{
			"componentGroupId": "cg-456",
			"configCloudName": "zscalerone.net",
			"enrollmentServer": "enroll.zscalerone.net",
			"nonceAssociationType": "SERVICE_EDGE_GRP",
			"tenantId": "tenant-002",
			"userCodes": ["ABC123", "DEF456"],
			"zcomponentId": "zc-002"
		}`

		var user OauthUser
		err := json.Unmarshal([]byte(apiResponse), &user)
		require.NoError(t, err)

		assert.Equal(t, "cg-456", user.ComponentGroupID)
		assert.Equal(t, "SERVICE_EDGE_GRP", user.NonceAssociationType)
		assert.Len(t, user.UserCodes, 2)
	})

	t.Run("UserCodeStatusRequest JSON marshaling", func(t *testing.T) {
		req := UserCodeStatusRequest{
			UserCodes: []string{"CODE1", "CODE2"},
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var unmarshaled UserCodeStatusRequest
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.UserCodes, 2)
	})
}

// TestOAuth2User_MockServerOperations tests operations
func TestOAuth2User_MockServerOperations(t *testing.T) {
	t.Run("POST verify user codes", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Contains(t, r.URL.Path, "/usercodes")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"componentGroupId": "cg-123",
				"tenantId": "tenant-001",
				"userCodes": ["CODE1"]
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/CONNECTOR_GRP/usercodes", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST verify user code status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Contains(t, r.URL.Path, "/usercodes/status")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"componentGroupId": "cg-123",
				"tenantId": "tenant-001"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/SERVICE_EDGE_GRP/usercodes/status", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestOAuth2User_SpecialCases tests edge cases
func TestOAuth2User_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Association types", func(t *testing.T) {
		types := []string{
			"CONNECTOR_GRP",
			"NP_ASSISTANT_GRP",
			"SITE_CONTROLLER_GRP",
			"SERVICE_EDGE_GRP",
		}

		for _, assocType := range types {
			user := OauthUser{
				NonceAssociationType: assocType,
				TenantID:             "tenant-001",
			}

			data, err := json.Marshal(user)
			require.NoError(t, err)

			assert.Contains(t, string(data), assocType)
		}
	})

	t.Run("Multiple user codes", func(t *testing.T) {
		user := OauthUser{
			TenantID: "tenant-001",
			UserCodes: []string{
				"CODE1", "CODE2", "CODE3", "CODE4", "CODE5",
			},
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		var unmarshaled OauthUser
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.UserCodes, 5)
	})

	t.Run("Cloud names", func(t *testing.T) {
		clouds := []string{
			"zscaler.net",
			"zscalerone.net",
			"zscalertwo.net",
			"zscloud.net",
		}

		for _, cloud := range clouds {
			user := OauthUser{
				ConfigCloudName: cloud,
			}

			data, err := json.Marshal(user)
			require.NoError(t, err)

			assert.Contains(t, string(data), cloud)
		}
	})
}

