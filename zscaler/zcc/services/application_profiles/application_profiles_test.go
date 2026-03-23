package application_profiles_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/application_profiles"
)

func TestGetApplicationProfiles(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	result, _, err := application_profiles.GetApplicationProfiles(context.Background(), service, "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("Error getting application profiles: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil response")
	}

	t.Logf("Total application profiles: %d", result.TotalCount)

	if result.TotalCount > 0 && len(result.Policies) == 0 {
		t.Error("TotalCount > 0 but no policies returned")
	}

	for i, p := range result.Policies {
		t.Logf("Profile %d: ID=%d, Name=%s, DeviceType=%s, Active=%d", i, p.ID, p.Name, p.DeviceType, p.Active)
		if p.ID == 0 {
			t.Errorf("Expected non-zero ID for profile at index %d", i)
		}
		if p.Name == "" {
			t.Errorf("Expected non-empty Name for profile at index %d", i)
		}
	}
}

func TestGetApplicationProfilesWithDeviceType(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	deviceTypes := []string{"windows", "ios", "android", "macos", "linux"}

	for _, dt := range deviceTypes {
		t.Run(fmt.Sprintf("deviceType=%s", dt), func(t *testing.T) {
			result, _, err := application_profiles.GetApplicationProfiles(context.Background(), service, "", "", dt, nil, nil)
			if err != nil {
				t.Fatalf("Error getting application profiles for deviceType %s: %v", dt, err)
			}

			if result == nil {
				t.Fatal("Expected non-nil response")
			}

			t.Logf("DeviceType=%s: TotalCount=%d, Returned=%d", dt, result.TotalCount, len(result.Policies))
		})
	}
}

func TestGetApplicationProfilesWithPagination(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	page := 1
	pageSize := 2

	result, _, err := application_profiles.GetApplicationProfiles(context.Background(), service, "", "", "", &page, &pageSize)
	if err != nil {
		t.Fatalf("Error getting application profiles with pagination: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil response")
	}

	t.Logf("Page %d (size %d): TotalCount=%d, Returned=%d", page, pageSize, result.TotalCount, len(result.Policies))

	if len(result.Policies) > pageSize {
		t.Errorf("Expected at most %d profiles, got %d", pageSize, len(result.Policies))
	}
}

func TestGetApplicationProfileByID(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	listResult, _, err := application_profiles.GetApplicationProfiles(context.Background(), service, "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("Error listing application profiles: %v", err)
	}

	if listResult == nil || len(listResult.Policies) == 0 {
		t.Log("No application profiles found to test GetByProfileID. Skipping.")
		return
	}

	firstProfile := listResult.Policies[0]
	profileID := fmt.Sprintf("%d", firstProfile.ID)

	result, _, err := application_profiles.GetByProfileID(context.Background(), service, profileID)
	if err != nil {
		t.Fatalf("Error getting application profile by ID %s: %v", profileID, err)
	}

	if result == nil {
		t.Fatalf("Expected non-nil response for profile ID %s", profileID)
	}

	t.Logf("Retrieved profile: ID=%d, Name=%s, DeviceType=%s, Active=%d", result.ID, result.Name, result.DeviceType, result.Active)

	if result.ID != firstProfile.ID {
		t.Errorf("Expected ID %d, got %d", firstProfile.ID, result.ID)
	}
	if result.Name == "" {
		t.Error("Expected non-empty Name")
	}
}

func TestGetApplicationProfileByName(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	listResult, _, err := application_profiles.GetApplicationProfiles(context.Background(), service, "", "", "", nil, nil)
	if err != nil {
		t.Fatalf("Error listing application profiles: %v", err)
	}

	if listResult == nil || len(listResult.Policies) == 0 {
		t.Log("No application profiles found to test GetByName. Skipping.")
		return
	}

	targetName := listResult.Policies[0].Name
	t.Logf("Searching for profile by name: %s", targetName)

	result, _, err := application_profiles.GetByName(context.Background(), service, targetName)
	if err != nil {
		t.Fatalf("Error getting application profile by name %q: %v", targetName, err)
	}

	if result == nil {
		t.Fatalf("Expected non-nil response for profile name %q", targetName)
	}

	t.Logf("Found profile: ID=%d, Name=%s, DeviceType=%s", result.ID, result.Name, result.DeviceType)

	if result.Name != targetName {
		t.Errorf("Expected Name %q, got %q", targetName, result.Name)
	}
}
