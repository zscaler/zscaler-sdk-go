// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/isolationprofile"
)

func TestIsolationProfile_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/isolation/profiles"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []isolationprofile.IsolationProfile{{ID: "iso-001", Name: "Test Profile"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := isolationprofile.GetByName(context.Background(), service, "Test Profile")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "iso-001", result.ID)
}

func TestIsolationProfile_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/isolation/profiles"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []isolationprofile.IsolationProfile{{ID: "iso-001"}, {ID: "iso-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := isolationprofile.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
