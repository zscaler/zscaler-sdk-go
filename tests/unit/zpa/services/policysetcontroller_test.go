// Package unit provides unit tests for ZPA Policy Set Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PolicySet represents the policy set structure for testing
type PolicySet struct {
	ID                string       `json:"id,omitempty"`
	Name              string       `json:"name,omitempty"`
	Description       string       `json:"description,omitempty"`
	PolicyType        string       `json:"policyType,omitempty"`
	Enabled           bool         `json:"enabled"`
	CreationTime      string       `json:"creationTime,omitempty"`
	ModifiedBy        string       `json:"modifiedBy,omitempty"`
	ModifiedTime      string       `json:"modifiedTime,omitempty"`
	MicroTenantID     string       `json:"microtenantId,omitempty"`
	MicroTenantName   string       `json:"microtenantName,omitempty"`
	Rules             []PolicyRule `json:"rules,omitempty"`
}

// PolicyRule represents a policy rule for testing
type PolicyRule struct {
	ID                       string        `json:"id,omitempty"`
	Name                     string        `json:"name,omitempty"`
	Description              string        `json:"description,omitempty"`
	Action                   string        `json:"action,omitempty"`
	ActionID                 string        `json:"actionId,omitempty"`
	PolicyType               string        `json:"policyType,omitempty"`
	RuleOrder                string        `json:"ruleOrder,omitempty"`
	Priority                 string        `json:"priority,omitempty"`
	Operator                 string        `json:"operator,omitempty"`
	CustomMsg                string        `json:"customMsg,omitempty"`
	PolicySetID              string        `json:"policySetId,omitempty"`
	ReauthIdleTimeout        string        `json:"reauthIdleTimeout,omitempty"`
	ReauthTimeout            string        `json:"reauthTimeout,omitempty"`
	ZpnCbiProfileID          string        `json:"zpnCbiProfileId,omitempty"`
	ZpnIsolationProfileID    string        `json:"zpnIsolationProfileId,omitempty"`
	ZpnInspectionProfileID   string        `json:"zpnInspectionProfileId,omitempty"`
	ZpnInspectionProfileName string        `json:"zpnInspectionProfileName,omitempty"`
	CreationTime             string        `json:"creationTime,omitempty"`
	ModifiedBy               string        `json:"modifiedBy,omitempty"`
	ModifiedTime             string        `json:"modifiedTime,omitempty"`
	MicroTenantID            string        `json:"microtenantId,omitempty"`
	Conditions               []Conditions  `json:"conditions,omitempty"`
	AppConnectorGroups       []IDNamePair  `json:"appConnectorGroups,omitempty"`
	AppServerGroups          []IDNamePair  `json:"appServerGroups,omitempty"`
	ServiceEdgeGroups        []IDNamePair  `json:"serviceEdgeGroups,omitempty"`
}

// Conditions represents policy conditions for testing
type Conditions struct {
	ID        string     `json:"id,omitempty"`
	Operator  string     `json:"operator,omitempty"`
	Negated   bool       `json:"negated"`
	Operands  []Operands `json:"operands,omitempty"`
}

// Operands represents policy operands for testing
type Operands struct {
	ID              string `json:"id,omitempty"`
	ObjectType      string `json:"objectType,omitempty"`
	LHS             string `json:"lhs,omitempty"`
	RHS             string `json:"rhs,omitempty"`
	IdpID           string `json:"idpId,omitempty"`
	Name            string `json:"name,omitempty"`
	EntryValuesJSON string `json:"entryValues,omitempty"`
}

