package shadowitreport

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
)

func TestGetAllCloudAppsLite(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	// Retrieve all cloud applications without optional parameters
	apps, err := GetAllCloudAppsLite(service, nil, nil)
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
	apps, err = GetAllCloudAppsLite(service, &pageNumber, &limit)
	if err != nil {
		t.Fatalf("Error getting all cloud applications with pagination: %v", err)
	}
	if len(apps) == 0 {
		t.Log("No cloud applications found with pagination")
		return
	}
}

func TestGetAllCustomTags(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	// Retrieve all cloud applications without optional parameters
	apps, err := GetAllCustomTags(service)
	if err != nil {
		t.Fatalf("Error getting all cloud applications: %v", err)
	}
	if len(apps) == 0 {
		t.Log("No cloud applications found")
		return
	}
}
