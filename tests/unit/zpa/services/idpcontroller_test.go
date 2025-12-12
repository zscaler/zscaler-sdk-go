// Package unit provides unit tests for ZPA IDP Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// IdpController represents the IDP controller structure for testing
type IdpController struct {
	ID                     string          `json:"id,omitempty"`
	Name                   string          `json:"name,omitempty"`
	Description            string          `json:"description,omitempty"`
	IdpEntityID            string          `json:"idpEntityId,omitempty"`
	LoginURL               string          `json:"loginUrl,omitempty"`
	AdminSPSigningCertID   string          `json:"adminSpSigningCertId,omitempty"`
	SignSamlRequest        string          `json:"signSamlRequest,omitempty"`
	UseCustomSPMetadata    bool            `json:"useCustomSpMetadata"`
	SsoType                []string        `json:"ssoType,omitempty"`
	DomainList             []string        `json:"domainList,omitempty"`
	Enabled                bool            `json:"enabled"`
	EnableScimBasedPolicy  bool            `json:"enableScimBasedPolicy"`
	DisableSamlBasedPolicy bool            `json:"disableSamlBasedPolicy"`
	EnableArbitraryAuthDomains bool        `json:"enableArbitraryAuthDomains"`
	ForceAuth              bool            `json:"forceAuth"`
	AutoProvision          string          `json:"autoProvision,omitempty"`
	CreationTime           string          `json:"creationTime,omitempty"`
	ModifiedBy             string          `json:"modifiedBy,omitempty"`
	ModifiedTime           string          `json:"modifiedTime,omitempty"`
	MicroTenantID          string          `json:"microtenantId,omitempty"`
	MicroTenantName        string          `json:"microtenantName,omitempty"`
	ReauthOnUserUpdate     bool            `json:"reauthOnUserUpdate"`
	RedirectBinding        bool            `json:"redirectBinding"`
	ScimEnabled            bool            `json:"scimEnabled"`
	ScimServiceProviderEndpoint string     `json:"scimServiceProviderEndpoint,omitempty"`
	ScimSharedSecretExists bool            `json:"scimSharedSecretExists"`
	AdminMetadata          AdminMetadata   `json:"adminMetadata,omitempty"`
	UserMetadata           UserMetadata    `json:"userMetadata,omitempty"`
}

type AdminMetadata struct {
	CertificateURL string `json:"certificateUrl,omitempty"`
	SpEntityID     string `json:"spEntityId,omitempty"`
	SpMetadataURL  string `json:"spMetadataUrl,omitempty"`
	SpPostURL      string `json:"spPostUrl,omitempty"`
}

type UserMetadata struct {
	CertificateURL string `json:"certificateUrl,omitempty"`
	SpEntityID     string `json:"spEntityId,omitempty"`
	SpMetadataURL  string `json:"spMetadataUrl,omitempty"`
	SpPostURL      string `json:"spPostUrl,omitempty"`
}

