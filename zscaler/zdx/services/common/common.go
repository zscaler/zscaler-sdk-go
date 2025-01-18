package common

import "fmt"

type Metric struct {
	Metric     string      `json:"metric,omitempty"`
	Unit       string      `json:"unit,omitempty"`
	DataPoints []DataPoint `json:"datapoints"`
}

type DataPoint struct {
	TimeStamp int     `json:"timestamp,omitempty"`
	Value     float64 `json:"value,omitempty"`
}

type GetFromToFilters struct {
	// The start time (in seconds) for the query. The value is entered in Unix Epoch. If not entered, returns the data for the last 2 hours.
	From int `json:"from,omitempty" url:"from,omitempty"`
	// The end time (in seconds) for the query. The value is entered in Unix Epoch. If not entered, returns the data for the last 2 hours.
	To             int      `json:"to,omitempty" url:"to,omitempty"`
	Loc            []int    `json:"loc,omitempty" url:"loc,omitempty"`
	Dept           []int    `json:"dept,omitempty" url:"dept,omitempty"`
	Geo            []string `json:"geo,omitempty" url:"geo,omitempty"`
	LocationGroups []string `json:"location_groups,omitempty" url:"location_groups,omitempty"`
	MetricName     string   `json:"metric_name,omitempty" url:"metric_name,omitempty"`
	Offset         string   `json:"offset,omitempty" url:"offset,omitempty"`
	// The number of items that must be returned per request from the list.
	Limit int `json:"limit,omitempty" url:"limit,omitempty"`
	// Search for a user name or email. The search results include active users for the first 1000 matches.
	Q string `json:"q,omitempty" url:"q,omitempty"`
}

// Centralized safe conversion function
func SafeCastToInt(value int64) (int, error) {
	minInt := int64(-1 << 31)      // Minimum value of int
	maxInt := int64((1 << 31) - 1) // Maximum value of int

	if value < minInt || value > maxInt {
		return 0, fmt.Errorf("value %d is out of range for int type", value)
	}
	return int(value), nil
}
