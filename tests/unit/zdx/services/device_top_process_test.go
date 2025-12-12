// Package services provides unit tests for ZDX device top process service
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	zdxcommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/devices"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestDeviceTopProcess_GetDeviceTopProcesses_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/devices/12345/deeptraces/trace-001/top-processes"

	server.On("GET", path, common.SuccessResponse([]devices.DeviceTopProcesses{
		{
			TimeStamp: 1699900000,
			TopProcesses: []devices.TopProcesses{
				{Category: "CPU", Processes: []devices.Processes{{ID: 1, Name: "chrome.exe"}}},
				{Category: "Memory", Processes: []devices.Processes{{ID: 2, Name: "teams.exe"}}},
			},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := devices.GetDeviceTopProcesses(context.Background(), service, 12345, "trace-001", zdxcommon.GetFromToFilters{})

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1699900000, result[0].TimeStamp)
	assert.Len(t, result[0].TopProcesses, 2)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestDeviceTopProcess_Structure(t *testing.T) {
	t.Parallel()

	t.Run("DeviceTopProcesses JSON marshaling", func(t *testing.T) {
		topProcs := devices.DeviceTopProcesses{
			TimeStamp: 1699900000,
			TopProcesses: []devices.TopProcesses{
				{
					Category: "CPU",
					Processes: []devices.Processes{
						{ID: 1, Name: "chrome.exe"},
						{ID: 2, Name: "teams.exe"},
						{ID: 3, Name: "outlook.exe"},
					},
				},
				{
					Category: "Memory",
					Processes: []devices.Processes{
						{ID: 4, Name: "chrome.exe"},
						{ID: 5, Name: "vscode.exe"},
					},
				},
			},
		}

		data, err := json.Marshal(topProcs)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"timestamp":1699900000`)
		assert.Contains(t, string(data), `"top_processes"`)
		assert.Contains(t, string(data), `"category":"CPU"`)
		assert.Contains(t, string(data), `"category":"Memory"`)
	})

	t.Run("DeviceTopProcesses JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"timestamp": 1699950000,
			"top_processes": [
				{
					"category": "Disk I/O",
					"processes": [
						{"id": 100, "name": "sqlserver.exe"},
						{"id": 101, "name": "mysqld.exe"}
					]
				}
			]
		}`

		var topProcs devices.DeviceTopProcesses
		err := json.Unmarshal([]byte(jsonData), &topProcs)
		require.NoError(t, err)

		assert.Equal(t, 1699950000, topProcs.TimeStamp)
		assert.Len(t, topProcs.TopProcesses, 1)
		assert.Equal(t, "Disk I/O", topProcs.TopProcesses[0].Category)
		assert.Len(t, topProcs.TopProcesses[0].Processes, 2)
	})

	t.Run("TopProcesses JSON marshaling", func(t *testing.T) {
		topProc := devices.TopProcesses{
			Category: "Network",
			Processes: []devices.Processes{
				{ID: 1, Name: "chrome.exe"},
				{ID: 2, Name: "slack.exe"},
				{ID: 3, Name: "zoom.exe"},
			},
		}

		data, err := json.Marshal(topProc)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"category":"Network"`)
		assert.Contains(t, string(data), `"processes"`)
	})

	t.Run("Processes JSON marshaling", func(t *testing.T) {
		process := devices.Processes{
			ID:   12345,
			Name: "java.exe",
		}

		data, err := json.Marshal(process)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"java.exe"`)
	})
}

func TestDeviceTopProcess_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse device top processes list response", func(t *testing.T) {
		jsonResponse := `[
			{
				"timestamp": 1699900000,
				"top_processes": [
					{
						"category": "CPU",
						"processes": [
							{"id": 1, "name": "chrome.exe"},
							{"id": 2, "name": "teams.exe"}
						]
					},
					{
						"category": "Memory",
						"processes": [
							{"id": 3, "name": "chrome.exe"},
							{"id": 4, "name": "vscode.exe"}
						]
					}
				]
			},
			{
				"timestamp": 1699903600,
				"top_processes": [
					{
						"category": "CPU",
						"processes": [
							{"id": 5, "name": "zoom.exe"},
							{"id": 6, "name": "chrome.exe"}
						]
					}
				]
			}
		]`

		var topProcs []devices.DeviceTopProcesses
		err := json.Unmarshal([]byte(jsonResponse), &topProcs)
		require.NoError(t, err)

		assert.Len(t, topProcs, 2)
		
		// First timestamp
		assert.Equal(t, 1699900000, topProcs[0].TimeStamp)
		assert.Len(t, topProcs[0].TopProcesses, 2)
		assert.Equal(t, "CPU", topProcs[0].TopProcesses[0].Category)
		assert.Len(t, topProcs[0].TopProcesses[0].Processes, 2)
		
		// Second timestamp
		assert.Equal(t, 1699903600, topProcs[1].TimeStamp)
		assert.Len(t, topProcs[1].TopProcesses, 1)
	})

	t.Run("Parse deep trace top processes", func(t *testing.T) {
		jsonResponse := `[
			{
				"timestamp": 1699900000,
				"top_processes": [
					{
						"category": "CPU",
						"processes": [
							{"id": 1001, "name": "System Idle Process"},
							{"id": 1002, "name": "Chrome"},
							{"id": 1003, "name": "Microsoft Teams"},
							{"id": 1004, "name": "Visual Studio Code"},
							{"id": 1005, "name": "Slack"}
						]
					},
					{
						"category": "Memory",
						"processes": [
							{"id": 2001, "name": "Chrome"},
							{"id": 2002, "name": "Visual Studio Code"},
							{"id": 2003, "name": "Microsoft Teams"}
						]
					},
					{
						"category": "Disk Read",
						"processes": [
							{"id": 3001, "name": "Windows Defender"},
							{"id": 3002, "name": "OneDrive"}
						]
					},
					{
						"category": "Disk Write",
						"processes": [
							{"id": 4001, "name": "Chrome"},
							{"id": 4002, "name": "Windows Update"}
						]
					}
				]
			}
		]`

		var topProcs []devices.DeviceTopProcesses
		err := json.Unmarshal([]byte(jsonResponse), &topProcs)
		require.NoError(t, err)

		assert.Len(t, topProcs, 1)
		assert.Len(t, topProcs[0].TopProcesses, 4)
		
		// Verify categories
		categories := make([]string, len(topProcs[0].TopProcesses))
		for i, tp := range topProcs[0].TopProcesses {
			categories[i] = tp.Category
		}
		assert.Contains(t, categories, "CPU")
		assert.Contains(t, categories, "Memory")
		assert.Contains(t, categories, "Disk Read")
		assert.Contains(t, categories, "Disk Write")
	})
}

