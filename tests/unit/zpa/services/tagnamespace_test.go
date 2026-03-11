// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/tag_controller/tag_namespace"
)

func TestTagNamespace_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	namespaceID := "ns-99999"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/namespace/" + namespaceID

	server.On("GET", path, common.SuccessResponse(tag_namespace.Namespace{
		ID:          namespaceID,
		Name:        "Test Namespace",
		Description: "Test description",
		Enabled:     true,
		Origin:      "CUSTOM",
		Type:        "STATIC",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, resp, err := tag_namespace.Get(context.Background(), service, namespaceID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, namespaceID, result.ID)
	assert.Equal(t, "Test Namespace", result.Name)
	assert.True(t, result.Enabled)
}

func TestTagNamespace_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	namespaceName := "Production Namespace"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/namespace/search"

	server.On("POST", path, common.SuccessResponse(map[string]interface{}{
		"list": []tag_namespace.Namespace{
			{ID: "ns-001", Name: "Other Namespace", Enabled: true},
			{ID: "ns-002", Name: namespaceName, Enabled: true},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := tag_namespace.GetByName(context.Background(), service, namespaceName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "ns-002", result.ID)
	assert.Equal(t, namespaceName, result.Name)
}

func TestTagNamespace_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/namespace"

	server.On("POST", path, common.SuccessResponse(tag_namespace.Namespace{
		ID:      "new-ns-123",
		Name:    "New Namespace",
		Enabled: true,
		Origin:  "CUSTOM",
		Type:    "STATIC",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := tag_namespace.Create(context.Background(), service, tag_namespace.Namespace{
		Name:    "New Namespace",
		Enabled: true,
		Origin:  "CUSTOM",
		Type:    "STATIC",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-ns-123", result.ID)
	assert.Equal(t, "New Namespace", result.Name)
}

func TestTagNamespace_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	namespaceID := "ns-99999"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/namespace/" + namespaceID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := tag_namespace.Update(context.Background(), service, namespaceID, &tag_namespace.Namespace{
		ID:          namespaceID,
		Name:        "Updated Namespace",
		Description: "Updated description",
	})

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestTagNamespace_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	namespaceID := "ns-99999"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/namespace/" + namespaceID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := tag_namespace.Delete(context.Background(), service, namespaceID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestTagNamespace_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/namespace/search"

	server.On("POST", path, common.SuccessResponse(map[string]interface{}{
		"list": []tag_namespace.Namespace{
			{ID: "ns-001", Name: "Namespace 1", Enabled: true},
			{ID: "ns-002", Name: "Namespace 2", Enabled: false},
			{ID: "ns-003", Name: "Namespace 3", Enabled: true},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := tag_namespace.GetAll(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 3)
}

func TestTagNamespace_UpdateStatus_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	namespaceID := "ns-99999"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/namespace/" + namespaceID + "/status"

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := tag_namespace.UpdateStatus(context.Background(), service, namespaceID, tag_namespace.UpdateStatusRequest{
		Enabled: false,
	})

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestTagNamespace_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/namespace/search"

	server.On("POST", path, common.SuccessResponse(map[string]interface{}{
		"list":       []tag_namespace.Namespace{},
		"totalPages": 0,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := tag_namespace.GetByName(context.Background(), service, "NonExistent")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no namespace named")
}

func TestTagNamespace_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	namespaceID := "nonexistent-id"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/namespace/" + namespaceID

	server.On("GET", path, common.MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"id": "resource.not.found", "message": "Resource not found"}`,
	})

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := tag_namespace.Get(context.Background(), service, namespaceID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
