// Package services provides unit tests for ZDX services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/devices"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestDevices_GetDevice_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	deviceID := "12345"
	path := "/zdx/v1/devices/" + deviceID

	server.On("GET", path, common.SuccessResponse(devices.DeviceDetail{
		ID:   12345,
		Name: "LAPTOP-ENG-001",
		Hardware: &devices.Hardware{
			HWModel: "MacBook Pro",
			HWMFG:   "Apple",
		},
		Software: &devices.Software{
			OSName: "macOS",
			OSVer:  "14.0",
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := devices.GetDevice(context.Background(), service, deviceID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 12345, result.ID)
	assert.Equal(t, "LAPTOP-ENG-001", result.Name)
	assert.Equal(t, "MacBook Pro", result.Hardware.HWModel)
}

func TestDevices_GetAllDevices_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/devices"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"devices": []devices.DeviceDetail{
			{ID: 1, Name: "Device 1"},
			{ID: 2, Name: "Device 2"},
		},
		"next_offset": nil,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	filters := devices.GetDevicesFilters{}
	result, _, err := devices.GetAllDevices(context.Background(), service, filters)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestDevices_Structure(t *testing.T) {
	t.Parallel()

	t.Run("DeviceDetail JSON marshaling", func(t *testing.T) {
		device := devices.DeviceDetail{
			ID:   12345,
			Name: "LAPTOP-ENG-001",
			Hardware: &devices.Hardware{
				HWModel:     "MacBook Pro",
				HWMFG:       "Apple",
				HWType:      "Laptop",
				TotMem:      16384,
				DiskSize:    512,
				DiskType:    "SSD",
				CPUMFG:      "Apple",
				CPUModel:    "M1 Pro",
				SpeedGHZ:    3.2,
				LogicalProc: 10,
				NumCores:    10,
			},
			Software: &devices.Software{
				OSName:        "macOS",
				OSVer:         "14.0",
				Hostname:      "laptop-eng-001",
				User:          "john.doe",
				ClientConnVer: "4.2.0",
				ZDXVer:        "3.5.0",
			},
		}

		data, err := json.Marshal(device)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"LAPTOP-ENG-001"`)
		assert.Contains(t, string(data), `"hw_model":"MacBook Pro"`)
		assert.Contains(t, string(data), `"os_name":"macOS"`)
	})

	t.Run("DeviceDetail JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 67890,
			"name": "DESKTOP-SALES-001",
			"hardware": {
				"hw_model": "ThinkPad X1 Carbon",
				"hw_mfg": "Lenovo",
				"hw_type": "Laptop",
				"tot_mem": 32768,
				"disk_size": 1024,
				"disk_type": "NVMe SSD",
				"cpu_mfg": "Intel",
				"cpu_model": "Core i7-1365U",
				"speed_ghz": 3.9,
				"logical_proc": 12,
				"num_cores": 6
			},
			"software": {
				"os_name": "Windows",
				"os_ver": "11 Pro",
				"hostname": "desktop-sales-001",
				"user": "jane.smith",
				"client_conn_ver": "4.2.1",
				"zdx_ver": "3.5.1"
			},
			"network": [
				{
					"net_type": "WiFi",
					"status": "Connected",
					"ipv4": "192.168.1.100",
					"gateway": "192.168.1.1",
					"ssid": "CorpWiFi"
				}
			]
		}`

		var device devices.DeviceDetail
		err := json.Unmarshal([]byte(jsonData), &device)
		require.NoError(t, err)

		assert.Equal(t, 67890, device.ID)
		assert.Equal(t, "DESKTOP-SALES-001", device.Name)
		assert.NotNil(t, device.Hardware)
		assert.Equal(t, "ThinkPad X1 Carbon", device.Hardware.HWModel)
		assert.Equal(t, 32768, device.Hardware.TotMem)
		assert.NotNil(t, device.Software)
		assert.Equal(t, "Windows", device.Software.OSName)
		assert.Len(t, device.Network, 1)
		assert.Equal(t, "CorpWiFi", device.Network[0].SSID)
	})

	t.Run("Hardware JSON marshaling", func(t *testing.T) {
		hardware := devices.Hardware{
			HWModel:     "Dell XPS 15",
			HWMFG:       "Dell",
			HWType:      "Laptop",
			HWSerial:    "ABC123XYZ",
			TotMem:      64000,
			GPU:         "NVIDIA RTX 4070",
			DiskSize:    2048,
			DiskModel:   "Samsung 990 Pro",
			DiskType:    "NVMe",
			CPUMFG:      "Intel",
			CPUModel:    "Core i9-13900H",
			SpeedGHZ:    5.4,
			LogicalProc: 24,
			NumCores:    14,
		}

		data, err := json.Marshal(hardware)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"hw_model":"Dell XPS 15"`)
		assert.Contains(t, string(data), `"gpu":"NVIDIA RTX 4070"`)
		assert.Contains(t, string(data), `"speed_ghz":5.4`)
	})

	t.Run("Network JSON marshaling", func(t *testing.T) {
		network := devices.Network{
			NetType:     "WiFi",
			Status:      "Connected",
			IPv4:        "10.0.0.100",
			IPv6:        "fe80::1",
			DNSSRVS:     "8.8.8.8,8.8.4.4",
			DNSSuffix:   "corp.example.com",
			Gateway:     "10.0.0.1",
			MAC:         "00:11:22:33:44:55",
			GUID:        "net-guid-123",
			WiFiAdapter: "Intel Wi-Fi 6E AX211",
			WiFiType:    "802.11ax",
			SSID:        "CorpWiFi-5G",
			Channel:     "36",
			BSSID:       "AA:BB:CC:DD:EE:FF",
		}

		data, err := json.Marshal(network)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"net_type":"WiFi"`)
		assert.Contains(t, string(data), `"ipv4":"10.0.0.100"`)
		assert.Contains(t, string(data), `"ssid":"CorpWiFi-5G"`)
		assert.Contains(t, string(data), `"wifi_type":"802.11ax"`)
	})

	t.Run("Software JSON marshaling", func(t *testing.T) {
		software := devices.Software{
			OSName:        "Windows",
			OSVer:         "11 Enterprise",
			Hostname:      "workstation-001",
			NetBios:       "WORKSTATION01",
			User:          "admin.user",
			ClientConnVer: "4.3.0.100",
			ZDXVer:        "3.6.0",
		}

		data, err := json.Marshal(software)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"os_name":"Windows"`)
		assert.Contains(t, string(data), `"os_ver":"11 Enterprise"`)
		assert.Contains(t, string(data), `"client_conn_ver":"4.3.0.100"`)
	})
}

func TestDevices_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse devices list response", func(t *testing.T) {
		jsonResponse := `{
			"devices": [
				{
					"id": 1,
					"name": "Device 1",
					"hardware": {"hw_model": "Model A"},
					"software": {"os_name": "Windows"}
				},
				{
					"id": 2,
					"name": "Device 2",
					"hardware": {"hw_model": "Model B"},
					"software": {"os_name": "macOS"}
				}
			],
			"next_offset": "page2"
		}`

		var response struct {
			Devices    []devices.DeviceDetail `json:"devices"`
			NextOffset string                 `json:"next_offset"`
		}
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Len(t, response.Devices, 2)
		assert.Equal(t, "page2", response.NextOffset)
		assert.Equal(t, "Device 1", response.Devices[0].Name)
		assert.Equal(t, "Windows", response.Devices[0].Software.OSName)
		assert.Equal(t, "macOS", response.Devices[1].Software.OSName)
	})

	t.Run("Parse single device response", func(t *testing.T) {
		jsonResponse := `{
			"id": 12345,
			"name": "Production Server",
			"hardware": {
				"hw_model": "PowerEdge R750",
				"hw_mfg": "Dell",
				"hw_type": "Server",
				"tot_mem": 131072,
				"disk_size": 4096,
				"num_cores": 32
			},
			"software": {
				"os_name": "Linux",
				"os_ver": "Ubuntu 22.04 LTS",
				"hostname": "prod-server-01"
			},
			"network": [
				{"net_type": "Ethernet", "status": "Connected", "ipv4": "10.10.10.10"},
				{"net_type": "Ethernet", "status": "Connected", "ipv4": "10.10.10.11"}
			]
		}`

		var device devices.DeviceDetail
		err := json.Unmarshal([]byte(jsonResponse), &device)
		require.NoError(t, err)

		assert.Equal(t, 12345, device.ID)
		assert.Equal(t, "Production Server", device.Name)
		assert.Equal(t, "PowerEdge R750", device.Hardware.HWModel)
		assert.Equal(t, 131072, device.Hardware.TotMem)
		assert.Equal(t, "Linux", device.Software.OSName)
		assert.Len(t, device.Network, 2)
	})
}

