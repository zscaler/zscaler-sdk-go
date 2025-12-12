// Package unit provides unit tests for ZPA Inspection Profile service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// InspectionProfile represents the inspection profile for testing
type InspectionProfile struct {
	ID                            string                         `json:"id,omitempty"`
	Name                          string                         `json:"name,omitempty"`
	Description                   string                         `json:"description,omitempty"`
	APIProfile                    bool                           `json:"apiProfile,omitempty"`
	OverrideAction                string                         `json:"overrideAction,omitempty"`
	ZSDefinedControlChoice        string                         `json:"zsDefinedControlChoice,omitempty"`
	GlobalControlActions          []string                       `json:"globalControlActions,omitempty"`
	IncarnationNumber             string                         `json:"incarnationNumber,omitempty"`
	ParanoiaLevel                 string                         `json:"paranoiaLevel,omitempty"`
	PredefinedControlsVersion     string                         `json:"predefinedControlsVersion,omitempty"`
	CheckControlDeploymentStatus  bool                           `json:"checkControlDeploymentStatus,omitempty"`
	ControlInfoResource           []ControlInfoResource          `json:"controlsInfo,omitempty"`
	CustomControls                []InspProfileCustomControl     `json:"customControls,omitempty"`
	WebSocketControls             []WebSocketControls            `json:"websocketControls,omitempty"`
	ThreatLabzControls            []ThreatLabzControls           `json:"threatlabzControls,omitempty"`
	CreationTime                  string                         `json:"creationTime,omitempty"`
	ModifiedBy                    string                         `json:"modifiedBy,omitempty"`
	ModifiedTime                  string                         `json:"modifiedTime,omitempty"`
}

// ControlInfoResource represents control info
type ControlInfoResource struct {
	ControlType string `json:"controlType,omitempty"`
	Count       string `json:"count,omitempty"`
}

// InspProfileCustomControl represents a custom control in the profile
type InspProfileCustomControl struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Action        string `json:"action,omitempty"`
	ActionValue   string `json:"actionValue,omitempty"`
	ControlNumber string `json:"controlNumber,omitempty"`
	ControlType   string `json:"controlType,omitempty"`
	Severity      string `json:"severity,omitempty"`
}

// WebSocketControls represents websocket controls
type WebSocketControls struct {
	ID                     string `json:"id,omitempty"`
	Name                   string `json:"name,omitempty"`
	Action                 string `json:"action,omitempty"`
	ActionValue            string `json:"actionValue,omitempty"`
	ControlNumber          string `json:"controlNumber,omitempty"`
	ControlType            string `json:"controlType,omitempty"`
	Severity               string `json:"severity,omitempty"`
	ZSDefinedControlChoice string `json:"zsDefinedControlChoice,omitempty"`
}

// ThreatLabzControls represents ThreatLabz controls
type ThreatLabzControls struct {
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	Description         string `json:"description,omitempty"`
	Enabled             bool   `json:"enabled,omitempty"`
	Action              string `json:"action,omitempty"`
	ControlGroup        string `json:"controlGroup,omitempty"`
	ControlNumber       string `json:"controlNumber,omitempty"`
	Severity            string `json:"severity,omitempty"`
	EngineVersion       string `json:"engineVersion,omitempty"`
	RuleDeploymentState string `json:"ruleDeploymentState,omitempty"`
	RulesetName         string `json:"rulesetName,omitempty"`
	RulesetVersion      string `json:"rulesetVersion,omitempty"`
}

