package appconnectorgroup

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appconnectorcontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig                = "/mgmtconfig/v1/admin/customers/"
	appConnectorGroupEndpoint = "/appConnectorGroup"
)

type AppConnectorGroup struct {
	CityCountry                   string                                `json:"cityCountry"`
	CountryCode                   string                                `json:"countryCode,omitempty"`
	CreationTime                  string                                `json:"creationTime,omitempty"`
	Description                   string                                `json:"description,omitempty"`
	DNSQueryType                  string                                `json:"dnsQueryType,omitempty"`
	Enabled                       bool                                  `json:"enabled"`
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
	AppServerGroup                []AppServerGroup                      `json:"serverGroups,omitempty"`
	Connectors                    []appconnectorcontroller.AppConnector `json:"connectors,omitempty"`
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

func Get(service *services.Service, appConnectorGroupID string) (*AppConnectorGroup, *http.Response, error) {
	v := new(AppConnectorGroup)
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+appConnectorGroupEndpoint, appConnectorGroupID)
	resp, err := service.Client.NewRequestDo("GET", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(service *services.Service, appConnectorGroupName string) (*AppConnectorGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appConnectorGroupEndpoint
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

func Create(service *services.Service, appConnectorGroup AppConnectorGroup) (*AppConnectorGroup, *http.Response, error) {
	v := new(AppConnectorGroup)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+appConnectorGroupEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, appConnectorGroup, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Update(service *services.Service, appConnectorGroupID string, appConnectorGroup *AppConnectorGroup) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appConnectorGroupEndpoint, appConnectorGroupID)
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, appConnectorGroup, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func Delete(service *services.Service, appConnectorGroupID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appConnectorGroupEndpoint, appConnectorGroupID)
	resp, err := service.Client.NewRequestDo("DELETE", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAll(service *services.Service) ([]AppConnectorGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appConnectorGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AppConnectorGroup](service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
