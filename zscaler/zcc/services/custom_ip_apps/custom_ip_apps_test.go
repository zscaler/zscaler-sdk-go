package custom_ip_apps_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/custom_ip_apps"
)

func TestGetCustomIPApps(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	result, _, err := custom_ip_apps.GetCustomIPApps(context.Background(), service, "", nil, nil)
	if err != nil {
		t.Fatalf("Error getting custom IP-based apps: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil response")
	}

	t.Logf("Total custom IP-based apps: %d", result.TotalCount)

	if result.TotalCount > 0 && len(result.CustomAppContracts) == 0 {
		t.Error("TotalCount > 0 but no apps returned")
	}

	for i, app := range result.CustomAppContracts {
		t.Logf("App %d: ID=%d, Name=%s, Active=%t", i, app.ID, app.AppName, app.Active)
		if app.ID == 0 {
			t.Errorf("Expected non-zero ID for app at index %d", i)
		}
		if app.AppName == "" {
			t.Errorf("Expected non-empty AppName for app at index %d", i)
		}
	}
}

func TestGetCustomIPAppByID(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	listResult, _, err := custom_ip_apps.GetCustomIPApps(context.Background(), service, "", nil, nil)
	if err != nil {
		t.Fatalf("Error listing custom IP-based apps: %v", err)
	}

	if listResult == nil || len(listResult.CustomAppContracts) == 0 {
		t.Log("No custom IP-based apps found to test GetByAppID. Skipping.")
		return
	}

	firstApp := listResult.CustomAppContracts[0]
	appID := fmt.Sprintf("%d", firstApp.ID)

	result, _, err := custom_ip_apps.GetByAppID(context.Background(), service, appID)
	if err != nil {
		t.Fatalf("Error getting custom IP-based app by ID %s: %v", appID, err)
	}

	if result == nil {
		t.Fatalf("Expected non-nil response for app ID %s", appID)
	}

	t.Logf("Retrieved app: ID=%d, Name=%s, Active=%t", result.ID, result.AppName, result.Active)

	if result.ID != firstApp.ID {
		t.Errorf("Expected ID %d, got %d", firstApp.ID, result.ID)
	}
	if result.AppName == "" {
		t.Error("Expected non-empty AppName")
	}
}

func TestGetCustomIPAppByName(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	listResult, _, err := custom_ip_apps.GetCustomIPApps(context.Background(), service, "", nil, nil)
	if err != nil {
		t.Fatalf("Error listing custom IP-based apps: %v", err)
	}

	if listResult == nil || len(listResult.CustomAppContracts) == 0 {
		t.Log("No custom IP-based apps found to test GetByName. Skipping.")
		return
	}

	targetName := listResult.CustomAppContracts[0].AppName
	t.Logf("Searching for app by name: %s", targetName)

	result, _, err := custom_ip_apps.GetByName(context.Background(), service, targetName)
	if err != nil {
		t.Fatalf("Error getting custom IP-based app by name %q: %v", targetName, err)
	}

	if result == nil {
		t.Fatalf("Expected non-nil response for app name %q", targetName)
	}

	t.Logf("Found app: ID=%d, Name=%s", result.ID, result.AppName)

	if result.AppName != targetName {
		t.Errorf("Expected AppName %q, got %q", targetName, result.AppName)
	}
}

func TestGetCustomIPAppsWithPagination(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	page := 1
	pageSize := 2

	result, _, err := custom_ip_apps.GetCustomIPApps(context.Background(), service, "", &page, &pageSize)
	if err != nil {
		t.Fatalf("Error getting custom IP-based apps with pagination: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil response")
	}

	t.Logf("Page %d (size %d): TotalCount=%d, Returned=%d", page, pageSize, result.TotalCount, len(result.CustomAppContracts))

	if len(result.CustomAppContracts) > pageSize {
		t.Errorf("Expected at most %d apps, got %d", pageSize, len(result.CustomAppContracts))
	}
}
