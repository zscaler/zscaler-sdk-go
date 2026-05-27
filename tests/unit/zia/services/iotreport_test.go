package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/iotreport"
)

const (
	iotDeviceTypesPath        = "/zia/api/v1/iotDiscovery/deviceTypes"
	iotCategoriesPath         = "/zia/api/v1/iotDiscovery/categories"
	iotClassificationsPath    = "/zia/api/v1/iotDiscovery/classifications"
	iotDeviceListPath         = "/zia/api/v1/iotDiscovery/deviceList"
)

func TestIOTReport_GetDeviceTypes_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", iotDeviceTypesPath, common.SuccessResponse(iotreport.CommonIOTReport{
		UUID: "device-type-uuid-1",
		Name: "Camera",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := iotreport.GetDeviceTypes(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "device-type-uuid-1", result.UUID)
	assert.Equal(t, "Camera", result.Name)
}

func TestIOTReport_GetIOTCategories_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", iotCategoriesPath, common.SuccessResponse(iotreport.CommonIOTReport{
		UUID:       "category-uuid-1",
		Name:       "Surveillance",
		ParentUuid: "device-type-uuid-1",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := iotreport.GetIOTCategories(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Surveillance", result.Name)
}

func TestIOTReport_GetIOTClassifications_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", iotClassificationsPath, common.SuccessResponse(iotreport.CommonIOTReport{
		UUID:       "classification-uuid-1",
		Name:       "Trusted",
		ParentUuid: "category-uuid-1",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := iotreport.GetIOTClassifications(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Trusted", result.Name)
}

func TestIOTReport_GetIOTDeviceList_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", iotDeviceListPath, common.SuccessResponse(iotreport.IOTDeviceList{
		CloudName:  "zscaler.net",
		CustomerID: 123456,
		Devices: []iotreport.Device{
			{
				LocationID:         "100",
				DeviceUUID:         "device-uuid-1",
				IPAddress:          "192.168.1.50",
				DeviceTypeUUID:     "device-type-uuid-1",
				AutoLabel:          "IP Camera",
				ClassificationUUID: "classification-uuid-1",
				CategoryUUID:       "category-uuid-1",
				FlowStartTime:      1699000000,
				FlowEndTime:        1699003600,
			},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := iotreport.GetIOTDeviceList(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "zscaler.net", result.CloudName)
	assert.Len(t, result.Devices, 1)
	assert.Equal(t, "192.168.1.50", result.Devices[0].IPAddress)
}

func TestIOTReport_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Device JSON marshaling", func(t *testing.T) {
		device := iotreport.Device{
			LocationID:  "100",
			DeviceUUID:  "device-uuid-1",
			IPAddress:   "192.168.1.50",
			AutoLabel:   "IP Camera",
			FlowStartTime: 1699000000,
		}

		data, err := json.Marshal(device)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"ipAddress":"192.168.1.50"`)
	})
}
