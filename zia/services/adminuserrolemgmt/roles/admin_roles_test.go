package roles

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestAdminRoles_data(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	roles, err := service.GetAllAdminRoles()
	if err != nil {
		t.Errorf("Error getting admin roles: %v", err)
		return
	}
	if len(roles) == 0 {
		t.Errorf("No admin roles found")
		return
	}
	name := roles[0].Name
	t.Log("Getting admin roles by name:" + name)
	role, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting admin roles by name: %v", err)
		return
	}
	if role.Name != name {
		t.Errorf("admin role name does not match: expected %s, got %s", name, role.Name)
		return
	}
	// Negative Test: Try to retrieve an admin role with a non-existent name
	nonExistentName := "ThisAdminRoleDoesNotExist"
	_, err = service.GetByName(nonExistentName)
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

	roles, err := service.GetAllAdminRoles()
	if err != nil {
		t.Errorf("Error getting admin role: %v", err)
		return
	}
	if len(roles) == 0 {
		t.Errorf("No admin role found")
		return
	}

	// Validate admin role
	for _, role := range roles {
		// Checking if essential fields are not empty
		if role.ID == 0 {
			t.Errorf("admin role ID is empty")
		}
		if role.Name == "" {
			t.Errorf("admin role Name is empty")
		}
	}
}
