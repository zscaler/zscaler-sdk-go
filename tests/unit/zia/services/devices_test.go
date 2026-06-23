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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/devices"
)

const devicesResourcePath = "/zia/api/v1/devices"

// =====================================================
// Get
// =====================================================

func TestDevicesResource_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesResourcePath+"/100", common.SuccessResponse(devices.Devices{
		ID:       100,
		Name:     "DESKTOP-ABC",
		Active:   true,
		Os:       "Windows",
		Hostname: "desktop-abc.corp.example.com",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devices.Get(context.Background(), service, 100)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 100, result.ID)
	assert.Equal(t, "DESKTOP-ABC", result.Name)
	assert.True(t, result.Active)
	assert.Equal(t, "Windows", result.Os)
}

func TestDevicesResource_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesResourcePath+"/999", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devices.Get(context.Background(), service, 999)

	require.Error(t, err)
	assert.Nil(t, result)
}

// =====================================================
// GetAll
// =====================================================

func TestDevicesResource_GetAll_NoOpts_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", devicesResourcePath, func(r *http.Request, _ []byte) common.MockResponse {
		// No filter query params should be present.
		q := r.URL.Query()
		assert.Empty(t, q["id"])
		assert.Empty(t, q.Get("search"))
		assert.Empty(t, q.Get("valid"))
		assert.Empty(t, q.Get("includeCbiDevices"))
		return common.SuccessResponse([]devices.Devices{
			{ID: 1, Name: "Device A"},
			{ID: 2, Name: "Device B"},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devices.GetAll(context.Background(), service, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Device A", result[0].Name)
}

func TestDevicesResource_GetAll_EmptyOpts_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", devicesResourcePath, func(r *http.Request, _ []byte) common.MockResponse {
		q := r.URL.Query()
		assert.Empty(t, q["id"])
		assert.Empty(t, q.Get("search"))
		assert.Empty(t, q.Get("valid"))
		assert.Empty(t, q.Get("includeCbiDevices"))
		return common.SuccessResponse([]devices.Devices{{ID: 3, Name: "Device C"}})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devices.GetAll(context.Background(), service, &devices.GetAllFilterOptions{})

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, 3, result[0].ID)
}

func TestDevicesResource_GetAll_WithIDFilter_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", devicesResourcePath, func(r *http.Request, _ []byte) common.MockResponse {
		// IDs are sent as repeated query params.
		assert.ElementsMatch(t, []string{"10", "20"}, r.URL.Query()["id"])
		return common.SuccessResponse([]devices.Devices{
			{ID: 10, Name: "Device 10"},
			{ID: 20, Name: "Device 20"},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devices.GetAll(context.Background(), service, &devices.GetAllFilterOptions{
		ID: []int{10, 20},
	})

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestDevicesResource_GetAll_WithSearchFilter_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", devicesResourcePath, func(r *http.Request, _ []byte) common.MockResponse {
		assert.Equal(t, "DESKTOP", r.URL.Query().Get("search"))
		return common.SuccessResponse([]devices.Devices{{ID: 1, Name: "DESKTOP-1"}})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	search := "DESKTOP"
	result, err := devices.GetAll(context.Background(), service, &devices.GetAllFilterOptions{
		Search: &search,
	})

	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestDevicesResource_GetAll_EmptySearchIgnored_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", devicesResourcePath, func(r *http.Request, _ []byte) common.MockResponse {
		// An empty search string should not be added to the query.
		assert.Empty(t, r.URL.Query().Get("search"))
		return common.SuccessResponse([]devices.Devices{{ID: 1, Name: "Device A"}})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	empty := ""
	result, err := devices.GetAll(context.Background(), service, &devices.GetAllFilterOptions{
		Search: &empty,
	})

	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestDevicesResource_GetAll_WithValidFilter_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", devicesResourcePath, func(r *http.Request, _ []byte) common.MockResponse {
		assert.Equal(t, "true", r.URL.Query().Get("valid"))
		return common.SuccessResponse([]devices.Devices{{ID: 1, Name: "Valid Device"}})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	valid := true
	result, err := devices.GetAll(context.Background(), service, &devices.GetAllFilterOptions{
		Valid: &valid,
	})

	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestDevicesResource_GetAll_WithIncludeCbiDevicesFilter_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", devicesResourcePath, func(r *http.Request, _ []byte) common.MockResponse {
		assert.Equal(t, "false", r.URL.Query().Get("includeCbiDevices"))
		return common.SuccessResponse([]devices.Devices{{ID: 1, Name: "Device A"}})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	include := false
	result, err := devices.GetAll(context.Background(), service, &devices.GetAllFilterOptions{
		IncludeCbiDevices: &include,
	})

	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestDevicesResource_GetAll_WithAllFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", devicesResourcePath, func(r *http.Request, _ []byte) common.MockResponse {
		q := r.URL.Query()
		assert.ElementsMatch(t, []string{"7"}, q["id"])
		assert.Equal(t, "laptop", q.Get("search"))
		assert.Equal(t, "true", q.Get("valid"))
		assert.Equal(t, "true", q.Get("includeCbiDevices"))
		return common.SuccessResponse([]devices.Devices{{ID: 7, Name: "Laptop"}})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	search := "laptop"
	valid := true
	include := true
	result, err := devices.GetAll(context.Background(), service, &devices.GetAllFilterOptions{
		ID:                []int{7},
		Search:            &search,
		Valid:             &valid,
		IncludeCbiDevices: &include,
	})

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, 7, result[0].ID)
}

