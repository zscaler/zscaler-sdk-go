package deeptrace

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/applications"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/devices"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestDeepTraceSession(t *testing.T) {
	name := "TestSession-" + acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	client, err := tests.NewZdxClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	// Step 1: Get all devices and retrieve the first device ID
	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()
	deviceFilters := devices.GetDevicesFilters{
		GetFromToFilters: common.GetFromToFilters{
			From: int(from),
			To:   int(to),
		},
	}

	device, _, err := devices.GetAllDevices(context.Background(), service, deviceFilters)
	if err != nil {
		t.Fatalf("Error getting all devices: %v", err)
	}

	if len(device) == 0 {
		t.Log("No devices found, skipping deep trace session test.")
		return
	}

	deviceID := device[0].ID

	// Step 2: Get all apps and retrieve the first app ID
	appFilters := common.GetFromToFilters{
		From: int(from),
		To:   int(to),
	}

	apps, _, err := applications.GetAllApps(context.Background(), service, appFilters)
	if err != nil {
		t.Fatalf("Error getting all apps: %v", err)
	}

	if len(apps) == 0 {
		t.Log("No apps found, skipping deep trace session test.")
		return
	}

	appID := apps[0].ID

	// Step 3: Get all web probes and retrieve the first web probe ID
	webProbes, _, err := devices.GetAllWebProbes(context.Background(), service, deviceID, appID, appFilters)
	if err != nil {
		t.Fatalf("Error getting all web probes: %v", err)
	}

	if len(webProbes) == 0 {
		t.Log("No web probes found, skipping deep trace session test.")
		return
	}

	webProbeID := webProbes[0].ID

	// Step 4: Get all cloud path probes and retrieve the first cloud path probe ID
	cloudPathProbes, _, err := devices.GetAllCloudPathProbes(context.Background(), service, deviceID, appID, appFilters)
	if err != nil {
		t.Fatalf("Error getting all cloud path probes: %v", err)
	}

	if len(cloudPathProbes) == 0 {
		t.Log("No cloud path probes found, skipping deep trace session test.")
		return
	}

	cloudPathProbeID := cloudPathProbes[0].ID

	// Step 5: Create a DeepTrace session
	payload := DeepTraceSessionPayload{
		SessionName:          name,
		AppID:                appID,
		WebProbeID:           webProbeID,
		CloudPathProbeID:     cloudPathProbeID,
		SessionLengthMinutes: 5,
		ProbeDevice:          true,
	}

	createdSession, resp, err := CreateDeepTraceSession(context.Background(), service, deviceID, payload)
	if err != nil {
		t.Fatalf("Error creating deep trace session: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code 200 or 201, got %d", resp.StatusCode)
	}

	traceID := createdSession.TraceID
	t.Logf("Created Deep Trace Session: %s", traceID)

	// Step 6: Get all deep traces
	deepTraces, resp, err := GetDeepTraces(context.Background(), service, deviceID)
	if err != nil {
		t.Fatalf("Error getting deep traces: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(deepTraces) == 0 {
		t.Log("No deep traces found.")
	} else {
		t.Logf("Retrieved %d deep traces", len(deepTraces))
		for _, trace := range deepTraces {
			t.Logf("Trace ID: %s, Status: %s", trace.TraceID, trace.Status)
		}
	}

	// Step 7: Get deep trace session
	traceSessionResp, err := GetDeepTraceSession(context.Background(), service, deviceID, traceID)
	if err != nil {
		t.Fatalf("Error getting deep trace session: %v", err)
	}

	if traceSessionResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, traceSessionResp.StatusCode)
	}

	t.Logf("Retrieved deep trace session: %s", traceID)

	// Step 8: Pause the test for 60 seconds
	t.Log("Pausing test for 60 seconds...")
	time.Sleep(60 * time.Second)

	// Step 9: Delete the deep trace session
	deleteResp, err := DeleteDeepTraceSession(context.Background(), service, deviceID, traceID)
	if err != nil {
		t.Fatalf("Error deleting deep trace session: %v", err)
	}

	if deleteResp.StatusCode != http.StatusOK && deleteResp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status code 200 or 204, got %d", deleteResp.StatusCode)
	}

	t.Logf("Deleted deep trace session: %s", traceID)
}
