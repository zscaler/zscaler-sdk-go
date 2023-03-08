package applicationsegment

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zpa/services/applicationsegmentpra"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/browseraccess"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/common"
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
	PassiveHealthEnabled      bool                                `json:"passiveHealthEnabled"`
	DoubleEncrypt             bool                                `json:"doubleEncrypt"`
	ConfigSpace               string                              `json:"configSpace,omitempty"`
	Applications              string                              `json:"applications,omitempty"`
	BypassType                string                              `json:"bypassType,omitempty"`
	HealthCheckType           string                              `json:"healthCheckType,omitempty"`
	IsCnameEnabled            bool                                `json:"isCnameEnabled"`
	IpAnchored                bool                                `json:"ipAnchored"`
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
	IsIncompleteDRConfig      bool                                `json:"isIncompleteDRConfig,omitempty"`
	UseInDrMode               bool                                `json:"useInDrMode"`
	TCPPortRanges             []string                            `json:"tcpPortRanges"`
	UDPPortRanges             []string                            `json:"udpPortRanges"`
	TCPAppPortRange           []common.NetworkPorts               `json:"tcpPortRange,omitempty"`
	UDPAppPortRange           []common.NetworkPorts               `json:"udpPortRange,omitempty"`
	ServerGroups              []AppServerGroups                   `json:"serverGroups,omitempty"`
	DefaultIdleTimeout        string                              `json:"defaultIdleTimeout,omitempty"`
	DefaultMaxAge             string                              `json:"defaultMaxAge,omitempty"`
	CommonAppsDto             applicationsegmentpra.CommonAppsDto `json:"commonAppsDto,omitempty"`
	ClientlessApps            []browseraccess.ClientlessApps      `json:"clientlessApps,omitempty"`
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

func (service *Service) Get(applicationID string) (*ApplicationSegmentResource, *http.Response, error) {
	v := new(ApplicationSegmentResource)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentEndpoint, applicationID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(appName string) (*ApplicationSegmentResource, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appSegmentEndpoint
	list, resp, err := common.GetAllPagesGeneric[ApplicationSegmentResource](service.Client, relativeURL, "")
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

func (service *Service) Create(appSegment ApplicationSegmentResource) (*ApplicationSegmentResource, *http.Response, error) {
	v := new(ApplicationSegmentResource)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+appSegmentEndpoint, nil, appSegment, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) Update(applicationId string, appSegmentRequest ApplicationSegmentResource) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentEndpoint, applicationId)
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, nil, appSegmentRequest, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (service *Service) Delete(applicationId string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentEndpoint, applicationId)
	resp, err := service.Client.NewRequestDo("DELETE", relativeURL, common.DeleteApplicationQueryParams{ForceDelete: true}, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (service *Service) GetAll() ([]ApplicationSegmentResource, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appSegmentEndpoint
	list, resp, err := common.GetAllPagesGeneric[ApplicationSegmentResource](service.Client, relativeURL, "")
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
