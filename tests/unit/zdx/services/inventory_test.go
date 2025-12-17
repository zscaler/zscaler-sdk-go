// Package services provides unit tests for ZDX inventory service
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/inventory"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestInventory_GetSoftware_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/inventory/software"

	server.On("GET", path, common.SuccessResponse(inventory.SoftwareOverviewResponse{
		Software: []inventory.SoftwareOverview{
			{SoftwareKey: "chrome", SoftwareName: "Google Chrome", Vendor: "Google", UserTotal: 500, DeviceTotal: 450},
			{SoftwareKey: "teams", SoftwareName: "Microsoft Teams", Vendor: "Microsoft", UserTotal: 400, DeviceTotal: 380},
			{SoftwareKey: "zoom", SoftwareName: "Zoom", Vendor: "Zoom Video", UserTotal: 350, DeviceTotal: 340},
		},
		NextOffset: "",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, nextOffset, _, err := inventory.GetSoftware(context.Background(), service, inventory.GetSoftwareFilters{})

	require.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "Google Chrome", result[0].SoftwareName)
	assert.Equal(t, "Google", result[0].Vendor)
	assert.Empty(t, nextOffset)
}

func TestInventory_GetSoftware_WithPagination_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/inventory/software"

	server.On("GET", path, common.SuccessResponse(inventory.SoftwareOverviewResponse{
		Software: []inventory.SoftwareOverview{
			{SoftwareKey: "software1", SoftwareName: "Software 1", Vendor: "Vendor 1"},
		},
		NextOffset: "page2token",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, nextOffset, _, err := inventory.GetSoftware(context.Background(), service, inventory.GetSoftwareFilters{})

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "page2token", nextOffset)
}

func TestInventory_GetSoftwareKey_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/inventory/software/chrome"

	server.On("GET", path, common.SuccessResponse(inventory.SoftwareKeyResponse{
		Software: []inventory.SoftwareUserList{
			{
				SoftwareKey:     "chrome",
				SoftwareName:    "Google Chrome",
				SoftwareVersion: "120.0.6099.130",
				OS:              "Windows 10",
				Vendor:          "Google",
				UserID:          1001,
				DeviceID:        2001,
				Hostname:        "LAPTOP-001",
				Username:        "john.doe",
				InstallDate:     "2024-01-15",
			},
			{
				SoftwareKey:     "chrome",
				SoftwareName:    "Google Chrome",
				SoftwareVersion: "120.0.6099.129",
				OS:              "macOS 14.2",
				Vendor:          "Google",
				UserID:          1002,
				DeviceID:        2002,
				Hostname:        "MACBOOK-001",
				Username:        "jane.smith",
				InstallDate:     "2024-01-14",
			},
		},
		NextOffset: "",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, nextOffset, _, err := inventory.GetSoftwareKey(context.Background(), service, "chrome", inventory.GetSoftwareFilters{})

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "120.0.6099.130", result[0].SoftwareVersion)
	assert.Equal(t, "Windows 10", result[0].OS)
	assert.Equal(t, "LAPTOP-001", result[0].Hostname)
	assert.Empty(t, nextOffset)
}

func TestInventory_GetSoftware_Empty_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/inventory/software"

	server.On("GET", path, common.SuccessResponse(inventory.SoftwareOverviewResponse{
		Software:   []inventory.SoftwareOverview{},
		NextOffset: "",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, _, err := inventory.GetSoftware(context.Background(), service, inventory.GetSoftwareFilters{})

	require.NoError(t, err)
	assert.Len(t, result, 0)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestInventory_Structure(t *testing.T) {
	t.Parallel()

	t.Run("SoftwareOverview JSON marshaling", func(t *testing.T) {
		software := inventory.SoftwareOverview{
			SoftwareKey:         "vscode",
			SoftwareName:        "Visual Studio Code",
			Vendor:              "Microsoft",
			SoftwareGroup:       "Development Tools",
			SoftwareInstallType: "MSI",
			UserTotal:           1000,
			DeviceTotal:         950,
		}

		data, err := json.Marshal(software)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"software_key":"vscode"`)
		assert.Contains(t, string(data), `"software_name":"Visual Studio Code"`)
		assert.Contains(t, string(data), `"vendor":"Microsoft"`)
		assert.Contains(t, string(data), `"user_total":1000`)
		assert.Contains(t, string(data), `"device_total":950`)
	})

	t.Run("SoftwareOverview JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"software_key": "slack",
			"software_name": "Slack",
			"vendor": "Salesforce",
			"software_group": "Communication",
			"sw_install_type": "EXE",
			"user_total": 500,
			"device_total": 480
		}`

		var software inventory.SoftwareOverview
		err := json.Unmarshal([]byte(jsonData), &software)
		require.NoError(t, err)

		assert.Equal(t, "slack", software.SoftwareKey)
		assert.Equal(t, "Slack", software.SoftwareName)
		assert.Equal(t, "Salesforce", software.Vendor)
		assert.Equal(t, 500, software.UserTotal)
	})

	t.Run("SoftwareUserList JSON marshaling", func(t *testing.T) {
		userList := inventory.SoftwareUserList{
			SoftwareKey:     "chrome",
			SoftwareName:    "Google Chrome",
			SoftwareVersion: "120.0.0.0",
			SoftwareGroup:   "Browsers",
			OS:              "Windows 11",
			Vendor:          "Google",
			UserID:          12345,
			DeviceID:        67890,
			Hostname:        "WORKSTATION-001",
			Username:        "alice.johnson",
			InstallDate:     "2024-01-20",
		}

		data, err := json.Marshal(userList)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"software_key":"chrome"`)
		assert.Contains(t, string(data), `"software_version":"120.0.0.0"`)
		assert.Contains(t, string(data), `"os":"Windows 11"`)
		assert.Contains(t, string(data), `"hostname":"WORKSTATION-001"`)
		assert.Contains(t, string(data), `"install_date":"2024-01-20"`)
	})

	t.Run("SoftwareUserList JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"software_key": "teams",
			"software_name": "Microsoft Teams",
			"software_version": "1.6.00.1381",
			"software_group": "Communication",
			"os": "macOS 14.2",
			"vendor": "Microsoft",
			"user_id": 11111,
			"device_id": 22222,
			"hostname": "MACBOOK-PRO",
			"username": "bob.williams",
			"install_date": "2024-01-18"
		}`

		var userList inventory.SoftwareUserList
		err := json.Unmarshal([]byte(jsonData), &userList)
		require.NoError(t, err)

		assert.Equal(t, "teams", userList.SoftwareKey)
		assert.Equal(t, "Microsoft Teams", userList.SoftwareName)
		assert.Equal(t, "1.6.00.1381", userList.SoftwareVersion)
		assert.Equal(t, "MACBOOK-PRO", userList.Hostname)
		assert.Equal(t, 11111, userList.UserID)
	})

	t.Run("SoftwareOverviewResponse JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"software": [
				{"software_key": "app1", "software_name": "Application 1", "user_total": 100},
				{"software_key": "app2", "software_name": "Application 2", "user_total": 200}
			],
			"next_offset": "nextpage"
		}`

		var response inventory.SoftwareOverviewResponse
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Len(t, response.Software, 2)
		assert.Equal(t, "nextpage", response.NextOffset)
		assert.Equal(t, "app1", response.Software[0].SoftwareKey)
	})

	t.Run("SoftwareKeyResponse JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"software": [
				{"software_key": "zoom", "hostname": "PC-001", "username": "user1"},
				{"software_key": "zoom", "hostname": "PC-002", "username": "user2"}
			],
			"next_offset": ""
		}`

		var response inventory.SoftwareKeyResponse
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Len(t, response.Software, 2)
		assert.Empty(t, response.NextOffset)
		assert.Equal(t, "PC-001", response.Software[0].Hostname)
	})
}

func TestInventory_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse software overview list", func(t *testing.T) {
		jsonResponse := `[
			{"software_key": "chrome", "software_name": "Google Chrome", "vendor": "Google", "user_total": 500},
			{"software_key": "firefox", "software_name": "Mozilla Firefox", "vendor": "Mozilla", "user_total": 200},
			{"software_key": "edge", "software_name": "Microsoft Edge", "vendor": "Microsoft", "user_total": 300}
		]`

		var software []inventory.SoftwareOverview
		err := json.Unmarshal([]byte(jsonResponse), &software)
		require.NoError(t, err)

		assert.Len(t, software, 3)
		assert.Equal(t, "Google Chrome", software[0].SoftwareName)
		assert.Equal(t, 500, software[0].UserTotal)
	})

	t.Run("Parse empty software list", func(t *testing.T) {
		jsonResponse := `[]`

		var software []inventory.SoftwareOverview
		err := json.Unmarshal([]byte(jsonResponse), &software)
		require.NoError(t, err)

		assert.Empty(t, software)
	})
}
