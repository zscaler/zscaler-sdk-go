package appconnectorgroup

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig                = "/zpa/mgmtconfig/v1/admin/customers/"
	appConnectorGroupEndpoint = "/appConnectorGroup"
)

type AppConnectorGroup struct {
	CityCountry                   string           `json:"cityCountry"`
	CountryCode                   string           `json:"countryCode,omitempty"`
	CreationTime                  string           `json:"creationTime,omitempty"`
	Description                   string           `json:"description,omitempty"`
	DNSQueryType                  string           `json:"dnsQueryType,omitempty"`
	Enabled                       bool             `json:"enabled"`
	GeoLocationID                 string           `json:"geoLocationId,omitempty"`
	ID                            string           `json:"id,omitempty"`
	Latitude                      string           `json:"latitude,omitempty"`
	Location                      string           `json:"location,omitempty"`
	Longitude                     string           `json:"longitude,omitempty"`
	ModifiedBy                    string           `json:"modifiedBy,omitempty"`
	ModifiedTime                  string           `json:"modifiedTime,omitempty"`
	Name                          string           `json:"name,omitempty"`
	OverrideVersionProfile        bool             `json:"overrideVersionProfile"`
	PRAEnabled                    bool             `json:"praEnabled"`
	WAFDisabled                   bool             `json:"wafDisabled"`
	UpgradeDay                    string           `json:"upgradeDay,omitempty"`
	UpgradeTimeInSecs             string           `json:"upgradeTimeInSecs,omitempty"`
	VersionProfileID              string           `json:"versionProfileId,omitempty"`
	VersionProfileName            string           `json:"versionProfileName,omitempty"`
	VersionProfileVisibilityScope string           `json:"versionProfileVisibilityScope,omitempty"`
	TCPQuickAckApp                bool             `json:"tcpQuickAckApp"`
	TCPQuickAckAssistant          bool             `json:"tcpQuickAckAssistant"`
	UseInDrMode                   bool             `json:"useInDrMode"`
	TCPQuickAckReadAssistant      bool             `json:"tcpQuickAckReadAssistant"`
	LSSAppConnectorGroup          bool             `json:"lssAppConnectorGroup"`
	MicroTenantID                 string           `json:"microtenantId,omitempty"`
	MicroTenantName               string           `json:"microtenantName,omitempty"`
	AppServerGroup                []AppServerGroup `json:"serverGroups,omitempty"`
	Connectors                    []Connector      `json:"connectors,omitempty"`
}

type Connector struct {
	ApplicationStartTime             string                 `json:"applicationStartTime,omitempty"`
	AppConnectorGroupID              string                 `json:"appConnectorGroupId,omitempty"`
	AppConnectorGroupName            string                 `json:"appConnectorGroupName,omitempty"`
	ControlChannelStatus             string                 `json:"controlChannelStatus,omitempty"`
	CreationTime                     string                 `json:"creationTime,omitempty"`
	CtrlBrokerName                   string                 `json:"ctrlBrokerName,omitempty"`
	CurrentVersion                   string                 `json:"currentVersion,omitempty"`
	Description                      string                 `json:"description,omitempty"`
	Enabled                          bool                   `json:"enabled,omitempty"`
	ExpectedUpgradeTime              string                 `json:"expectedUpgradeTime,omitempty"`
	ExpectedVersion                  string                 `json:"expectedVersion,omitempty"`
	Fingerprint                      string                 `json:"fingerprint,omitempty"`
	ID                               string                 `json:"id,omitempty"`
	IPACL                            string                 `json:"ipAcl,omitempty"`
	IssuedCertID                     string                 `json:"issuedCertId,omitempty"`
	LastBrokerConnectTime            string                 `json:"lastBrokerConnectTime,omitempty"`
	LastBrokerConnectTimeDuration    string                 `json:"lastBrokerConnectTimeDuration,omitempty"`
	LastBrokerDisconnectTime         string                 `json:"lastBrokerDisconnectTime,omitempty"`
	LastBrokerDisconnectTimeDuration string                 `json:"lastBrokerDisconnectTimeDuration,omitempty"`
	LastUpgradeTime                  string                 `json:"lastUpgradeTime,omitempty"`
	Latitude                         string                 `json:"latitude,omitempty"`
	Location                         string                 `json:"location,omitempty"`
	Longitude                        string                 `json:"longitude,omitempty"`
	ModifiedBy                       string                 `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                 `json:"modifiedTime,omitempty"`
	Name                             string                 `json:"name,omitempty"`
	ProvisioningKeyID                string                 `json:"provisioningKeyId"`
	ProvisioningKeyName              string                 `json:"provisioningKeyName"`
	Platform                         string                 `json:"platform,omitempty"`
	PreviousVersion                  string                 `json:"previousVersion,omitempty"`
	PrivateIP                        string                 `json:"privateIp,omitempty"`
	PublicIP                         string                 `json:"publicIp,omitempty"`
	SargeVersion                     string                 `json:"sargeVersion,omitempty"`
	EnrollmentCert                   map[string]interface{} `json:"enrollmentCert,omitempty"`
	UpgradeAttempt                   string                 `json:"upgradeAttempt,omitempty"`
	UpgradeStatus                    string                 `json:"upgradeStatus,omitempty"`
}

type AppServerGroup struct {
	ConfigSpace      string `json:"configSpace,omitempty"`
	CreationTime     string `json:"creationTime,omitempty"`
	Description      string `json:"description,omitempty"`
	Enabled          bool   `json:"enabled,omitempty"`
	ID               string `json:"id,omitempty"`
	DynamicDiscovery bool   `json:"dynamicDiscovery,omitempty"`
	ModifiedBy       string `json:"modifiedBy,omitempty"`
	ModifiedTime     string `json:"modifiedTime,omitempty"`
	Name             string `json:"name,omitempty"`
}

func Get(service *zscaler.Service, appConnectorGroupID string) (*AppConnectorGroup, *http.Response, error) {
	v := new(AppConnectorGroup)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+appConnectorGroupEndpoint, appConnectorGroupID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(service *zscaler.Service, appConnectorGroupName string) (*AppConnectorGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + appConnectorGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AppConnectorGroup](service.Client, relativeURL, common.Filter{Search: appConnectorGroupName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, appConnectorGroupName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no app connector group named '%s' was found", appConnectorGroupName)
}

func Create(service *zscaler.Service, appConnectorGroup AppConnectorGroup) (*AppConnectorGroup, *http.Response, error) {
	v := new(AppConnectorGroup)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.GetCustomerID()+appConnectorGroupEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, appConnectorGroup, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Update(service *zscaler.Service, appConnectorGroupID string, appConnectorGroup *AppConnectorGroup) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+appConnectorGroupEndpoint, appConnectorGroupID)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, appConnectorGroup, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func Delete(service *zscaler.Service, appConnectorGroupID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+appConnectorGroupEndpoint, appConnectorGroupID)
	resp, err := service.Client.NewRequestDo("DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAll(service *zscaler.Service) ([]AppConnectorGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + appConnectorGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AppConnectorGroup](service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
