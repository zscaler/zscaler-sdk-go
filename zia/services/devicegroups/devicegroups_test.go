package devicegroups

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func TestCaseSensitivityOfGetByName(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	// Assuming a group with the name "IOS" exists
	knownName := "IOS"

	// Case variations to test
	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
		cases.Title(language.English).String(knownName),
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve group with name variation: %s", variation)
		group, err := service.GetDeviceGroupByName(variation)
		if err != nil {
			t.Errorf("Error getting device group with name variation '%s': %v", variation, err)
			continue
		}

		// Check if the group's actual name matches the known name
		if group.Name != knownName {
			t.Errorf("Expected group name to be '%s' for variation '%s', but got '%s'", knownName, variation, group.Name)
		}
	}
}

func TestDeviceGroupFields(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := New(client)

	// Retrieve all device groups
	groups, err := service.GetAllDevicesGroups()
	if err != nil {
		t.Fatalf("Error getting all device groups: %v", err)
	}
	if len(groups) == 0 {
		t.Fatalf("No device groups found for testing")
	}

	// Check the first device group
	firstGroup := groups[0]
	if firstGroup.GroupType != "ZCC_OS" {
		t.Errorf("Group Type field is incorrect, expected 'ZCC_OS', got '%s'", firstGroup.GroupType)
	}
	if firstGroup.OSType != "IOS" {
		t.Errorf("OS Type field is incorrect, expected 'IOS', got '%s'", firstGroup.OSType)
	}
	if !firstGroup.Predefined {
		t.Errorf("Predefined field is not set to true as expected")
	}
}
