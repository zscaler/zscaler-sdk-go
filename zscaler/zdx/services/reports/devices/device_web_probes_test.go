package devices

import (
	"net/http"
	"testing"
	"time"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zdx/services"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zdx/services/common"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zdx/services/reports/applications"
)

func TestGetAllWebProbes(t *testing.T) {
	client, err := tests.NewZdxClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	// Define a time filter for the last 2 hours
	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()
	filters := common.GetFromToFilters{
		From: int(from),
		To:   int(to),
	}

	// Invoke GetAllDevices to retrieve the ID of the first device
	devices, _, err := GetAllDevices(service, GetDevicesFilters{
		GetFromToFilters: filters,
	})
	if err != nil {
		t.Fatalf("Error getting all devices: %v", err)
	}

	if len(devices) == 0 {
		t.Log("No devices found, skipping GetAllWebProbes test.")
		return
	}

	firstDeviceID := devices[0].ID

	// Invoke GetAllApps to retrieve the ID of the first app
	apps, _, err := applications.GetAllApps(service, filters)
	if err != nil {
		t.Fatalf("Error getting all apps: %v", err)
	}

	if len(apps) == 0 {
		t.Log("No apps found, skipping GetAllWebProbes test.")
		return
	}

	firstAppID := apps[0].ID

	// Call GetAllWebProbes with the first device's ID and first app's ID
	webProbes, resp, err := GetAllWebProbes(service, firstDeviceID, firstAppID, filters)
	if err != nil {
		t.Fatalf("Error getting web probes: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(webProbes) == 0 {
		t.Log("No web probes found.")
	} else {
		t.Logf("Retrieved %d web probes", len(webProbes))
		for _, probe := range webProbes {
			t.Logf("Probe ID: %d, Name: %s, NumProbes: %d, AvgScore: %f, AvgPFT: %f", probe.ID, probe.Name, probe.NumProbes, probe.AvgScore, probe.AvgPFT)
		}
	}
}

func TestGetWebProbes(t *testing.T) {
	client, err := tests.NewZdxClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	// Define a time filter for the last 2 hours
	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()
	filters := common.GetFromToFilters{
		From: int(from),
		To:   int(to),
	}

	// Invoke GetAllDevices to retrieve the ID of the first device
	devices, _, err := GetAllDevices(service, GetDevicesFilters{
		GetFromToFilters: filters,
	})
	if err != nil {
		t.Fatalf("Error getting all devices: %v", err)
	}

	if len(devices) == 0 {
		t.Log("No devices found, skipping GetWebProbes test.")
		return
	}

	firstDeviceID := devices[0].ID

	// Invoke GetAllApps to retrieve the ID of the first app
	apps, _, err := applications.GetAllApps(service, filters)
	if err != nil {
		t.Fatalf("Error getting all apps: %v", err)
	}

	if len(apps) == 0 {
		t.Log("No apps found, skipping GetWebProbes test.")
		return
	}

	firstAppID := apps[0].ID

	// Invoke GetAllWebProbes to retrieve the ID of the first web probe
	webProbes, _, err := GetAllWebProbes(service, firstDeviceID, firstAppID, filters)
	if err != nil {
		t.Fatalf("Error getting web probes: %v", err)
	}

	if len(webProbes) == 0 {
		t.Log("No web probes found, skipping GetWebProbes test.")
		return
	}

	firstProbeID := webProbes[0].ID

	// Call GetWebProbes with the first device's ID, first app's ID, and first web probe's ID
	webProbeMetrics, resp, err := GetWebProbes(service, firstDeviceID, firstAppID, firstProbeID, filters)
	if err != nil {
		t.Fatalf("Error getting web probe metrics: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(webProbeMetrics) == 0 {
		t.Log("No web probe metrics found.")
	} else {
		t.Logf("Retrieved %d data points for web probe", len(webProbeMetrics))
		for _, metric := range webProbeMetrics {
			t.Logf("Metric: %s, Unit: %s", metric.Metric, metric.Unit)
			for _, dp := range metric.DataPoints {
				t.Logf("Timestamp: %d, Value: %f", dp.TimeStamp, dp.Value)
			}
		}
	}
}