// TestIdpController_Structure tests the struct definitions
func TestIdpController_Structure(t *testing.T) {
	t.Parallel()

	t.Run("IdpController JSON marshaling", func(t *testing.T) {
		idp := IdpController{
			ID:                    "idp-123",
			Name:                  "Okta IDP",
			Description:           "Okta Identity Provider",
			IdpEntityID:           "https://okta.example.com/entity",
			LoginURL:              "https://okta.example.com/login",
			SignSamlRequest:       "1",
			SsoType:               []string{"USER", "ADMIN"},
			DomainList:            []string{"example.com", "test.com"},
			Enabled:               true,
			EnableScimBasedPolicy: true,
			ForceAuth:             false,
			AutoProvision:         "SCIM_HYBRID",
			ScimEnabled:           true,
		}

		data, err := json.Marshal(idp)
		require.NoError(t, err)

		var unmarshaled IdpController
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, idp.ID, unmarshaled.ID)
		assert.Equal(t, idp.Name, unmarshaled.Name)
		assert.Equal(t, idp.IdpEntityID, unmarshaled.IdpEntityID)
		assert.ElementsMatch(t, idp.SsoType, unmarshaled.SsoType)
		assert.ElementsMatch(t, idp.DomainList, unmarshaled.DomainList)
	})

	t.Run("IdpController JSON unmarshaling from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "idp-456",
			"name": "Azure AD",
			"description": "Azure Active Directory",
			"idpEntityId": "https://login.microsoftonline.com/tenant-id",
			"loginUrl": "https://login.microsoftonline.com/tenant-id/saml2",
			"signSamlRequest": "1",
			"ssoType": ["USER"],
			"domainList": ["company.com"],
			"enabled": true,
			"enableScimBasedPolicy": true,
			"disableSamlBasedPolicy": false,
			"enableArbitraryAuthDomains": false,
			"forceAuth": true,
			"autoProvision": "SCIM_HYBRID",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com",
			"reauthOnUserUpdate": true,
			"redirectBinding": true,
			"scimEnabled": true,
			"scimServiceProviderEndpoint": "https://scim.example.com/Users",
			"scimSharedSecretExists": true,
			"adminMetadata": {
				"certificateUrl": "https://admin.example.com/cert",
				"spEntityId": "admin-entity-id",
				"spMetadataUrl": "https://admin.example.com/metadata",
				"spPostUrl": "https://admin.example.com/sso"
			},
			"userMetadata": {
				"certificateUrl": "https://user.example.com/cert",
				"spEntityId": "user-entity-id",
				"spMetadataUrl": "https://user.example.com/metadata",
				"spPostUrl": "https://user.example.com/sso"
			}
		}`

		var idp IdpController
		err := json.Unmarshal([]byte(apiResponse), &idp)
		require.NoError(t, err)

		assert.Equal(t, "idp-456", idp.ID)
		assert.Equal(t, "Azure AD", idp.Name)
		assert.True(t, idp.Enabled)
		assert.True(t, idp.ForceAuth)
		assert.True(t, idp.ScimEnabled)
		assert.True(t, idp.ReauthOnUserUpdate)
		assert.NotEmpty(t, idp.AdminMetadata.SpEntityID)
		assert.NotEmpty(t, idp.UserMetadata.SpEntityID)
	})

	t.Run("AdminMetadata structure", func(t *testing.T) {
		metadata := AdminMetadata{
			CertificateURL: "https://cert.example.com",
			SpEntityID:     "sp-entity-123",
			SpMetadataURL:  "https://metadata.example.com",
			SpPostURL:      "https://sso.example.com/post",
		}

		data, err := json.Marshal(metadata)
		require.NoError(t, err)

		var unmarshaled AdminMetadata
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, metadata.SpEntityID, unmarshaled.SpEntityID)
	})
}

// TestIdpController_ResponseParsing tests parsing of various API responses
func TestIdpController_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse IDP list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "Okta", "enabled": true, "ssoType": ["USER", "ADMIN"]},
				{"id": "2", "name": "Azure AD", "enabled": true, "ssoType": ["USER"]},
				{"id": "3", "name": "OneLogin", "enabled": false, "ssoType": ["USER"]}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []IdpController `json:"list"`
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

// TestIdpController_MockServerOperations tests CRUD operations with mock server
func TestIdpController_MockServerOperations(t *testing.T) {
	t.Run("GET IDP by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/idp/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "idp-123",
				"name": "Mock IDP",
				"enabled": true,
				"ssoType": ["USER"]
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v2/admin/customers/123/idp/idp-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all IDPs", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "IDP A", "enabled": true},
					{"id": "2", "name": "IDP B", "enabled": true}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v2/admin/customers/123/idp")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET IDP by SSO type", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "USER", r.URL.Query().Get("ssoType"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "User IDP", "enabled": true, "ssoType": ["USER"]}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/idp?ssoType=USER")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestIdpController_ErrorHandling tests error scenarios
func TestIdpController_ErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("404 IDP Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"code": "NOT_FOUND", "message": "IDP not found"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/idp/nonexistent")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

// TestIdpController_SpecialCases tests edge cases
func TestIdpController_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("IDP with SCIM enabled", func(t *testing.T) {
		idp := IdpController{
			ID:                          "idp-123",
			Name:                        "SCIM IDP",
			ScimEnabled:                 true,
			ScimServiceProviderEndpoint: "https://scim.example.com",
			ScimSharedSecretExists:      true,
			EnableScimBasedPolicy:       true,
		}

		data, err := json.Marshal(idp)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"scimEnabled":true`)
		assert.Contains(t, string(data), `"enableScimBasedPolicy":true`)
	})

	t.Run("IDP with multiple SSO types", func(t *testing.T) {
		idp := IdpController{
			ID:      "idp-123",
			Name:    "Multi SSO IDP",
			SsoType: []string{"USER", "ADMIN"},
		}

		data, err := json.Marshal(idp)
		require.NoError(t, err)

		var unmarshaled IdpController
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.SsoType, 2)
		assert.Contains(t, unmarshaled.SsoType, "USER")
		assert.Contains(t, unmarshaled.SsoType, "ADMIN")
	})

	t.Run("IDP with force auth", func(t *testing.T) {
		idp := IdpController{
			ID:        "idp-123",
			Name:      "Force Auth IDP",
			ForceAuth: true,
		}

		data, err := json.Marshal(idp)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"forceAuth":true`)
	})

	t.Run("IDP with arbitrary auth domains", func(t *testing.T) {
		idp := IdpController{
			ID:                         "idp-123",
			Name:                       "Arbitrary Domains IDP",
			EnableArbitraryAuthDomains: true,
			DomainList:                 []string{"*"},
		}

		data, err := json.Marshal(idp)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"enableArbitraryAuthDomains":true`)
	})

	t.Run("Auto provision modes", func(t *testing.T) {
		modes := []string{"SCIM_ONLY", "SAML_ONLY", "SCIM_HYBRID", "DISABLED"}

		for _, mode := range modes {
			idp := IdpController{
				ID:            "idp-" + mode,
				Name:          mode + " IDP",
				AutoProvision: mode,
			}

			data, err := json.Marshal(idp)
			require.NoError(t, err)

			assert.Contains(t, string(data), mode)
		}
	})
}

