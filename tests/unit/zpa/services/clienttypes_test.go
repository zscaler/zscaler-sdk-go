// Package unit provides unit tests for ZPA Client Types service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/clienttypes"
)

func TestClientTypes_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ClientTypes JSON marshaling", func(t *testing.T) {
		types := clienttypes.ClientTypes{
			ZPNClientTypeExplorer:         "zpn_client_type_exporter",
			ZPNClientTypeBrowserIsolation: "zpn_client_type_browser_isolation",
			ZPNClientTypeIPAnchoring:      "zpn_client_type_ip_anchoring",
			ZPNClientTypeEdgeConnector:    "zpn_client_type_edge_connector",
			ZPNClientTypeMachineTunnel:    "zpn_client_type_machine_tunnel",
			ZPNClientTypeZAPP:             "zpn_client_type_zapp",
		}

		data, err := json.Marshal(types)
		require.NoError(t, err)

		var unmarshaled clienttypes.ClientTypes
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, types.ZPNClientTypeZAPP, unmarshaled.ZPNClientTypeZAPP)
	})
}

func TestClientTypes_MockServerOperations(t *testing.T) {
	t.Run("GET all client types", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"zpn_client_type_zapp": "zpn_client_type_zapp"}]`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/clientTypes")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
