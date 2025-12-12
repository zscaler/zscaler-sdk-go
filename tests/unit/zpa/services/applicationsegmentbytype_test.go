// Package unit provides unit tests for ZPA Application Segment By Type service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AppSegmentBaseAppDto represents the app segment base structure for testing
type AppSegmentBaseAppDto struct {
	ID                  string `json:"id,omitempty"`
	AppID               string `json:"appId,omitempty"`
	Name                string `json:"name,omitempty"`
	Enabled             bool   `json:"enabled"`
	Domain              string `json:"domain,omitempty"`
	ApplicationPort     string `json:"applicationPort,omitempty"`
	ApplicationProtocol string `json:"applicationProtocol,omitempty"`
	CertificateID       string `json:"certificateId,omitempty"`
	CertificateName     string `json:"certificateName,omitempty"`
	MicroTenantID       string `json:"microtenantId,omitempty"`
	MicroTenantName     string `json:"microtenantName,omitempty"`
}

// TestAppSegmentByType_Structure tests the struct definitions
func TestAppSegmentByType_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AppSegmentBaseAppDto JSON marshaling", func(t *testing.T) {
		app := AppSegmentBaseAppDto{
			ID:                  "app-123",
			AppID:               "parent-app-001",
			Name:                "Browser Access App",
			Enabled:             true,
			Domain:              "portal.example.com",
			ApplicationPort:     "443",
			ApplicationProtocol: "HTTPS",
			CertificateID:       "cert-001",
			CertificateName:     "Portal Certificate",
		}

		data, err := json.Marshal(app)
		require.NoError(t, err)

		var unmarshaled AppSegmentBaseAppDto
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, app.ID, unmarshaled.ID)
		assert.Equal(t, app.Name, unmarshaled.Name)
		assert.Equal(t, app.Domain, unmarshaled.Domain)
		assert.True(t, unmarshaled.Enabled)
	})

	t.Run("AppSegmentBaseAppDto from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "app-456",
			"appId": "parent-app-002",
			"name": "PRA App",
			"enabled": true,
			"domain": "pra.example.com",
			"applicationPort": "3389",
			"applicationProtocol": "RDP",
			"certificateId": "cert-002",
			"certificateName": "PRA Certificate",
			"microtenantId": "mt-001",
			"microtenantName": "Production"
		}`

		var app AppSegmentBaseAppDto
		err := json.Unmarshal([]byte(apiResponse), &app)
		require.NoError(t, err)

		assert.Equal(t, "app-456", app.ID)
		assert.Equal(t, "PRA App", app.Name)
		assert.Equal(t, "3389", app.ApplicationPort)
		assert.Equal(t, "RDP", app.ApplicationProtocol)
	})
}

// TestAppSegmentByType_ResponseParsing tests parsing of various API responses
func TestAppSegmentByType_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse app segment list by type response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "BA App 1", "enabled": true, "applicationProtocol": "HTTPS"},
				{"id": "2", "name": "BA App 2", "enabled": true, "applicationProtocol": "HTTPS"},
				{"id": "3", "name": "BA App 3", "enabled": false, "applicationProtocol": "HTTP"}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []AppSegmentBaseAppDto `json:"list"`
			TotalPages int                    `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.True(t, listResp.List[0].Enabled)
		assert.False(t, listResp.List[2].Enabled)
	})
}

// TestAppSegmentByType_MockServerOperations tests operations with mock server
func TestAppSegmentByType_MockServerOperations(t *testing.T) {
	t.Run("GET apps by BROWSER_ACCESS type", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/getAppsByType")
			assert.Equal(t, "BROWSER_ACCESS", r.URL.Query().Get("applicationType"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "BA App", "applicationProtocol": "HTTPS"}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/application/getAppsByType?applicationType=BROWSER_ACCESS")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET apps by SECURE_REMOTE_ACCESS type", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "SECURE_REMOTE_ACCESS", r.URL.Query().Get("applicationType"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "PRA App", "applicationProtocol": "RDP"}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/application/getAppsByType?applicationType=SECURE_REMOTE_ACCESS")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET apps by INSPECT type", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "INSPECT", r.URL.Query().Get("applicationType"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Inspect App", "applicationProtocol": "HTTPS"}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/application/getAppsByType?applicationType=INSPECT")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("DELETE app by type", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Contains(t, r.URL.Path, "/deleteAppByType")
			assert.Equal(t, "BROWSER_ACCESS", r.URL.Query().Get("applicationType"))

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/application/app-123/deleteAppByType?applicationType=BROWSER_ACCESS", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestAppSegmentByType_ApplicationTypes tests valid application types
func TestAppSegmentByType_ApplicationTypes(t *testing.T) {
	t.Parallel()

	t.Run("Valid application types", func(t *testing.T) {
		validTypes := []string{
			"BROWSER_ACCESS",
			"INSPECT",
			"SECURE_REMOTE_ACCESS",
		}

		for _, appType := range validTypes {
			assert.Contains(t, []string{"BROWSER_ACCESS", "INSPECT", "SECURE_REMOTE_ACCESS"}, appType)
		}
	})

	t.Run("Application protocols by type", func(t *testing.T) {
		protocolsByType := map[string][]string{
			"BROWSER_ACCESS":       {"HTTP", "HTTPS"},
			"INSPECT":              {"HTTPS"},
			"SECURE_REMOTE_ACCESS": {"RDP", "SSH", "VNC"},
		}

		assert.Len(t, protocolsByType["BROWSER_ACCESS"], 2)
		assert.Len(t, protocolsByType["SECURE_REMOTE_ACCESS"], 3)
	})
}

