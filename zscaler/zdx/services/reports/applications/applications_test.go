package applications

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
)

func TestGetAllApps(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Define a time filter for the last 2 hours
	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()
	filters := common.GetFromToFilters{
		From: int(from),
		To:   int(to),
	}

	apps, resp, err := GetAllApps(context.Background(), service, filters)
	if err != nil {
		t.Fatalf("Error getting all apps: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(apps) == 0 {
		t.Log("No apps found.")
	} else {
		t.Logf("Retrieved %d apps", len(apps))
		for _, app := range apps {
			t.Logf("App ID: %d, Name: %s, Score: %f", app.ID, app.Name, app.Score)
		}
	}
}

func TestGetApp(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Define a time filter for the last 2 hours
	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()
	filters := common.GetFromToFilters{
		From: int(from),
		To:   int(to),
	}

	apps, _, err := GetAllApps(context.Background(), service, filters)
	if err != nil {
		t.Fatalf("Error getting all apps: %v", err)
	}

	if len(apps) == 0 {
		t.Log("No apps found, skipping GetApp test.")
		return
	}

	firstAppID := apps[0].ID

	// Get the specific app by ID
	app, resp, err := GetApp(context.Background(), service, strconv.Itoa(firstAppID), filters)
	if err != nil {
		t.Fatalf("Error getting app: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if app.ID == 0 {
		t.Log("No app found with the specified ID.")
	} else {
		t.Logf("Retrieved app ID: %d, Name: %s, Score: %f", app.ID, app.Name, app.Score)
	}
}
