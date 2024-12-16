package inventory

import "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"

type GetSoftwareFilters struct {
	common.GetFromToFilters
	// The Zscaler location (ID). You can add multiple location IDs.
	Loc []int `json:"loc,omitempty" url:"loc,omitempty"`
	// The department (ID). You can add multiple department IDs.
	Dept []int `json:"dept,omitempty" url:"dept,omitempty"`
	// The active geolocation (ID). You can add multiple active geolocation IDs.
	Geo []string `json:"geo,omitempty" url:"geo,omitempty"`
	// User IDs filter
	UserIDs []int `json:"userids,omitempty" url:"userids,omitempty"`
	// Device IDs filter
	DeviceIDs []int `json:"deviceids,omitempty" url:"deviceids,omitempty"`
	// Software Key (required for softwareKeyEndpoint)
	SoftwareKey string `json:"software_key,omitempty" url:"software_key,omitempty"`
}
