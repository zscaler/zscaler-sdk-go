package endpoint_resource

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	dlpEndpointResourceEndpoint = "/zia/api/v1/dlpEndpointResource"
)

type EndpointResource struct {
	ID               int              `json:"id,omitempty"`
	Name             string           `json:"name,omitempty"`
	Channel          string           `json:"channel,omitempty"`
	IsPredefined     bool             `json:"isPredefined,omitempty"`
	NetworkDriveType string           `json:"networkDriveType,omitempty"`
	Description      string           `json:"description,omitempty"`
	ServerName       string           `json:"serverName,omitempty"`
	AppID            int              `json:"appId,omitempty"`
	NetworkDrives    []NetworkDrive   `json:"networkDrives,omitempty"`
	Printer          Printer          `json:"printer,omitempty"`
	RemovableStorage RemovableStorage `json:"removableStorage,omitempty"`
	Application      Application      `json:"application,omitempty"`
}

type NetworkDrive struct {
	NetworkPath string `json:"networkPath,omitempty"`
}

type Printer struct {
	Unc       string `json:"unc,omitempty"`
	IpAddress string `json:"ipAddress,omitempty"`
	Domain    string `json:"domain,omitempty"`
}

type RemovableStorage struct {
	VendorId     string `json:"vendorId,omitempty"`
	ProductId    string `json:"productId,omitempty"`
	SerialNumber string `json:"serialNumber,omitempty"`
}

type Application struct {
	OsType           string `json:"osType,omitempty"`
	FileName         string `json:"fileName,omitempty"`
	OriginalFileName string `json:"originalFileName,omitempty"`
	BundleID         string `json:"bundleID,omitempty"`
	DigitallySigned  bool   `json:"digitallySigned,omitempty"`
}

func Create(ctx context.Context, service *zscaler.Service, dlpRule *EndpointResource) (*EndpointResource, *http.Response, error) {
	resp, err := service.Client.Create(ctx, dlpEndpointResourceEndpoint, *dlpRule)
	if err != nil {
		return nil, nil, err
	}
	createdEndpointResource, ok := resp.(*EndpointResource)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a endpoint dlp rule pointer")
	}
	service.Client.GetLogger().Printf("[DEBUG]returning new endpoint dlp resource from create: %d", createdEndpointResource.ID)
	return createdEndpointResource, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleID int, dlpRule *EndpointResource) (*EndpointResource, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", dlpEndpointResourceEndpoint, ruleID), *dlpRule)
	if err != nil {
		return nil, nil, err
	}
	updatedOutboundEmail, _ := resp.(*EndpointResource)

	service.Client.GetLogger().Printf("[DEBUG]returning updates outbound email dlp policy from update: %d", updatedOutboundEmail.ID)
	return updatedOutboundEmail, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, profileID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", dlpEndpointResourceEndpoint, profileID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
