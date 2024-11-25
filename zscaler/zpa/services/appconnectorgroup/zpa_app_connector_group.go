package appconnectorgroup

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorcontroller"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig                = "/zpa/mgmtconfig/v1/admin/customers/"
	appConnectorGroupEndpoint = "/appConnectorGroup"
)

type AppConnectorGroup struct {
	CityCountry                   string                                `json:"cityCountry"`
	CountryCode                   string                                `json:"countryCode,omitempty"`
	CreationTime                  string                                `json:"creationTime,omitempty"`
	Description                   string                                `json:"description,omitempty"`
	DNSQueryType                  string                                `json:"dnsQueryType,omitempty"`
	Enabled                       bool                                  `json:"enabled"`
	ConnectorGroupType            string                                `json:"connectorGroupType,omitempty"`
	GeoLocationID                 string                                `json:"geoLocationId,omitempty"`
	ID                            string                                `json:"id,omitempty"`
	Latitude                      string                                `json:"latitude,omitempty"`
	Location                      string                                `json:"location,omitempty"`
	Longitude                     string                                `json:"longitude,omitempty"`
	ModifiedBy                    string                                `json:"modifiedBy,omitempty"`
	ModifiedTime                  string                                `json:"modifiedTime,omitempty"`
	Name                          string                                `json:"name,omitempty"`
	OverrideVersionProfile        bool                                  `json:"overrideVersionProfile"`
	PRAEnabled                    bool                                  `json:"praEnabled"`
	WAFDisabled                   bool                                  `json:"wafDisabled"`
	UpgradeDay                    string                                `json:"upgradeDay,omitempty"`
	UpgradeTimeInSecs             string                                `json:"upgradeTimeInSecs,omitempty"`
	VersionProfileID              string                                `json:"versionProfileId,omitempty"`
	VersionProfileName            string                                `json:"versionProfileName,omitempty"`
	VersionProfileVisibilityScope string                                `json:"versionProfileVisibilityScope,omitempty"`
	TCPQuickAckApp                bool                                  `json:"tcpQuickAckApp"`
	TCPQuickAckAssistant          bool                                  `json:"tcpQuickAckAssistant"`
	UseInDrMode                   bool                                  `json:"useInDrMode"`
	TCPQuickAckReadAssistant      bool                                  `json:"tcpQuickAckReadAssistant"`
	LSSAppConnectorGroup          bool                                  `json:"lssAppConnectorGroup"`
	MicroTenantID                 string                                `json:"microtenantId,omitempty"`
	MicroTenantName               string                                `json:"microtenantName,omitempty"`
	SiteID                        string                                `json:"siteId,omitempty"`
	SiteName                      string                                `json:"siteName,omitempty"`
	AppServerGroup                []AppServerGroup                      `json:"serverGroups,omitempty"`
	Connectors                    []appconnectorcontroller.AppConnector `json:"connectors,omitempty"`
	NPAssistantGroup              NPAssistantGroup                      `json:"npAssistantGroup,omitempty"`
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

type NPAssistantGroup struct {
	AppConnectorGroupID string      `json:"appConnectorGroupId,omitempty"`
	CreationTime        string      `json:"creationTime,omitempty"`
	ID                  string      `json:"id,omitempty"`
	LanSubnets          []LanSubnet `json:"lanSubnets,omitempty"`
	ModifiedBy          string      `json:"modifiedBy,omitempty"`
	ModifiedTime        string      `json:"modifiedTime,omitempty"`
}

type LanSubnet struct {
	AppConnectorGroupID string        `json:"appConnectorGroupId,omitempty"`
	CreationTime        string        `json:"creationTime,omitempty"`
	Description         string        `json:"description,omitempty"`
	FQDNs               []string      `json:"fqdns,omitempty"`
	ID                  string        `json:"id,omitempty"`
	ModifiedBy          string        `json:"modifiedBy,omitempty"`
	ModifiedTime        string        `json:"modifiedTime,omitempty"`
	Name                string        `json:"name,omitempty"`
	NPDnsNsRecord       NPDnsNsRecord `json:"npDnsNsRecord,omitempty"`
	NPServerIPs         []string      `json:"npserverips,omitempty"`
	OldAuditString      string        `json:"oldAuditString,omitempty"`
	Subnet              string        `json:"subnet,omitempty"`
}

type NPDnsNsRecord struct {
	CreationTime  string   `json:"creationTime,omitempty"`
	FQDN          []string `json:"fqdn,omitempty"`
	ID            string   `json:"id,omitempty"`
	ModifiedBy    string   `json:"modifiedBy,omitempty"`
	ModifiedTime  string   `json:"modifiedTime,omitempty"`
	Name          string   `json:"name,omitempty"`
	NameserverIPs []string `json:"nameserverIps,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, appConnectorGroupID string) (*AppConnectorGroup, *http.Response, error) {
	v := new(AppConnectorGroup)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+appConnectorGroupEndpoint, appConnectorGroupID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, appConnectorGroupName string) (*AppConnectorGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + appConnectorGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AppConnectorGroup](ctx, service.Client, relativeURL, common.Filter{Search: appConnectorGroupName, MicroTenantID: service.MicroTenantID()})
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

func Create(ctx context.Context, service *zscaler.Service, appConnectorGroup AppConnectorGroup) (*AppConnectorGroup, *http.Response, error) {
	v := new(AppConnectorGroup)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+appConnectorGroupEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, appConnectorGroup, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, appConnectorGroupID string, appConnectorGroup *AppConnectorGroup) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+appConnectorGroupEndpoint, appConnectorGroupID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, appConnectorGroup, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, appConnectorGroupID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+appConnectorGroupEndpoint, appConnectorGroupID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]AppConnectorGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + appConnectorGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AppConnectorGroup](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
