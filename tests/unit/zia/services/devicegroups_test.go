// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/devicegroups"
)

const (
	deviceGroupsPath = "/zia/api/v1/deviceGroups"
	devicesPath      = "/zia/api/v1/deviceGroups/devices"
)

// sampleIOSDeviceGroup mirrors integration test expectations (GroupType ZCC_OS, OSType IOS, predefined).
func sampleIOSDeviceGroup(name string) devicegroups.DeviceGroups {
	return devicegroups.DeviceGroups{
		ID:          1,
		Name:        name,
		GroupType:   "ZCC_OS",
		Description: "Predefined iOS device group",
		OSType:      "IOS",
		Predefined:  true,
		DeviceCount: 150,
	}
}

// sampleDevice mirrors typical device fields exercised in integration tests.
func sampleDevice(name string) devicegroups.Devices {
	return devicegroups.Devices{
		ID:              100,
		Name:            name,
		DeviceGroupType: "ZCC_OS",
		DeviceModel:     "iPhone 15",
		OSType:          "IOS",
		OSVersion:       "17.0",
		Description:     "Corporate iPhone",
		OwnerUserId:     500,
		OwnerName:       "john@company.com",
		HostName:        "iphone-001",
	}
}

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestDeviceGroups_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", deviceGroupsPath, common.SuccessResponse([]devicegroups.DeviceGroups{
		sampleIOSDeviceGroup("IOS"),
		{ID: 2, Name: "ANDROID_OS", GroupType: "ZCC_OS", OSType: "ANDROID_OS", Predefined: true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetAllDevicesGroups(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "ZCC_OS", result[0].GroupType)
	assert.Equal(t, "IOS", result[0].OSType)
	assert.True(t, result[0].Predefined)
}

func TestDeviceGroups_GetAll_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", deviceGroupsPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetAllDevicesGroups(context.Background(), service)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDeviceGroups_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", deviceGroupsPath, common.SuccessResponse([]devicegroups.DeviceGroups{
		sampleIOSDeviceGroup("IOS"),
		{ID: 2, Name: "WINDOWS_OS", GroupType: "ZCC_OS", OSType: "WINDOWS_OS", Predefined: true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetDeviceGroupByName(context.Background(), service, "IOS")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "IOS", result.Name)
	assert.Equal(t, "ZCC_OS", result.GroupType)
	assert.True(t, result.Predefined)
}

func TestDeviceGroups_GetByName_CaseInsensitive_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", deviceGroupsPath, common.SuccessResponse([]devicegroups.DeviceGroups{
		sampleIOSDeviceGroup("IOS"),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetDeviceGroupByName(context.Background(), service, "ios")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "IOS", result.Name)
}

func TestDeviceGroups_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", deviceGroupsPath, common.SuccessResponse([]devicegroups.DeviceGroups{
		sampleIOSDeviceGroup("IOS"),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetDeviceGroupByName(context.Background(), service, "ThisDeviceGroupDoesNotExist")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no device group found with name: ThisDeviceGroupDoesNotExist")
}

