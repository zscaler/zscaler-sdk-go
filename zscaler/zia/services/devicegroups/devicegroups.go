package devicegroups

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	deviceGroupEndpoint = "/zia/api/v1/deviceGroups"
	devicesEndpoint     = "/zia/api/v1/deviceGroups/devices"
	devicesLiteEndpoint = "/zia/api/v1/deviceGroups/devices/lite"
)

// DeviceGroups represents a device group returned by /deviceGroups endpoint
type DeviceGroups struct {
	// The unique identifier for the device group
	ID int `json:"id"`

	// The device group name
	Name string `json:"name,omitempty"`

	// The device group type
	GroupType string `json:"groupType,omitempty"`

	// The device group's description
	Description string `json:"description,omitempty"`

	// The operating system (OS)
	OSType string `json:"osType,omitempty"`

	// Indicates whether this is a predefined device group. If this value is set to true, the group is predefined
	Predefined bool `json:"predefined"`

	DeviceNames string `json:"deviceNames,omitempty"`

	// The number of devices within the group
	DeviceCount int `json:"deviceCount,omitempty"`
}

// Devices represents a device returned by /deviceGroups/devices endpoint
type Devices struct {
	// The unique identifier for the device
	ID int `json:"id"`

	// The device name
	Name string `json:"name,omitempty"`

	// The device group type
	DeviceGroupType string `json:"deviceGroupType,omitempty"`

	// The device model
	DeviceModel string `json:"deviceModel,omitempty"`

	// The operating system (OS)
	OSType string `json:"osType,omitempty"`

	// The operating system version
	OSVersion string `json:"osVersion,omitempty"`

	// The device's description
	Description string `json:"description,omitempty"`

	// The unique identifier of the device owner (i.e., user)
	OwnerUserId int `json:"ownerUserId,omitempty"`

	// The device owner's user name
	OwnerName string `json:"ownerName,omitempty"`

	// The hostname of the device
	HostName string `json:"hostName,omitempty"`
}

// =============================================================================
// Endpoint 1: /deviceGroups - Gets a list of device groups
// Note: This endpoint does NOT support pagination
// =============================================================================

// GetAllDeviceGroupsOptions contains optional parameters for GetAllDeviceGroups
type GetAllDeviceGroupsOptions struct {
	// Include or exclude device information
	IncludeDeviceInfo bool
	// Include or exclude Zscaler Client Connector and Cloud Browser Isolation-related device groups
	IncludePseudoGroups bool
	// Include or exclude IoT device groups
	IncludeIOTGroups bool
}

// GetAllDeviceGroups retrieves all device groups.
// Note: This endpoint does NOT support pagination.
func GetAllDeviceGroups(ctx context.Context, service *zscaler.Service, opts *GetAllDeviceGroupsOptions) ([]DeviceGroups, error) {
	queryParams := url.Values{}

	if opts != nil {
		if opts.IncludeDeviceInfo {
			queryParams.Set("includeDeviceInfo", "true")
		}
		if opts.IncludePseudoGroups {
			queryParams.Set("includePseudoGroups", "true")
		}
		if opts.IncludeIOTGroups {
			queryParams.Set("includeIOTGroups", "true")
		}
	}

	endpoint := deviceGroupEndpoint
	if len(queryParams) > 0 {
		endpoint = fmt.Sprintf("%s?%s", deviceGroupEndpoint, queryParams.Encode())
	}

	var deviceGroups []DeviceGroups
	err := service.Client.Read(ctx, endpoint, &deviceGroups)
	return deviceGroups, err
}

// GetDeviceGroupByName retrieves a device group by name.
func GetDeviceGroupByName(ctx context.Context, service *zscaler.Service, deviceGroupName string) (*DeviceGroups, error) {
	deviceGroups, err := GetAllDeviceGroups(ctx, service, nil)
	if err != nil {
		return nil, err
	}
	for _, deviceGroup := range deviceGroups {
		if strings.EqualFold(deviceGroup.Name, deviceGroupName) {
			return &deviceGroup, nil
		}
	}
	return nil, fmt.Errorf("no device group found with name: %s", deviceGroupName)
}

