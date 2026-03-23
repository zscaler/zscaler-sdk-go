package predefined_ip_apps_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/predefined_ip_apps"
)

func TestGetPredefinedIPApps(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	result, _, err := predefined_ip_apps.GetPredefinedIPApps(context.Background(), service, "", nil, nil)
	if err != nil {
		t.Fatalf("Error getting predefined IP-based apps: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil response")
	}

	t.Logf("Total predefined IP-based apps: %d", result.TotalCount)

	if result.TotalCount > 0 && len(result.AppServiceContracts) == 0 {
		t.Error("TotalCount > 0 but no apps returned")
	}

	for i, app := range result.AppServiceContracts {
		t.Logf("App %d: ID=%d, AppSvcId=%d, Name=%s, Active=%t", i, app.ID, app.AppSvcId, app.AppName, app.Active)
		if app.ID == 0 {
			t.Errorf("Expected non-zero ID for app at index %d", i)
		}
		if app.AppName == "" {
			t.Errorf("Expected non-empty AppName for app at index %d", i)
		}
	}
}

func TestGetPredefinedIPAppByID(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	listResult, _, err := predefined_ip_apps.GetPredefinedIPApps(context.Background(), service, "", nil, nil)
	if err != nil {
		t.Fatalf("Error listing predefined IP-based apps: %v", err)
	}

	if listResult == nil || len(listResult.AppServiceContracts) == 0 {
		t.Log("No predefined IP-based apps found to test GetByAppID. Skipping.")
		return
	}

	firstApp := listResult.AppServiceContracts[0]
	appID := fmt.Sprintf("%d", firstApp.ID)

	result, _, err := predefined_ip_apps.GetByAppID(context.Background(), service, appID)
	if err != nil {
		t.Fatalf("Error getting predefined IP-based app by ID %s: %v", appID, err)
	}

	if result == nil {
		t.Fatalf("Expected non-nil response for app ID %s", appID)
	}

	t.Logf("Retrieved app: ID=%d, AppSvcId=%d, Name=%s, Active=%t", result.ID, result.AppSvcId, result.AppName, result.Active)

	if result.ID != firstApp.ID {
		t.Errorf("Expected ID %d, got %d", firstApp.ID, result.ID)
	}
	if result.AppName == "" {
		t.Error("Expected non-empty AppName")
	}
	if result.UID == "" {
		t.Error("Expected non-empty UID")
	}
}

func TestGetPredefinedIPAppByName(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	listResult, _, err := predefined_ip_apps.GetPredefinedIPApps(context.Background(), service, "", nil, nil)
	if err != nil {
		t.Fatalf("Error listing predefined IP-based apps: %v", err)
	}

	if listResult == nil || len(listResult.AppServiceContracts) == 0 {
		t.Log("No predefined IP-based apps found to test GetByName. Skipping.")
		return
	}

	targetName := listResult.AppServiceContracts[0].AppName
	t.Logf("Searching for app by name: %s", targetName)

	result, _, err := predefined_ip_apps.GetByName(context.Background(), service, targetName)
	if err != nil {
		t.Fatalf("Error getting predefined IP-based app by name %q: %v", targetName, err)
	}

	if result == nil {
		t.Fatalf("Expected non-nil response for app name %q", targetName)
	}

	t.Logf("Found app: ID=%d, Name=%s, AppSvcId=%d", result.ID, result.AppName, result.AppSvcId)

	if result.AppName != targetName {
		t.Errorf("Expected AppName %q, got %q", targetName, result.AppName)
	}
}

func TestGetPredefinedIPAppsWithSearch(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	listResult, _, err := predefined_ip_apps.GetPredefinedIPApps(context.Background(), service, "", nil, nil)
	if err != nil {
		t.Fatalf("Error listing predefined IP-based apps: %v", err)
	}

	if listResult == nil || len(listResult.AppServiceContracts) == 0 {
		t.Log("No predefined IP-based apps found to test search. Skipping.")
		return
	}

	searchTerm := listResult.AppServiceContracts[0].AppName
	t.Logf("Searching with term: %s", searchTerm)

	searchResult, _, err := predefined_ip_apps.GetPredefinedIPApps(context.Background(), service, searchTerm, nil, nil)
	if err != nil {
		t.Fatalf("Error searching predefined IP-based apps: %v", err)
	}

	if searchResult == nil {
		t.Fatal("Expected non-nil search response")
	}

	t.Logf("Search returned %d apps for term %q", len(searchResult.AppServiceContracts), searchTerm)
}

func TestGetPredefinedIPAppsWithPagination(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	page := 1
	pageSize := 3

	result, _, err := predefined_ip_apps.GetPredefinedIPApps(context.Background(), service, "", &page, &pageSize)
	if err != nil {
		t.Fatalf("Error getting predefined IP-based apps with pagination: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil response")
	}

	t.Logf("Page %d (size %d): TotalCount=%d, Returned=%d", page, pageSize, result.TotalCount, len(result.AppServiceContracts))

	if len(result.AppServiceContracts) > pageSize {
		t.Errorf("Expected at most %d apps, got %d", pageSize, len(result.AppServiceContracts))
	}
}
