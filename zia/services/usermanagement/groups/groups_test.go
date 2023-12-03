package groups

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestAccGroupManagement(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	groups, err := service.GetAllGroups()
	if err != nil {
		t.Errorf("Error getting groups: %v", err)
		return
	}
	if len(groups) == 0 {
		t.Errorf("No groups found")
		return
	}
	name := groups[0].Name
	t.Log("Getting groups by name:" + name)
	group, err := service.GetGroupByName(name)
	if err != nil {
		t.Errorf("Error getting groups by name: %v", err)
		return
	}
	if group.Name != name {
		t.Errorf("group name does not match: expected %s, got %s", name, group.Name)
		return
	}
	// Negative Test: Try to retrieve a group with a non-existent name
	nonExistentName := "ThisGroupDoesNotExist"
	_, err = service.GetGroupByName(nonExistentName)
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

	groups, err := service.GetAllGroups()
	if err != nil {
		t.Errorf("Error getting group: %v", err)
		return
	}
	if len(groups) == 0 {
		t.Errorf("No group found")
		return
	}

	// Validate group
	for _, group := range groups {
		// Checking if essential fields are not empty
		if group.ID == 0 {
			t.Errorf("group ID is empty")
		}
		if group.Name == "" {
			t.Errorf("group Name is empty")
		}
	}
}

func TestAllFieldsGroups(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)
	groups, err := service.GetAllGroups()
	if err != nil {
		t.Errorf("Error getting all Groups: %v", err)
		return
	}

	if len(groups) == 0 {
		t.Errorf("No SCIM Group found")
		return
	}

	specificID := groups[0].ID
	group, err := service.GetGroups(specificID)
	if err != nil {
		t.Errorf("Error getting group by ID: %v", err)
		return
	}

	// Now check each field
	if group.ID == 0 {
		t.Errorf("ID is empty")
	}
	if group.Name == "" {
		t.Errorf("Name is empty")
	}
}
