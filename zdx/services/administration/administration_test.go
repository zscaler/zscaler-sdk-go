package administration

import (
	"fmt"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestAdministrationDepartments(t *testing.T) {
	client, err := tests.NewZdxClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := New(client)

	// Fetch all departments
	allDepartments, err := service.GetAllDepartments()
	if err != nil {
		t.Fatalf("Error fetching all departments: %v", err)
	}
	if len(allDepartments) == 0 {
		t.Logf("No departments found. Test considered successful.")
		return // Exit the test successfully as no departments are available
	}
	firstDepartmentID := allDepartments[0].ID

	// Fetch all locations
	allLocations, err := service.GetAllLocations()
	if err != nil {
		t.Fatalf("Error fetching all locations: %v", err)
	}
	if len(allLocations) == 0 {
		t.Logf("No locations found. Test considered successful.")
		return // Exit the test successfully as no locations are available
	}
	firstLocationID := allLocations[0].ID

	// Example test for GetDepartments with first department ID
	departmentsFilters := GetDepartmentsFilters{
		Loc:    []int{firstLocationID},
		Search: "SearchTerm", // Use an appropriate search term or leave it empty
	}
	departments, resp, err := service.GetDepartments(fmt.Sprintf("%d", firstDepartmentID), departmentsFilters)
	if err != nil {
		t.Errorf("Error getting departments: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got: %v", resp.StatusCode)
	}
	if departments == nil {
		t.Error("Expected non-nil departments")
	}

	// Additional assertions based on expected behavior

	// Example test for GetLocations with first location ID
	locationsFilters := GetLocationsFilters{
		Loc:    []int{firstDepartmentID},
		Search: "SearchTerm", // Use an appropriate search term or leave it empty
	}
	locations, resp, err := service.GetLocations(fmt.Sprintf("%d", firstLocationID), locationsFilters)
	if err != nil {
		t.Errorf("Error getting locations: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got: %v", resp.StatusCode)
	}
	if locations == nil {
		t.Error("Expected non-nil locations")
	}

}