// IDNamePair represents a simple ID/Name pair for testing
type IDNamePair struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// TestPolicySetController_Structure tests the struct definitions
func TestPolicySetController_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PolicySet JSON marshaling", func(t *testing.T) {
		policySet := PolicySet{
			ID:          "ps-123",
			Name:        "Access Policy",
			Description: "Test access policy",
			PolicyType:  "ACCESS_POLICY",
			Enabled:     true,
		}

		data, err := json.Marshal(policySet)
		require.NoError(t, err)

		var unmarshaled PolicySet
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, policySet.ID, unmarshaled.ID)
		assert.Equal(t, policySet.PolicyType, unmarshaled.PolicyType)
	})

	t.Run("PolicyRule JSON marshaling", func(t *testing.T) {
		rule := PolicyRule{
			ID:          "pr-123",
			Name:        "Allow Engineering",
			Description: "Allow engineering team access",
			Action:      "ALLOW",
			RuleOrder:   "1",
			Priority:    "1",
			Operator:    "AND",
			PolicyType:  "ACCESS_POLICY",
			Conditions: []Conditions{
				{
					ID:       "cond-1",
					Operator: "OR",
					Negated:  false,
					Operands: []Operands{
						{
							ObjectType: "APP",
							LHS:        "id",
							RHS:        "app-123",
						},
					},
				},
			},
		}

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		var unmarshaled PolicyRule
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, rule.ID, unmarshaled.ID)
		assert.Equal(t, rule.Action, unmarshaled.Action)
		assert.Len(t, unmarshaled.Conditions, 1)
	})

	t.Run("PolicyRule with all policy types", func(t *testing.T) {
		policyTypes := []string{
			"ACCESS_POLICY",
			"TIMEOUT_POLICY",
			"REAUTH_POLICY",
			"CLIENT_FORWARDING_POLICY",
			"BYPASS_POLICY",
			"ISOLATION_POLICY",
			"INSPECTION_POLICY",
			"CREDENTIAL_POLICY",
			"CAPABILITIES_POLICY",
			"CLIENTLESS_SESSION_PROTECTION_POLICY",
			"REDIRECTION_POLICY",
		}

		for _, policyType := range policyTypes {
			rule := PolicyRule{
				ID:         "pr-" + policyType,
				Name:       policyType + " Rule",
				PolicyType: policyType,
				Action:     "ALLOW",
			}

			data, err := json.Marshal(rule)
			require.NoError(t, err)

			var unmarshaled PolicyRule
			err = json.Unmarshal(data, &unmarshaled)
			require.NoError(t, err)

			assert.Equal(t, policyType, unmarshaled.PolicyType)
		}
	})
}

// TestPolicySetController_ResponseParsing tests parsing of various API responses
func TestPolicySetController_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse policy set with rules", func(t *testing.T) {
		response := `{
			"id": "ps-789",
			"name": "Access Policy",
			"description": "Main access policy",
			"policyType": "ACCESS_POLICY",
			"enabled": true,
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"rules": [
				{
					"id": "rule-1",
					"name": "Allow All Apps",
					"action": "ALLOW",
					"ruleOrder": "1",
					"conditions": []
				},
				{
					"id": "rule-2",
					"name": "Block Untrusted",
					"action": "DENY",
					"ruleOrder": "2",
					"conditions": []
				}
			]
		}`

		var policySet PolicySet
		err := json.Unmarshal([]byte(response), &policySet)
		require.NoError(t, err)

		assert.Equal(t, "ps-789", policySet.ID)
		assert.Equal(t, "ACCESS_POLICY", policySet.PolicyType)
		assert.Len(t, policySet.Rules, 2)
		assert.Equal(t, "ALLOW", policySet.Rules[0].Action)
		assert.Equal(t, "DENY", policySet.Rules[1].Action)
	})

	t.Run("Parse complex policy rule with conditions", func(t *testing.T) {
		response := `{
			"id": "rule-complex",
			"name": "Complex Rule",
			"action": "ALLOW",
			"operator": "AND",
			"conditions": [
				{
					"id": "cond-1",
					"operator": "OR",
					"negated": false,
					"operands": [
						{
							"objectType": "APP",
							"lhs": "id",
							"rhs": "app-001"
						},
						{
							"objectType": "APP",
							"lhs": "id",
							"rhs": "app-002"
						}
					]
				},
				{
					"id": "cond-2",
					"operator": "OR",
					"negated": false,
					"operands": [
						{
							"objectType": "SAML",
							"lhs": "department",
							"rhs": "Engineering",
							"idpId": "idp-001"
						}
					]
				}
			]
		}`

		var rule PolicyRule
		err := json.Unmarshal([]byte(response), &rule)
		require.NoError(t, err)

		assert.Equal(t, "Complex Rule", rule.Name)
		assert.Equal(t, "AND", rule.Operator)
		assert.Len(t, rule.Conditions, 2)
		assert.Len(t, rule.Conditions[0].Operands, 2)
		assert.Equal(t, "SAML", rule.Conditions[1].Operands[0].ObjectType)
	})
}

