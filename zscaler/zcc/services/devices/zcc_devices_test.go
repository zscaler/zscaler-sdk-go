package devices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zcc/services/devices"
)

func TestGetDevices(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Define test cases
	testCases := []struct {
		username string
		osType   string
	}{
		{"", ""}, // No filters
		{"", ""}, // Filter by username only
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("GetAll with filters - username: %s, osType: %s", tc.username, tc.osType), func(t *testing.T) {
			t.Logf("Testing with filters - username: %s, osType: %s", tc.username, tc.osType)

			devicesList, err := devices.GetAll(context.Background(), service, tc.username, tc.osType)
			if err != nil {
				t.Errorf("Error while getting devices with username=%s, osType=%s: %v", tc.username, tc.osType, err)
				return
			}

			// Log the raw response for debugging
			t.Logf("Raw devices list response: %+v", devicesList)

			// Check if the response slice is not nil
			if devicesList == nil {
				t.Errorf("Expected non-nil slice of devices")
				return
			}

			// Log the number of devices returned for the given filters
			t.Logf("Number of devices returned for username=%s, osType=%s: %d", tc.username, tc.osType, len(devicesList))

			// Check specific fields in the returned structure if necessary
			if len(devicesList) > 0 {
				device := devicesList[0] // Check the first device as an example
				if device.Udid == "" {
					t.Errorf("Expected non-empty UDID for the first device")
				}
				if device.CompanyName == "" {
					t.Errorf("Expected non-empty CompanyName for the first device")
				}
				t.Logf("First device details: %+v", device)
			} else {
				t.Logf("No devices returned for username=%s, osType=%s", tc.username, tc.osType)
			}
		})
	}
}
