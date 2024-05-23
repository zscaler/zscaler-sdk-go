package devicegroups

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	deviceGroupEndpoint = "/deviceGroups"
	devicesEndpoint     = "/deviceGroups/devices"
)

type DeviceGroups struct {
	// The unique identifer for the device group
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
	Predefined  bool   `json:"predefined"`
	DeviceNames string `json:"deviceNames,omitempty"`

	// The number of devices within the group
	DeviceCount int `json:"deviceCount,omitempty"`
}

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

func (service *Service) GetDeviceGroupByName(deviceGroupName string) (*DeviceGroups, error) {
	var deviceGroups []DeviceGroups
	err := service.Client.Read(deviceGroupEndpoint, &deviceGroups)
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

func (service *Service) GetIncludeDeviceInfo(includeDeviceInfo, includePseudoGroups bool) ([]DeviceGroups, error) {
	queryParams := url.Values{}
	if includeDeviceInfo {
		queryParams.Set("includeDeviceInfo", "true")
	}
	if includePseudoGroups {
		queryParams.Set("includePseudoGroups", "true")
	}

	endpoint := fmt.Sprintf("%s?%s", deviceGroupEndpoint, queryParams.Encode())
	var deviceInfos []DeviceGroups
	err := service.Client.Read(endpoint, &deviceInfos)
	if err != nil {
		return nil, err
	}
	return deviceInfos, nil
}

func (service *Service) GetAllDevicesGroups() ([]DeviceGroups, error) {
	var owners []DeviceGroups
	err := common.ReadAllPages(service.Client, deviceGroupEndpoint, &owners)
	return owners, err
}

func (service *Service) GetDevicesByID(deviceId int) (*Devices, error) {
	devices, err := service.GetAllDevices()
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		if device.ID == deviceId {
			return &device, nil
		}
	}

	return nil, fmt.Errorf("no device found with ID: %d", deviceId)
}

// Get Devices by Name.
func (service *Service) GetDevicesByName(deviceName string) (*Devices, error) {
	var devices []Devices
	// We are assuming this device name will be in the firsy 1000 obejcts
	err := common.ReadAllPages(service.Client, fmt.Sprintf("%s?page=1&pageSize=1000", devicesEndpoint), &devices)
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

func (service *Service) GetDevicesByModel(deviceModel string) (*Devices, error) {
	var models []Devices
	err := common.ReadAllPages(service.Client, fmt.Sprintf("%s?model=%s", devicesEndpoint, url.QueryEscape(deviceModel)), &models)
	if err != nil {
		return nil, err
	}
	for _, model := range models {
		if strings.EqualFold(model.DeviceModel, deviceModel) {
			return &model, nil
		}
	}
	return nil, fmt.Errorf("no device found with model: %s", deviceModel)
}

func (service *Service) GetDevicesByOwner(ownerName string) (*Devices, error) {
	var owners []Devices
	err := common.ReadAllPages(service.Client, fmt.Sprintf("%s?owner=%s", devicesEndpoint, url.QueryEscape(ownerName)), &owners)
	if err != nil {
		return nil, err
	}
	for _, owner := range owners {
		if strings.EqualFold(owner.OwnerName, ownerName) {
			return &owner, nil
		}
	}
	return nil, fmt.Errorf("no device found for owner: %s", ownerName)
}

func (service *Service) GetDevicesByOSType(osTypeName string) (*Devices, error) {
	var osTypes []Devices
	err := common.ReadAllPages(service.Client, fmt.Sprintf("%s?osType=%s", devicesEndpoint, url.QueryEscape(osTypeName)), &osTypes)
	if err != nil {
		return nil, err
	}
	for _, osType := range osTypes {
		if strings.EqualFold(osType.OSType, osTypeName) {
			return &osType, nil
		}
	}
	return nil, fmt.Errorf("no device found for type: %s", osTypeName)
}

func (service *Service) GetDevicesByOSVersion(osVersionName string) (*Devices, error) {
	var osVersions []Devices
	err := common.ReadAllPages(service.Client, fmt.Sprintf("%s?osVersion=%s", devicesEndpoint, url.QueryEscape(osVersionName)), &osVersions)
	if err != nil {
		return nil, err
	}
	for _, osVersion := range osVersions {
		if strings.EqualFold(osVersion.OSVersion, osVersionName) {
			return &osVersion, nil
		}
	}
	return nil, fmt.Errorf("no device found for version: %s", osVersionName)
}

func (service *Service) GetAllDevices() ([]Devices, error) {
	var owners []Devices
	err := common.ReadAllPages(service.Client, devicesEndpoint, &owners)
	return owners, err
}
