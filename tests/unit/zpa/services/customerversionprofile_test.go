// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/customerversionprofile"
)

func TestCustomerVersionProfile_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileName := "Default"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/visible/versionProfiles"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []customerversionprofile.CustomerVersionProfile{
			{ID: "profile-001", Name: "Other Profile"},
			{ID: "profile-002", Name: profileName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := customerversionprofile.GetByName(context.Background(), service, profileName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "profile-002", result.ID)
	assert.Equal(t, profileName, result.Name)
}

func TestCustomerVersionProfile_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/visible/versionProfiles"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []customerversionprofile.CustomerVersionProfile{
			{ID: "profile-001", Name: "Profile 1"},
			{ID: "profile-002", Name: "Profile 2"},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := customerversionprofile.GetAll(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 2)
}

func TestCustomerVersionProfile_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileID := "profile-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/visible/versionProfiles/" + profileID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateProfile := &customerversionprofile.CustomerVersionProfile{
		ID:   profileID,
		Name: "Updated Profile",
	}

	resp, err := customerversionprofile.Update(context.Background(), service, profileID, updateProfile)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

