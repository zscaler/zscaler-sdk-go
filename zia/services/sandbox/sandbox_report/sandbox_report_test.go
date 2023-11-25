package sandbox_report

import (
	"fmt"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestGetRatingQuota(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := &Service{Client: client}

	quotas, err := service.GetRatingQuota()
	if err != nil {
		t.Errorf("Error getting Rating Quotas: %v", err)
	}
	if len(quotas) == 0 {
		t.Error("Expected non-empty Rating Quotas, got empty")
	}
	// Add additional checks for expected data, e.g., check the values of quotas[0] if needed.
}

func TestGetReportMD5Hash(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := &Service{Client: client}

	// Replace with an actual MD5 hash of known malware. This is just an example.
	md5Hashes := []string{"c0202cf6aeab8437c638533d14563d35", "42914d6d213a20a2684064be5c80ffa9"}

	for _, md5Hash := range md5Hashes {
		for _, details := range []string{"full", "summary"} {
			t.Run(fmt.Sprintf("MD5Hash=%s-Details=%s", md5Hash, details), func(t *testing.T) {
				report, err := service.GetReportMD5Hash(md5Hash, details)
				if err != nil {
					t.Errorf("Error getting MD5 Hash Report: %v", err)
				}
				if report == nil {
					t.Error("Expected MD5 Hash Report, got nil")
				}
			})
		}
	}
}
