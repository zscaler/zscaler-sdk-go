package applications

const (
	appsEndpoint = "/apps"
)

/*
https://help.zscaler.com/zdx/reports#/apps-get
Lists all active applications configured for a tenant.
The endpoint gets each application’s ZDX score (default for the last 2 hours), most impacted location, and the total number of users impacted.
To learn more, see About the ZDX Dashboard.
*/

type Apps struct {
	ID              int               `json:"id"`
	Name            string            `json:"name,omitempty"`
	Score           int               `json:"score,omitempty"`
	MostImpactedGeo []MostImpactedGeo `json:"most_impacted_geo"`
}

type MostImpactedGeo struct {
	ID      int    `json:"id"`
	City    string `json:"city,omitempty"`
	Region  string `json:"region,omitempty"`
	Country string `json:"country,omitempty"`
	GeoType string `json:"geo_type,omitempty"`
}

type Stats struct {
	ID      int    `json:"id"`
	City    string `json:"city,omitempty"`
	Region  string `json:"region,omitempty"`
	Country string `json:"country,omitempty"`
	GeoType string `json:"geo_type,omitempty"`
}
