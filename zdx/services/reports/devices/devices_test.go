package devices

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services/common"
)

func TestGetAllDevices(t *testing.T) {
	client, err := tests.NewZdxClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	// Define a time filter for the last 2 hours
	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()
	filters := GetDevicesFilters{
		GetFromToFilters: common.GetFromToFilters{
			From: int(from),
			To:   int(to),
		},
	}

	// Call GetAllDevices
	devices, resp, err := GetAllDevices(service, filters)
	if err != nil {
		t.Fatalf("Error getting all devices: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(devices) == 0 {
		t.Log("No devices found.")
	} else {
		t.Logf("Retrieved %d devices", len(devices))
		for _, device := range devices {
			if device.Software != nil {
				t.Logf("Device ID: %d, Name: %s, User ID: %s", device.ID, device.Name, device.Software.User)
			} else {
				t.Logf("Device ID: %d, Name: %s, User ID: <nil>", device.ID, device.Name)
			}
		}
	}
}

func TestGetDevice(t *testing.T) {
	client, err := tests.NewZdxClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	// Define a time filter for the last 2 hours
	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()
	filters := GetDevicesFilters{
		GetFromToFilters: common.GetFromToFilters{
			From: int(from),
			To:   int(to),
		},
	}

	// Invoke GetAllDevices to retrieve the ID of the first device
	devices, _, err := GetAllDevices(service, filters)
	if err != nil {
		t.Fatalf("Error getting all devices: %v", err)
	}

	if len(devices) == 0 {
		t.Log("No devices found, skipping GetDevice test.")
		return
	}

	firstDeviceID := devices[0].ID

	// Call GetDevice with the first device's ID
	device, resp, err := GetDevice(service, strconv.Itoa(firstDeviceID))
	if err != nil {
		t.Fatalf("Error getting device: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if device.ID == 0 {
		t.Log("No device found with the specified ID.")
	} else {
		if device.Software != nil {
			t.Logf("Retrieved device ID: %d, Name: %s, User ID: %s", device.ID, device.Name, device.Software.User)
		} else {
			t.Logf("Retrieved device ID: %d, Name: %s, User ID: <nil>", device.ID, device.Name)
		}
	}
}