// TestInspectionProfile_Structure tests the struct definitions
func TestInspectionProfile_Structure(t *testing.T) {
	t.Parallel()

	t.Run("InspectionProfile JSON marshaling", func(t *testing.T) {
		profile := InspectionProfile{
			ID:                        "ip-123",
			Name:                      "Web App Protection",
			Description:               "Profile for web application protection",
			ParanoiaLevel:             "1",
			PredefinedControlsVersion: "OWASP_CRS/3.3.0",
			ZSDefinedControlChoice:    "ALL",
			GlobalControlActions:      []string{"BLOCK", "REDIRECT"},
			ControlInfoResource: []ControlInfoResource{
				{ControlType: "PREDEFINED", Count: "100"},
				{ControlType: "CUSTOM", Count: "5"},
			},
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		var unmarshaled InspectionProfile
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, profile.ID, unmarshaled.ID)
		assert.Equal(t, profile.Name, unmarshaled.Name)
		assert.Equal(t, "1", unmarshaled.ParanoiaLevel)
		assert.Len(t, unmarshaled.ControlInfoResource, 2)
	})

	t.Run("InspectionProfile from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "ip-456",
			"name": "API Protection Profile",
			"description": "Profile for API protection",
			"apiProfile": true,
			"overrideAction": "BLOCK",
			"zsDefinedControlChoice": "OWASP",
			"globalControlActions": ["BLOCK"],
			"incarnationNumber": "1",
			"paranoiaLevel": "2",
			"predefinedControlsVersion": "OWASP_CRS/3.3.0",
			"checkControlDeploymentStatus": true,
			"controlsInfo": [
				{"controlType": "PREDEFINED", "count": "150"},
				{"controlType": "WEBSOCKET", "count": "10"},
				{"controlType": "THREATLABZ", "count": "25"}
			],
			"customControls": [
				{"id": "cc-1", "name": "Custom Rule 1", "action": "BLOCK"}
			],
			"websocketControls": [
				{"id": "ws-1", "name": "WS Control 1", "action": "PASS"}
			],
			"threatlabzControls": [
				{"id": "tl-1", "name": "TL Control 1", "enabled": true}
			],
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var profile InspectionProfile
		err := json.Unmarshal([]byte(apiResponse), &profile)
		require.NoError(t, err)

		assert.Equal(t, "ip-456", profile.ID)
		assert.True(t, profile.APIProfile)
		assert.Equal(t, "OWASP", profile.ZSDefinedControlChoice)
		assert.Len(t, profile.CustomControls, 1)
		assert.Len(t, profile.WebSocketControls, 1)
		assert.Len(t, profile.ThreatLabzControls, 1)
	})
}

// TestInspectionProfile_MockServerOperations tests CRUD operations
func TestInspectionProfile_MockServerOperations(t *testing.T) {
	t.Run("GET profile by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/inspectionProfile/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "ip-123", "name": "Mock Profile"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/inspectionProfile/ip-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all profiles", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"list": [{"id": "1"}, {"id": "2"}], "totalPages": 1}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/inspectionProfile")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create profile", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "new-ip", "name": "New Profile"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/inspectionProfile", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT associate all predefined controls", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Contains(t, r.URL.Path, "/associateAllPredefinedControls")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/inspectionProfile/ip-123/associateAllPredefinedControls", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("PUT deassociate all predefined controls", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Contains(t, r.URL.Path, "/deAssociateAllPredefinedControls")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/inspectionProfile/ip-123/deAssociateAllPredefinedControls", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("PATCH profile", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PATCH", r.Method)
			assert.Contains(t, r.URL.Path, "/patch")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PATCH", server.URL+"/inspectionProfile/ip-123/patch", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE profile", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/inspectionProfile/ip-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestInspectionProfile_SpecialCases tests edge cases
func TestInspectionProfile_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("ZS defined control choices", func(t *testing.T) {
		choices := []string{"ALL", "OWASP", "NONE"}

		for _, choice := range choices {
			profile := InspectionProfile{
				ID:                     "ip-" + choice,
				Name:                   choice + " Profile",
				ZSDefinedControlChoice: choice,
			}

			data, err := json.Marshal(profile)
			require.NoError(t, err)

			assert.Contains(t, string(data), choice)
		}
	})

	t.Run("API profile", func(t *testing.T) {
		profile := InspectionProfile{
			ID:         "ip-api",
			Name:       "API Profile",
			APIProfile: true,
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"apiProfile":true`)
	})

	t.Run("ThreatLabz controls", func(t *testing.T) {
		tl := ThreatLabzControls{
			ID:                  "tl-123",
			Name:                "Malware Detection",
			Enabled:             true,
			Action:              "BLOCK",
			RuleDeploymentState: "DEPLOYED",
			RulesetName:         "ThreatLabz Ruleset",
			RulesetVersion:      "1.0.0",
		}

		data, err := json.Marshal(tl)
		require.NoError(t, err)

		var unmarshaled ThreatLabzControls
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.True(t, unmarshaled.Enabled)
		assert.Equal(t, "DEPLOYED", unmarshaled.RuleDeploymentState)
	})
}

