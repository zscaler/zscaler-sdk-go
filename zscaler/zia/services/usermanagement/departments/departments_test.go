package departments

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestAccDepartmentManagement(t *testing.T) {
	// Step 1: Create the general ZIA client
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Step 2: Create a department-specific service using the ZIA client
	departmentService := New(service.Client)

	// Step 3: Test fetching all departments
	departments, err := departmentService.GetAll() // Use departmentService, not the general service
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
	department, err := departmentService.GetDepartmentsByName(name) // Use departmentService
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
	_, err = departmentService.GetDepartmentsByName(nonExistentName) // Use departmentService
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}

func TestResponseFormatValidation(t *testing.T) {
	// Step 1: Create the general ZIA client
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Step 2: Create a department-specific service using the ZIA client
	departmentService := New(service.Client)

	// Step 3: Fetch all departments
	departments, err := departmentService.GetAll()
	if err != nil {
		t.Errorf("Error getting department: %v", err)
		return
	}
	if len(departments) == 0 {
		t.Errorf("No department found")
		return
	}

	// Step 4: Validate the departments' response
	for _, department := range departments {
		// Checking if essential fields are not empty
		if department.ID == 0 {
			t.Errorf("department ID is empty")
		}
		if department.Name == "" {
			t.Errorf("department Name is empty")
		}
	}
}

func TestAllFieldsDepartments(t *testing.T) {
	// Step 1: Create the general ZIA client
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Step 2: Create a department-specific service using the ZIA client
	departmentService := New(service.Client)

	// Step 3: Fetch all departments
	departments, err := departmentService.GetAll()
	if err != nil {
		t.Errorf("Error getting all departments: %v", err)
		return
	}

	if len(departments) == 0 {
		t.Errorf("No department found")
		return
	}

	// Step 4: Get a specific department by ID
	specificID := departments[0].ID
	department, err := departmentService.GetDepartments(specificID)
	if err != nil {
		t.Errorf("Error getting department by ID: %v", err)
		return
	}

	// Step 5: Validate all fields
	if department.ID == 0 {
		t.Errorf("ID is empty")
	}
	if department.Name == "" {
		t.Errorf("Name is empty")
	}
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	// Step 1: Create the general ZIA client
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Step 2: Create a department-specific service using the ZIA client
	departmentService := New(service.Client)

	// Assuming a department with the name "Engineering" exists
	knownName := "Engineering"

	// Case variations to test
	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
		cases.Title(language.English).String(knownName),
	}

	// Step 3: Test retrieving the department with different case variations
	for _, variation := range variations {
		t.Logf("Attempting to retrieve department with name variation: %s", variation)
		department, err := departmentService.GetDepartmentsByName(variation)
		if err != nil {
			t.Errorf("Error getting department with name variation '%s': %v", variation, err)
			continue
		}

		// Check if the department's actual name matches the known name
		if department.Name != knownName {
			t.Errorf("Expected department name to be '%s' for variation '%s', but got '%s'", knownName, variation, department.Name)
		}
	}
}
