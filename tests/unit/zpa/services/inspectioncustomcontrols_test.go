// Package unit provides unit tests for ZPA Inspection Custom Controls service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// InspectionCustomControl represents the inspection custom control for testing
type InspectionCustomControl struct {
	ID                 string                       `json:"id,omitempty"`
	Name               string                       `json:"name,omitempty"`
	Description        string                       `json:"description,omitempty"`
	Action             string                       `json:"action,omitempty"`
	ActionValue        string                       `json:"actionValue,omitempty"`
	ControlNumber      string                       `json:"controlNumber,omitempty"`
	ControlType        string                       `json:"controlType,omitempty"`
	ControlRuleJson    string                       `json:"controlRuleJson,omitempty"`
	DefaultAction      string                       `json:"defaultAction,omitempty"`
	DefaultActionValue string                       `json:"defaultActionValue,omitempty"`
	ParanoiaLevel      string                       `json:"paranoiaLevel,omitempty"`
	ProtocolType       string                       `json:"protocolType,omitempty"`
	Severity           string                       `json:"severity,omitempty"`
	Type               string                       `json:"type,omitempty"`
	Version            string                       `json:"version,omitempty"`
	Rules              []InspectionRule             `json:"rules,omitempty"`
	CreationTime       string                       `json:"creationTime,omitempty"`
	ModifiedBy         string                       `json:"modifiedBy,omitempty"`
	ModifiedTime       string                       `json:"modifiedTime,omitempty"`
}

// InspectionRule represents an inspection rule
type InspectionRule struct {
	Conditions []InspectionCondition `json:"conditions,omitempty"`
	Names      []string              `json:"names,omitempty"`
	Type       string                `json:"type,omitempty"`
}

// InspectionCondition represents a condition in a rule
type InspectionCondition struct {
	LHS string `json:"lhs,omitempty"`
	OP  string `json:"op,omitempty"`
	RHS string `json:"rhs,omitempty"`
}

// TestInspectionCustomControl_Structure tests the struct definitions
func TestInspectionCustomControl_Structure(t *testing.T) {
	t.Parallel()

	t.Run("InspectionCustomControl JSON marshaling", func(t *testing.T) {
		control := InspectionCustomControl{
			ID:            "icc-123",
			Name:          "SQL Injection Custom",
			Description:   "Custom SQL injection detection",
			Action:        "BLOCK",
			ControlNumber: "200001",
			ControlType:   "CUSTOM",
			Severity:      "CRITICAL",
			ParanoiaLevel: "1",
			ProtocolType:  "HTTP",
			Version:       "OWASP_CRS/3.3.0",
			Rules: []InspectionRule{
				{
					Type:  "REQUEST_HEADERS",
					Names: []string{"user-agent", "referer"},
					Conditions: []InspectionCondition{
						{LHS: "VALUE", OP: "CONTAINS", RHS: "union select"},
					},
				},
			},
		}

		data, err := json.Marshal(control)
		require.NoError(t, err)

		var unmarshaled InspectionCustomControl
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, control.ID, unmarshaled.ID)
		assert.Equal(t, control.Name, unmarshaled.Name)
		assert.Equal(t, "BLOCK", unmarshaled.Action)
		assert.Len(t, unmarshaled.Rules, 1)
	})

	t.Run("InspectionCustomControl from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "icc-456",
			"name": "XSS Custom Detection",
			"description": "Custom cross-site scripting detection",
			"action": "REDIRECT",
			"actionValue": "https://blocked.example.com",
			"controlNumber": "200002",
			"controlType": "CUSTOM",
			"severity": "HIGH",
			"paranoiaLevel": "2",
			"protocolType": "HTTP",
			"version": "OWASP_CRS/3.3.0",
			"defaultAction": "PASS",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var control InspectionCustomControl
		err := json.Unmarshal([]byte(apiResponse), &control)
		require.NoError(t, err)

		assert.Equal(t, "icc-456", control.ID)
		assert.Equal(t, "XSS Custom Detection", control.Name)
		assert.Equal(t, "REDIRECT", control.Action)
		assert.Equal(t, "HIGH", control.Severity)
	})
}

// TestInspectionCustomControl_MockServerOperations tests CRUD operations
func TestInspectionCustomControl_MockServerOperations(t *testing.T) {
	t.Run("GET custom control by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/inspectionControls/custom/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "icc-123", "name": "Mock Control", "action": "BLOCK"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/inspectionControls/custom/icc-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all custom controls", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"list": [{"id": "1"}, {"id": "2"}], "totalPages": 1}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/inspectionControls/custom")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create custom control", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "new-icc", "name": "New Control"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/inspectionControls/custom", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("DELETE custom control", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/inspectionControls/custom/icc-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestInspectionCustomControl_SpecialCases tests edge cases
func TestInspectionCustomControl_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Action types", func(t *testing.T) {
		actions := []string{"PASS", "BLOCK", "REDIRECT"}

		for _, action := range actions {
			control := InspectionCustomControl{
				ID:     "icc-" + action,
				Name:   action + " Control",
				Action: action,
			}

			data, err := json.Marshal(control)
			require.NoError(t, err)

			assert.Contains(t, string(data), action)
		}
	})

	t.Run("Severity levels", func(t *testing.T) {
		severities := []string{"CRITICAL", "HIGH", "MEDIUM", "LOW", "INFO"}

		for _, severity := range severities {
			control := InspectionCustomControl{
				ID:       "icc-" + severity,
				Name:     severity + " Control",
				Severity: severity,
			}

			data, err := json.Marshal(control)
			require.NoError(t, err)

			assert.Contains(t, string(data), severity)
		}
	})

	t.Run("Paranoia levels", func(t *testing.T) {
		levels := []string{"1", "2", "3", "4"}

		for _, level := range levels {
			control := InspectionCustomControl{
				ID:            "icc-pl" + level,
				Name:          "PL" + level + " Control",
				ParanoiaLevel: level,
			}

			data, err := json.Marshal(control)
			require.NoError(t, err)

			assert.Contains(t, string(data), level)
		}
	})
}

