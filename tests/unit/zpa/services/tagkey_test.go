// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	tag_key "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/tag_controller/tag_key"
)

const testNamespaceID = "ns-12345"

func TestTagKey_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	tagKeyID := "tk-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/namespace/" + testNamespaceID + "/tagKey/" + tagKeyID

	server.On("GET", path, common.SuccessResponse(tag_key.TagKey{
		ID:          tagKeyID,
		Name:        "Test Tag Key",
		Description: "Test description",
		Enabled:     true,
		Origin:      "CUSTOM",
		Type:        "STATIC",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, resp, err := tag_key.Get(context.Background(), service, testNamespaceID, tagKeyID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, tagKeyID, result.ID)
	assert.Equal(t, "Test Tag Key", result.Name)
	assert.True(t, result.Enabled)
}

func TestTagKey_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	tagKeyName := "Production Key"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/namespace/" + testNamespaceID + "/tagKey/search"

	server.On("POST", path, common.SuccessResponse(map[string]interface{}{
		"list": []tag_key.TagKey{
			{ID: "tk-001", Name: "Other Key", Enabled: true},
			{ID: "tk-002", Name: tagKeyName, Enabled: true},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := tag_key.GetByName(context.Background(), service, testNamespaceID, tagKeyName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "tk-002", result.ID)
	assert.Equal(t, tagKeyName, result.Name)
}

func TestTagKey_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/namespace/" + testNamespaceID + "/tagKey"

	server.On("POST", path, common.SuccessResponse(tag_key.TagKey{
		ID:      "new-tk-123",
		Name:    "New Tag Key",
		Enabled: true,
		Origin:  "CUSTOM",
		Type:    "STATIC",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := tag_key.Create(context.Background(), service, testNamespaceID, tag_key.TagKey{
		Name:      "New Tag Key",
		Enabled:   true,
		Origin:    "CUSTOM",
		Type:      "STATIC",
		TagValues: []tag_key.TagValue{{Name: "value1"}, {Name: "value2"}},
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-tk-123", result.ID)
	assert.Equal(t, "New Tag Key", result.Name)
}

func TestTagKey_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	tagKeyID := "tk-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/namespace/" + testNamespaceID + "/tagKey/" + tagKeyID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := tag_key.Update(context.Background(), service, testNamespaceID, tagKeyID, &tag_key.TagKey{
		ID:         tagKeyID,
		Name:       "Updated Tag Key",
		TagValues:  []tag_key.TagValue{},
	})

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestTagKey_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	tagKeyID := "tk-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/namespace/" + testNamespaceID + "/tagKey/" + tagKeyID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := tag_key.Delete(context.Background(), service, testNamespaceID, tagKeyID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestTagKey_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/namespace/" + testNamespaceID + "/tagKey/search"

	server.On("POST", path, common.SuccessResponse(map[string]interface{}{
		"list": []tag_key.TagKey{
			{ID: "tk-001", Name: "Key 1", Enabled: true},
			{ID: "tk-002", Name: "Key 2", Enabled: false},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := tag_key.GetAll(context.Background(), service, testNamespaceID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 2)
}

func TestTagKey_BulkUpdateStatus_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/namespace/" + testNamespaceID + "/tagKey/bulkUpdateStatus"

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := tag_key.BulkUpdateStatus(context.Background(), service, testNamespaceID, tag_key.BulkUpdateStatusRequest{
		Enabled:   false,
		TagKeyIDs: []string{"tk-001", "tk-002"},
	})

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestTagKey_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/namespace/" + testNamespaceID + "/tagKey/search"

	server.On("POST", path, common.SuccessResponse(map[string]interface{}{
		"list":       []tag_key.TagKey{},
		"totalPages": 0,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := tag_key.GetByName(context.Background(), service, testNamespaceID, "NonExistent")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no tag key named")
}

func TestTagKey_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	tagKeyID := "nonexistent-id"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/namespace/" + testNamespaceID + "/tagKey/" + tagKeyID

	server.On("GET", path, common.MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"id": "resource.not.found", "message": "Resource not found"}`,
	})

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := tag_key.Get(context.Background(), service, testNamespaceID, tagKeyID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
