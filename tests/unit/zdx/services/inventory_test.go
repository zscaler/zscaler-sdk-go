// Package services provides unit tests for ZDX services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/inventory"
)

func TestInventory_Structure(t *testing.T) {
	t.Parallel()

	t.Run("SoftwareOverview JSON marshaling", func(t *testing.T) {
		software := inventory.SoftwareOverview{
			SoftwareKey:         "chrome-browser",
			SoftwareName:        "Google Chrome",
			Vendor:              "Google LLC",
			SoftwareGroup:       "Browsers",
			SoftwareInstallType: "MSI",
			UserTotal:           1500,
			DeviceTotal:         2000,
		}

		data, err := json.Marshal(software)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"software_key":"chrome-browser"`)
		assert.Contains(t, string(data), `"software_name":"Google Chrome"`)
		assert.Contains(t, string(data), `"vendor":"Google LLC"`)
		assert.Contains(t, string(data), `"user_total":1500`)
	})

	t.Run("SoftwareOverview JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"software_key": "vscode",
			"software_name": "Visual Studio Code",
			"vendor": "Microsoft Corporation",
			"software_group": "Development Tools",
			"sw_install_type": "EXE",
			"user_total": 500,
			"device_total": 600
		}`

		var software inventory.SoftwareOverview
		err := json.Unmarshal([]byte(jsonData), &software)
		require.NoError(t, err)

		assert.Equal(t, "vscode", software.SoftwareKey)
		assert.Equal(t, "Visual Studio Code", software.SoftwareName)
		assert.Equal(t, "Microsoft Corporation", software.Vendor)
		assert.Equal(t, "Development Tools", software.SoftwareGroup)
		assert.Equal(t, 500, software.UserTotal)
		assert.Equal(t, 600, software.DeviceTotal)
	})

	t.Run("SoftwareUserList JSON marshaling", func(t *testing.T) {
		software := inventory.SoftwareUserList{
			SoftwareKey:     "slack",
			SoftwareName:    "Slack",
			SoftwareVersion: "4.35.126",
			SoftwareGroup:   "Communication",
			OS:              "Windows 11",
			Vendor:          "Slack Technologies",
			UserID:          12345,
			DeviceID:        67890,
			Hostname:        "DESKTOP-001",
			Username:        "john.doe",
			InstallDate:     "2024-01-15",
		}

		data, err := json.Marshal(software)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"software_key":"slack"`)
		assert.Contains(t, string(data), `"software_version":"4.35.126"`)
		assert.Contains(t, string(data), `"hostname":"DESKTOP-001"`)
		assert.Contains(t, string(data), `"install_date":"2024-01-15"`)
	})

	t.Run("SoftwareUserList JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"software_key": "zoom",
			"software_name": "Zoom",
			"software_version": "5.16.0",
			"software_group": "Communication",
			"os": "macOS 14.0",
			"vendor": "Zoom Video Communications",
			"user_id": 1001,
			"device_id": 2001,
			"hostname": "MACBOOK-001",
			"username": "jane.smith",
			"install_date": "2024-02-01"
		}`

		var software inventory.SoftwareUserList
		err := json.Unmarshal([]byte(jsonData), &software)
		require.NoError(t, err)

		assert.Equal(t, "zoom", software.SoftwareKey)
		assert.Equal(t, "5.16.0", software.SoftwareVersion)
		assert.Equal(t, "macOS 14.0", software.OS)
		assert.Equal(t, 1001, software.UserID)
		assert.Equal(t, "MACBOOK-001", software.Hostname)
	})

	t.Run("SoftwareOverviewResponse JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"software": [
				{"software_key": "app1", "software_name": "App 1", "user_total": 100},
				{"software_key": "app2", "software_name": "App 2", "user_total": 200}
			],
			"next_offset": "offset123"
		}`

		var response inventory.SoftwareOverviewResponse
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Len(t, response.Software, 2)
		assert.Equal(t, "offset123", response.NextOffset)
		assert.Equal(t, "App 1", response.Software[0].SoftwareName)
	})

	t.Run("SoftwareKeyResponse JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"software": [
				{"software_key": "chrome", "software_version": "120.0.1", "hostname": "PC-001"},
				{"software_key": "chrome", "software_version": "119.0.5", "hostname": "PC-002"}
			],
			"next_offset": ""
		}`

		var response inventory.SoftwareKeyResponse
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Len(t, response.Software, 2)
		assert.Empty(t, response.NextOffset)
		assert.Equal(t, "120.0.1", response.Software[0].SoftwareVersion)
	})
}

func TestInventory_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse software overview list", func(t *testing.T) {
		jsonResponse := `{
			"software": [
				{
					"software_key": "microsoft-office",
					"software_name": "Microsoft Office 365",
					"vendor": "Microsoft Corporation",
					"software_group": "Productivity",
					"sw_install_type": "MSI",
					"user_total": 5000,
					"device_total": 6000
				},
				{
					"software_key": "adobe-reader",
					"software_name": "Adobe Acrobat Reader",
					"vendor": "Adobe Inc.",
					"software_group": "PDF Tools",
					"user_total": 4500,
					"device_total": 5500
				}
			],
			"next_offset": "page2"
		}`

		var response inventory.SoftwareOverviewResponse
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Len(t, response.Software, 2)
		assert.Equal(t, "page2", response.NextOffset)
		
		// Check first software
		assert.Equal(t, "microsoft-office", response.Software[0].SoftwareKey)
		assert.Equal(t, "Microsoft Office 365", response.Software[0].SoftwareName)
		assert.Equal(t, 5000, response.Software[0].UserTotal)
		
		// Check second software
		assert.Equal(t, "Adobe Acrobat Reader", response.Software[1].SoftwareName)
		assert.Equal(t, "Adobe Inc.", response.Software[1].Vendor)
	})

	t.Run("Parse software installations by key", func(t *testing.T) {
		jsonResponse := `{
			"software": [
				{
					"software_key": "chrome",
					"software_name": "Google Chrome",
					"software_version": "120.0.6099.130",
					"os": "Windows 11",
					"hostname": "LAPTOP-ENG-001",
					"username": "engineer1",
					"install_date": "2024-01-10"
				},
				{
					"software_key": "chrome",
					"software_name": "Google Chrome",
					"software_version": "119.0.6045.200",
					"os": "Windows 10",
					"hostname": "DESKTOP-SALES-001",
					"username": "sales1",
					"install_date": "2023-12-15"
				}
			],
			"next_offset": ""
		}`

		var response inventory.SoftwareKeyResponse
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Len(t, response.Software, 2)
		assert.Empty(t, response.NextOffset)
		
		// Verify different versions
		assert.Equal(t, "120.0.6099.130", response.Software[0].SoftwareVersion)
		assert.Equal(t, "Windows 11", response.Software[0].OS)
		assert.Equal(t, "119.0.6045.200", response.Software[1].SoftwareVersion)
		assert.Equal(t, "Windows 10", response.Software[1].OS)
	})
}

