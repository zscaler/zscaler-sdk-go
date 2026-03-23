package process_based_apps_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/process_based_apps"
)

func TestGetProcessBasedApps(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	result, _, err := process_based_apps.GetProcessBasedApps(context.Background(), service, "", nil, nil)
	if err != nil {
		t.Fatalf("Error getting process-based apps: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil response")
	}

	t.Logf("Total process-based apps: %d", result.TotalCount)

	if result.TotalCount > 0 && len(result.AppIdentities) == 0 {
		t.Error("TotalCount > 0 but no apps returned")
	}

	for i, app := range result.AppIdentities {
		t.Logf("App %d: ID=%d, Name=%s, MatchingCriteria=%d", i, app.ID, app.AppName, app.MatchingCriteria)
		if app.ID == 0 {
			t.Errorf("Expected non-zero ID for app at index %d", i)
		}
		if app.AppName == "" {
			t.Errorf("Expected non-empty AppName for app at index %d", i)
		}
	}
}

func TestGetProcessBasedAppByID(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	listResult, _, err := process_based_apps.GetProcessBasedApps(context.Background(), service, "", nil, nil)
	if err != nil {
		t.Fatalf("Error listing process-based apps: %v", err)
	}

	if listResult == nil || len(listResult.AppIdentities) == 0 {
		t.Log("No process-based apps found to test GetByAppID. Skipping.")
		return
	}

	firstApp := listResult.AppIdentities[0]
	appID := fmt.Sprintf("%d", firstApp.ID)

	result, _, err := process_based_apps.GetByAppID(context.Background(), service, appID)
	if err != nil {
		t.Fatalf("Error getting process-based app by ID %s: %v", appID, err)
	}

	if result == nil {
		t.Fatalf("Expected non-nil response for app ID %s", appID)
	}

	t.Logf("Retrieved app: ID=%d, Name=%s, MatchingCriteria=%d", result.ID, result.AppName, result.MatchingCriteria)

	if result.ID != firstApp.ID {
		t.Errorf("Expected ID %d, got %d", firstApp.ID, result.ID)
	}
	if result.AppName == "" {
		t.Error("Expected non-empty AppName")
	}

	t.Logf("FilePaths: %v", result.FilePaths)
	t.Logf("FileNames: %v", result.FileNames)
	t.Logf("SignaturePayload: %s", result.SignaturePayload)
	t.Logf("CertificatePayload: %s", result.CertificatePayload)
}

func TestGetProcessBasedAppByName(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	listResult, _, err := process_based_apps.GetProcessBasedApps(context.Background(), service, "", nil, nil)
	if err != nil {
		t.Fatalf("Error listing process-based apps: %v", err)
	}

	if listResult == nil || len(listResult.AppIdentities) == 0 {
		t.Log("No process-based apps found to test GetByName. Skipping.")
		return
	}

	targetName := listResult.AppIdentities[0].AppName
	t.Logf("Searching for app by name: %s", targetName)

	result, _, err := process_based_apps.GetByName(context.Background(), service, targetName)
	if err != nil {
		t.Fatalf("Error getting process-based app by name %q: %v", targetName, err)
	}

	if result == nil {
		t.Fatalf("Expected non-nil response for app name %q", targetName)
	}

	t.Logf("Found app: ID=%d, Name=%s", result.ID, result.AppName)

	if result.AppName != targetName {
		t.Errorf("Expected AppName %q, got %q", targetName, result.AppName)
	}
}

func TestGetProcessBasedAppsWithPagination(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	page := 1
	pageSize := 2

	result, _, err := process_based_apps.GetProcessBasedApps(context.Background(), service, "", &page, &pageSize)
	if err != nil {
		t.Fatalf("Error getting process-based apps with pagination: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil response")
	}

	t.Logf("Page %d (size %d): TotalCount=%d, Returned=%d", page, pageSize, result.TotalCount, len(result.AppIdentities))

	if len(result.AppIdentities) > pageSize {
		t.Errorf("Expected at most %d apps, got %d", pageSize, len(result.AppIdentities))
	}
}
