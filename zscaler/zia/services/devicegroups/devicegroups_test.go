package devicegroups

import (
	"context"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// =============================================================================
// Tests for /deviceGroups endpoint (GetAllDeviceGroups)
// =============================================================================

func TestGetAllDeviceGroups(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "devicegroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// Test GetAllDeviceGroups without options
	t.Run("WithoutOptions", func(t *testing.T) {
		groups, err := GetAllDeviceGroups(context.Background(), service, nil)
		if err != nil {
			t.Errorf("Error getting device groups: %v", err)
			return
		}
		if len(groups) == 0 {
			t.Log("No device groups found")
			return
		}
		t.Logf("Found %d device groups", len(groups))
	})

	// Test GetAllDeviceGroups with includeDeviceInfo=true
	t.Run("WithIncludeDeviceInfo", func(t *testing.T) {
		groups, err := GetAllDeviceGroups(context.Background(), service, &GetAllDeviceGroupsOptions{
			IncludeDeviceInfo: true,
		})
		if err != nil {
			t.Errorf("Error getting device groups with includeDeviceInfo=true: %v", err)
			return
		}
		t.Logf("Found %d device groups with includeDeviceInfo=true", len(groups))
	})

	// Test GetAllDeviceGroups with includePseudoGroups=true
	t.Run("WithIncludePseudoGroups", func(t *testing.T) {
		groups, err := GetAllDeviceGroups(context.Background(), service, &GetAllDeviceGroupsOptions{
			IncludePseudoGroups: true,
		})
		if err != nil {
			t.Errorf("Error getting device groups with includePseudoGroups=true: %v", err)
			return
		}
		t.Logf("Found %d device groups with includePseudoGroups=true", len(groups))
	})

	// Note: includeIOTGroups requires ZT_IOT_STD SKU & ZT_IOT_VIS SKU subscription
	// Skipping test for includeIOTGroups to avoid NOT_SUBSCRIBED errors

	// Test GetAllDeviceGroups with includeDeviceInfo and includePseudoGroups
	t.Run("WithDeviceInfoAndPseudoGroups", func(t *testing.T) {
		groups, err := GetAllDeviceGroups(context.Background(), service, &GetAllDeviceGroupsOptions{
			IncludeDeviceInfo:   true,
			IncludePseudoGroups: true,
		})
		if err != nil {
			t.Errorf("Error getting device groups with includeDeviceInfo and includePseudoGroups: %v", err)
			return
		}
		t.Logf("Found %d device groups with includeDeviceInfo and includePseudoGroups", len(groups))
	})
}

func TestGetDeviceGroupByName(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "devicegroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// First get all groups to find a valid name
	groups, err := GetAllDeviceGroups(context.Background(), service, nil)
	if err != nil {
		t.Fatalf("Error getting device groups: %v", err)
	}
	if len(groups) == 0 {
		t.Skip("No device groups found for testing")
	}

	name := groups[0].Name
	t.Logf("Getting device group by name: %s", name)

	group, err := GetDeviceGroupByName(context.Background(), service, name)
	if err != nil {
		t.Errorf("Error getting device group by name: %v", err)
		return
	}
	if group.Name != name {
		t.Errorf("Device group name does not match: expected %s, got %s", name, group.Name)
	}

	// Negative Test: Try to retrieve a device group with a non-existent name
	nonExistentName := "ThisDeviceGroupDoesNotExist"
	_, err = GetDeviceGroupByName(context.Background(), service, nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
	}
}

func TestDeviceGroupByNameCaseSensitivity(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "devicegroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

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
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "devicegroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// Retrieve all device groups
	groups, err := GetAllDeviceGroups(context.Background(), service, nil)
	if err != nil {
		t.Fatalf("Error getting all device groups: %v", err)
	}
	if len(groups) == 0 {
		t.Skip("No device groups found for testing")
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

// =============================================================================
// Tests for /deviceGroups/devices endpoint (GetAllDevices)
// =============================================================================

func TestGetAllDevices(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "devicegroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// Test GetAllDevices without options
	t.Run("WithoutOptions", func(t *testing.T) {
		devices, err := GetAllDevices(context.Background(), service, nil)
		if err != nil {
			t.Errorf("Error getting devices: %v", err)
			return
		}
		if len(devices) == 0 {
			t.Log("No devices found")
			return
		}
		t.Logf("Found %d devices", len(devices))
	})
}

func TestGetDevicesByID(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "devicegroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// First get all devices to find a valid ID
	devices, err := GetAllDevices(context.Background(), service, nil)
	if err != nil {
		t.Fatalf("Error getting devices: %v", err)
	}
	if len(devices) == 0 {
		t.Skip("No devices found for testing")
	}

	firstDevice := devices[0]
	deviceByID, err := GetDevicesByID(context.Background(), service, firstDevice.ID)
	if err != nil {
		t.Errorf("Error getting device by ID: %v", err)
		return
	}
	if deviceByID == nil || deviceByID.ID != firstDevice.ID {
		t.Errorf("Device ID does not match: expected %d, got %d", firstDevice.ID, deviceByID.ID)
	}
}

func TestGetDevicesByName(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "devicegroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// First get all devices to find a valid name
	devices, err := GetAllDevices(context.Background(), service, nil)
	if err != nil {
		t.Fatalf("Error getting devices: %v", err)
	}
	if len(devices) == 0 {
		t.Skip("No devices found for testing")
	}

	firstDevice := devices[0]
	deviceByName, err := GetDevicesByName(context.Background(), service, firstDevice.Name)
	if err != nil {
		t.Errorf("Error getting device by name: %v", err)
		return
	}
	if deviceByName == nil || deviceByName.Name != firstDevice.Name {
		t.Errorf("Device name does not match: expected %s, got %s", firstDevice.Name, deviceByName.Name)
	}
}

func TestGetDevicesByModel(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "devicegroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// First get all devices to find a valid model
	devices, err := GetAllDevices(context.Background(), service, nil)
	if err != nil {
		t.Fatalf("Error getting devices: %v", err)
	}
	if len(devices) == 0 {
		t.Skip("No devices found for testing")
	}

	firstDevice := devices[0]
	if firstDevice.DeviceModel == "" {
		t.Skip("First device has no model for testing")
	}

	deviceByModel, err := GetDevicesByModel(context.Background(), service, firstDevice.DeviceModel)
	if err != nil {
		t.Errorf("Error getting device by model: %v", err)
		return
	}
	if deviceByModel == nil || deviceByModel.DeviceModel != firstDevice.DeviceModel {
		t.Errorf("Device model does not match: expected %s, got %s", firstDevice.DeviceModel, deviceByModel.DeviceModel)
	}
}

func TestGetDevicesByOwner(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "devicegroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// First get all devices to find a valid owner
	devices, err := GetAllDevices(context.Background(), service, nil)
	if err != nil {
		t.Fatalf("Error getting devices: %v", err)
	}
	if len(devices) == 0 {
		t.Skip("No devices found for testing")
	}

	firstDevice := devices[0]
	if firstDevice.OwnerName == "" {
		t.Skip("First device has no owner for testing")
	}

	deviceByOwner, err := GetDevicesByOwner(context.Background(), service, firstDevice.OwnerName)
	if err != nil {
		t.Errorf("Error getting device by owner: %v", err)
		return
	}
	if deviceByOwner == nil || deviceByOwner.OwnerName != firstDevice.OwnerName {
		t.Errorf("Device owner does not match: expected %s, got %s", firstDevice.OwnerName, deviceByOwner.OwnerName)
	}
}

func TestGetDevicesByOSType(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "devicegroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// First get all devices to find a valid OS type
	devices, err := GetAllDevices(context.Background(), service, nil)
	if err != nil {
		t.Fatalf("Error getting devices: %v", err)
	}
	if len(devices) == 0 {
		t.Skip("No devices found for testing")
	}

	firstDevice := devices[0]
	if firstDevice.OSType == "" {
		t.Skip("First device has no OS type for testing")
	}

	deviceByOSType, err := GetDevicesByOSType(context.Background(), service, firstDevice.OSType)
	if err != nil {
		t.Errorf("Error getting device by OS type: %v", err)
		return
	}
	if deviceByOSType == nil || deviceByOSType.OSType != firstDevice.OSType {
		t.Errorf("Device OS type does not match: expected %s, got %s", firstDevice.OSType, deviceByOSType.OSType)
	}
}

func TestGetDevicesByOSVersion(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "devicegroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// First get all devices to find a valid OS version
	devices, err := GetAllDevices(context.Background(), service, nil)
	if err != nil {
		t.Fatalf("Error getting devices: %v", err)
	}
	if len(devices) == 0 {
		t.Skip("No devices found for testing")
	}

	firstDevice := devices[0]
	if firstDevice.OSVersion == "" {
		t.Skip("First device has no OS version for testing")
	}

	deviceByOSVersion, err := GetDevicesByOSVersion(context.Background(), service, firstDevice.OSVersion)
	if err != nil {
		t.Errorf("Error getting device by OS version: %v", err)
		return
	}
	if deviceByOSVersion == nil || deviceByOSVersion.OSVersion != firstDevice.OSVersion {
		t.Errorf("Device OS version does not match: expected %s, got %s", firstDevice.OSVersion, deviceByOSVersion.OSVersion)
	}
}

// =============================================================================
// Tests for /deviceGroups/devices/lite endpoint (GetAllDevicesLite)
// =============================================================================

func TestGetAllDevicesLite(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "devicegroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// Test GetAllDevicesLite without options
	t.Run("WithoutOptions", func(t *testing.T) {
		devices, err := GetAllDevicesLite(context.Background(), service, nil)
		if err != nil {
			t.Errorf("Error getting devices lite: %v", err)
			return
		}
		if len(devices) == 0 {
			t.Log("No devices found in lite endpoint")
			return
		}
		t.Logf("Found %d devices in lite endpoint", len(devices))
	})

	// Test GetAllDevicesLite with includeAll=true
	t.Run("WithIncludeAll", func(t *testing.T) {
		devices, err := GetAllDevicesLite(context.Background(), service, &GetAllDevicesLiteOptions{
			IncludeAll: true,
		})
		if err != nil {
			t.Errorf("Error getting devices lite with includeAll=true: %v", err)
			return
		}
		t.Logf("Found %d devices in lite endpoint with includeAll=true", len(devices))
	})
}

// =============================================================================
// Tests for deprecated functions (backward compatibility)
// =============================================================================

func TestDeprecatedGetAllDevicesGroups(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "devicegroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// Test deprecated GetAllDevicesGroups function still works
	groups, err := GetAllDevicesGroups(context.Background(), service)
	if err != nil {
		t.Errorf("Error calling deprecated GetAllDevicesGroups: %v", err)
		return
	}
	t.Logf("Deprecated GetAllDevicesGroups returned %d groups", len(groups))
}

func TestDeprecatedGetIncludeDeviceInfo(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "devicegroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// Test deprecated GetIncludeDeviceInfo function still works
	groups, err := GetIncludeDeviceInfo(context.Background(), service, true, true)
	if err != nil {
		t.Errorf("Error calling deprecated GetIncludeDeviceInfo: %v", err)
		return
	}
	t.Logf("Deprecated GetIncludeDeviceInfo returned %d groups", len(groups))
}
