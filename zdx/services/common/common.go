package common

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
	To         int      `json:"to,omitempty" url:"to,omitempty"`
	Loc        []int    `json:"loc,omitempty" url:"loc,omitempty"`
	Dept       []int    `json:"dept,omitempty" url:"dept,omitempty"`
	Geo        []string `json:"geo,omitempty" url:"geo,omitempty"`
	MetricName string   `json:"metric_name,omitempty" url:"metric_name,omitempty"`
}
