package integration

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
	"github.com/zscaler/zscaler-sdk-go/zia/services/usermanagement"
)

func TestAccDepartmentManagement(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := usermanagement.New(client)

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
}
