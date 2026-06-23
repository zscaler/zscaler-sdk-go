package devices

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	devicesEndpoint = "/zia/api/v1/devices"
)

type Devices struct {
	ID                         int                      `json:"id,omitempty"`
	Name                       string                   `json:"name,omitempty"`
	Active                     bool                     `json:"active,omitempty"`
	Version                    string                   `json:"version,omitempty"`
	Hostname                   string                   `json:"hostname,omitempty"`
	Vendor                     string                   `json:"vendor,omitempty"`
	Model                      string                   `json:"model,omitempty"`
	Locale                     string                   `json:"locale,omitempty"`
	Os                         string                   `json:"os,omitempty"`
	Udid                       string                   `json:"udid,omitempty"`
	HardwareId                 string                   `json:"hardwareId,omitempty"`
	MacAddress                 string                   `json:"macAddress,omitempty"`
	User                       *common.IDNameExtensions `json:"user,omitempty"`
	FirstRegistrationTimestamp int                      `json:"firstRegistrationTimestamp,omitempty"`
	LastRegistrationTimestamp  int                      `json:"lastRegistrationTimestamp,omitempty"`
	UnRegistrationTimestamp    int                      `json:"unRegistrationTimestamp,omitempty"`
	Deleted                    bool                     `json:"deleted,omitempty"`
	Rooted                     bool                     `json:"rooted,omitempty"`
}

type GetAllFilterOptions struct {
	// Filters the list based on the device IDs.
	ID []int

	// Search string to filter devices.
	Search *string

	// Filters the list for valid devices.
	Valid *bool

	// Filters to include CBI devices in the search result.
	IncludeCbiDevices *bool
}

func Get(ctx context.Context, service *zscaler.Service, deviceID int) (*Devices, error) {
	var device Devices
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", devicesEndpoint, deviceID), &device)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning device from Get: %d", device.ID)
	return &device, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, deviceName string) (*Devices, error) {
	// Use GetAll to leverage the single API call and built-in pagination
	devices, err := GetAll(ctx, service, nil)
	if err != nil {
		return nil, err
	}
	// Search for exact match (case-insensitive)
	for _, device := range devices {
		if strings.EqualFold(device.Name, deviceName) {
			return &device, nil
		}
	}
	return nil, fmt.Errorf("no device found with name: %s", deviceName)
}

// GetAll retrieves a list of all the devices registered with Zscaler.
//
// The optional filtering parameters (id, search, valid, includeCbiDevices)
// are appended as query parameters. The page and pageSize parameters are
// handled internally by common.ReadAllPages, which paginates through all
// results, so they should not be set here.
func GetAll(ctx context.Context, service *zscaler.Service, opts *GetAllFilterOptions) ([]Devices, error) {
	var devices []Devices
	endpoint := devicesEndpoint

	queryParams := url.Values{}
	if opts != nil {
		for _, id := range opts.ID {
			queryParams.Add("id", strconv.Itoa(id))
		}
		if opts.Search != nil && *opts.Search != "" {
			queryParams.Set("search", *opts.Search)
		}
		if opts.Valid != nil {
			queryParams.Set("valid", strconv.FormatBool(*opts.Valid))
		}
		if opts.IncludeCbiDevices != nil {
			queryParams.Set("includeCbiDevices", strconv.FormatBool(*opts.IncludeCbiDevices))
		}
	}
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	err := common.ReadAllPages(ctx, service.Client, endpoint, &devices)
	return devices, err
}
