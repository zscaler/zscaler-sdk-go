package devices

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
)

func TestGetEvents(t *testing.T) {
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
		t.Log("No devices found, skipping GetEvents test.")
		return
	}

	firstDeviceID := devices[0].ID

	// Call GetEvents with the first device's ID
	deviceEvents, resp, err := GetEvents(context.Background(), service, firstDeviceID, common.GetFromToFilters{From: int(from), To: int(to)})
	if err != nil {
		t.Fatalf("Error getting events: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(deviceEvents) == 0 {
		t.Log("No events found.")
	} else {
		t.Logf("Retrieved events for device ID: %d", firstDeviceID)
		for _, deviceEvent := range deviceEvents {
			for _, event := range deviceEvent.Events {
				t.Logf("Event Category: %s, Name: %s, DisplayName: %s, Prev: %s, Curr: %s", event.Category, event.Name, event.DisplayName, event.Prev, event.Curr)
			}
		}
	}
}
