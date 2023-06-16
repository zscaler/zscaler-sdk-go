package cloudiacalerts

/*
const (
	cloudAlertsEndpoint = "/v1/alerts/iac"
)

type CloudAlerts struct {
	Kind                   string        `json:"kind,omitempty"`
	AlertID                string        `json:"alert_id,omitempty"`
	AlertType              string        `json:"alert_type,omitempty"`
	AlertDescription       string        `json:"alert_description,omitempty"`
	AlertStatus            string        `json:"alert_status,omitempty"`
	AlertStatusDescription string        `json:"alert_status_description,omitempty"`
	AlertURL               string        `json:"alert_url,omitempty"`
	CSP                    string        `json:"csp,omitempty"`
	Severity               string        `json:"severity,omitempty"`
	RiskLevel              string        `json:"risk_level,omitempty"`
	RiskScore              string        `json:"risk_score,omitempty"`
	ThreatCategory         string        `json:"threat_category,omitempty"`
	MitreAttack            string        `json:"mitre_attack,omitempty"`
	Theme                  string        `json:"theme,omitempty"`
	AssetCategory          string        `json:"asset_category,omitempty"`
	AssetType              string        `json:"asset_type,omitempty"`
	Age                    string        `json:"age,omitempty"`
	AuditProcedure         string        `json:"audit_procedure,omitempty"`
	Recommendations        string        `json:"recommendations,omitempty"`
	RemediationProcedure   string        `json:"remediation_procedure,omitempty"`
	AlertFocus             string        `json:"alert_focus,omitempty"`
	Reason                 string        `json:"reason,omitempty"`
	UpdatedBy              string        `json:"updated_by,omitempty"`
	Created                string        `json:"created,omitempty"`
	Updated                string        `json:"updated,omitempty"`
	TotalCount             int64         `json:"total_count,omitempty"`
	PrevLink               string        `json:"prev_link,omitempty"`
	NextLink               string        `json:"next_link,omitempty"`
	Policy                 Policy        `json:"policy,omitempty"`
	CloudResource          CloudResource `json:"cloud_resource,omitempty"`
	Identity               Identity      `json:"identity,omitempty"`
	Compliance             Compliance    `json:"compliance,omitempty"`
	BusinessUnit           BusinessUnit  `json:"business_unit,omitempty"`
	CloudAccount           CloudAccount  `json:"cloud_account,omitempty"`
	Kubernetes             Kubernetes    `json:"Kubernetes,omitempty"`
}

type Policy struct {
	PolicyID          string `json:"policy_id,omitempty"`
	PolicySource      string `json:"policy_source,omitempty"`
	PolicyName        string `json:"policy_name,omitempty"`
	PolicyDescription string `json:"policy_description,omitempty"`
}

type CloudResource struct {
	ResourceID     string `json:"resource_id,omitempty"`
	ResourceName   string `json:"resource_name,omitempty"`
	ResourceType   string `json:"resource_type,omitempty"`
	ResourceRegion string `json:"resource_region,omitempty"`
}

type Identity struct {
	IdentityID     string `json:"identity_id,omitempty"`
	IdentityName   string `json:"identity_name,omitempty"`
	IdentityType   string `json:"identity_type,omitempty"`
	IdentitySource string `json:"identity_source,omitempty"`
}

type Compliance struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Number  string `json:"number,omitempty"`
	Type    string `json:"type,omitempty"`
	Domains string `json:"domains,omitempty"`
}

type BusinessUnit struct {
	BusinessUnitName string `json:"business_unit_name,omitempty"`
}

type CloudAccount struct {
	AccountID      string `json:"account_id,omitempty"`
	AccountName    string `json:"account_name,omitempty"`
	Organization   string `json:"organization,omitempty"`
	OrganizationID string `json:"organization_id,omitempty"`
}

type Kubernetes struct {
	ClusteType  string `json:"cluster_type,omitempty"`
	ClusterName string `json:"cluster_name,omitempty"`
}
*/
