package devices

import (
	"net/http"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/troubleshooting/deeptrace"
)

func TestGetDeviceTopProcesses(t *testing.T) {
	client, err := tests.NewZdxClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	// Step 1: Get all devices and retrieve the first device ID
	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()
	filters := common.GetFromToFilters{
		From: int(from),
		To:   int(to),
	}

	devices, _, err := GetAllDevices(service, GetDevicesFilters{
		GetFromToFilters: filters,
	})
	if err != nil {
		t.Fatalf("Error getting all devices: %v", err)
	}

	if len(devices) == 0 {
		t.Log("No devices found, skipping test.")
		return
	}

	deviceID := devices[0].ID

	// Step 2: Get all deep traces for the device
	deepTraces, _, err := deeptrace.GetDeepTraces(service, deviceID)
	if err != nil {
		t.Fatalf("Error getting deep traces: %v", err)
	}

	if len(deepTraces) == 0 {
		t.Log("No deep traces found, skipping test.")
		return
	}

	// Step 3: Retrieve the first trace ID from the list of deep traces
	traceID := deepTraces[0].TraceID

	// Step 4: Call GetDeviceTopProcesses with the device ID and trace ID
	topProcesses, resp, err := GetDeviceTopProcesses(service, deviceID, traceID, filters)
	if err != nil {
		t.Fatalf("Error getting device top processes: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(topProcesses) == 0 {
		t.Log("No top processes found.")
	} else {
		for _, process := range topProcesses {
			t.Logf("Timestamp: %d", process.TimeStamp)
			for _, topProcess := range process.TopProcesses {
				t.Logf("Category: %s", topProcess.Category)
				for _, proc := range topProcess.Processes {
					t.Logf("Process ID: %d, Name: %s", proc.ID, proc.Name)
				}
			}
		}
	}
}