func TestDeviceGroups_GetByName_APIError_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", deviceGroupsPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetDeviceGroupByName(context.Background(), service, "IOS")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDevices_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesPath, common.SuccessResponse([]devicegroups.Devices{
		sampleDevice("Device 1"),
		func() devicegroups.Devices {
			d := sampleDevice("Device 2")
			d.ID = 101
			d.DeviceModel = "Galaxy S24"
			d.OSType = "ANDROID_OS"
			return d
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetAllDevices(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestDevices_GetAll_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetAllDevices(context.Background(), service)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDevices_GetByID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	device := sampleDevice("Device 1")

	server.On("GET", devicesPath, common.SuccessResponse([]devicegroups.Devices{
		device,
		func() devicegroups.Devices {
			d := sampleDevice("Device 2")
			d.ID = 101
			return d
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetDevicesByID(context.Background(), service, 100)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 100, result.ID)
	assert.Equal(t, "iPhone 15", result.DeviceModel)
}

func TestDevices_GetByID_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesPath, common.SuccessResponse([]devicegroups.Devices{
		sampleDevice("Device 1"),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetDevicesByID(context.Background(), service, 9999)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no device found with ID: 9999")
}

func TestDevices_GetByID_APIError_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetDevicesByID(context.Background(), service, 100)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDevices_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesPath, common.SuccessResponse([]devicegroups.Devices{
		sampleDevice("Device 1"),
		func() devicegroups.Devices {
			d := sampleDevice("Device 2")
			d.ID = 101
			return d
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetDevicesByName(context.Background(), service, "Device 1")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Device 1", result.Name)
}

func TestDevices_GetByName_CaseInsensitive_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesPath, common.SuccessResponse([]devicegroups.Devices{
		sampleDevice("Device 1"),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetDevicesByName(context.Background(), service, "device 1")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Device 1", result.Name)
}

func TestDevices_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesPath, common.SuccessResponse([]devicegroups.Devices{
		sampleDevice("Device 1"),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetDevicesByName(context.Background(), service, "missing-device")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDevices_GetByModel_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesPath, common.SuccessResponse([]devicegroups.Devices{
		sampleDevice("Device 1"),
		func() devicegroups.Devices {
			d := sampleDevice("Device 2")
			d.ID = 101
			d.DeviceModel = "Galaxy S24"
			return d
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetDevicesByModel(context.Background(), service, "iPhone 15")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "iPhone 15", result.DeviceModel)
}

func TestDevices_GetByModel_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesPath, common.SuccessResponse([]devicegroups.Devices{
		sampleDevice("Device 1"),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetDevicesByModel(context.Background(), service, "Unknown Model")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDevices_GetByOwner_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesPath, common.SuccessResponse([]devicegroups.Devices{
		sampleDevice("Device 1"),
		func() devicegroups.Devices {
			d := sampleDevice("Device 2")
			d.ID = 101
			d.OwnerName = "jane@company.com"
			return d
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetDevicesByOwner(context.Background(), service, "john@company.com")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "john@company.com", result.OwnerName)
}

func TestDevices_GetByOwner_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesPath, common.SuccessResponse([]devicegroups.Devices{
		sampleDevice("Device 1"),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetDevicesByOwner(context.Background(), service, "nobody@company.com")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDevices_GetByOSType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesPath, common.SuccessResponse([]devicegroups.Devices{
		sampleDevice("Device 1"),
		func() devicegroups.Devices {
			d := sampleDevice("Device 2")
			d.ID = 101
			d.OSType = "ANDROID_OS"
			return d
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetDevicesByOSType(context.Background(), service, "IOS")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "IOS", result.OSType)
}

func TestDevices_GetByOSType_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesPath, common.SuccessResponse([]devicegroups.Devices{
		sampleDevice("Device 1"),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetDevicesByOSType(context.Background(), service, "UNKNOWN_OS")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDevices_GetByOSVersion_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesPath, common.SuccessResponse([]devicegroups.Devices{
		sampleDevice("Device 1"),
		func() devicegroups.Devices {
			d := sampleDevice("Device 2")
			d.ID = 101
			d.OSVersion = "14.0"
			return d
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetDevicesByOSVersion(context.Background(), service, "17.0")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "17.0", result.OSVersion)
}

func TestDevices_GetByOSVersion_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesPath, common.SuccessResponse([]devicegroups.Devices{
		sampleDevice("Device 1"),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetDevicesByOSVersion(context.Background(), service, "99.0")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDeviceGroups_GetIncludeDeviceInfo_BothTrue_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", deviceGroupsPath, func(r *http.Request, _ []byte) common.MockResponse {
		assert.Contains(t, r.URL.RawQuery, "includeDeviceInfo=true")
		assert.Contains(t, r.URL.RawQuery, "includePseudoGroups=true")
		return common.SuccessResponse([]devicegroups.DeviceGroups{
			sampleIOSDeviceGroup("IOS"),
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetIncludeDeviceInfo(context.Background(), service, true, true)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 1)
}

func TestDeviceGroups_GetIncludeDeviceInfo_IncludeDeviceInfoOnly_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", deviceGroupsPath, func(r *http.Request, _ []byte) common.MockResponse {
		assert.Contains(t, r.URL.RawQuery, "includeDeviceInfo=true")
		assert.NotContains(t, r.URL.RawQuery, "includePseudoGroups=true")
		return common.SuccessResponse([]devicegroups.DeviceGroups{
			sampleIOSDeviceGroup("IOS"),
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetIncludeDeviceInfo(context.Background(), service, true, false)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestDeviceGroups_GetIncludeDeviceInfo_IncludePseudoGroupsOnly_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", deviceGroupsPath, func(r *http.Request, _ []byte) common.MockResponse {
		assert.Contains(t, r.URL.RawQuery, "includePseudoGroups=true")
		assert.NotContains(t, r.URL.RawQuery, "includeDeviceInfo=true")
		return common.SuccessResponse([]devicegroups.DeviceGroups{
			sampleIOSDeviceGroup("IOS"),
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetIncludeDeviceInfo(context.Background(), service, false, true)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestDeviceGroups_GetIncludeDeviceInfo_None_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", deviceGroupsPath, func(r *http.Request, _ []byte) common.MockResponse {
		assert.Empty(t, r.URL.RawQuery)
		return common.SuccessResponse([]devicegroups.DeviceGroups{
			sampleIOSDeviceGroup("IOS"),
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetIncludeDeviceInfo(context.Background(), service, false, false)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestDeviceGroups_GetIncludeDeviceInfo_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", deviceGroupsPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetIncludeDeviceInfo(context.Background(), service, true, true)

	require.Error(t, err)
	assert.Nil(t, result)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestDeviceGroups_Structure(t *testing.T) {
	t.Parallel()

	t.Run("DeviceGroups JSON marshaling", func(t *testing.T) {
		group := sampleIOSDeviceGroup("IOS")

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"name":"IOS"`)
		assert.Contains(t, string(data), `"groupType":"ZCC_OS"`)
		assert.Contains(t, string(data), `"osType":"IOS"`)
		assert.Contains(t, string(data), `"predefined":true`)
	})

	t.Run("DeviceGroups JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Mobile Devices",
			"groupType": "ZCC_OS",
			"description": "All mobile devices",
			"osType": "IOS",
			"predefined": true,
			"deviceCount": 150
		}`

		var group devicegroups.DeviceGroups
		err := json.Unmarshal([]byte(jsonData), &group)
		require.NoError(t, err)

		assert.Equal(t, 54321, group.ID)
		assert.Equal(t, "Mobile Devices", group.Name)
		assert.True(t, group.Predefined)
	})

	t.Run("Devices JSON marshaling", func(t *testing.T) {
		device := sampleDevice("iphone-001")

		data, err := json.Marshal(device)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"name":"iphone-001"`)
		assert.Contains(t, string(data), `"deviceModel":"iPhone 15"`)
		assert.Contains(t, string(data), `"osType":"IOS"`)
		assert.Contains(t, string(data), `"ownerName":"john@company.com"`)
	})
}
