package inventory

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services"
)

const (
	softwareEndpoint    = "v1/inventory/software"
	softwareKeyEndpoint = "v1/inventory/software"
)

type SoftwareOverviewResponse struct {
	Software   []SoftwareOverview `json:"software"`
	NextOffset string             `json:"next_offset,omitempty"`
}

type SoftwareKeyResponse struct {
	Software   []SoftwareUserList `json:"software"`
	NextOffset string             `json:"next_offset,omitempty"`
}

type SoftwareOverview struct {
	SoftwareKey         string `json:"software_key,omitempty"`
	SoftwareName        string `json:"software_name,omitempty"`
	Vendor              string `json:"vendor,omitempty"`
	SoftwareGroup       string `json:"software_group,omitempty"`
	SoftwareInstallType string `json:"sw_install_type,omitempty"`
	UserTotal           int    `json:"user_total,omitempty"`
	DeviceTotal         int    `json:"device_total,omitempty"`
}

type SoftwareUserList struct {
	SoftwareKey     string `json:"software_key,omitempty"`
	SoftwareName    string `json:"software_name,omitempty"`
	SoftwareVersion string `json:"software_version,omitempty"`
	SoftwareGroup   string `json:"software_group,omitempty"`
	OS              string `json:"os,omitempty"`
	Vendor          string `json:"vendor,omitempty"`
	UserID          int    `json:"user_id,omitempty"`
	DeviceID        int    `json:"device_id,omitempty"`
	Hostname        string `json:"hostname,omitempty"`
	Username        string `json:"username,omitempty"`
	InstallDate     string `json:"install_date,omitempty"`
}

func GetSoftware(ctx context.Context, service *services.Service, filters GetSoftwareFilters) ([]SoftwareOverview, string, *http.Response, error) {
	var response SoftwareOverviewResponse
	path := softwareEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "GET", path, filters, nil, &response)
	if err != nil {
		return nil, "", nil, err
	}
	return response.Software, response.NextOffset, resp, nil
}

func GetSoftwareKey(ctx context.Context, service *services.Service, softwareKey string, filters GetSoftwareFilters) ([]SoftwareUserList, string, *http.Response, error) {
	var response SoftwareKeyResponse
	path := fmt.Sprintf("%v/%v", softwareKeyEndpoint, softwareKey)
	resp, err := service.Client.NewRequestDo(ctx, "GET", path, filters, nil, &response)
	if err != nil {
		return nil, "", nil, err
	}
	return response.Software, response.NextOffset, resp, nil
}
