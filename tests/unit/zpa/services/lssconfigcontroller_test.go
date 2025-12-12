// Package unit provides unit tests for ZPA LSS Config Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// LSSConfig represents the LSS configuration structure for testing
type LSSConfig struct {
	ID            string      `json:"id,omitempty"`
	Name          string      `json:"name,omitempty"`
	Enabled       bool        `json:"enabled"`
	Filter        []string    `json:"filter,omitempty"`
	Format        string      `json:"format,omitempty"`
	LSSHost       string      `json:"lssHost,omitempty"`
	LSSPort       string      `json:"lssPort,omitempty"`
	SourceLogType string      `json:"sourceLogType,omitempty"`
	UseTLS        bool        `json:"useTls"`
	CreationTime  string      `json:"creationTime,omitempty"`
	ModifiedBy    string      `json:"modifiedBy,omitempty"`
	ModifiedTime  string      `json:"modifiedTime,omitempty"`
	Config        LSSLogConfig `json:"config,omitempty"`
}

// LSSLogConfig represents LSS log configuration
type LSSLogConfig struct {
	AuditMessage   string   `json:"auditMessage,omitempty"`
	Description    string   `json:"description,omitempty"`
	Enabled        bool     `json:"enabled"`
	Filter         []string `json:"filter,omitempty"`
	Format         string   `json:"format,omitempty"`
	ID             string   `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	LSSHost        string   `json:"lssHost,omitempty"`
	LSSPort        string   `json:"lssPort,omitempty"`
	SourceLogType  string   `json:"sourceLogType,omitempty"`
	UseTLS         bool     `json:"useTls"`
}

// LSSLogFormat represents log format options
type LSSLogFormat struct {
	CSV  string `json:"csv,omitempty"`
	JSON string `json:"json,omitempty"`
	TSV  string `json:"tsv,omitempty"`
}

// TestLSSConfig_Structure tests the struct definitions
func TestLSSConfig_Structure(t *testing.T) {
	t.Parallel()

	t.Run("LSSConfig JSON marshaling", func(t *testing.T) {
		lss := LSSConfig{
			ID:            "lss-123",
			Name:          "Audit Log Receiver",
			Enabled:       true,
			LSSHost:       "siem.example.com",
			LSSPort:       "514",
			SourceLogType: "zpn_trans_log",
			UseTLS:        true,
			Format:        "json",
			Filter:        []string{"ZPN_STATUS_AUTH_FAILED", "ZPN_STATUS_DISCONNECTED"},
		}

		data, err := json.Marshal(lss)
		require.NoError(t, err)

		var unmarshaled LSSConfig
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, lss.ID, unmarshaled.ID)
		assert.Equal(t, lss.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
		assert.True(t, unmarshaled.UseTLS)
		assert.Equal(t, "json", unmarshaled.Format)
	})

	t.Run("LSSConfig JSON unmarshaling from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "lss-456",
			"name": "User Activity Log",
			"enabled": true,
			"lssHost": "splunk.example.com",
			"lssPort": "9997",
			"sourceLogType": "zpn_auth_log",
			"useTls": true,
			"format": "csv",
			"filter": ["ZPN_STATUS_AUTH_FAILED"],
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com",
			"config": {
				"id": "config-001",
				"name": "Auth Config",
				"enabled": true,
				"lssHost": "splunk.example.com",
				"lssPort": "9997",
				"sourceLogType": "zpn_auth_log",
				"useTls": true
			}
		}`

		var lss LSSConfig
		err := json.Unmarshal([]byte(apiResponse), &lss)
		require.NoError(t, err)

		assert.Equal(t, "lss-456", lss.ID)
		assert.Equal(t, "User Activity Log", lss.Name)
		assert.True(t, lss.Enabled)
		assert.Equal(t, "9997", lss.LSSPort)
		assert.True(t, lss.UseTLS)
		assert.NotEmpty(t, lss.Config.ID)
	})
}

// TestLSSConfig_ResponseParsing tests parsing of various API responses
func TestLSSConfig_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse LSS config list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "Config A", "enabled": true, "sourceLogType": "zpn_trans_log"},
				{"id": "2", "name": "Config B", "enabled": true, "sourceLogType": "zpn_auth_log"},
				{"id": "3", "name": "Config C", "enabled": false, "sourceLogType": "zpn_ast_auth_log"}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []LSSConfig `json:"list"`
			TotalPages int         `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.True(t, listResp.List[0].Enabled)
		assert.False(t, listResp.List[2].Enabled)
	})
}

// TestLSSConfig_MockServerOperations tests CRUD operations with mock server
func TestLSSConfig_MockServerOperations(t *testing.T) {
	t.Run("GET LSS config by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/lssConfig/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "lss-123",
				"name": "Mock LSS Config",
				"enabled": true
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v2/admin/customers/123/lssConfig/lss-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all LSS configs", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Config A", "enabled": true},
					{"id": "2", "name": "Config B", "enabled": true}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v2/admin/customers/123/lssConfig")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create LSS config", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "new-lss-456",
				"name": "New LSS Config",
				"enabled": true
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/lssConfig", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update LSS config", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/lssConfig/lss-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE LSS config", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/lssConfig/lss-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestLSSConfig_SpecialCases tests edge cases
func TestLSSConfig_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Source log types", func(t *testing.T) {
		logTypes := []string{
			"zpn_trans_log",
			"zpn_auth_log",
			"zpn_ast_auth_log",
			"zpn_http_trans_log",
			"zpn_audit_log",
			"zpn_sys_auth_log",
			"zpn_ast_comprehensive_stats",
			"zpn_pbroker_comprehensive_stats",
			"zpn_waf_http_exchanges_log",
		}

		for _, logType := range logTypes {
			lss := LSSConfig{
				ID:            "lss-" + logType,
				Name:          logType + " Config",
				SourceLogType: logType,
			}

			data, err := json.Marshal(lss)
			require.NoError(t, err)

			assert.Contains(t, string(data), logType)
		}
	})

	t.Run("Log formats", func(t *testing.T) {
		formats := []string{"json", "csv", "tsv"}

		for _, format := range formats {
			lss := LSSConfig{
				ID:     "lss-format-" + format,
				Name:   format + " Format Config",
				Format: format,
			}

			data, err := json.Marshal(lss)
			require.NoError(t, err)

			assert.Contains(t, string(data), format)
		}
	})

	t.Run("TLS configuration", func(t *testing.T) {
		lss := LSSConfig{
			ID:      "lss-123",
			Name:    "TLS Config",
			UseTLS:  true,
			LSSHost: "secure.example.com",
			LSSPort: "6514",
		}

		data, err := json.Marshal(lss)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"useTls":true`)
	})

	t.Run("Multiple filters", func(t *testing.T) {
		lss := LSSConfig{
			ID:     "lss-123",
			Name:   "Filtered Config",
			Filter: []string{
				"ZPN_STATUS_AUTH_FAILED",
				"ZPN_STATUS_DISCONNECTED",
				"ZPN_STATUS_TIMEOUT",
				"ZPN_STATUS_UNREACHABLE",
			},
		}

		data, err := json.Marshal(lss)
		require.NoError(t, err)

		var unmarshaled LSSConfig
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.Filter, 4)
	})
}

