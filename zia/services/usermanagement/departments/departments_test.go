package departments

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestAccDepartmentManagement(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	departments, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting departments: %v", err)
		return
	}
	if len(departments) == 0 {
		t.Errorf("No departments found")
		return
	}
	name := departments[0].Name
	t.Log("Getting departments by name:" + name)
	department, err := service.GetDepartmentsByName(name)
	if err != nil {
		t.Errorf("Error getting departments by name: %v", err)
		return
	}
	if department.Name != name {
		t.Errorf("department name does not match: expected %s, got %s", name, department.Name)
		return
	}

	// Negative Test: Try to retrieve a department with a non-existent name
	nonExistentName := "ThisDepartmentDoesNotExist"
	_, err = service.GetDepartmentsByName(nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}

func TestResponseFormatValidation(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	departments, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting department: %v", err)
		return
	}
	if len(departments) == 0 {
		t.Errorf("No department found")
		return
	}

	// Validate department
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
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)
	departments, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting all departments: %v", err)
		return
	}

	if len(departments) == 0 {
		t.Errorf("No department found")
		return
	}

	specificID := departments[0].ID
	department, err := service.GetDepartments(specificID)
	if err != nil {
		t.Errorf("Error getting department by ID: %v", err)
		return
	}

	// Now check each field
	if department.ID == 0 {
		t.Errorf("ID is empty")
	}
	if department.Name == "" {
		t.Errorf("Name is empty")
	}
	if department.Comments == "" {
		t.Errorf("Name is empty")
	}
}
