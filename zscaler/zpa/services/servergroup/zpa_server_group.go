package servergroup

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/appservercontroller"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig          = "/zpa/mgmtconfig/v1/admin/customers/"
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

type ApplicationServer struct {
	Address           string   `json:"address,omitempty"`
	AppServerGroupIds []string `json:"appServerGroupIds,omitempty"`
	ConfigSpace       string   `json:"configSpace,omitempty"`
	CreationTime      string   `json:"creationTime,omitempty"`
	Description       string   `json:"description,omitempty"`
	Enabled           bool     `json:"enabled"`
	ID                string   `json:"id,omitempty"`
	ModifiedBy        string   `json:"modifiedBy,omitempty"`
	ModifiedTime      string   `json:"modifiedTime,omitempty"`
	Name              string   `json:"name"`
}

func Get(ctx context.Context, service *zscaler.Service, groupID string) (*ServerGroup, *http.Response, error) {
	v := new(ServerGroup)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+serverGroupEndpoint, groupID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, serverGroupName string) (*ServerGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + serverGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ServerGroup](ctx, service.Client, relativeURL, common.Filter{Search: serverGroupName, MicroTenantID: service.MicroTenantID()})
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

func Create(ctx context.Context, service *zscaler.Service, serverGroup *ServerGroup) (*ServerGroup, *http.Response, error) {
	v := new(ServerGroup)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+serverGroupEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, serverGroup, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, groupId string, serverGroup *ServerGroup) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+serverGroupEndpoint, groupId)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, serverGroup, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, groupId string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+serverGroupEndpoint, groupId)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]ServerGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + serverGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ServerGroup](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
