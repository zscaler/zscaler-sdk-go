package devices

import (
	"net/http"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/applications"
)

func TestGetQualityMetrics(t *testing.T) {
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
		t.Log("No devices found, skipping GetQualityMetrics test.")
		return
	}

	firstDeviceID := devices[0].ID

	// Invoke GetAllApps to retrieve the ID of the first app
	apps, _, err := applications.GetAllApps(service, filters)
	if err != nil {
		t.Fatalf("Error getting all apps: %v", err)
	}

	if len(apps) == 0 {
		t.Log("No apps found, skipping GetQualityMetrics test.")
		return
	}

	firstAppID := apps[0].ID

	// Call GetQualityMetrics with the first device's ID and first app's ID
	qualityMetrics, resp, err := GetQualityMetrics(service, firstDeviceID, firstAppID, filters)
	if err != nil {
		t.Fatalf("Error getting quality metrics: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(qualityMetrics) == 0 {
		t.Log("No quality metrics found.")
	} else {
		for _, qm := range qualityMetrics {
			t.Logf("Retrieved quality metrics for Meet ID: %s, Meet Session ID: %s, Meet Subject: %s", qm.MeetID, qm.MeetSessionID, qm.MeetSubject)
			for _, metric := range qm.Metrics {
				t.Logf("Metric: %s, Unit: %s", metric.Metric, metric.Unit)
				for _, dp := range metric.DataPoints {
					t.Logf("Timestamp: %d, Value: %f", dp.TimeStamp, dp.Value)
				}
			}
		}
	}
}
