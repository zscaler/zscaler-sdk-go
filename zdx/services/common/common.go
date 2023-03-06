package common

type DataPoint struct {
	TimeStamp string `json:"timestamp,omitempty"`
	Value     string `json:"value,omitempty"`
}

type GetFromToFilters struct {
	// The start time (in seconds) for the query. The value is entered in Unix Epoch. If not entered, returns the data for the last 2 hours.
	From int `json:"from,omitempty" url:"from,omitempty"`
	// The end time (in seconds) for the query. The value is entered in Unix Epoch. If not entered, returns the data for the last 2 hours.
	To int `json:"to,omitempty" url:"to,omitempty"`
}

type Metric struct {
	Metric     string      `json:"metric,omitempty"`
	Unit       string      `json:"unit,omitempty"`
	DataPoints []DataPoint `json:"datapoints"`
}
