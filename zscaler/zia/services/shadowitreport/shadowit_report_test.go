package shadowitreport

import (
	"context"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestGetAllCloudAppsLite(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "shadowitreport", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	ctx := context.Background() // Create a context

	// Retrieve all cloud applications without optional parameters
	apps, err := GetAllCloudAppsLite(ctx, service, nil, nil) // Pass context
	if err != nil {
		t.Fatalf("Error getting all cloud applications: %v", err)
	}
	if len(apps) == 0 {
		t.Log("No cloud applications found")
		return
	}

	// Retrieve all cloud applications with optional parameters
	pageNumber := 1
	limit := 10
	apps, err = GetAllCloudAppsLite(ctx, service, &pageNumber, &limit) // Pass context
	if err != nil {
		t.Fatalf("Error getting all cloud applications with pagination: %v", err)
	}
	if len(apps) == 0 {
		t.Log("No cloud applications found with pagination")
		return
	}
}

func TestGetAllCustomTags(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "shadowitreport", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	ctx := context.Background() // Create a context

	// Retrieve all custom tags without optional parameters
	apps, err := GetAllCustomTags(ctx, service) // Pass context
	if err != nil {
		t.Fatalf("Error getting all custom tags: %v", err)
	}
	if len(apps) == 0 {
		t.Log("No custom tags found")
		return
	}
}
