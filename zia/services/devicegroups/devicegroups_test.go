package devicegroups

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestDeviceGroup_data(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	groups, err := service.GetAllDevicesGroups()
	if err != nil {
		t.Errorf("Error getting device groups: %v", err)
		return
	}
	if len(groups) == 0 {
		t.Errorf("No device group found")
		return
	}
	name := groups[0].Name
	t.Log("Getting device group by name:" + name)
	server, err := service.GetDeviceGroupByName(name)
	if err != nil {
		t.Errorf("Error getting device group by name: %v", err)
		return
	}
	if server.Name != name {
		t.Errorf("device group name does not match: expected %s, got %s", name, server.Name)
		return
	}
	// Negative Test: Try to retrieve an device group with a non-existent name
	nonExistentName := "ThisDeviceGroupDoesNotExist"
	_, err = service.GetDeviceGroupByName(nonExistentName)
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

	groups, err := service.GetAllDevicesGroups()
	if err != nil {
		t.Errorf("Error getting device group: %v", err)
		return
	}
	if len(groups) == 0 {
		t.Errorf("No device group found")
		return
	}

	// Validate device group
	for _, group := range groups {
		// Checking if essential fields are not empty
		if group.ID == 0 {
			t.Errorf("device group ID is empty")
		}
		if group.Name == "" {
			t.Errorf("device group Name is empty")
		}
	}
}
