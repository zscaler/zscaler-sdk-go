package iotreport

import (
	"context"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	deviceTypesEndpoint     = "/zia/api/v1/iotDiscovery/deviceTypes"
	categoriesEndpoint      = "/zia/api/v1/iotDiscovery/categories"
	classificationsEndpoint = "/zia/api/v1/iotDiscovery/classifications"
	deviceListEndpoint      = "/zia/api/v1/iotDiscovery/deviceList"
)

type CommonIOTReport struct {
	// The universally unique identifier (UUID) of the item.
	UUID string `json:"uuid"`

	//The name of the item.
	Name string `json:"name,omitempty"`

	// The parent UUID of the item. This field is not applicable for the device types endpoint.
	ParentUuid string `json:"parent_uuid"`
}

type IOTDeviceList struct {
	// The Zscaler cloud on which your organization is provisioned.
	CloudName string `json:"cloudName"`

	// The unique identifier for your organization.
	CustomerID int `json:"customerId"`

	// A list of DeviceInfo objects.
	Devices []Device `json:"devices"`
}

type Device struct {
	// The location ID where the device resides.
	LocationID string `json:"locationId"`

	// The universally unique identifier (UUID) of the device identified by the Zscaler AI/ML.
	DeviceUUID string `json:"deviceUuid"`

	// The IP address of the device.
	IPAddress string `json:"ipAddress"`

	// The UUID of the device type identified by the Zscaler AI/Ml.
	DeviceTypeUUID string `json:"deviceTypeUuid"`

	// A label generated by the Zscaler AI/ML engine to describe the device's prominent characteristics.
	AutoLabel string `json:"autoLabel"`

	// The UUID for the device classification.
	ClassificationUUID string `json:"classificationUuid"`

	// The UUID for the device category.
	CategoryUUID string `json:"categoryUuid"`

	// The start timestamp at which Zscaler AI/ML engine starts evaluating the device's weblog records. It's noted in epoch seconds.
	FlowStartTime int `json:"flowStartTime"`

	// The end timestamp at which Zscaler AI/ML engine stops evaluating the device's weblog records. It's noted in epoch seconds.
	FlowEndTime int `json:"flowEndTime"`
}

// Retrieve the mapping between device type universally unique identifier (UUID) values and the device type names for all the device types supported by the Zscaler AI/ML.
func GetDeviceTypes(ctx context.Context, service *zscaler.Service) (*CommonIOTReport, error) {
	var deviceTypes CommonIOTReport
	err := service.Client.Read(ctx, deviceTypesEndpoint, &deviceTypes)
	if err != nil {
		return nil, err
	}
	return &deviceTypes, nil
}

// Retrieve the mapping between the device category universally unique identifier (UUID) values and the category names for all the device categories supported by the Zscaler AI/ML. The parent of device category is device type.
func GetIOTCategories(ctx context.Context, service *zscaler.Service) (*CommonIOTReport, error) {
	var categories CommonIOTReport
	err := service.Client.Read(ctx, categoriesEndpoint, &categories)
	if err != nil {
		return nil, err
	}
	return &categories, nil
}

// Retrieve the mapping between the device classification universally unique identifier (UUID) values and the classification names for all the device classifications supported by Zscaler AI/ML. The parent of device classification is device category.
func GetIOTClassifications(ctx context.Context, service *zscaler.Service) (*CommonIOTReport, error) {
	var categories CommonIOTReport
	err := service.Client.Read(ctx, classificationsEndpoint, &categories)
	if err != nil {
		return nil, err
	}
	return &categories, nil
}

// Retrieve a list of discovered devices with the following key contexts, IP address, location, ML auto-label, classification, category, and type.
func GetIOTDeviceList(ctx context.Context, service *zscaler.Service) (*IOTDeviceList, error) {
	var deviceList IOTDeviceList
	err := service.Client.Read(ctx, deviceListEndpoint, &deviceList)
	if err != nil {
		return nil, err
	}
	return &deviceList, nil
}
