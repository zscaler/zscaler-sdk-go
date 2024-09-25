package devicegroups

import (
	"context"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestDevices_data(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Retrieve all devices
	devices, err := GetAllDevices(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting devices: %v", err)
		return
	}
	if len(devices) == 0 {
		t.Log("No devices found")
		return
	}

	// Use the first device for testing
	firstDevice := devices[0]

	// Test GetDevicesByID
	t.Run("GetDevicesByID", func(t *testing.T) {
		deviceByID, err := GetDevicesByID(context.Background(), service, firstDevice.ID)
		if err != nil {
			t.Errorf("Error getting device by ID: %v", err)
			return
		}
		if deviceByID == nil || deviceByID.ID != firstDevice.ID {
			t.Errorf("Device ID does not match: expected %d, got %d", firstDevice.ID, deviceByID.ID)
		}
	})

	// Test GetDevicesByName
	t.Run("GetDevicesByName", func(t *testing.T) {
		deviceByName, err := GetDevicesByName(context.Background(), service, firstDevice.Name)
		if err != nil {
			t.Errorf("Error getting device by name: %v", err)
			return
		}
		if deviceByName == nil || deviceByName.Name != firstDevice.Name {
			t.Errorf("Device name does not match: expected %s, got %s", firstDevice.Name, deviceByName.Name)
		}
	})

	// Test GetDevicesByModel
	t.Run("GetDevicesByModel", func(t *testing.T) {
		deviceByModel, err := GetDevicesByModel(context.Background(), service, firstDevice.DeviceModel)
		if err != nil {
			t.Errorf("Error getting device by model: %v", err)
			return
		}
		if deviceByModel == nil || deviceByModel.DeviceModel != firstDevice.DeviceModel {
			t.Errorf("Device model does not match: expected %s, got %s", firstDevice.DeviceModel, deviceByModel.DeviceModel)
		}
	})

	// Test GetDevicesByOwner
	t.Run("GetDevicesByOwner", func(t *testing.T) {
		deviceByOwner, err := GetDevicesByOwner(context.Background(), service, firstDevice.OwnerName)
		if err != nil {
			t.Errorf("Error getting device by owner: %v", err)
			return
		}
		if deviceByOwner == nil || deviceByOwner.OwnerName != firstDevice.OwnerName {
			t.Errorf("Device owner does not match: expected %s, got %s", firstDevice.OwnerName, deviceByOwner.OwnerName)
		}
	})

	// Test GetDevicesByOSType
	t.Run("GetDevicesByOSType", func(t *testing.T) {
		deviceByOSType, err := GetDevicesByOSType(context.Background(), service, firstDevice.OSType)
		if err != nil {
			t.Errorf("Error getting device by OS type: %v", err)
			return
		}
		if deviceByOSType == nil || deviceByOSType.OSType != firstDevice.OSType {
			t.Errorf("Device OS type does not match: expected %s, got %s", firstDevice.OSType, deviceByOSType.OSType)
		}
	})

	// Test GetDevicesByOSVersion
	t.Run("GetDevicesByOSVersion", func(t *testing.T) {
		deviceByOSVersion, err := GetDevicesByOSVersion(context.Background(), service, firstDevice.OSVersion)
		if err != nil {
			t.Errorf("Error getting device by OS version: %v", err)
			return
		}
		if deviceByOSVersion == nil || deviceByOSVersion.OSVersion != firstDevice.OSVersion {
			t.Errorf("Device OS version does not match: expected %s, got %s", firstDevice.OSVersion, deviceByOSVersion.OSVersion)
		}
	})
}

func TestDeviceGroup_data(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	groups, err := GetAllDevicesGroups(context.Background(), service)
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
	server, err := GetDeviceGroupByName(context.Background(), service, name)
	if err != nil {
		t.Errorf("Error getting device group by name: %v", err)
		return
	}
	if server.Name != name {
		t.Errorf("device group name does not match: expected %s, got %s", name, server.Name)
		return
	}
	// Negative Test: Try to retrieve a device group with a non-existent name
	nonExistentName := "ThisDeviceGroupDoesNotExist"
	_, err = GetDeviceGroupByName(context.Background(), service, nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}

	// Test the GetIncludeDeviceInfo function with includeDeviceInfo=true and includePseudoGroups=false
	deviceInfos, err := GetIncludeDeviceInfo(context.Background(), service, true, false)
	if err != nil {
		t.Errorf("Error getting device info with includeDeviceInfo=true: %v", err)
		return
	}
	if len(deviceInfos) == 0 {
		t.Log("No device info found with includeDeviceInfo=true")
	} else {
		t.Logf("Found %d device infos with includeDeviceInfo=true", len(deviceInfos))
	}

	// Test the GetIncludeDeviceInfo function with includeDeviceInfo=false and includePseudoGroups=true
	deviceInfos, err = GetIncludeDeviceInfo(context.Background(), service, false, true)
	if err != nil {
		t.Errorf("Error getting device info with includePseudoGroups=true: %v", err)
		return
	}
	if len(deviceInfos) == 0 {
		t.Log("No device info found with includePseudoGroups=true")
	} else {
		t.Logf("Found %d device infos with includePseudoGroups=true", len(deviceInfos))
	}

	// Test the GetIncludeDeviceInfo function with both includeDeviceInfo and includePseudoGroups set to true
	deviceInfos, err = GetIncludeDeviceInfo(context.Background(), service, true, true)
	if err != nil {
		t.Errorf("Error getting device info with both includeDeviceInfo and includePseudoGroups=true: %v", err)
		return
	}
	if len(deviceInfos) == 0 {
		t.Log("No device info found with both includeDeviceInfo and includePseudoGroups=true")
	} else {
		t.Logf("Found %d device infos with both includeDeviceInfo and includePseudoGroups=true", len(deviceInfos))
	}

	// Test the GetIncludeDeviceInfo function with both includeDeviceInfo and includePseudoGroups set to false
	deviceInfos, err = GetIncludeDeviceInfo(context.Background(), service, false, false)
	if err != nil {
		t.Errorf("Error getting device info with both includeDeviceInfo and includePseudoGroups=false: %v", err)
		return
	}
	if len(deviceInfos) == 0 {
		t.Log("No device info found with both includeDeviceInfo and includePseudoGroups=false")
	} else {
		t.Logf("Found %d device infos with both includeDeviceInfo and includePseudoGroups=false", len(deviceInfos))
	}
}

func TestResponseFormatValidation(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	groups, err := GetAllDevicesGroups(context.Background(), service)
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
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

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
		group, err := GetDeviceGroupByName(context.Background(), service, variation)
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
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Retrieve all device groups
	groups, err := GetAllDevicesGroups(context.Background(), service)
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
