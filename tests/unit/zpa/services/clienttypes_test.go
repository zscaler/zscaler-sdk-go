// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/clienttypes"
)

func TestClientTypes_GetAllClientTypes_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// Correct path: /zpa/mgmtconfig/v1/admin/customers/{customerId}/clientTypes
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/clientTypes"

	server.On("GET", path, common.SuccessResponse(clienttypes.ClientTypes{
		ZPNClientTypeExplorer:      "zpn_client_type_exporter",
		ZPNClientTypeMachineTunnel: "zpn_client_type_machine_tunnel",
		ZPNClientTypeIPAnchoring:   "zpn_client_type_ip_anchoring",
		ZPNClientTypeEdgeConnector: "zpn_client_type_edge_connector",
		ZPNClientTypeZAPP:          "zpn_client_type_zapp",
		ZPNClientTypeSlogger:       "zpn_client_type_slogger",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := clienttypes.GetAllClientTypes(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.ZPNClientTypeExplorer)
}
