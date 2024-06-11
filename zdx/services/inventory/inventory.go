package inventory

import "net/http"

const (
	softwareEndpoint    = "v1/inventory/software"
	softwareKeyEndpoint = "v1/inventory/softwares"
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
	SoftwareInstallType string `json:"software_install_type,omitempty"`
	UserTotal           string `json:"user_total,omitempty"`
	DeviceTotal         string `json:"device_total,omitempty"`
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

func (service *Service) GetSoftware(filters GetSoftwareFilters) ([]SoftwareOverview, string, *http.Response, error) {
	var response SoftwareOverviewResponse
	path := softwareEndpoint
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, &response)
	if err != nil {
		return nil, "", nil, err
	}
	return response.Software, response.NextOffset, resp, nil
}

func (service *Service) GetSoftwareKey(filters GetSoftwareFilters) ([]SoftwareUserList, string, *http.Response, error) {
	var response SoftwareKeyResponse
	path := softwareKeyEndpoint
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, &response)
	if err != nil {
		return nil, "", nil, err
	}
	return response.Software, response.NextOffset, resp, nil
}
