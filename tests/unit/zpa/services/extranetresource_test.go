// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	zpacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/extranet_resource"
)

func TestExtranetResource_GetPartner_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/extranetPartner"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []zpacommon.CommonSummary{{ID: "ext-001"}, {ID: "ext-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := extranet_resource.GetExtranetResourcePartner(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
