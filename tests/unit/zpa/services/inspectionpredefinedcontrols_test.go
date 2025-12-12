// Package unit provides unit tests for ZPA Inspection Predefined Controls service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PredefinedControls represents predefined controls for testing
type PredefinedControls struct {
	ID                 string `json:"id,omitempty"`
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	Action             string `json:"action,omitempty"`
	ActionValue        string `json:"actionValue,omitempty"`
	Attachment         string `json:"attachment,omitempty"`
	ControlGroup       string `json:"controlGroup,omitempty"`
	ControlType        string `json:"controlType,omitempty"`
	ControlNumber      string `json:"controlNumber,omitempty"`
	DefaultAction      string `json:"defaultAction,omitempty"`
	DefaultActionValue string `json:"defaultActionValue,omitempty"`
	ParanoiaLevel      string `json:"paranoiaLevel,omitempty"`
	ProtocolType       string `json:"protocolType,omitempty"`
	Severity           string `json:"severity,omitempty"`
	Version            string `json:"version,omitempty"`
	CreationTime       string `json:"creationTime,omitempty"`
	ModifiedBy         string `json:"modifiedBy,omitempty"`
	ModifiedTime       string `json:"modifiedTime,omitempty"`
}

// ControlGroupItem represents a control group for testing
type ControlGroupItem struct {
	ControlGroup                 string               `json:"controlGroup,omitempty"`
	PredefinedInspectionControls []PredefinedControls `json:"predefinedInspectionControls,omitempty"`
	DefaultGroup                 bool                 `json:"defaultGroup,omitempty"`
}

// TestPredefinedControls_Structure tests the struct definitions
func TestPredefinedControls_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PredefinedControls JSON marshaling", func(t *testing.T) {
		control := PredefinedControls{
			ID:            "pc-123",
			Name:          "SQL Injection Attack",
			Description:   "Detects SQL injection attempts",
			Action:        "BLOCK",
			ControlGroup:  "SQL_INJECTION",
			ControlNumber: "942100",
			Severity:      "CRITICAL",
			ParanoiaLevel: "1",
			Version:       "OWASP_CRS/3.3.0",
		}

		data, err := json.Marshal(control)
		require.NoError(t, err)

		var unmarshaled PredefinedControls
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, control.ID, unmarshaled.ID)
		assert.Equal(t, control.Name, unmarshaled.Name)
		assert.Equal(t, "SQL_INJECTION", unmarshaled.ControlGroup)
	})

	t.Run("PredefinedControls from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "pc-456",
			"name": "Cross-Site Scripting",
			"description": "Detects XSS attempts",
			"action": "BLOCK",
			"controlGroup": "XSS",
			"controlType": "PREDEFINED",
			"controlNumber": "941100",
			"defaultAction": "BLOCK",
			"paranoiaLevel": "1",
			"protocolType": "HTTP",
			"severity": "CRITICAL",
			"version": "OWASP_CRS/3.3.0",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var control PredefinedControls
		err := json.Unmarshal([]byte(apiResponse), &control)
		require.NoError(t, err)

		assert.Equal(t, "pc-456", control.ID)
		assert.Equal(t, "XSS", control.ControlGroup)
		assert.Equal(t, "CRITICAL", control.Severity)
	})

	t.Run("ControlGroupItem JSON marshaling", func(t *testing.T) {
		group := ControlGroupItem{
			ControlGroup: "SQL_INJECTION",
			DefaultGroup: true,
			PredefinedInspectionControls: []PredefinedControls{
				{ID: "pc-1", Name: "SQL Injection 1"},
				{ID: "pc-2", Name: "SQL Injection 2"},
			},
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled ControlGroupItem
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.True(t, unmarshaled.DefaultGroup)
		assert.Len(t, unmarshaled.PredefinedInspectionControls, 2)
	})
}

// TestPredefinedControls_MockServerOperations tests operations
func TestPredefinedControls_MockServerOperations(t *testing.T) {
	t.Run("GET predefined control by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/inspectionControls/predefined/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "pc-123", "name": "Mock Control"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/inspectionControls/predefined/pc-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all predefined controls with version", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.NotEmpty(t, r.URL.Query().Get("version"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `[
				{
					"controlGroup": "SQL_INJECTION",
					"defaultGroup": true,
					"predefinedInspectionControls": [
						{"id": "1", "name": "SQL Control 1"}
					]
				}
			]`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/inspectionControls/predefined?version=OWASP_CRS/3.3.0")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET controls by group", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Query().Get("search"), "controlGroup")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `[
				{
					"controlGroup": "XSS",
					"predefinedInspectionControls": [
						{"id": "1", "name": "XSS Control 1"}
					]
				}
			]`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/inspectionControls/predefined?version=OWASP_CRS/3.3.0&search=controlGroup+EQ+XSS")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestPredefinedControls_SpecialCases tests edge cases
func TestPredefinedControls_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Control groups", func(t *testing.T) {
		groups := []string{
			"SQL_INJECTION",
			"XSS",
			"LOCAL_FILE_INCLUSION",
			"REMOTE_FILE_INCLUSION",
			"REMOTE_CODE_EXECUTION",
			"PROTOCOL_ATTACK",
			"DATA_LEAKAGE",
		}

		for _, group := range groups {
			control := PredefinedControls{
				ID:           "pc-" + group,
				Name:         group + " Control",
				ControlGroup: group,
			}

			data, err := json.Marshal(control)
			require.NoError(t, err)

			assert.Contains(t, string(data), group)
		}
	})

	t.Run("OWASP versions", func(t *testing.T) {
		versions := []string{
			"OWASP_CRS/3.3.0",
			"OWASP_CRS/3.2.0",
			"OWASP_CRS/3.1.0",
		}

		for _, version := range versions {
			control := PredefinedControls{
				ID:      "pc-version",
				Name:    version + " Control",
				Version: version,
			}

			data, err := json.Marshal(control)
			require.NoError(t, err)

			assert.Contains(t, string(data), version)
		}
	})
}

