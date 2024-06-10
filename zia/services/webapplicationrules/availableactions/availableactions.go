package availableactions

const (
	availableActionsEndpoint = "/webApplicationRules/availableActions"
)

type AvailableActions struct {
	CloudApps []CloudApps `json:"cloudApps,omitempty"`
	Type      string      `json:"type,omitempty"`
}
type CloudApps struct {
	Val                 int    `json:"val,omitempty"`
	Name                string `json:"name,omitempty"`
	WebApplicationClass string `json:"webApplicationClass,omitempty"`
	BackendName         string `json:"backendName,omitempty"`
	OriginalName        string `json:"originalName,omitempty"`
	Deprecated          bool   `json:"deprecated,omitempty"`
	Misc                bool   `json:"misc,omitempty"`
	AppNotReady         bool   `json:"appNotReady,omitempty"`
	UnderMigration      bool   `json:"underMigration,omitempty"`
	AppCatModified      bool   `json:"appCatModified,omitempty"`
}
