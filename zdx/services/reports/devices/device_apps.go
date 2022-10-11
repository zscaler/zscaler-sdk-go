package devices

const (
	deviceAppsEndpoint = "/apps"
)

/*
https://help.zscaler.com/zdx/reports#/apps-get
Lists all active applications configured for a tenant.
The endpoint gets each application’s ZDX score (default for the last 2 hours), most impacted location, and the total number of users impacted.
To learn more, see About the ZDX Dashboard.
*/
type Model struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Score int    `json:"score,omitempty"`
}
