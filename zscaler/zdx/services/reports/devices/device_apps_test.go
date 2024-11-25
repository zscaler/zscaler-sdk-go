package devices

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zdx/services"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zdx/services/common"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zdx/services/reports/applications"
)

func TestGetDeviceApp(t *testing.T) {
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
		t.Log("No devices found, skipping GetDeviceApp test.")
		return
	}

	firstDeviceID := devices[0].ID

	// Invoke GetAllApps to retrieve the ID of the first app
	apps, _, err := applications.GetAllApps(service, common.GetFromToFilters{From: int(from), To: int(to)})
	if err != nil {
		t.Fatalf("Error getting all apps: %v", err)
	}

	if len(apps) == 0 {
		t.Log("No apps found, skipping GetDeviceApp test.")
		return
	}

	firstAppID := apps[0].ID

	// Call GetDeviceApp with the first device's ID and first app's ID
	deviceApp, resp, err := GetDeviceApp(service, strconv.Itoa(firstDeviceID), strconv.Itoa(firstAppID), common.GetFromToFilters{From: int(from), To: int(to)})
	if err != nil {
		t.Fatalf("Error getting device app: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if deviceApp.ID == 0 {
		t.Log("No app found for the specified device and app ID.")
	} else {
		t.Logf("Retrieved app for device ID: %d, App ID: %d, App Name: %s, Score: %f", firstDeviceID, deviceApp.ID, deviceApp.Name, deviceApp.Score)
	}
}

func TestGetDeviceAllApps(t *testing.T) {
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
		t.Log("No devices found, skipping GetDeviceAllApps test.")
		return
	}

	firstDeviceID := devices[0].ID

	// Call GetDeviceAllApps with the first device's ID
	deviceApps, resp, err := GetDeviceAllApps(service, strconv.Itoa(firstDeviceID), common.GetFromToFilters{From: int(from), To: int(to)})
	if err != nil {
		t.Fatalf("Error getting all apps for device: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(deviceApps) == 0 {
		t.Log("No apps found for the specified device.")
	} else {
		t.Logf("Retrieved %d apps for device ID: %d", len(deviceApps), firstDeviceID)
		for _, app := range deviceApps {
			t.Logf("App ID: %d, App Name: %s, Score: %f", app.ID, app.Name, app.Score)
		}
	}
}
