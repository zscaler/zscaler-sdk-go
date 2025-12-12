// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/postureprofile"
)

func TestPostureProfile_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileID := "posture-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/posture/" + profileID

	server.On("GET", path, common.SuccessResponse(postureprofile.PostureProfile{
		ID:          profileID,
		Name:        "Test Posture Profile",
		Platform:    "windows",
		PostureType: "DEVICE",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, resp, err := postureprofile.Get(context.Background(), service, profileID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, profileID, result.ID)
	assert.Equal(t, "Test Posture Profile", result.Name)
}

func TestPostureProfile_GetByPostureUDID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	postureUDID := "udid-12345"
	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/posture"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []postureprofile.PostureProfile{
			{ID: "pp-001", Name: "Other Profile", PostureudID: "udid-001"},
			{ID: "pp-002", Name: "Target Profile", PostureudID: postureUDID},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := postureprofile.GetByPostureUDID(context.Background(), service, postureUDID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "pp-002", result.ID)
	assert.Equal(t, postureUDID, result.PostureudID)
}

func TestPostureProfile_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileName := "Production Profile"
	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/posture"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []postureprofile.PostureProfile{
			{ID: "pp-001", Name: "Other Profile"},
			{ID: "pp-002", Name: profileName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := postureprofile.GetByName(context.Background(), service, profileName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "pp-002", result.ID)
	assert.Equal(t, profileName, result.Name)
}

func TestPostureProfile_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/posture"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []postureprofile.PostureProfile{
			{ID: "pp-001", Name: "Profile 1", Platform: "windows"},
			{ID: "pp-002", Name: "Profile 2", Platform: "macos"},
			{ID: "pp-003", Name: "Profile 3", Platform: "linux"},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := postureprofile.GetAll(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 3)
}

func TestPostureProfile_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/posture"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []postureprofile.PostureProfile{},
		"totalPages": 0,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := postureprofile.GetByName(context.Background(), service, "NonExistent")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no posture profile named")
}

func TestPostureProfile_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileID := "nonexistent-id"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/posture/" + profileID

	server.On("GET", path, common.MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"id": "resource.not.found", "message": "Resource not found"}`,
	})

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := postureprofile.Get(context.Background(), service, profileID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
