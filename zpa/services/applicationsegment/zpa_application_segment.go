package applicationsegment

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegmentpra"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/browseraccess"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig              = "/mgmtconfig/v1/admin/customers/"
	appSegmentEndpoint      = "/application"
	applicationTypeEndpoint = "/application/getAppsByType"
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
	SharedMicrotenantDetails  SharedMicrotenantDetails            `json:"sharedMicrotenantDetails,omitempty"`
	InconsistentConfigDetails InconsistentConfigDetails           `json:"inconsistentConfigDetails,omitempty"`
}

type InconsistentConfigDetails struct {
	Application          []common.CommonConfigDetails `json:"application,omitempty"`
	SegmentGroup         []common.CommonConfigDetails `json:"segmentGroup,omitempty"`
	BaCertificate        []common.CommonConfigDetails `json:"baCertificate,omitempty"`
	BranchConnectorGroup []common.CommonConfigDetails `json:"branchConnectorGroup,omitempty"`
	CloudConnectorGroup  []common.CommonConfigDetails `json:"cloudConnectorGroup,omitempty"`
	IDP                  []common.CommonConfigDetails `json:"idp,omitempty"`
	Location             []common.CommonConfigDetails `json:"location,omitempty"`
	MachineGroup         []common.CommonConfigDetails `json:"machineGroup,omitempty"`
	PostureProfile       []common.CommonConfigDetails `json:"postureProfile,omitempty"`
	SamlAttributes       []common.CommonConfigDetails `json:"samlAttributes,omitempty"`
	ScimAttributes       []common.CommonConfigDetails `json:"scimAttributes,omitempty"`
	ServerGroup          []common.CommonConfigDetails `json:"serverGroup,omitempty"`
	SRAApplication       []common.CommonConfigDetails `json:"sraApplication,omitempty"`
	TrustedNetwork       []common.CommonConfigDetails `json:"trustedNetwork,omitempty"`
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

func (service *Service) Get(applicationID string) (*ApplicationSegmentResource, *http.Response, error) {
	v := new(ApplicationSegmentResource)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentEndpoint, applicationID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(appName string) (*ApplicationSegmentResource, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appSegmentEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ApplicationSegmentResource](service.Client, relativeURL, common.Filter{Search: appName, MicroTenantID: service.microTenantID})
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

func (service *Service) GetByApplicationType(applicationType string, expandAll bool) ([]ApplicationSegmentResource, *http.Response, error) {
	if applicationType != "BROWSER_ACCESS" && applicationType != "SECURE_REMOTE_ACCESS" && applicationType != "INSPECT" {
		return nil, nil, fmt.Errorf("invalid applicationType '%s'. Valid types are 'BROWSER_ACCESS', 'SECURE_REMOTE_ACCESS', 'INSPECT'", applicationType)
	}
	// Constructing the query parameters as part of the URL
	relativeURL := fmt.Sprintf("%s%s%s?applicationType=%s&expandAll=%t&page=1&pagesize=20",
		mgmtConfig, service.Client.Config.CustomerID, applicationTypeEndpoint, applicationType, expandAll)
	filter := common.Filter{} // Initialize an empty filter or with minimal required fields

	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ApplicationSegmentResource](service.Client, relativeURL, filter)
	if err != nil {
		return nil, nil, err
	}

	return list, resp, nil
}

func (service *Service) Create(appSegment ApplicationSegmentResource) (*ApplicationSegmentResource, *http.Response, error) {
	v := new(ApplicationSegmentResource)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+appSegmentEndpoint, common.Filter{MicroTenantID: service.microTenantID}, appSegment, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) Update(applicationId string, appSegmentRequest ApplicationSegmentResource) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentEndpoint, applicationId)
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, appSegmentRequest, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (service *Service) Delete(applicationId string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentEndpoint, applicationId)
	resp, err := service.Client.NewRequestDo("DELETE", relativeURL, common.DeleteApplicationQueryParams{ForceDelete: true, MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (service *Service) GetAll() ([]ApplicationSegmentResource, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appSegmentEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ApplicationSegmentResource](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
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
