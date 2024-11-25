package administration

import (
	"net/http"
	"testing"
	"time"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zdx/services"
)

func TestGetDepartments(t *testing.T) {
	client, err := tests.NewZdxClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()

	filters := GetDepartmentsFilters{
		From:   int(from),
		To:     int(to),
		Search: "A000",
	}

	departments, resp, err := GetDepartments(service, filters)
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
	client, err := tests.NewZdxClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()

	filters := GetLocationsFilters{
		From:   int(from),
		To:     int(to),
		Search: "Road Warrior",
	}

	locations, resp, err := GetLocations(service, filters)
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
