// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/platforms"
)

func TestPlatforms_GetAllPlatforms_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/platform"

	server.On("GET", path, common.SuccessResponse(platforms.Platforms{
		Linux:   "true",
		Windows: "true",
		MacOS:   "true",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := platforms.GetAllPlatforms(context.Background(), service)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "true", result.Linux)
}
