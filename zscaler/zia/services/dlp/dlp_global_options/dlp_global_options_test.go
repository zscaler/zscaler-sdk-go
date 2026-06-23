package dlp_global_options

import (
	"context"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

// TestDLPGlobalOptionsGet validates that the global options endpoint returns
// data without error.
func TestDLPGlobalOptionsGet(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	options, err := GetDLPGlobalOptions(context.Background(), service)
	if err != nil {
		t.Fatalf("Error retrieving DLP global options: %v", err)
	}
	if options == nil {
		t.Fatal("Expected DLP global options, got nil")
	}
}

// TestDLPGlobalOptionsUpdate toggles a boolean option, verifies the change, and
// then restores the original value so the test is non-destructive.
func TestDLPGlobalOptionsUpdate(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	ctx := context.Background()

	original, err := GetDLPGlobalOptions(ctx, service)
	if err != nil {
		t.Fatalf("Error retrieving initial DLP global options: %v", err)
	}

	// Always restore the original configuration at the end of the test.
	defer func() {
		if _, _, err := UpdateDLPGlobalOptions(ctx, service, *original); err != nil {
			t.Logf("Warning: failed to restore original DLP global options: %v", err)
		}
	}()

	// Toggle a low-risk boolean flag.
	updated := *original
	updated.ExemptUrlEncodedData = !original.ExemptUrlEncodedData

	result, _, err := UpdateDLPGlobalOptions(ctx, service, updated)
	if err != nil {
		t.Fatalf("Error updating DLP global options: %v", err)
	}
	if result == nil {
		t.Fatal("Expected updated DLP global options, got nil")
	}

	// Verify the change took effect.
	verify, err := GetDLPGlobalOptions(ctx, service)
	if err != nil {
		t.Fatalf("Error fetching updated DLP global options: %v", err)
	}
	if verify.ExemptUrlEncodedData != updated.ExemptUrlEncodedData {
		t.Errorf("Expected exemptUrlEncodedData=%v, got %v",
			updated.ExemptUrlEncodedData, verify.ExemptUrlEncodedData)
	}
}
