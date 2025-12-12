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

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/zpnClientTypes"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"zpnClientTypesMap": map[string]interface{}{
			"zpn_client_type_ip_anchoring":      "IP Anchoring Client",
			"zpn_client_type_edge_connector":    "Cloud Connector",
			"zpn_client_type_browser_isolation": "Browser Isolation",
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := clienttypes.GetAllClientTypes(context.Background(), service)

	require.NoError(t, err)
	assert.NotNil(t, result)
}
