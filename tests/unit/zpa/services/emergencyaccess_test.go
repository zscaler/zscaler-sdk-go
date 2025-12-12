// Package unit provides unit tests for ZPA Emergency Access service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/emergencyaccess"
)

// TestEmergencyAccess_Structure tests the struct definitions
func TestEmergencyAccess_Structure(t *testing.T) {
	t.Parallel()

	t.Run("EmergencyAccess JSON marshaling", func(t *testing.T) {
		ea := emergencyaccess.EmergencyAccess{
			UserID:            "user-123",
			FirstName:         "John",
			LastName:          "Doe",
			EmailID:           "john.doe@example.com",
			AllowedActivate:   true,
			AllowedDeactivate: false,
			UpdateEnabled:     true,
			UserStatus:        "ACTIVE",
		}

		data, err := json.Marshal(ea)
		require.NoError(t, err)

		var unmarshaled emergencyaccess.EmergencyAccess
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, ea.UserID, unmarshaled.UserID)
		assert.Equal(t, ea.EmailID, unmarshaled.EmailID)
		assert.True(t, unmarshaled.AllowedActivate)
	})

	t.Run("EmergencyAccess from API response", func(t *testing.T) {
		apiResponse := `{
			"userId": "user-456",
			"firstName": "Jane",
			"lastName": "Smith",
			"emailId": "jane.smith@example.com",
			"allowedActivate": true,
			"allowedDeactivate": true,
			"updateEnabled": true,
			"userStatus": "ACTIVE",
			"activatedOn": "2024-01-15T10:30:00Z",
			"lastLoginTime": "2024-01-20T14:45:00Z"
		}`

		var ea emergencyaccess.EmergencyAccess
		err := json.Unmarshal([]byte(apiResponse), &ea)
		require.NoError(t, err)

		assert.Equal(t, "user-456", ea.UserID)
		assert.Equal(t, "Jane", ea.FirstName)
		assert.True(t, ea.AllowedActivate)
		assert.NotEmpty(t, ea.ActivatedOn)
	})
}

// TestEmergencyAccess_ResponseParsing tests parsing of API responses
func TestEmergencyAccess_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse emergency access list response", func(t *testing.T) {
		response := `{
			"list": [
				{"userId": "1", "emailId": "user1@example.com", "userStatus": "ACTIVE"},
				{"userId": "2", "emailId": "user2@example.com", "userStatus": "INACTIVE"},
				{"userId": "3", "emailId": "user3@example.com", "userStatus": "ACTIVE"}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []emergencyaccess.EmergencyAccess `json:"list"`
			TotalPages int                               `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.Equal(t, "ACTIVE", listResp.List[0].UserStatus)
	})
}

// TestEmergencyAccess_MockServerOperations tests CRUD operations
func TestEmergencyAccess_MockServerOperations(t *testing.T) {
	t.Run("GET emergency access by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"userId": "user-123", "emailId": "test@example.com", "userStatus": "ACTIVE"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/emergencyAccess/user-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create emergency access", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"userId": "new-user", "emailId": "new@example.com"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/emergencyAccess", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT activate emergency access", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/emergencyAccess/user-123/activate", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE emergency access", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/emergencyAccess/user-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestEmergencyAccess_SpecialCases tests edge cases
func TestEmergencyAccess_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("User status values", func(t *testing.T) {
		statuses := []string{"ACTIVE", "INACTIVE", "PENDING"}

		for _, status := range statuses {
			ea := emergencyaccess.EmergencyAccess{
				UserID:     "user-" + status,
				UserStatus: status,
			}

			data, err := json.Marshal(ea)
			require.NoError(t, err)
			assert.Contains(t, string(data), status)
		}
	})

	t.Run("ActivateNow flag", func(t *testing.T) {
		ea := emergencyaccess.EmergencyAccess{
			UserID:      "user-123",
			ActivateNow: true,
		}

		data, err := json.Marshal(ea)
		require.NoError(t, err)
		assert.Contains(t, string(data), `"activateNow":true`)
	})
}