// =============================================================================
// Endpoint 2: /deviceGroups/devices - Gets a list of devices
// This endpoint supports pagination (page, pageSize up to 1000)
// =============================================================================

// GetAllDevicesOptions contains optional parameters for GetAllDevices
type GetAllDevicesOptions struct {
	// The device group name
	Name string
	// The device models
	Model string
	// The device owners
	Owner string
	// The device's operating system (IOS, ANDROID_SAMSUNG, WINDOWS, MAC, LINUX, ANDROID_NON_SAMSUNG, ANY)
	OSType string
	// The device's operating system version
	OSVersion string
	// The unique identifier for the device group
	DeviceGroupID int
	// Used to list devices for specific users (array of user IDs)
	UserIDs []int
	// Used to match against all device attribute information
	SearchAll string
	// Used to include or exclude Cloud Browser Isolation devices
	IncludeAll bool
}

// GetAllDevices retrieves all devices with optional filtering.
// This endpoint supports pagination.
func GetAllDevices(ctx context.Context, service *zscaler.Service, opts *GetAllDevicesOptions) ([]Devices, error) {
	queryParams := url.Values{}

	if opts != nil {
		if opts.Name != "" {
			queryParams.Set("name", opts.Name)
		}
		if opts.Model != "" {
			queryParams.Set("model", opts.Model)
		}
		if opts.Owner != "" {
			queryParams.Set("owner", opts.Owner)
		}
		if opts.OSType != "" {
			queryParams.Set("osType", opts.OSType)
		}
		if opts.OSVersion != "" {
			queryParams.Set("osVersion", opts.OSVersion)
		}
		if opts.DeviceGroupID != 0 {
			queryParams.Set("deviceGroupId", fmt.Sprintf("%d", opts.DeviceGroupID))
		}
		if len(opts.UserIDs) > 0 {
			for _, userID := range opts.UserIDs {
				queryParams.Add("userIds", fmt.Sprintf("%d", userID))
			}
		}
		if opts.SearchAll != "" {
			queryParams.Set("search_all", opts.SearchAll)
		}
		if opts.IncludeAll {
			queryParams.Set("includeAll", "true")
		}
	}

	endpoint := devicesEndpoint
	if len(queryParams) > 0 {
		endpoint = fmt.Sprintf("%s?%s", devicesEndpoint, queryParams.Encode())
	}

	var devices []Devices
	err := common.ReadAllPages(ctx, service.Client, endpoint, &devices)
	return devices, err
}

// GetDevicesByID retrieves a device by ID.
func GetDevicesByID(ctx context.Context, service *zscaler.Service, deviceID int) (*Devices, error) {
	devices, err := GetAllDevices(ctx, service, nil)
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		if device.ID == deviceID {
			return &device, nil
		}
	}
	return nil, fmt.Errorf("no device found with ID: %d", deviceID)
}

// GetDevicesByName retrieves a device by name.
func GetDevicesByName(ctx context.Context, service *zscaler.Service, deviceName string) (*Devices, error) {
	// Use the API's name filter for efficiency
	devices, err := GetAllDevices(ctx, service, &GetAllDevicesOptions{Name: deviceName})
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		if strings.EqualFold(device.Name, deviceName) {
			return &device, nil
		}
	}
	return nil, fmt.Errorf("no device found with name: %s", deviceName)
}

// GetDevicesByModel retrieves a device by model.
func GetDevicesByModel(ctx context.Context, service *zscaler.Service, deviceModel string) (*Devices, error) {
	devices, err := GetAllDevices(ctx, service, &GetAllDevicesOptions{Model: deviceModel})
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		if strings.EqualFold(device.DeviceModel, deviceModel) {
			return &device, nil
		}
	}
	return nil, fmt.Errorf("no device found with model: %s", deviceModel)
}

