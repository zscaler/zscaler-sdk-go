// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbiprofilecontroller"
)

func TestCBIProfileController_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileID := "profile-12345"
	// Correct path: /zpa/cbiconfig/cbi/api/customers/{customerId}/profiles/{id} (plural)
	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/profiles/" + profileID

	server.On("GET", path, common.SuccessResponse(cbiprofilecontroller.IsolationProfile{
		ID:   profileID,
		Name: "Test Profile",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := cbiprofilecontroller.Get(context.Background(), service, profileID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, profileID, result.ID)
}

func TestCBIProfileController_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// Correct path: /zpa/cbiconfig/cbi/api/customers/{customerId}/profiles (plural)
	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/profiles"

	server.On("GET", path, common.SuccessResponse([]cbiprofilecontroller.IsolationProfile{{ID: "profile-001"}, {ID: "profile-002"}}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := cbiprofilecontroller.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestCBIProfileController_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// Create uses singular endpoint "/profile"
	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/profile"

	server.On("POST", path, common.SuccessResponse(cbiprofilecontroller.IsolationProfile{
		ID:   "new-profile-123",
		Name: "New Profile",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newProfile := &cbiprofilecontroller.IsolationProfile{
		Name: "New Profile",
	}

	result, _, err := cbiprofilecontroller.Create(context.Background(), service, newProfile)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-profile-123", result.ID)
}

func TestCBIProfileController_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileID := "profile-12345"
	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/profiles/" + profileID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateProfile := &cbiprofilecontroller.IsolationProfile{
		ID:   profileID,
		Name: "Updated Profile",
	}

	resp, err := cbiprofilecontroller.Update(context.Background(), service, profileID, updateProfile)

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestCBIProfileController_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileID := "profile-12345"
	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/profiles/" + profileID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := cbiprofilecontroller.Delete(context.Background(), service, profileID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
