package serviceedgegroup

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/serviceedgecontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/trustednetwork"
)

const (
	mgmtConfig               = "/mgmtconfig/v1/admin/customers/"
	serviceEdgeGroupEndpoint = "/serviceEdgeGroup"
)

type ServiceEdgeGroup struct {
	CityCountry                   string                                        `json:"cityCountry,omitempty"`
	CountryCode                   string                                        `json:"countryCode,omitempty"`
	CreationTime                  string                                        `json:"creationTime,omitempty"`
	Description                   string                                        `json:"description,omitempty"`
	Enabled                       bool                                          `json:"enabled"`
	GeoLocationID                 string                                        `json:"geoLocationId,omitempty"`
	GraceDistanceEnabled          bool                                          `json:"graceDistanceEnabled"`
	GraceDistanceValue            string                                        `json:"graceDistanceValue,omitempty"`
	GraceDistanceValueUnit        string                                        `json:"graceDistanceValueUnit,omitempty"`
	ID                            string                                        `json:"id,omitempty"`
	IsPublic                      string                                        `json:"isPublic,omitempty"`
	Latitude                      string                                        `json:"latitude,omitempty"`
	Location                      string                                        `json:"location,omitempty"`
	Longitude                     string                                        `json:"longitude,omitempty"`
	ModifiedBy                    string                                        `json:"modifiedBy,omitempty"`
	ModifiedTime                  string                                        `json:"modifiedTime,omitempty"`
	Name                          string                                        `json:"name,omitempty"`
	UseInDrMode                   bool                                          `json:"useInDrMode"`
	OverrideVersionProfile        bool                                          `json:"overrideVersionProfile"`
	ServiceEdges                  []serviceedgecontroller.ServiceEdgeController `json:"serviceEdges,omitempty"`
	TrustedNetworks               []trustednetwork.TrustedNetwork               `json:"trustedNetworks,omitempty"`
	UpgradeDay                    string                                        `json:"upgradeDay,omitempty"`
	UpgradeTimeInSecs             string                                        `json:"upgradeTimeInSecs,omitempty"`
	VersionProfileID              string                                        `json:"versionProfileId,omitempty"`
	VersionProfileName            string                                        `json:"versionProfileName,omitempty"`
	VersionProfileVisibilityScope string                                        `json:"versionProfileVisibilityScope,omitempty"`
	ObjectType                    string                                        `json:"objectType,omitempty"`
	ScopeName                     string                                        `json:"scopeName,omitempty"`
	RestrictedEntity              bool                                          `json:"restrictedEntity,omitempty"`
	AltCloud                      string                                        `json:"altCloud,omitempty"`
	MicroTenantID                 string                                        `json:"microtenantId,omitempty"`
	MicroTenantName               string                                        `json:"microtenantName,omitempty"`
	SiteID                        string                                        `json:"siteId,omitempty"`
	SiteName                      string                                        `json:"siteName,omitempty"`
}

func Get(service *services.Service, serviceEdgeGroupID string) (*ServiceEdgeGroup, *http.Response, error) {
	v := new(ServiceEdgeGroup)
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+serviceEdgeGroupEndpoint, serviceEdgeGroupID)
	resp, err := service.Client.NewRequestDo("GET", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(service *services.Service, serviceEdgeGroupName string) (*ServiceEdgeGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + serviceEdgeGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ServiceEdgeGroup](service.Client, relativeURL, common.Filter{Search: serviceEdgeGroupName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, serviceEdgeGroupName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no service edge group named '%s' was found", serviceEdgeGroupName)
}

func Create(service *services.Service, serviceEdge ServiceEdgeGroup) (*ServiceEdgeGroup, *http.Response, error) {
	v := new(ServiceEdgeGroup)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+serviceEdgeGroupEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, serviceEdge, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Update(service *services.Service, serviceEdgeGroupID string, serviceEdge *ServiceEdgeGroup) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+serviceEdgeGroupEndpoint, serviceEdgeGroupID)
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, serviceEdge, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Delete(service *services.Service, serviceEdgeGroupID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+serviceEdgeGroupEndpoint, serviceEdgeGroupID)
	resp, err := service.Client.NewRequestDo("DELETE", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAll(service *services.Service) ([]ServiceEdgeGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + serviceEdgeGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ServiceEdgeGroup](service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
