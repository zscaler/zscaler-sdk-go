// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_profile"
)

func TestInspectionProfile_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileID := "profile-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/inspectionProfile/" + profileID

	server.On("GET", path, common.SuccessResponse(inspection_profile.InspectionProfile{
		ID:   profileID,
		Name: "Test Inspection Profile",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := inspection_profile.Get(context.Background(), service, profileID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, profileID, result.ID)
}

func TestInspectionProfile_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/inspectionProfile"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []inspection_profile.InspectionProfile{{ID: "profile-001"}, {ID: "profile-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := inspection_profile.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
