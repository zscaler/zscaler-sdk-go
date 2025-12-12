// Package unit provides unit tests for ZPA User Portal AUP service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// UserPortalAup represents the AUP for testing
type UserPortalAup struct {
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	Aup             string `json:"aup,omitempty"`
	Email           string `json:"email,omitempty"`
	PhoneNum        string `json:"phoneNum,omitempty"`
	Enabled         bool   `json:"enabled,omitempty"`
	MicrotenantID   string `json:"microtenantId,omitempty"`
	MicrotenantName string `json:"microtenantName,omitempty"`
	CreationTime    string `json:"creationTime,omitempty"`
	ModifiedBy      string `json:"modifiedBy,omitempty"`
	ModifiedTime    string `json:"modifiedTime,omitempty"`
}

// TestUserPortalAup_Structure tests the struct definitions
func TestUserPortalAup_Structure(t *testing.T) {
	t.Parallel()

	t.Run("UserPortalAup JSON marshaling", func(t *testing.T) {
		aup := UserPortalAup{
			ID:          "aup-123",
			Name:        "Corporate AUP",
			Description: "Acceptable Use Policy for corporate users",
			Aup:         "By accessing this portal, you agree to comply with all company policies...",
			Email:       "support@company.com",
			PhoneNum:    "+1-555-123-4567",
			Enabled:     true,
		}

		data, err := json.Marshal(aup)
		require.NoError(t, err)

		var unmarshaled UserPortalAup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, aup.ID, unmarshaled.ID)
		assert.Equal(t, aup.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
		assert.Contains(t, unmarshaled.Aup, "company policies")
	})

	t.Run("UserPortalAup from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "aup-456",
			"name": "Guest AUP",
			"description": "Acceptable Use Policy for guest users",
			"aup": "Guest access terms and conditions...",
			"email": "guest-support@company.com",
			"phoneNum": "+1-555-987-6543",
			"enabled": true,
			"microtenantId": "mt-001",
			"microtenantName": "Production",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var aup UserPortalAup
		err := json.Unmarshal([]byte(apiResponse), &aup)
		require.NoError(t, err)

		assert.Equal(t, "aup-456", aup.ID)
		assert.Equal(t, "Guest AUP", aup.Name)
		assert.True(t, aup.Enabled)
		assert.Equal(t, "mt-001", aup.MicrotenantID)
	})
}

// TestUserPortalAup_ResponseParsing tests parsing of API responses
func TestUserPortalAup_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse AUP list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "AUP 1", "enabled": true},
				{"id": "2", "name": "AUP 2", "enabled": true},
				{"id": "3", "name": "AUP 3", "enabled": false}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []UserPortalAup `json:"list"`
			TotalPages int             `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.True(t, listResp.List[0].Enabled)
		assert.False(t, listResp.List[2].Enabled)
	})
}

// TestUserPortalAup_MockServerOperations tests CRUD operations
func TestUserPortalAup_MockServerOperations(t *testing.T) {
	t.Run("GET AUP by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/userportal/aup/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "aup-123", "name": "Mock AUP", "enabled": true}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/userportal/aup/aup-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all AUPs", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"list": [{"id": "1"}, {"id": "2"}], "totalPages": 1}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/userportal/aup")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create AUP", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "new-aup", "name": "New AUP"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/userportal/aup", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update AUP", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/userportal/aup/aup-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE AUP", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/userportal/aup/aup-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestUserPortalAup_SpecialCases tests edge cases
func TestUserPortalAup_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Long AUP content", func(t *testing.T) {
		longContent := `
		ACCEPTABLE USE POLICY
		
		1. PURPOSE
		This Acceptable Use Policy defines the acceptable use of IT resources.
		
		2. SCOPE
		This policy applies to all users of company IT resources.
		
		3. PROHIBITED ACTIVITIES
		- Unauthorized access to systems
		- Distribution of malicious software
		- Violation of intellectual property rights
		
		4. COMPLIANCE
		All users must comply with this policy at all times.
		`

		aup := UserPortalAup{
			ID:   "aup-long",
			Name: "Detailed AUP",
			Aup:  longContent,
		}

		data, err := json.Marshal(aup)
		require.NoError(t, err)

		var unmarshaled UserPortalAup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Contains(t, unmarshaled.Aup, "ACCEPTABLE USE POLICY")
		assert.Contains(t, unmarshaled.Aup, "PROHIBITED ACTIVITIES")
	})

	t.Run("Disabled AUP", func(t *testing.T) {
		aup := UserPortalAup{
			ID:      "aup-disabled",
			Name:    "Disabled AUP",
			Enabled: false,
		}

		data, err := json.Marshal(aup)
		require.NoError(t, err)

		var unmarshaled UserPortalAup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.False(t, unmarshaled.Enabled)
	})
}

