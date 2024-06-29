package users

import "github.com/zscaler/zscaler-sdk-go/v2/zdx/services/common"

type GetUsersFilters struct {
	common.GetFromToFilters
	// The start time (in seconds) for the query. The value is entered in Unix Epoch. If not entered, returns the data for the last 2 hours.
	From int `json:"from,omitempty" url:"from,omitempty"`
	// The end time (in seconds) for the query. The value is entered in Unix Epoch. If not entered, returns the data for the last 2 hours.
	To int `json:"to,omitempty" url:"to,omitempty"`
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
	// Search for a user name or email. The search results include active users for the first 1000 matches.
	Q string `json:"q,omitempty" url:"q,omitempty"`
}