// GetDevicesByOwner retrieves a device by owner name.
func GetDevicesByOwner(ctx context.Context, service *zscaler.Service, ownerName string) (*Devices, error) {
	devices, err := GetAllDevices(ctx, service, &GetAllDevicesOptions{Owner: ownerName})
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		if strings.EqualFold(device.OwnerName, ownerName) {
			return &device, nil
		}
	}
	return nil, fmt.Errorf("no device found for owner: %s", ownerName)
}

// GetDevicesByOSType retrieves a device by OS type.
func GetDevicesByOSType(ctx context.Context, service *zscaler.Service, osTypeName string) (*Devices, error) {
	devices, err := GetAllDevices(ctx, service, &GetAllDevicesOptions{OSType: osTypeName})
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		if strings.EqualFold(device.OSType, osTypeName) {
			return &device, nil
		}
	}
	return nil, fmt.Errorf("no device found for type: %s", osTypeName)
}

// GetDevicesByOSVersion retrieves a device by OS version.
func GetDevicesByOSVersion(ctx context.Context, service *zscaler.Service, osVersionName string) (*Devices, error) {
	devices, err := GetAllDevices(ctx, service, &GetAllDevicesOptions{OSVersion: osVersionName})
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		if strings.EqualFold(device.OSVersion, osVersionName) {
			return &device, nil
		}
	}
	return nil, fmt.Errorf("no device found for version: %s", osVersionName)
}

// =============================================================================
// Endpoint 3: /deviceGroups/devices/lite - Gets a lite list of devices
// This endpoint supports pagination (page, pageSize up to 1000)
// =============================================================================

// GetAllDevicesLiteOptions contains optional parameters for GetAllDevicesLite
type GetAllDevicesLiteOptions struct {
	// The device group name
	Name string
	// Used to list devices for specific users (array of user IDs)
	UserIDs []int
	// Used to include or exclude Cloud Browser Isolation devices
	IncludeAll bool
}

// GetAllDevicesLite retrieves all devices (lite version) with optional filtering.
// This endpoint supports pagination.
func GetAllDevicesLite(ctx context.Context, service *zscaler.Service, opts *GetAllDevicesLiteOptions) ([]Devices, error) {
	queryParams := url.Values{}

	if opts != nil {
		if opts.Name != "" {
			queryParams.Set("name", opts.Name)
		}
		if len(opts.UserIDs) > 0 {
			for _, userID := range opts.UserIDs {
				queryParams.Add("userIds", fmt.Sprintf("%d", userID))
			}
		}
		if opts.IncludeAll {
			queryParams.Set("includeAll", "true")
		}
	}

	endpoint := devicesLiteEndpoint
	if len(queryParams) > 0 {
		endpoint = fmt.Sprintf("%s?%s", devicesLiteEndpoint, queryParams.Encode())
	}

	var devices []Devices
	err := common.ReadAllPages(ctx, service.Client, endpoint, &devices)
	return devices, err
}

// =============================================================================
// Deprecated functions - kept for backward compatibility
// =============================================================================

// Deprecated: Use GetAllDeviceGroups instead.
func GetAllDevicesGroups(ctx context.Context, service *zscaler.Service) ([]DeviceGroups, error) {
	return GetAllDeviceGroups(ctx, service, nil)
}

// Deprecated: Use GetAllDeviceGroups with GetAllDeviceGroupsOptions instead.
func GetIncludeDeviceInfo(ctx context.Context, service *zscaler.Service, includeDeviceInfo, includePseudoGroups bool) ([]DeviceGroups, error) {
	return GetAllDeviceGroups(ctx, service, &GetAllDeviceGroupsOptions{
		IncludeDeviceInfo:   includeDeviceInfo,
		IncludePseudoGroups: includePseudoGroups,
	})
}
