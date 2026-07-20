// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource"
)

func TestEndpointResource_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("POST", dlpEndpointResourceBasePath, func(_ *http.Request, body []byte) common.MockResponse {
		assert.Contains(t, string(body), "Finance Printer")
		return common.SuccessResponse(endpoint_resource.EndpointResource{
			ID:      501,
			Name:    "Finance Printer",
			Channel: "PRINTING",
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := endpoint_resource.Create(context.Background(), service, &endpoint_resource.EndpointResource{
		Name:    "Finance Printer",
		Channel: "PRINTING",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 501, result.ID)
	assert.Equal(t, "Finance Printer", result.Name)
}

func TestEndpointResource_Create_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", dlpEndpointResourceBasePath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := endpoint_resource.Create(context.Background(), service, &endpoint_resource.EndpointResource{
		Name: "Finance Printer",
	})

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestEndpointResource_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceBasePath + "/501"

	server.On("PUT", path, common.SuccessResponse(endpoint_resource.EndpointResource{
		ID:   501,
		Name: "Updated Printer",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := endpoint_resource.Update(context.Background(), service, 501, &endpoint_resource.EndpointResource{
		ID:   501,
		Name: "Updated Printer",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Printer", result.Name)
}

func TestEndpointResource_Update_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", dlpEndpointResourceBasePath+"/501", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := endpoint_resource.Update(context.Background(), service, 501, &endpoint_resource.EndpointResource{ID: 501})

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestEndpointResource_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceBasePath + "/501"
	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = endpoint_resource.Delete(context.Background(), service, 501)

	require.NoError(t, err)
}

func TestEndpointResource_Delete_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("DELETE", dlpEndpointResourceBasePath+"/999", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = endpoint_resource.Delete(context.Background(), service, 999)

	require.Error(t, err)
}
