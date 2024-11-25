package sandbox_report

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
)

func TestGetRatingQuota(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	quotas, err := GetRatingQuota(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting Rating Quotas: %v", err)
	}
	if len(quotas) == 0 {
		t.Error("Expected non-empty Rating Quotas, got empty")
	}
	// Add additional checks for expected data, e.g., check the values of quotas[0] if needed.
}

func TestGetReportMD5Hash(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Replace with an actual MD5 hash of known malware. This is just an example.
	md5Hashes := []string{"F69CA01D65E6C8F9E3540029E5F6AB92"}

	for _, md5Hash := range md5Hashes {
		for _, details := range []string{"full", "summary"} {
			t.Run(fmt.Sprintf("MD5Hash=%s-Details=%s", md5Hash, details), func(t *testing.T) {
				report, err := GetReportMD5Hash(context.Background(), service, md5Hash, details)

				if err != nil {
					if strings.Contains(err.Error(), "md5 is unknown or analysis has yet not been completed.Please try again later") {
						t.Logf("Known error encountered: %v", err)
						return
					}
					t.Errorf("Error getting MD5 Hash Report: %v", err)
				}

				if report == nil {
					t.Error("Expected MD5 Hash Report, got nil")
				}
			})
		}
	}
}
