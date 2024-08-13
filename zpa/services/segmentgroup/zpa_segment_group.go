package segmentgroup

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfigV1         = "/mgmtconfig/v1/admin/customers/"
	mgmtConfigV2         = "/mgmtconfig/v2/admin/customers/"
	segmentGroupEndpoint = "/segmentGroup"
)

type SegmentGroup struct {
	ID                  string             `json:"id,omitempty"`
	Name                string             `json:"name"`
	Description         string             `json:"description,omitempty"`
	Enabled             bool               `json:"enabled"`
	ConfigSpace         string             `json:"configSpace,omitempty"`
	CreationTime        string             `json:"creationTime,omitempty"`
	ModifiedBy          string             `json:"modifiedBy,omitempty"`
	ModifiedTime        string             `json:"modifiedTime,omitempty"`
	PolicyMigrated      bool               `json:"policyMigrated"`
	TcpKeepAliveEnabled string             `json:"tcpKeepAliveEnabled,omitempty"`
	MicroTenantID       string             `json:"microtenantId,omitempty"`
	MicroTenantName     string             `json:"microtenantName,omitempty"`
	AddedApps           string             `json:"addedApps,omitempty"`
	DeletedApps         string             `json:"deletedApps,omitempty"`
	Applications        []Application      `json:"applications"`
	ApplicationNames    []ApplicationNames `json:"applicationNames,omitempty"`
}

type Application struct {
	BypassType           string           `json:"bypassType,omitempty"`
	ConfigSpace          string           `json:"configSpace,omitempty"`
	CreationTime         string           `json:"creationTime,omitempty"`
	DefaultIdleTimeout   string           `json:"defaultIdleTimeout,omitempty"`
	DefaultMaxAge        string           `json:"defaultMaxAge,omitempty"`
	Description          string           `json:"description,omitempty"`
	DomainName           string           `json:"domainName,omitempty"`
	DomainNames          []string         `json:"domainNames,omitempty"`
	DoubleEncrypt        bool             `json:"doubleEncrypt"`
	Enabled              bool             `json:"enabled"`
	HealthCheckType      string           `json:"healthCheckType,omitempty"`
	ID                   string           `json:"id,omitempty"`
	IPAnchored           bool             `json:"ipAnchored"`
	LogFeatures          []string         `json:"logFeatures,omitempty"`
	ModifiedBy           string           `json:"modifiedBy,omitempty"`
	ModifiedTime         string           `json:"modifiedTime,omitempty"`
	Name                 string           `json:"name"`
	PassiveHealthEnabled bool             `json:"passiveHealthEnabled"`
	ServerGroup          []AppServerGroup `json:"serverGroups,omitempty"`
	TCPPortRanges        interface{}      `json:"tcpPortRanges,omitempty"`
	TCPPortsIn           interface{}      `json:"tcpPortsIn,omitempty"`
	TCPPortsOut          interface{}      `json:"tcpPortsOut,omitempty"`
	UDPPortRanges        interface{}      `json:"udpPortRangesg,omitempty"`
}

type AppServerGroup struct {
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

type ApplicationNames struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

func Get(service *services.Service, segmentGroupID string) (*SegmentGroup, *http.Response, error) {
	v := new(SegmentGroup)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfigV1+service.Client.Config.CustomerID+segmentGroupEndpoint, segmentGroupID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(service *services.Service, segmentName string) (*SegmentGroup, *http.Response, error) {
	relativeURL := mgmtConfigV1 + service.Client.Config.CustomerID + segmentGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[SegmentGroup](service.Client, relativeURL, common.Filter{Search: segmentName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, segmentName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no application named '%s' was found", segmentName)
}

func Create(service *services.Service, segmentGroup *SegmentGroup) (*SegmentGroup, *http.Response, error) {
	v := new(SegmentGroup)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfigV1+service.Client.Config.CustomerID+segmentGroupEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, segmentGroup, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(service *services.Service, segmentGroupId string, segmentGroupRequest *SegmentGroup) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfigV1+service.Client.Config.CustomerID+segmentGroupEndpoint, segmentGroupId)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, segmentGroupRequest, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func UpdateV2(service *services.Service, segmentGroupId string, segmentGroupRequest *SegmentGroup) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfigV2+service.Client.Config.CustomerID+segmentGroupEndpoint, segmentGroupId)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, segmentGroupRequest, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(service *services.Service, segmentGroupId string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfigV1+service.Client.Config.CustomerID+segmentGroupEndpoint, segmentGroupId)
	resp, err := service.Client.NewRequestDo("DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(service *services.Service) ([]SegmentGroup, *http.Response, error) {
	relativeURL := mgmtConfigV1 + service.Client.Config.CustomerID + segmentGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[SegmentGroup](service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
