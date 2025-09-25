package provisioning_url

import "github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"

type ProvisioningURL struct {
	ID             int                       `json:"id,omitempty"`
	Name           string                    `json:"name,omitempty"`
	Desc           string                    `json:"desc,omitempty"`
	ProvUrl        string                    `json:"provUrl,omitempty"`
	ProvUrlType    string                    `json:"provUrlType,omitempty"`
	ProvUrlData    ProvUrlData               `json:"provUrlData,omitempty"`
	UsedInEcGroups []common.IDNameExtensions `json:"usedInEcGroups,omitempty"`
	Status         string                    `json:"status,omitempty"`
	LastModUid     *common.IDNameExtensions  `json:"lastModUid,omitempty"`
	LastModTime    int                       `json:"lastModTime,omitempty"`
}

type ProvUrlData struct {
	ZsCloudDomain      string                         `json:"zsCloudDomain,omitempty"`
	OrgID              int                            `json:"orgId,omitempty"`
	ConfigServer       string                         `json:"configServer,omitempty"`
	RegistrationServer string                         `json:"registrationServer,omitempty"`
	ApiServer          string                         `json:"apiServer,omitempty"`
	PacServer          string                         `json:"pacServer,omitempty"`
	LocationTemplate   LocationTemplate               `json:"locationTemplate,omitempty"`
	CloudProvider      *common.CommonIDNameExternalID `json:"cloudProvider,omitempty"`
	CloudProviderType  string                         `json:"cloudProviderType,omitempty"`
	FormFactor         string                         `json:"formFactor,omitempty"`
	HyperVisors        string                         `json:"hyperVisors,omitempty"`
	Location           *common.CommonIDNameExternalID `json:"location,omitempty"`
	BcGroup            BcGroup                        `json:"bcGroup,omitempty"`
}

type LocationTemplate struct {
	ID          int                            `json:"id,omitempty"`
	Name        string                         `json:"name,omitempty"`
	Desc        string                         `json:"desc,omitempty"`
	Template    Template                       `json:"template,omitempty"`
	Editable    bool                           `json:"editable,omitempty"`
	LastModUid  *common.CommonIDNameExternalID `json:"lastModUid,omitempty"`
	LastModTime int                            `json:"lastModTime,omitempty"`
}

type Template struct {
	TemplatePrefix                      string                         `json:"templatePrefix,omitempty"`
	XffForwardEnabled                   bool                           `json:"xffForwardEnabled,omitempty"`
	AuthRequired                        bool                           `json:"authRequired,omitempty"`
	CautionEnabled                      bool                           `json:"cautionEnabled,omitempty"`
	AupEnabled                          bool                           `json:"aupEnabled,omitempty"`
	AupTimeoutInDays                    int                            `json:"aupTimeoutInDays,omitempty"`
	OfwEnabled                          bool                           `json:"ofwEnabled,omitempty"`
	IpsControl                          bool                           `json:"ipsControl,omitempty"`
	EnforceBandwidthControl             bool                           `json:"enforceBandwidthControl,omitempty"`
	UpBandwidth                         int                            `json:"upBandwidth,omitempty"`
	DnBandwidth                         int                            `json:"dnBandwidth,omitempty"`
	DisplayTimeUnit                     string                         `json:"displayTimeUnit,omitempty"`
	IdleTimeInMinutes                   int                            `json:"idleTimeInMinutes,omitempty"`
	SurrogateIpEnforcedForKnownBrowsers bool                           `json:"surrogateIPEnforcedForKnownBrowsers,omitempty"`
	SurrogateRefreshTimeUnit            string                         `json:"surrogateRefreshTimeUnit,omitempty"`
	SurrogateRefreshTimeInMinutes       int                            `json:"surrogateRefreshTimeInMinutes,omitempty"`
	Surrogate                           bool                           `json:"surrogateIP,omitempty"`
	Editable                            bool                           `json:"editable,omitempty"`
	LastModUid                          *common.CommonIDNameExternalID `json:"lastModUid,omitempty"`
}

type BcGroup struct {
	ID                    int                            `json:"id,omitempty"`
	Name                  string                         `json:"name,omitempty"`
	Desc                  string                         `json:"desc,omitempty"`
	DeployType            string                         `json:"deployType,omitempty"`
	Status                []string                       `json:"status,omitempty"`
	Platform              string                         `json:"platform,omitempty"`
	AwsAvailabilityZone   string                         `json:"awsAvailabilityZone,omitempty"`
	AzureAvailabilityZone string                         `json:"azureAvailabilityZone,omitempty"`
	Location              *common.CommonIDNameExternalID `json:"location,omitempty"`
	MaxEcCount            int                            `json:"maxEcCount,omitempty"`
	ProvTemplate          *common.CommonIDNameExternalID `json:"provTemplate,omitempty"`
	TunnelMode            string                         `json:"tunnelMode,omitempty"`
	EcVMs                 []EcVM                         `json:"ecVMs,omitempty"`
}

type EcVM struct {
	ID                int          `json:"id,omitempty"`
	Name              string       `json:"name,omitempty"`
	Status            []string     `json:"status,omitempty"`
	OperationalStatus string       `json:"operationalStatus,omitempty"`
	FormFactor        string       `json:"formFactor,omitempty"`
	ManagementNw      ManagementNw `json:"managementNw,omitempty"`
	EcInstances       []EcInstance `json:"ecInstances,omitempty"`
	CityGeoId         int          `json:"cityGeoId,omitempty"`
	NatIp             string       `json:"natIp,omitempty"`
	ZiaGateway        string       `json:"ziaGateway,omitempty"`
	ZpaBroker         string       `json:"zpaBroker,omitempty"`
	BuildVersion      string       `json:"buildVersion,omitempty"`
	LastUpgradeTime   int          `json:"lastUpgradeTime,omitempty"`
	UpgradeStatus     int          `json:"upgradeStatus,omitempty"`
	UpgradeStartTime  int          `json:"upgradeStartTime,omitempty"`
	UpgradeEndTime    int          `json:"upgradeEndTime,omitempty"`
	UpgradeDayOfWeek  int          `json:"upgradeDayOfWeek,omitempty"`
}

type ManagementNw struct {
	ID             int    `json:"id,omitempty"`
	IpStart        string `json:"ipStart,omitempty"`
	IpEnd          string `json:"ipEnd,omitempty"`
	Netmask        string `json:"netmask,omitempty"`
	DefaultGateway string `json:"defaultGateway,omitempty"`
	NwType         string `json:"nwType,omitempty"`
	DNS            DNS    `json:"dns,omitempty"`
}

type DNS struct {
	ID      int      `json:"id,omitempty"`
	Ips     []string `json:"ips,omitempty"`
	DNSType string   `json:"dnsType,omitempty"`
}

type EcInstance struct {
	ID           int        `json:"id,omitempty"`
	InstanceType string     `json:"instanceType,omitempty"`
	ServiceIps   ServiceIps `json:"serviceIps,omitempty"`
	LbIpAddr     ServiceIps `json:"lbIpAddr,omitempty"`
	OutGwIp      string     `json:"outGwIp,omitempty"`
	NatIp        string     `json:"natIp,omitempty"`
	DnsIp        []string   `json:"dnsIp,omitempty"`
}

type ServiceIps struct {
	IpStart string `json:"ipStart,omitempty"`
	IpEnd   string `json:"ipEnd,omitempty"`
}
