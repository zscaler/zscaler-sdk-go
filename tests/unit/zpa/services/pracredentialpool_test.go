// Package unit provides unit tests for ZPA PRA Credential Pool service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CredentialPool represents the PRA credential pool for testing
type CredentialPool struct {
	ID                     string                    `json:"id,omitempty"`
	Name                   string                    `json:"name,omitempty"`
	CredentialType         string                    `json:"credentialType,omitempty"`
	PRACredentials         []CredentialPoolRef       `json:"credentials"`
	CredentialMappingCount string                    `json:"credentialMappingCount,omitempty"`
	CreationTime           string                    `json:"creationTime,omitempty"`
	ModifiedBy             string                    `json:"modifiedBy,omitempty"`
	ModifiedTime           string                    `json:"modifiedTime,omitempty"`
	MicroTenantID          string                    `json:"microtenantId,omitempty"`
	MicroTenantName        string                    `json:"microtenantName,omitempty"`
}

// CredentialPoolRef represents a credential reference in the pool
type CredentialPoolRef struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// TestCredentialPool_Structure tests the struct definitions
func TestCredentialPool_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CredentialPool JSON marshaling", func(t *testing.T) {
		pool := CredentialPool{
			ID:             "pool-123",
			Name:           "SSH Credential Pool",
			CredentialType: "SSH",
			PRACredentials: []CredentialPoolRef{
				{ID: "cred-001", Name: "SSH Key 1"},
				{ID: "cred-002", Name: "SSH Key 2"},
				{ID: "cred-003", Name: "SSH Key 3"},
			},
			CredentialMappingCount: "3",
		}

		data, err := json.Marshal(pool)
		require.NoError(t, err)

		var unmarshaled CredentialPool
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, pool.ID, unmarshaled.ID)
		assert.Equal(t, pool.Name, unmarshaled.Name)
		assert.Equal(t, "SSH", unmarshaled.CredentialType)
		assert.Len(t, unmarshaled.PRACredentials, 3)
	})

	t.Run("CredentialPool from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "pool-456",
			"name": "RDP Credential Pool",
			"credentialType": "RDP",
			"credentials": [
				{"id": "cred-004", "name": "RDP Admin 1"},
				{"id": "cred-005", "name": "RDP Admin 2"}
			],
			"credentialMappingCount": "2",
			"creationTime": "1609459200000",
			"modifiedBy": "admin@example.com",
			"modifiedTime": "1612137600000",
			"microtenantId": "mt-001",
			"microtenantName": "Production"
		}`

		var pool CredentialPool
		err := json.Unmarshal([]byte(apiResponse), &pool)
		require.NoError(t, err)

		assert.Equal(t, "pool-456", pool.ID)
		assert.Equal(t, "RDP Credential Pool", pool.Name)
		assert.Equal(t, "RDP", pool.CredentialType)
		assert.Len(t, pool.PRACredentials, 2)
		assert.Equal(t, "2", pool.CredentialMappingCount)
	})
}

// TestCredentialPool_ResponseParsing tests parsing of API responses
func TestCredentialPool_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse credential pool list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "SSH Pool", "credentialType": "SSH"},
				{"id": "2", "name": "RDP Pool", "credentialType": "RDP"},
				{"id": "3", "name": "VNC Pool", "credentialType": "VNC"}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []CredentialPool `json:"list"`
			TotalPages int              `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.Equal(t, "SSH", listResp.List[0].CredentialType)
	})
}

// TestCredentialPool_MockServerOperations tests CRUD operations
func TestCredentialPool_MockServerOperations(t *testing.T) {
	t.Run("GET credential pool by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/credential-pool/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "pool-123", "name": "Mock Pool", "credentialType": "SSH"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/credential-pool/pool-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all credential pools", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"list": [{"id": "1"}, {"id": "2"}], "totalPages": 1}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/credential-pool")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create credential pool", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/credential-pool", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("PUT update credential pool", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/credential-pool/pool-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE credential pool", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/credential-pool/pool-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestCredentialPool_SpecialCases tests edge cases
func TestCredentialPool_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Pool with many credentials", func(t *testing.T) {
		credentials := make([]CredentialPoolRef, 20)
		for i := 0; i < 20; i++ {
			credentials[i] = CredentialPoolRef{
				ID:   "cred-" + string(rune('0'+i)),
				Name: "Credential " + string(rune('0'+i)),
			}
		}

		pool := CredentialPool{
			ID:                     "pool-large",
			Name:                   "Large Pool",
			CredentialType:         "SSH",
			PRACredentials:         credentials,
			CredentialMappingCount: "20",
		}

		data, err := json.Marshal(pool)
		require.NoError(t, err)

		var unmarshaled CredentialPool
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.PRACredentials, 20)
	})

	t.Run("Empty credential pool", func(t *testing.T) {
		pool := CredentialPool{
			ID:                     "pool-empty",
			Name:                   "Empty Pool",
			CredentialType:         "SSH",
			PRACredentials:         []CredentialPoolRef{},
			CredentialMappingCount: "0",
		}

		data, err := json.Marshal(pool)
		require.NoError(t, err)

		var unmarshaled CredentialPool
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.PRACredentials, 0)
	})
}

