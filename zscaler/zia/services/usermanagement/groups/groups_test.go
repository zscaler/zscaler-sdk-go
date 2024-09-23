package groups

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestAccGroupManagement(t *testing.T) {
	// Step 1: Create the general ZIA client
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Step 2: Create a group-specific service using the ZIA client
	groupService := New(service.Client)

	// Step 3: Fetch all groups
	groups, err := groupService.GetAllGroups()
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

	// Step 4: Get a group by name
	group, err := groupService.GetGroupByName(name)
	if err != nil {
		t.Errorf("Error getting groups by name: %v", err)
		return
	}
	if group.Name != name {
		t.Errorf("Group name does not match: expected %s, got %s", name, group.Name)
		return
	}

	// Step 5: Negative test: Try to retrieve a group with a non-existent name
	nonExistentName := "ThisGroupDoesNotExist"
	_, err = groupService.GetGroupByName(nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}
