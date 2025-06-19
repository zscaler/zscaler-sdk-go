package administration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestGetDepartments(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()

	filters := GetDepartmentsFilters{
		From:   int(from),
		To:     int(to),
		Search: "A000",
	}

	departments, resp, err := GetDepartments(context.Background(), service, filters)
	if err != nil {
		t.Fatalf("Error getting departments: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(departments) == 0 {
		t.Log("No departments found within the specified time range.")
	} else {
		t.Logf("Retrieved %d departments", len(departments))
		for _, department := range departments {
			t.Logf("Department ID: %d, Name: %s", department.ID, department.Name)
		}
	}
}

func TestGetLocations(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()

	filters := GetLocationsFilters{
		From:   int(from),
		To:     int(to),
		Search: "Road Warrior",
	}

	locations, resp, err := GetLocations(context.Background(), service, filters)
	if err != nil {
		t.Fatalf("Error getting locations: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(locations) == 0 {
		t.Log("No locations found within the specified time range.")
	} else {
		t.Logf("Retrieved %d locations", len(locations))
		for _, location := range locations {
			t.Logf("Location ID: %d, Name: %s", location.ID, location.Name)
		}
	}
}
