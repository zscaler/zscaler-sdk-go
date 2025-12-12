// Package unit provides unit tests for ZPA PRA Credential service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Credential represents the PRA credential for testing
type Credential struct {
	ID                      string `json:"id,omitempty"`
	Name                    string `json:"name,omitempty"`
	Description             string `json:"description,omitempty"`
	LastCredentialResetTime string `json:"lastCredentialResetTime,omitempty"`
	CredentialType          string `json:"credentialType,omitempty"`
	Passphrase              string `json:"passphrase,omitempty"`
	Password                string `json:"password,omitempty"`
	PrivateKey              string `json:"privateKey,omitempty"`
	UserDomain              string `json:"userDomain,omitempty"`
	UserName                string `json:"userName,omitempty"`
	CreationTime            string `json:"creationTime,omitempty"`
	ModifiedBy              string `json:"modifiedBy,omitempty"`
	ModifiedTime            string `json:"modifiedTime,omitempty"`
	MicroTenantID           string `json:"microtenantId,omitempty"`
	MicroTenantName         string `json:"microtenantName,omitempty"`
	TargetMicrotenantId     string `json:"targetMicrotenantId,omitempty"`
}

// TestCredential_Structure tests the struct definitions
func TestCredential_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Credential JSON marshaling", func(t *testing.T) {
		cred := Credential{
			ID:             "cred-123",
			Name:           "SSH Credential",
			Description:    "SSH key for server access",
			CredentialType: "SSH",
			UserName:       "admin",
			PrivateKey:     "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----",
			Passphrase:     "encrypted",
		}

		data, err := json.Marshal(cred)
		require.NoError(t, err)

		var unmarshaled Credential
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, cred.ID, unmarshaled.ID)
		assert.Equal(t, cred.Name, unmarshaled.Name)
		assert.Equal(t, "SSH", unmarshaled.CredentialType)
		assert.Equal(t, "admin", unmarshaled.UserName)
	})

	t.Run("Credential from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "cred-456",
			"name": "RDP Credential",
			"description": "Windows server credential",
			"credentialType": "RDP",
			"userName": "Administrator",
			"password": "encrypted_password",
			"userDomain": "CORP",
			"lastCredentialResetTime": "1612137600000",
			"creationTime": "1609459200000",
			"modifiedBy": "admin@example.com",
			"modifiedTime": "1612137600000",
			"microtenantId": "mt-001",
			"microtenantName": "Production"
		}`

		var cred Credential
		err := json.Unmarshal([]byte(apiResponse), &cred)
		require.NoError(t, err)

		assert.Equal(t, "cred-456", cred.ID)
		assert.Equal(t, "RDP Credential", cred.Name)
		assert.Equal(t, "RDP", cred.CredentialType)
		assert.Equal(t, "CORP", cred.UserDomain)
	})
}

// TestCredential_ResponseParsing tests parsing of API responses
func TestCredential_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse credential list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "SSH Cred", "credentialType": "SSH"},
				{"id": "2", "name": "RDP Cred", "credentialType": "RDP"},
				{"id": "3", "name": "VNC Cred", "credentialType": "VNC"}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []Credential `json:"list"`
			TotalPages int          `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.Equal(t, "SSH", listResp.List[0].CredentialType)
		assert.Equal(t, "VNC", listResp.List[2].CredentialType)
	})
}

// TestCredential_MockServerOperations tests CRUD operations
func TestCredential_MockServerOperations(t *testing.T) {
	t.Run("GET credential by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/credential/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "cred-123", "name": "Mock Credential", "credentialType": "SSH"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/credential/cred-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all credentials", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"list": [{"id": "1"}, {"id": "2"}], "totalPages": 1}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/credential")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create credential", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "new-cred", "name": "New Credential"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/credential", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update credential", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/credential/cred-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE credential", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/credential/cred-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("POST move credential to microtenant", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Contains(t, r.URL.Path, "/move")
			assert.NotEmpty(t, r.URL.Query().Get("targetMicrotenantId"))

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/credential/cred-123/move?targetMicrotenantId=mt-002", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestCredential_SpecialCases tests edge cases
func TestCredential_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Credential types", func(t *testing.T) {
		types := []string{"SSH", "RDP", "VNC"}

		for _, credType := range types {
			cred := Credential{
				ID:             "cred-" + credType,
				Name:           credType + " Credential",
				CredentialType: credType,
			}

			data, err := json.Marshal(cred)
			require.NoError(t, err)

			assert.Contains(t, string(data), credType)
		}
	})

	t.Run("SSH credential with private key", func(t *testing.T) {
		cred := Credential{
			ID:             "cred-ssh",
			Name:           "SSH Key Credential",
			CredentialType: "SSH",
			UserName:       "root",
			PrivateKey:     "-----BEGIN RSA PRIVATE KEY-----\nMIIE...\n-----END RSA PRIVATE KEY-----",
			Passphrase:     "keypassphrase",
		}

		data, err := json.Marshal(cred)
		require.NoError(t, err)

		var unmarshaled Credential
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Contains(t, unmarshaled.PrivateKey, "RSA PRIVATE KEY")
	})

	t.Run("RDP credential with domain", func(t *testing.T) {
		cred := Credential{
			ID:             "cred-rdp",
			Name:           "RDP Domain Credential",
			CredentialType: "RDP",
			UserName:       "Administrator",
			UserDomain:     "CORP.EXAMPLE.COM",
			Password:       "encrypted",
		}

		data, err := json.Marshal(cred)
		require.NoError(t, err)

		var unmarshaled Credential
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "CORP.EXAMPLE.COM", unmarshaled.UserDomain)
	})
}

