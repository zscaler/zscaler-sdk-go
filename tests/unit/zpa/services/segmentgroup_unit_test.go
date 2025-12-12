// Package unit provides unit tests for ZPA services
// These tests actually call the SDK functions and generate real code coverage
package unit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
)

// TestSegmentGroup_Get tests the segmentgroup.Get function
// This test actually calls the SDK function and generates coverage
func TestSegmentGroup_Get_SDK(t *testing.T) {
	// Create mock server
	server := common.NewTestServer()
	defer server.Close()

	segmentGroupID := "sg-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/segmentGroup/" + segmentGroupID

	// Register mock response using the MockResponse type from mocks.go
	server.On("GET", path, common.SuccessResponse(segmentgroup.SegmentGroup{
		ID:          segmentGroupID,
		Name:        "Test Segment Group",
		Description: "Test description",
		Enabled:     true,
		ConfigSpace: "DEFAULT",
	}))

	// Create test service - this uses the real SDK client with mocked HTTP
	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)
	require.NotNil(t, service)

	// Call actual SDK function - this generates coverage!
	result, resp, err := segmentgroup.Get(context.Background(), service, segmentGroupID)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, segmentGroupID, result.ID)
	assert.Equal(t, "Test Segment Group", result.Name)
	assert.True(t, result.Enabled)
}

// TestSegmentGroup_GetByName_SDK tests the segmentgroup.GetByName function
func TestSegmentGroup_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupName := "Production Group"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/segmentGroup"

	// Register mock response for list endpoint (GetByName calls GetAll internally)
	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []segmentgroup.SegmentGroup{
			{
				ID:          "sg-001",
				Name:        "Other Group",
				Enabled:     true,
				ConfigSpace: "DEFAULT",
			},
			{
				ID:          "sg-002",
				Name:        groupName,
				Enabled:     true,
				ConfigSpace: "DEFAULT",
			},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	// Call SDK function
	result, _, err := segmentgroup.GetByName(context.Background(), service, groupName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "sg-002", result.ID)
	assert.Equal(t, groupName, result.Name)
}

// TestSegmentGroup_Create_SDK tests the segmentgroup.Create function
func TestSegmentGroup_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/segmentGroup"

	// Register mock response
	server.On("POST", path, common.SuccessResponse(segmentgroup.SegmentGroup{
		ID:          "new-sg-123",
		Name:        "New Test Group",
		Description: "Created via unit test",
		Enabled:     true,
		ConfigSpace: "DEFAULT",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	// Create segment group
	newGroup := &segmentgroup.SegmentGroup{
		Name:        "New Test Group",
		Description: "Created via unit test",
		Enabled:     true,
	}

	result, _, err := segmentgroup.Create(context.Background(), service, newGroup)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-sg-123", result.ID)
	assert.Equal(t, "New Test Group", result.Name)

	// Verify request was made
	req := server.LastRequest()
	assert.Equal(t, "POST", req.Method)
}

// TestSegmentGroup_Update_SDK tests the segmentgroup.Update function
func TestSegmentGroup_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	segmentGroupID := "sg-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/segmentGroup/" + segmentGroupID

	// Register mock response for PUT
	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	// Update segment group
	updateGroup := &segmentgroup.SegmentGroup{
		ID:          segmentGroupID,
		Name:        "Updated Group Name",
		Description: "Updated description",
		Enabled:     false,
	}

	resp, err := segmentgroup.Update(context.Background(), service, segmentGroupID, updateGroup)

	require.NoError(t, err)
	assert.NotNil(t, resp)

	// Verify request
	req := server.LastRequest()
	assert.Equal(t, "PUT", req.Method)
}

// TestSegmentGroup_Delete_SDK tests the segmentgroup.Delete function
func TestSegmentGroup_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	segmentGroupID := "sg-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/segmentGroup/" + segmentGroupID

	// Register mock response
	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	// Delete segment group
	resp, err := segmentgroup.Delete(context.Background(), service, segmentGroupID)

	require.NoError(t, err)
	assert.NotNil(t, resp)

	// Verify request
	req := server.LastRequest()
	assert.Equal(t, "DELETE", req.Method)
}

// TestSegmentGroup_GetAll_SDK tests the segmentgroup.GetAll function
func TestSegmentGroup_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/segmentGroup"

	// Register mock response
	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []segmentgroup.SegmentGroup{
			{ID: "sg-001", Name: "Group 1", Enabled: true},
			{ID: "sg-002", Name: "Group 2", Enabled: false},
			{ID: "sg-003", Name: "Group 3", Enabled: true},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	// Get all segment groups
	result, _, err := segmentgroup.GetAll(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 3)
	assert.Equal(t, "sg-001", result[0].ID)
	assert.Equal(t, "Group 1", result[0].Name)
}

// TestSegmentGroup_GetByName_NotFound_SDK tests error handling when group not found
func TestSegmentGroup_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/segmentGroup"

	// Register empty list response
	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []segmentgroup.SegmentGroup{},
		"totalPages": 0,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	// Try to find non-existent group
	result, _, err := segmentgroup.GetByName(context.Background(), service, "NonExistent")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no application named")
}

// TestSegmentGroup_UpdateV2_SDK tests the segmentgroup.UpdateV2 function
func TestSegmentGroup_UpdateV2_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	segmentGroupID := "sg-12345"
	// Note: UpdateV2 uses v2 endpoint
	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/segmentGroup/" + segmentGroupID

	// Register mock response for PUT
	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	// Update segment group using V2 API
	updateGroup := &segmentgroup.SegmentGroup{
		ID:          segmentGroupID,
		Name:        "Updated Group V2",
		Description: "Updated via V2 API",
		Enabled:     true,
	}

	resp, err := segmentgroup.UpdateV2(context.Background(), service, segmentGroupID, updateGroup)

	require.NoError(t, err)
	assert.NotNil(t, resp)

	// Verify V2 endpoint was called
	req := server.LastRequest()
	assert.Equal(t, "PUT", req.Method)
	assert.Contains(t, req.Path, "/v2/")
}

// TestSegmentGroup_NotFound_SDK tests 404 error handling
func TestSegmentGroup_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	segmentGroupID := "nonexistent-id"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/segmentGroup/" + segmentGroupID

	// Register 404 response
	server.On("GET", path, common.MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"id": "resource.not.found", "message": "Resource not found"}`,
	})

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	// Try to get non-existent segment group
	result, _, err := segmentgroup.Get(context.Background(), service, segmentGroupID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
