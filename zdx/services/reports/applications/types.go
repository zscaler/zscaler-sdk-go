package applications

import "github.com/zscaler/zscaler-sdk-go/v2/zdx/services/common"

type GetAppsFilters struct {
	common.GetFromToFilters
	// The Zscaler location (ID). You can add multiple location IDs.
	Loc []int `json:"loc,omitempty" url:"loc,omitempty"`
	// The department (ID). You can add multiple department IDs.
	Dept []int `json:"dept,omitempty" url:"dept,omitempty"`
	// The active geolocation (ID). You can add multiple active geolocation IDs.
	Geo []string `json:"geo,omitempty" url:"geo,omitempty"`
}
