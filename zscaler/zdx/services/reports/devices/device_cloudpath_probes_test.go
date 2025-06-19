package devices

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/applications"
)

func TestGetAllCloudPathProbes(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Define a time filter for the last 2 hours
	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()
	filters := common.GetFromToFilters{
		From: int(from),
		To:   int(to),
	}

	// Invoke GetAllDevices to retrieve the ID of the first device
	devices, _, err := GetAllDevices(context.Background(), service, GetDevicesFilters{
		GetFromToFilters: filters,
	})
	if err != nil {
		t.Fatalf("Error getting all devices: %v", err)
	}

	if len(devices) == 0 {
		t.Log("No devices found, skipping GetAllCloudPathProbes test.")
		return
	}

	firstDeviceID := devices[0].ID

	// Invoke GetAllApps to retrieve the ID of the first app
	apps, _, err := applications.GetAllApps(context.Background(), service, filters)
	if err != nil {
		t.Fatalf("Error getting all apps: %v", err)
	}

	if len(apps) == 0 {
		t.Log("No apps found, skipping GetAllCloudPathProbes test.")
		return
	}

	firstAppID := apps[0].ID

	// Call GetAllCloudPathProbes with the first device's ID and first app's ID
	cloudPathProbes, resp, err := GetAllCloudPathProbes(context.Background(), service, firstDeviceID, firstAppID, filters)
	if err != nil {
		t.Fatalf("Error getting cloud path probes: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(cloudPathProbes) == 0 {
		t.Log("No cloud path probes found.")
	} else {
		t.Logf("Retrieved %d cloud path probes", len(cloudPathProbes))
		for _, probe := range cloudPathProbes {
			t.Logf("Probe ID: %d, Name: %s, NumProbes: %d", probe.ID, probe.Name, probe.NumProbes)
			for _, latency := range probe.AverageLatency {
				t.Logf("LegSRC: %s, LegDst: %s, Latency: %f", latency.LegSRC, latency.LegDst, latency.Latency)
			}
		}
	}
}

func TestGetDeviceAppCloudPathProbe(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

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
	devices, _, err := GetAllDevices(context.Background(), service, filters)
	if err != nil {
		t.Fatalf("Error getting all devices: %v", err)
	}

	if len(devices) == 0 {
		t.Log("No devices found, skipping GetDeviceAppCloudPathProbe test.")
		return
	}

	firstDeviceID := devices[0].ID

	// Invoke GetAllApps to retrieve the ID of the first app
	apps, _, err := applications.GetAllApps(context.Background(), service, common.GetFromToFilters{From: int(from), To: int(to)})
	if err != nil {
		t.Fatalf("Error getting all apps: %v", err)
	}

	if len(apps) == 0 {
		t.Log("No apps found, skipping GetDeviceAppCloudPathProbe test.")
		return
	}

	firstAppID := apps[0].ID

	// Invoke GetAllCloudPathProbes to retrieve the ID of the first probe
	probes, _, err := GetAllCloudPathProbes(context.Background(), service, firstDeviceID, firstAppID, common.GetFromToFilters{From: int(from), To: int(to)})
	if err != nil {
		t.Fatalf("Error getting all cloud path probes: %v", err)
	}

	if len(probes) == 0 {
		t.Log("No cloud path probes found, skipping GetDeviceAppCloudPathProbe test.")
		return
	}

	firstProbeID := probes[0].ID

	// Call GetDeviceAppCloudPathProbe with the first device's ID, first app's ID, and first probe's ID
	networkStats, resp, err := GetDeviceAppCloudPathProbe(context.Background(), service, firstDeviceID, firstAppID, firstProbeID, common.GetFromToFilters{From: int(from), To: int(to)})
	if err != nil {
		t.Fatalf("Error getting device app cloud path probe: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(networkStats) == 0 {
		t.Log("No network stats found.")
	} else {
		for _, stat := range networkStats {
			t.Logf("Retrieved network stats: LegSRC: %s, LegDst: %s, Stats: %+v", stat.LegSRC, stat.LegDst, stat.Stats)
		}
	}
}

func TestGetCloudPathAppDevice(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

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
	devices, _, err := GetAllDevices(context.Background(), service, filters)
	if err != nil {
		t.Fatalf("Error getting all devices: %v", err)
	}

	if len(devices) == 0 {
		t.Log("No devices found, skipping GetCloudPathAppDevice test.")
		return
	}

	firstDeviceID := devices[0].ID

	// Invoke GetAllApps to retrieve the ID of the first app
	apps, _, err := applications.GetAllApps(context.Background(), service, common.GetFromToFilters{From: int(from), To: int(to)})
	if err != nil {
		t.Fatalf("Error getting all apps: %v", err)
	}

	if len(apps) == 0 {
		t.Log("No apps found, skipping GetCloudPathAppDevice test.")
		return
	}

	firstAppID := apps[0].ID

	// Invoke GetAllCloudPathProbes to retrieve the ID of the first probe
	probes, _, err := GetAllCloudPathProbes(context.Background(), service, firstDeviceID, firstAppID, common.GetFromToFilters{From: int(from), To: int(to)})
	if err != nil {
		t.Fatalf("Error getting all cloud path probes: %v", err)
	}

	if len(probes) == 0 {
		t.Log("No cloud path probes found, skipping GetCloudPathAppDevice test.")
		return
	}

	firstProbeID := probes[0].ID

	// Call GetCloudPathAppDevice with the first device's ID, first app's ID, and first probe's ID
	cloudPathProbes, resp, err := GetCloudPathAppDevice(context.Background(), service, firstDeviceID, firstAppID, firstProbeID, common.GetFromToFilters{From: int(from), To: int(to)})
	if err != nil {
		t.Fatalf("Error getting cloud path app device: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(cloudPathProbes) == 0 {
		t.Log("No cloud path probes found.")
	} else {
		for _, probe := range cloudPathProbes {
			t.Logf("Retrieved cloud path probe: TimeStamp: %d", probe.TimeStamp)
			for _, cloudPath := range probe.CloudPath {
				t.Logf("CloudPath: SRC: %s, DST: %s, NumHops: %d, Latency: %v, Loss: %v, NumUnrespHops: %d, TunnelType: %d",
					cloudPath.SRC, cloudPath.DST, cloudPath.NumHops, cloudPath.Latency, cloudPath.Loss, cloudPath.NumUnrespHops, cloudPath.TunnelType)
				for _, hop := range cloudPath.Hops {
					t.Logf("Hops: IP: %s, GWMac: %s, GWMacVendor: %s, PktSent: %d, PktRcvd: %d, LatencyMin: %d, LatencyMax: %d, LatencyAvg: %d, LatencyDiff: %d",
						hop.IP, hop.GWMac, hop.GWMacVendor, hop.PktSent, hop.PktRcvd, hop.LatencyMin, hop.LatencyMax, hop.LatencyAvg, hop.LatencyDiff)
				}
			}
		}
	}
}
