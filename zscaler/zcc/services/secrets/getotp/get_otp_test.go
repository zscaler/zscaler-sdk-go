package getotp

import (
	"context"
	"fmt"
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zcc/services/devices"
)

func TestGetOtp(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	deviceList, err := devices.GetAll(context.Background(), service, "", "")
	if err != nil {
		t.Errorf("Error getting devices: %v", err)
		return
	}
	if len(deviceList) == 0 {
		t.Log("No devices found to test. Passing the test.")
		return
	}

	// Extract the UDID from the first device
	udid := deviceList[0].Udid

	// Define test cases
	testCases := []struct {
		udid string
	}{
		{udid}, // Use the UDID from the first device
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(fmt.Sprintf("udid=%s", tc.udid), func(t *testing.T) {
			otpResponse, err := GetOtp(context.Background(), service, tc.udid)
			if err != nil {
				t.Fatalf("Error retrieving OTP for UDID %s: %v", tc.udid, err)
			}

			// Log the raw response for debugging
			t.Logf("Raw OTP response: %+v", otpResponse)

			// Check if the response is not nil
			if otpResponse == nil {
				t.Errorf("Expected non-nil response for UDID %s", tc.udid)
				return
			}

			// Ensure at least one OTP field is populated
			if otpResponse.Otp == "" && otpResponse.ExitOtp == "" && otpResponse.LogoutOtp == "" &&
				otpResponse.RevertOtp == "" && otpResponse.UninstallOtp == "" &&
				otpResponse.ZdpDisableOtp == "" && otpResponse.ZdxDisableOtp == "" &&
				otpResponse.ZiaDisableOtp == "" && otpResponse.ZpaDisableOtp == "" {
				t.Errorf("Expected at least one non-empty OTP for UDID %s", tc.udid)
			}
		})
	}
}
