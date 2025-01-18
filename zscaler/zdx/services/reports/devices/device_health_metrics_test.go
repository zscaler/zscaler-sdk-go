package devices

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
)

func TestGetHealthMetrics(t *testing.T) {
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
	devices, _, err := GetAllDevices(context.Background(), service, filters)
	if err != nil {
		t.Fatalf("Error getting all devices: %v", err)
	}

	if len(devices) == 0 {
		t.Log("No devices found, skipping GetHealthMetrics test.")
		return
	}

	firstDeviceID := devices[0].ID

	// Call GetHealthMetrics with the first device's ID and first app's ID
	healthMetrics, resp, err := GetHealthMetrics(context.Background(), service, firstDeviceID, common.GetFromToFilters{From: int(from), To: int(to)})
	if err != nil {
		t.Fatalf("Error getting health metrics: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(healthMetrics) == 0 {
		t.Log("No health metrics found.")
	} else {
		for _, metric := range healthMetrics {
			t.Logf("Retrieved health metrics: Category: %s", metric.Category)
			for _, instance := range metric.Instances {
				t.Logf("Instance: %+v", instance)
			}
		}
	}
}
