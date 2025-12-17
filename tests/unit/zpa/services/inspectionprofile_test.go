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

func TestInspectionProfile_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileName := "Production Profile"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/inspectionProfile"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []inspection_profile.InspectionProfile{
			{ID: "profile-001", Name: "Other Profile"},
			{ID: "profile-002", Name: profileName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := inspection_profile.GetByName(context.Background(), service, profileName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "profile-002", result.ID)
	assert.Equal(t, profileName, result.Name)
}

func TestInspectionProfile_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/inspectionProfile"

	server.On("POST", path, common.SuccessResponse(inspection_profile.InspectionProfile{
		ID:   "new-profile-123",
		Name: "New Profile",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newProfile := inspection_profile.InspectionProfile{
		Name: "New Profile",
	}

	result, _, err := inspection_profile.Create(context.Background(), service, newProfile)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-profile-123", result.ID)
}

func TestInspectionProfile_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileID := "profile-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/inspectionProfile/" + profileID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateProfile := &inspection_profile.InspectionProfile{
		ID:   profileID,
		Name: "Updated Profile",
	}

	resp, err := inspection_profile.Update(context.Background(), service, profileID, updateProfile)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestInspectionProfile_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileID := "profile-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/inspectionProfile/" + profileID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := inspection_profile.Delete(context.Background(), service, profileID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestInspectionProfile_Patch_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileID := "profile-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/inspectionProfile/" + profileID + "/patch"

	server.On("PATCH", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	patchProfile := &inspection_profile.InspectionProfile{
		ID:   profileID,
		Name: "Patched Profile",
	}

	resp, err := inspection_profile.Patch(context.Background(), service, profileID, patchProfile)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
