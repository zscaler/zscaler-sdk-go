// Package unit provides unit tests for ZPA Client Types service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ClientTypes represents all client types for testing
type ClientTypes struct {
	ZPNClientTypeExplorer         string `json:"zpn_client_type_exporter"`
	ZPNClientTypeNoAuth           string `json:"zpn_client_type_exporter_noauth"`
	ZPNClientTypeBrowserIsolation string `json:"zpn_client_type_browser_isolation"`
	ZPNClientTypeMachineTunnel    string `json:"zpn_client_type_machine_tunnel"`
	ZPNClientTypeIPAnchoring      string `json:"zpn_client_type_ip_anchoring"`
	ZPNClientTypeEdgeConnector    string `json:"zpn_client_type_edge_connector"`
	ZPNClientTypeZAPP             string `json:"zpn_client_type_zapp"`
	ZPNClientTypeSlogger          string `json:"zpn_client_type_slogger"`
	ZPNClientTypeBranchConnector  string `json:"zpn_client_type_branch_connector"`
	ZPNClientTypePartner          string `json:"zpn_client_type_zapp_partner"`
	ZPNClientTypeVDI              string `json:"zpn_client_type_vdi"`
	ZPNClientTypeZIAInspection    string `json:"zpn_client_type_zia_inspection"`
}

// TestClientTypes_Structure tests the struct definitions
func TestClientTypes_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ClientTypes JSON marshaling", func(t *testing.T) {
		types := ClientTypes{
			ZPNClientTypeExplorer:         "zpn_client_type_exporter",
			ZPNClientTypeNoAuth:           "zpn_client_type_exporter_noauth",
			ZPNClientTypeBrowserIsolation: "zpn_client_type_browser_isolation",
			ZPNClientTypeMachineTunnel:    "zpn_client_type_machine_tunnel",
			ZPNClientTypeIPAnchoring:      "zpn_client_type_ip_anchoring",
			ZPNClientTypeEdgeConnector:    "zpn_client_type_edge_connector",
			ZPNClientTypeZAPP:             "zpn_client_type_zapp",
			ZPNClientTypeSlogger:          "zpn_client_type_slogger",
			ZPNClientTypeBranchConnector:  "zpn_client_type_branch_connector",
			ZPNClientTypePartner:          "zpn_client_type_zapp_partner",
			ZPNClientTypeVDI:              "zpn_client_type_vdi",
			ZPNClientTypeZIAInspection:    "zpn_client_type_zia_inspection",
		}

		data, err := json.Marshal(types)
		require.NoError(t, err)

		var unmarshaled ClientTypes
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, types.ZPNClientTypeZAPP, unmarshaled.ZPNClientTypeZAPP)
		assert.Equal(t, types.ZPNClientTypeBrowserIsolation, unmarshaled.ZPNClientTypeBrowserIsolation)
	})

	t.Run("ClientTypes from API response", func(t *testing.T) {
		apiResponse := `{
			"zpn_client_type_exporter": "Web Browser",
			"zpn_client_type_exporter_noauth": "Web Browser (Unauthenticated)",
			"zpn_client_type_browser_isolation": "Cloud Browser Isolation",
			"zpn_client_type_machine_tunnel": "Machine Tunnel",
			"zpn_client_type_ip_anchoring": "ZIA Service Edge",
			"zpn_client_type_edge_connector": "Cloud Connector",
			"zpn_client_type_zapp": "Client Connector",
			"zpn_client_type_slogger": "LSS Receiver",
			"zpn_client_type_branch_connector": "Branch Connector",
			"zpn_client_type_zapp_partner": "Client Connector Partner",
			"zpn_client_type_vdi": "VDI Client",
			"zpn_client_type_zia_inspection": "ZIA Inspection"
		}`

		var types ClientTypes
		err := json.Unmarshal([]byte(apiResponse), &types)
		require.NoError(t, err)

		assert.Equal(t, "Client Connector", types.ZPNClientTypeZAPP)
		assert.Equal(t, "Cloud Browser Isolation", types.ZPNClientTypeBrowserIsolation)
		assert.Equal(t, "Branch Connector", types.ZPNClientTypeBranchConnector)
		assert.Equal(t, "VDI Client", types.ZPNClientTypeVDI)
	})
}

// TestClientTypes_MockServerOperations tests operations
func TestClientTypes_MockServerOperations(t *testing.T) {
	t.Run("GET all client types", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/clientTypes")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"zpn_client_type_zapp": "Client Connector",
				"zpn_client_type_browser_isolation": "Cloud Browser Isolation",
				"zpn_client_type_machine_tunnel": "Machine Tunnel",
				"zpn_client_type_branch_connector": "Branch Connector"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/clientTypes")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestClientTypes_AllTypes tests all supported client types
func TestClientTypes_AllTypes(t *testing.T) {
	t.Parallel()

	t.Run("All client type fields", func(t *testing.T) {
		typeFields := []string{
			"zpn_client_type_exporter",
			"zpn_client_type_exporter_noauth",
			"zpn_client_type_browser_isolation",
			"zpn_client_type_machine_tunnel",
			"zpn_client_type_ip_anchoring",
			"zpn_client_type_edge_connector",
			"zpn_client_type_zapp",
			"zpn_client_type_slogger",
			"zpn_client_type_branch_connector",
			"zpn_client_type_zapp_partner",
			"zpn_client_type_vdi",
			"zpn_client_type_zia_inspection",
		}

		types := ClientTypes{
			ZPNClientTypeExplorer:         "exporter",
			ZPNClientTypeNoAuth:           "exporter_noauth",
			ZPNClientTypeBrowserIsolation: "browser_isolation",
			ZPNClientTypeMachineTunnel:    "machine_tunnel",
			ZPNClientTypeIPAnchoring:      "ip_anchoring",
			ZPNClientTypeEdgeConnector:    "edge_connector",
			ZPNClientTypeZAPP:             "zapp",
			ZPNClientTypeSlogger:          "slogger",
			ZPNClientTypeBranchConnector:  "branch_connector",
			ZPNClientTypePartner:          "zapp_partner",
			ZPNClientTypeVDI:              "vdi",
			ZPNClientTypeZIAInspection:    "zia_inspection",
		}

		data, err := json.Marshal(types)
		require.NoError(t, err)

		for _, field := range typeFields {
			assert.Contains(t, string(data), field)
		}
	})
}

