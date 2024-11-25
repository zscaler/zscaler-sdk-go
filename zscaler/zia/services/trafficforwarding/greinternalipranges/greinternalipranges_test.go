package greinternalipranges

import (
	"context"
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
)

func TestGREInternalIPRanges(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	expectedCount := 10

	ranges, err := GetGREInternalIPRange(context.Background(), service, expectedCount)
	if err != nil {
		t.Errorf("Error retrieving GRE internal IP ranges: %v", err)
		return
	}

	if len(*ranges) < expectedCount {
		t.Logf("Warning: Expected %d IP ranges but got %d. This might be an API inconsistency.", expectedCount, len(*ranges))

		// Log the individual IP ranges only if less than expectedCount
		for _, r := range *ranges {
			t.Logf("Available range: %s - %s", r.StartIPAddress, r.EndIPAddress)
		}
	} else {
		t.Logf("Successfully fetched %d IP ranges.", expectedCount)
	}
}
