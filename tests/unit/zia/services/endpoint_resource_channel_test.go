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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource_channel"
)

const dlpEndpointResourceBasePath = "/zia/api/v1/dlpEndpointResource"

// =====================================================
// GetChannelList
// =====================================================

func TestEndpointResourceChannel_GetChannelList_NoOpts_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceBasePath + "/" + string(endpoint_resource_channel.ChannelPrinting)

	server.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
		q := r.URL.Query()
		// No filters set, so neither sortOrder nor name should be present.
		assert.Empty(t, q.Get("sortOrder"))
		assert.Empty(t, q.Get("name"))
		return common.SuccessResponse([]endpoint_resource.EndpointResource{
			{ID: 1, Name: "Printer A", Channel: string(endpoint_resource_channel.ChannelPrinting)},
			{ID: 2, Name: "Printer B", Channel: string(endpoint_resource_channel.ChannelPrinting)},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_resource_channel.GetChannelList(context.Background(), service, endpoint_resource_channel.ChannelPrinting, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Printer A", result[0].Name)
}

func TestEndpointResourceChannel_GetChannelList_WithFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceBasePath + "/" + string(endpoint_resource_channel.ChannelNetworkDriveTransfer)

	server.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
		q := r.URL.Query()
		assert.Equal(t, "asc", q.Get("sortOrder"))
		assert.Equal(t, "share", q.Get("name"))
		return common.SuccessResponse([]endpoint_resource.EndpointResource{
			{ID: 10, Name: "Network Share", Channel: string(endpoint_resource_channel.ChannelNetworkDriveTransfer)},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	opts := &endpoint_resource_channel.GetChannelFilterOptions{
		SortOrder: "asc",
		Name:      "share",
	}
	result, err := endpoint_resource_channel.GetChannelList(context.Background(), service, endpoint_resource_channel.ChannelNetworkDriveTransfer, opts)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, 10, result[0].ID)
}

func TestEndpointResourceChannel_GetChannelList_Empty_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceBasePath + "/" + string(endpoint_resource_channel.ChannelRemovableDriveTransfer)
	server.On("GET", path, common.SuccessResponse([]endpoint_resource.EndpointResource{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_resource_channel.GetChannelList(context.Background(), service, endpoint_resource_channel.ChannelRemovableDriveTransfer, nil)

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestEndpointResourceChannel_GetChannelList_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceBasePath + "/" + string(endpoint_resource_channel.ChannelPersonalCloudStorage)
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_resource_channel.GetChannelList(context.Background(), service, endpoint_resource_channel.ChannelPersonalCloudStorage, nil)

	require.Error(t, err)
	assert.Empty(t, result)
}

// =====================================================
// GetChannel
// =====================================================

func TestEndpointResourceChannel_GetChannel_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceBasePath + "/" + string(endpoint_resource_channel.ChannelPrinting) + "/123"

	server.On("GET", path, common.SuccessResponse(endpoint_resource.EndpointResource{
		ID:      123,
		Name:    "Printer 123",
		Channel: string(endpoint_resource_channel.ChannelPrinting),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_resource_channel.GetChannel(context.Background(), service, endpoint_resource_channel.ChannelPrinting, 123)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 123, result.ID)
	assert.Equal(t, "Printer 123", result.Name)
	assert.Equal(t, string(endpoint_resource_channel.ChannelPrinting), result.Channel)
}

func TestEndpointResourceChannel_GetChannel_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceBasePath + "/" + string(endpoint_resource_channel.ChannelPrinting) + "/999"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_resource_channel.GetChannel(context.Background(), service, endpoint_resource_channel.ChannelPrinting, 999)

	require.Error(t, err)
	assert.Nil(t, result)
}
