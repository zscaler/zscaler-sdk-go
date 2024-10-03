package servergroup

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appservercontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig          = "/mgmtconfig/v1/admin/customers/"
	serverGroupEndpoint = "/serverGroup"
)

type ServerGroup struct {
	ID                 string                                  `json:"id,omitempty"`
	Enabled            bool                                    `json:"enabled"`
	Name               string                                  `json:"name,omitempty"`
	Description        string                                  `json:"description,omitempty"`
	IpAnchored         bool                                    `json:"ipAnchored"`
	ConfigSpace        string                                  `json:"configSpace,omitempty"`
	DynamicDiscovery   bool                                    `json:"dynamicDiscovery"`
	CreationTime       string                                  `json:"creationTime,omitempty"`
	ModifiedBy         string                                  `json:"modifiedBy,omitempty"`
	ModifiedTime       string                                  `json:"modifiedTime,omitempty"`
	MicroTenantID      string                                  `json:"microtenantId,omitempty"`
	MicroTenantName    string                                  `json:"microtenantName,omitempty"`
	AppConnectorGroups []appconnectorgroup.AppConnectorGroup   `json:"appConnectorGroups"`
	Servers            []appservercontroller.ApplicationServer `json:"servers"`
	Applications       []Applications                          `json:"applications"`
}

type Applications struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type AppConnectorGroups struct {
	Citycountry           string            `json:"cityCountry,omitempty"`
	CountryCode           string            `json:"countryCode,omitempty"`
	CreationTime          string            `json:"creationTime,omitempty"`
	Description           string            `json:"description,omitempty"`
	DnsqueryType          string            `json:"dnsQueryType,omitempty"`
	Enabled               bool              `json:"enabled"`
	GeolocationID         string            `json:"geoLocationId,omitempty"`
	ID                    string            `json:"id,omitempty"`
	Latitude              string            `json:"latitude,omitempty"`
	Location              string            `json:"location,omitempty"`
	Longitude             string            `json:"longitude,omitempty"`
	ModifiedBy            string            `json:"modifiedBy,omitempty"`
	ModifiedTime          string            `json:"modifiedTime,omitempty"`
	Name                  string            `json:"name"`
	SiemAppconnectorGroup bool              `json:"siemAppConnectorGroup"`
	UpgradeDay            string            `json:"upgradeDay,omitempty"`
	UpgradeTimeinSecs     string            `json:"upgradeTimeInSecs,omitempty"`
	VersionProfileID      string            `json:"versionProfileId,omitempty"`
	AppServerGroups       []AppServerGroups `json:"serverGroups,omitempty"`
	Connectors            []Connectors      `json:"connectors,omitempty"`
}

type Connectors struct {
	ApplicationStartTime     string                 `json:"applicationStartTime,omitempty"`
	AppConnectorGroupID      string                 `json:"appConnectorGroupId,omitempty"`
	AppConnectorGroupName    string                 `json:"appConnectorGroupName,omitempty"`
	ControlChannelStatus     string                 `json:"controlChannelStatus,omitempty"`
	CreationTime             string                 `json:"creationTime,omitempty"`
	CtrlBrokerName           string                 `json:"ctrlBrokerName,omitempty"`
	CurrentVersion           string                 `json:"currentVersion,omitempty"`
	Description              string                 `json:"description,omitempty"`
	Enabled                  bool                   `json:"enabled"`
	ExpectedUpgradeTime      string                 `json:"expectedUpgradeTime,omitempty"`
	ExpectedVersion          string                 `json:"expectedVersion,omitempty"`
	Fingerprint              string                 `json:"fingerprint,omitempty"`
	ID                       string                 `json:"id,omitempty"`
	IPACL                    []string               `json:"ipAcl,omitempty"`
	IssuedCertID             string                 `json:"issuedCertId,omitempty"`
	LastBrokerConnecttime    string                 `json:"lastBrokerConnectTime,omitempty"`
	LastBrokerDisconnectTime string                 `json:"lastBrokerDisconnectTime,omitempty"`
	LastUpgradeTime          string                 `json:"lastUpgradeTime,omitempty"`
	Latitude                 float64                `json:"latitude,omitempty"`
	Location                 string                 `json:"location,omitempty"`
	Longitude                float64                `json:"longitude,omitempty"`
	ModifiedBy               string                 `json:"modifiedBy,omitempty"`
	ModifiedTime             string                 `json:"modifiedTime,omitempty"`
	Name                     string                 `json:"name"`
	Platform                 string                 `json:"platform,omitempty"`
	PreviousVersion          string                 `json:"previousVersion,omitempty"`
	PrivateIP                string                 `json:"privateIp,omitempty"`
	PublicIP                 string                 `json:"publicIp,omitempty"`
	SigningCert              map[string]interface{} `json:"signingCert,omitempty"`
	UpgradeAttempt           string                 `json:"upgradeAttempt,omitempty"`
	UpgradeStatus            string                 `json:"upgradeStatus,omitempty"`
}

type AppServerGroups struct {
	ConfigSpace      string `json:"configSpace,omitempty"`
	CreationTime     string `json:"creationTime,omitempty"`
	Description      string `json:"description,omitempty"`
	Enabled          bool   `json:"enabled"`
	ID               string `json:"id,omitempty"`
	DynamicDiscovery bool   `json:"dynamicDiscovery"`
	ModifiedBy       string `json:"modifiedBy,omitempty"`
	ModifiedTime     string `json:"modifiedTime,omitempty"`
	Name             string `json:"name"`
}

func Get(service *services.Service, groupID string) (*ServerGroup, *http.Response, error) {
	v := new(ServerGroup)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+serverGroupEndpoint, groupID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(service *services.Service, serverGroupName string) (*ServerGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + serverGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ServerGroup](service.Client, relativeURL, common.Filter{Search: serverGroupName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, serverGroupName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no server group named '%s' was found", serverGroupName)
}

func Create(service *services.Service, serverGroup *ServerGroup) (*ServerGroup, *http.Response, error) {
	v := new(ServerGroup)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+serverGroupEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, serverGroup, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(service *services.Service, groupId string, serverGroup *ServerGroup) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+serverGroupEndpoint, groupId)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, serverGroup, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(service *services.Service, groupId string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+serverGroupEndpoint, groupId)
	resp, err := service.Client.NewRequestDo("DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(service *services.Service) ([]ServerGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + serverGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ServerGroup](service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