// TestPolicySetController_MockServerOperations tests CRUD operations with mock server
func TestPolicySetController_MockServerOperations(t *testing.T) {
	t.Run("GET policy set by type", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "ps-123",
				"name": "Access Policy",
				"policyType": "ACCESS_POLICY",
				"enabled": true,
				"rules": []
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/policySet/policyType/ACCESS_POLICY")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET policy rule by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/rule/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "rule-123",
				"name": "Test Rule",
				"action": "ALLOW",
				"ruleOrder": "1"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/policySet/ps-123/rule/rule-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create policy rule", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "new-rule-456",
				"name": "New Policy Rule",
				"action": "ALLOW",
				"ruleOrder": "1"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/rule", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update policy rule", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/rule/rule-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("PUT reorder policy rules", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Contains(t, r.URL.Path, "/reorder")

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/policySet/ps-123/rule/rule-123/reorder/2", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE policy rule", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/rule/rule-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestPolicySetController_ErrorHandling tests error scenarios
func TestPolicySetController_ErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("404 Policy Set Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"code": "NOT_FOUND", "message": "Policy set not found"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/policySet/nonexistent")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("400 Invalid Rule Configuration", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code": "INVALID_REQUEST", "message": "At least one condition is required"}`))
		}))
		defer server.Close()

		resp, _ := http.Post(server.URL+"/rule", "application/json", nil)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

// TestPolicySetController_SpecialCases tests edge cases
func TestPolicySetController_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Rule with negated condition", func(t *testing.T) {
		rule := PolicyRule{
			ID:     "rule-123",
			Name:   "Negated Rule",
			Action: "DENY",
			Conditions: []Conditions{
				{
					Operator: "OR",
					Negated:  true,
					Operands: []Operands{
						{
							ObjectType: "APP_GROUP",
							LHS:        "id",
							RHS:        "group-123",
						},
					},
				},
			},
		}

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		var unmarshaled PolicyRule
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.True(t, unmarshaled.Conditions[0].Negated)
	})

	t.Run("Timeout policy rule", func(t *testing.T) {
		rule := PolicyRule{
			ID:                "rule-timeout",
			Name:              "Idle Timeout Rule",
			PolicyType:        "TIMEOUT_POLICY",
			Action:            "RE_AUTH",
			ReauthIdleTimeout: "600",
			ReauthTimeout:     "172800",
		}

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"reauthIdleTimeout":"600"`)
		assert.Contains(t, string(data), `"reauthTimeout":"172800"`)
	})

	t.Run("Isolation policy rule", func(t *testing.T) {
		rule := PolicyRule{
			ID:                    "rule-isolation",
			Name:                  "CBI Isolation Rule",
			PolicyType:            "ISOLATION_POLICY",
			Action:                "ISOLATE",
			ZpnCbiProfileID:       "cbi-profile-123",
			ZpnIsolationProfileID: "iso-profile-123",
		}

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"zpnCbiProfileId":"cbi-profile-123"`)
		assert.Contains(t, string(data), `"zpnIsolationProfileId":"iso-profile-123"`)
	})

	t.Run("Inspection policy rule", func(t *testing.T) {
		rule := PolicyRule{
			ID:                       "rule-inspection",
			Name:                     "DLP Inspection Rule",
			PolicyType:               "INSPECTION_POLICY",
			Action:                   "INSPECT",
			ZpnInspectionProfileID:   "insp-profile-123",
			ZpnInspectionProfileName: "DLP Profile",
		}

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"zpnInspectionProfileId":"insp-profile-123"`)
	})

	t.Run("Rule with app connector groups", func(t *testing.T) {
		rule := PolicyRule{
			ID:     "rule-acg",
			Name:   "ACG Rule",
			Action: "ALLOW",
			AppConnectorGroups: []IDNamePair{
				{ID: "acg-1", Name: "Connector Group 1"},
				{ID: "acg-2", Name: "Connector Group 2"},
			},
		}

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		var unmarshaled PolicyRule
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.AppConnectorGroups, 2)
	})

	t.Run("Custom message in deny rule", func(t *testing.T) {
		rule := PolicyRule{
			ID:        "rule-deny",
			Name:      "Deny With Message",
			Action:    "DENY",
			CustomMsg: "Access denied. Please contact IT support.",
		}

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		assert.Contains(t, string(data), "Access denied")
	})
}

