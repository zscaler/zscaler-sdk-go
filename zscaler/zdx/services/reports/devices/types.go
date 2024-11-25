package devices

import "github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zdx/services/common"

type GetDevicesFilters struct {
	common.GetFromToFilters
	// The user IDs.
	UserIDs []int `json:"userids,omitempty" url:"userids,omitempty"`
	// Emails
	Emails []string `json:"emails,omitempty" url:"emails,omitempty"`
	// The Zscaler location (ID). You can add multiple location IDs.
	Loc []int `json:"loc,omitempty" url:"loc,omitempty"`
	// The department (ID). You can add multiple department IDs.
	Dept []int `json:"dept,omitempty" url:"dept,omitempty"`
	// The active geolocation (ID). You can add multiple active geolocation IDs.
	Geo []string `json:"geo,omitempty" url:"geo,omitempty"`
	// The next_offset value from the last request. You must enter this value to get the next batch from the list. When the next_offset value becomes null, the list is complete.
	Offset string `json:"offset,omitempty" url:"offset,omitempty"`
	// The number of items that must be returned per request from the list.
	Limit int `json:"limit,omitempty" url:"limit,omitempty"`
}

type GeoLocationFilter struct {
	common.GetFromToFilters
	// The parent geo ID.
	ParentGeoID string `json:"parent_geo_id,omitempty" url:"parent_geo_id,omitempty"`
	// The search string used to support search by name.
	Search string `json:"search,omitempty" url:"search,omitempty"`
}
