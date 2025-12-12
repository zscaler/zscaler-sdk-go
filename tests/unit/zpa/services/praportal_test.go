// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praportal"
)

func TestPRAPortal_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	portalID := "portal-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/praPortal/" + portalID

	server.On("GET", path, common.SuccessResponse(praportal.PRAPortal{
		ID:   portalID,
		Name: "Test Portal",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := praportal.Get(context.Background(), service, portalID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, portalID, result.ID)
}

func TestPRAPortal_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/praPortal"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []praportal.PRAPortal{{ID: "portal-001"}, {ID: "portal-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := praportal.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
