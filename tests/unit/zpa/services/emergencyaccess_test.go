// Package unit provides unit tests for ZPA Emergency Access service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// EmergencyAccessUser represents the emergency access user structure for testing
type EmergencyAccessUser struct {
	ID                    string `json:"id,omitempty"`
	UserID                string `json:"userId,omitempty"`
	FirstName             string `json:"firstName,omitempty"`
	LastName              string `json:"lastName,omitempty"`
	EmailID               string `json:"emailId,omitempty"`
	ActivateTime          string `json:"activateTime,omitempty"`
	DeactivateTime        string `json:"deactivateTime,omitempty"`
	IsEnabled             bool   `json:"isEnabled"`
	CreationTime          string `json:"creationTime,omitempty"`
	ModifiedBy            string `json:"modifiedBy,omitempty"`
	ModifiedTime          string `json:"modifiedTime,omitempty"`
	MicroTenantID         string `json:"microtenantId,omitempty"`
	MicroTenantName       string `json:"microtenantName,omitempty"`
}

// TestEmergencyAccessUser_Structure tests the struct definitions
func TestEmergencyAccessUser_Structure(t *testing.T) {
	t.Parallel()

	t.Run("EmergencyAccessUser JSON marshaling", func(t *testing.T) {
		user := EmergencyAccessUser{
			ID:             "ea-123",
			UserID:         "user-001",
			FirstName:      "John",
			LastName:       "Doe",
			EmailID:        "john.doe@example.com",
			ActivateTime:   "1609459200000",
			DeactivateTime: "1612137600000",
			IsEnabled:      true,
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		var unmarshaled EmergencyAccessUser
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, user.ID, unmarshaled.ID)
		assert.Equal(t, user.EmailID, unmarshaled.EmailID)
		assert.True(t, unmarshaled.IsEnabled)
	})

	t.Run("EmergencyAccessUser JSON unmarshaling from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "ea-456",
			"userId": "user-002",
			"firstName": "Jane",
			"lastName": "Smith",
			"emailId": "jane.smith@example.com",
			"activateTime": "1609459200000",
			"deactivateTime": "1612137600000",
			"isEnabled": true,
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com",
			"microtenantId": "mt-001",
			"microtenantName": "Production"
		}`

		var user EmergencyAccessUser
		err := json.Unmarshal([]byte(apiResponse), &user)
		require.NoError(t, err)

		assert.Equal(t, "ea-456", user.ID)
		assert.Equal(t, "Jane", user.FirstName)
		assert.Equal(t, "Smith", user.LastName)
		assert.True(t, user.IsEnabled)
		assert.NotEmpty(t, user.ActivateTime)
	})
}

// TestEmergencyAccessUser_ResponseParsing tests parsing of various API responses
func TestEmergencyAccessUser_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse emergency access user list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "emailId": "user1@example.com", "isEnabled": true},
				{"id": "2", "emailId": "user2@example.com", "isEnabled": true},
				{"id": "3", "emailId": "user3@example.com", "isEnabled": false}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []EmergencyAccessUser `json:"list"`
			TotalPages int                   `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.True(t, listResp.List[0].IsEnabled)
		assert.False(t, listResp.List[2].IsEnabled)
	})
}

// TestEmergencyAccessUser_MockServerOperations tests CRUD operations with mock server
func TestEmergencyAccessUser_MockServerOperations(t *testing.T) {
	t.Run("GET emergency access user by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/emergencyAccess/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "ea-123",
				"emailId": "mock@example.com",
				"isEnabled": true
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/emergencyAccess/ea-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all emergency access users", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "emailId": "user1@example.com", "isEnabled": true},
					{"id": "2", "emailId": "user2@example.com", "isEnabled": true}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/emergencyAccess")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create emergency access user", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "new-ea-456",
				"emailId": "newuser@example.com",
				"isEnabled": true
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/emergencyAccess", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update emergency access user", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/emergencyAccess/ea-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE emergency access user", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/emergencyAccess/ea-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestEmergencyAccessUser_SpecialCases tests edge cases
func TestEmergencyAccessUser_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Activation and deactivation times", func(t *testing.T) {
		user := EmergencyAccessUser{
			ID:             "ea-123",
			EmailID:        "user@example.com",
			ActivateTime:   "1609459200000",  // 2021-01-01
			DeactivateTime: "1612137600000",  // 2021-02-01
			IsEnabled:      true,
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		var unmarshaled EmergencyAccessUser
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.NotEmpty(t, unmarshaled.ActivateTime)
		assert.NotEmpty(t, unmarshaled.DeactivateTime)
	})

	t.Run("Disabled emergency access user", func(t *testing.T) {
		user := EmergencyAccessUser{
			ID:        "ea-123",
			EmailID:   "disabled@example.com",
			IsEnabled: false,
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"isEnabled":false`)
	})

	t.Run("Microtenant association", func(t *testing.T) {
		user := EmergencyAccessUser{
			ID:              "ea-123",
			EmailID:         "tenant@example.com",
			MicroTenantID:   "mt-001",
			MicroTenantName: "Production",
			IsEnabled:       true,
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"microtenantId":"mt-001"`)
	})
}