func TestDevicesResource_GetAll_Empty_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesResourcePath, common.SuccessResponse([]devices.Devices{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devices.GetAll(context.Background(), service, nil)

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestDevicesResource_GetAll_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesResourcePath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devices.GetAll(context.Background(), service, nil)

	require.Error(t, err)
	assert.Empty(t, result)
}

// =====================================================
// GetByName
// =====================================================

func TestDevicesResource_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesResourcePath, common.SuccessResponse([]devices.Devices{
		{ID: 1, Name: "Device A"},
		{ID: 2, Name: "Device B"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devices.GetByName(context.Background(), service, "Device B")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.ID)
	assert.Equal(t, "Device B", result.Name)
}

func TestDevicesResource_GetByName_CaseInsensitive_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesResourcePath, common.SuccessResponse([]devices.Devices{
		{ID: 1, Name: "Device A"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devices.GetByName(context.Background(), service, "device a")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
}

func TestDevicesResource_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesResourcePath, common.SuccessResponse([]devices.Devices{
		{ID: 1, Name: "Device A"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devices.GetByName(context.Background(), service, "NonExistent")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no device found with name")
}

func TestDevicesResource_GetByName_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", devicesResourcePath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devices.GetByName(context.Background(), service, "Device A")

	require.Error(t, err)
	assert.Nil(t, result)
}

// =====================================================
// Structure Tests
// =====================================================

func TestDevicesResource_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Devices JSON marshaling", func(t *testing.T) {
		device := devices.Devices{
			ID:         42,
			Name:       "MOBILE-1",
			Active:     true,
			Version:    "1.2.3",
			Hostname:   "mobile-1.corp.example.com",
			Vendor:     "Apple",
			Model:      "iPhone",
			Os:         "iOS",
			Udid:       "udid-123",
			MacAddress: "00:11:22:33:44:55",
		}

		data, err := json.Marshal(device)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":42`)
		assert.Contains(t, string(data), `"name":"MOBILE-1"`)
		assert.Contains(t, string(data), `"active":true`)
		assert.Contains(t, string(data), `"os":"iOS"`)
		assert.Contains(t, string(data), `"macAddress":"00:11:22:33:44:55"`)
	})

	t.Run("Devices JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 7,
			"name": "DESKTOP-7",
			"active": true,
			"os": "Windows",
			"user": {"id": 99, "name": "user@example.com"},
			"deleted": false,
			"rooted": false
		}`

		var device devices.Devices
		err := json.Unmarshal([]byte(jsonData), &device)
		require.NoError(t, err)

		assert.Equal(t, 7, device.ID)
		assert.Equal(t, "DESKTOP-7", device.Name)
		assert.True(t, device.Active)
		require.NotNil(t, device.User)
		assert.Equal(t, 99, device.User.ID)
		assert.Equal(t, "user@example.com", device.User.Name)
	})
}
