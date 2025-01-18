package departments

import (
	"context"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestAccDepartmentManagement(t *testing.T) {
	// Step 1: Create the general ZIA client
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Step 3: Test fetching all departments
	departments, err := GetAll(context.Background(), service) // Pass context and service
	if err != nil {
		t.Errorf("Error getting departments: %v", err)
		return
	}
	if len(departments) == 0 {
		t.Errorf("No departments found")
		return
	}

	// Step 4: Test getting departments by name
	name := departments[0].Name
	t.Log("Getting departments by name: " + name)
	department, err := GetDepartmentsByName(context.Background(), service, name) // Pass context and service
	if err != nil {
		t.Errorf("Error getting department by name: %v", err)
		return
	}
	if department.Name != name {
		t.Errorf("Department name does not match: expected %s, got %s", name, department.Name)
		return
	}

	// Step 5: Negative test: Try to retrieve a non-existent department
	nonExistentName := "ThisDepartmentDoesNotExist"
	_, err = GetDepartmentsByName(context.Background(), service, nonExistentName) // Pass context and service
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}
