package applicationsegment

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegmentpra"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/browseraccess"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig         = "/mgmtconfig/v1/admin/customers/"
	appSegmentEndpoint = "/application"
)

type ApplicationSegmentResource struct {
	ID                        string                              `json:"id,omitempty"`
	DomainNames               []string                            `json:"domainNames,omitempty"`
	Name                      string                              `json:"name,omitempty"`
	Description               string                              `json:"description,omitempty"`
	Enabled                   bool                                `json:"enabled"`
	ADPEnabled                bool                                `json:"adpEnabled"`
	PassiveHealthEnabled      bool                                `json:"passiveHealthEnabled"`
	DoubleEncrypt             bool                                `json:"doubleEncrypt"`
	ConfigSpace               string                              `json:"configSpace,omitempty"`
	Applications              string                              `json:"applications,omitempty"`
	BypassType                string                              `json:"bypassType,omitempty"`
	BypassOnReauth            bool                                `json:"bypassOnReauth,omitempty"`
	HealthCheckType           string                              `json:"healthCheckType,omitempty"`
	IsCnameEnabled            bool                                `json:"isCnameEnabled"`
	IpAnchored                bool                                `json:"ipAnchored"`
	FQDNDnsCheck              bool                                `json:"fqdnDnsCheck"`
	HealthReporting           string                              `json:"healthReporting,omitempty"`
	SelectConnectorCloseToApp bool                                `json:"selectConnectorCloseToApp"`
	IcmpAccessType            string                              `json:"icmpAccessType,omitempty"`
	AppRecommendationId       string                              `json:"appRecommendationId,omitempty"`
	SegmentGroupID            string                              `json:"segmentGroupId"`
	SegmentGroupName          string                              `json:"segmentGroupName,omitempty"`
	CreationTime              string                              `json:"creationTime,omitempty"`
	ModifiedBy                string                              `json:"modifiedBy,omitempty"`
	ModifiedTime              string                              `json:"modifiedTime,omitempty"`
	TCPKeepAlive              string                              `json:"tcpKeepAlive,omitempty"`
	IsIncompleteDRConfig      bool                                `json:"isIncompleteDRConfig"`
	UseInDrMode               bool                                `json:"useInDrMode"`
	InspectTrafficWithZia     bool                                `json:"inspectTrafficWithZia"`
	MicroTenantID             string                              `json:"microtenantId,omitempty"`
	MicroTenantName           string                              `json:"microtenantName,omitempty"`
	MatchStyle                string                              `json:"matchStyle,omitempty"`
	TCPPortRanges             []string                            `json:"tcpPortRanges"`
	UDPPortRanges             []string                            `json:"udpPortRanges"`
	TCPAppPortRange           []common.NetworkPorts               `json:"tcpPortRange,omitempty"`
	UDPAppPortRange           []common.NetworkPorts               `json:"udpPortRange,omitempty"`
	ServerGroups              []AppServerGroups                   `json:"serverGroups"`
	DefaultIdleTimeout        string                              `json:"defaultIdleTimeout,omitempty"`
	DefaultMaxAge             string                              `json:"defaultMaxAge,omitempty"`
	CommonAppsDto             applicationsegmentpra.CommonAppsDto `json:"commonAppsDto,omitempty"`
	ClientlessApps            []browseraccess.ClientlessApps      `json:"clientlessApps,omitempty"`
	ShareToMicrotenants       []string                            `json:"shareToMicrotenants"`
	SharedMicrotenantDetails  SharedMicrotenantDetails            `json:"sharedMicrotenantDetails,omitempty"`
}

type SharedMicrotenantDetails struct {
	SharedFromMicrotenant SharedFromMicrotenant `json:"sharedFromMicrotenant,omitempty"`
	SharedToMicrotenants  []SharedToMicrotenant `json:"sharedToMicrotenants,omitempty"`
}

type SharedFromMicrotenant struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type SharedToMicrotenant struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
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

func Get(service *services.Service, applicationID string) (*ApplicationSegmentResource, *http.Response, error) {
	v := new(ApplicationSegmentResource)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentEndpoint, applicationID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(service *services.Service, appName string) (*ApplicationSegmentResource, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appSegmentEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ApplicationSegmentResource](service.Client, relativeURL, common.Filter{Search: appName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, appName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no application segment named '%s' was found", appName)
}

func Create(service *services.Service, appSegment ApplicationSegmentResource) (*ApplicationSegmentResource, *http.Response, error) {
	v := new(ApplicationSegmentResource)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+appSegmentEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, appSegment, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Update(service *services.Service, appID string, appSegmentRequest ApplicationSegmentResource) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentEndpoint, appID)
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, appSegmentRequest, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func Delete(service *services.Service, appID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentEndpoint, appID)
	resp, err := service.Client.NewRequestDo("DELETE", relativeURL, common.DeleteApplicationQueryParams{ForceDelete: true, MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAll(service *services.Service) ([]ApplicationSegmentResource, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appSegmentEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ApplicationSegmentResource](service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	result := []ApplicationSegmentResource{}
	// filter apps
	for _, item := range list {
		if len(item.ClientlessApps) == 0 && (len(item.CommonAppsDto.AppsConfig) == 0 || !common.InList(item.CommonAppsDto.AppsConfig[0].AppTypes, "SECURE_REMOTE_ACCESS") && !common.InList(item.CommonAppsDto.AppsConfig[0].AppTypes, "INSPECT")) {
			result = append(result, item)
		}
	}
	return result, resp, nil
}
