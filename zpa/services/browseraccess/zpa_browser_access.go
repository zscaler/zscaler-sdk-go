package browseraccess

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig              = "/mgmtconfig/v1/admin/customers/"
	browserAccessEndpoint   = "/application"
	applicationTypeEndpoint = "/application/getAppsByType"
)

type BrowserAccess struct {
	ID                        string                   `json:"id,omitempty"`
	Name                      string                   `json:"name,omitempty"`
	Description               string                   `json:"description,omitempty"`
	SegmentGroupID            string                   `json:"segmentGroupId,omitempty"`
	SegmentGroupName          string                   `json:"segmentGroupName,omitempty"`
	BypassType                string                   `json:"bypassType,omitempty"`
	BypassOnReauth            bool                     `json:"bypassOnReauth,omitempty"`
	AppRecommendationId       string                   `json:"appRecommendationId,omitempty"`
	MatchStyle                string                   `json:"matchStyle,omitempty"`
	ConfigSpace               string                   `json:"configSpace,omitempty"`
	DomainNames               []string                 `json:"domainNames,omitempty"`
	Enabled                   bool                     `json:"enabled"`
	PassiveHealthEnabled      bool                     `json:"passiveHealthEnabled"`
	FQDNDnsCheck              bool                     `json:"fqdnDnsCheck"`
	SelectConnectorCloseToApp bool                     `json:"selectConnectorCloseToApp"`
	DoubleEncrypt             bool                     `json:"doubleEncrypt"`
	HealthCheckType           string                   `json:"healthCheckType,omitempty"`
	IsCnameEnabled            bool                     `json:"isCnameEnabled"`
	IPAnchored                bool                     `json:"ipAnchored"`
	TCPKeepAlive              string                   `json:"tcpKeepAlive,omitempty"`
	IsIncompleteDRConfig      bool                     `json:"isIncompleteDRConfig"`
	UseInDrMode               bool                     `json:"useInDrMode"`
	InspectTrafficWithZia     bool                     `json:"inspectTrafficWithZia"`
	MicroTenantID             string                   `json:"microtenantId,omitempty"`
	MicroTenantName           string                   `json:"microtenantName,omitempty"`
	HealthReporting           string                   `json:"healthReporting,omitempty"`
	ICMPAccessType            string                   `json:"icmpAccessType,omitempty"`
	CreationTime              string                   `json:"creationTime,omitempty"`
	ModifiedBy                string                   `json:"modifiedBy,omitempty"`
	ModifiedTime              string                   `json:"modifiedTime,omitempty"`
	TCPPortRanges             []string                 `json:"tcpPortRanges,omitempty"`
	UDPPortRanges             []string                 `json:"udpPortRanges,omitempty"`
	TCPAppPortRange           []common.NetworkPorts    `json:"tcpPortRange,omitempty"`
	UDPAppPortRange           []common.NetworkPorts    `json:"udpPortRange,omitempty"`
	ClientlessApps            []ClientlessApps         `json:"clientlessApps,omitempty"`
	AppServerGroups           []AppServerGroups        `json:"serverGroups,omitempty"`
	SharedMicrotenantDetails  SharedMicrotenantDetails `json:"sharedMicrotenantDetails,omitempty"`
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

type ClientlessApps struct {
	AllowOptions        bool   `json:"allowOptions"`
	AppID               string `json:"appId,omitempty"`
	ApplicationPort     string `json:"applicationPort,omitempty"`
	ApplicationProtocol string `json:"applicationProtocol,omitempty"`
	CertificateID       string `json:"certificateId,omitempty"`
	CertificateName     string `json:"certificateName,omitempty"`
	Cname               string `json:"cname,omitempty"`
	CreationTime        string `json:"creationTime,omitempty"`
	Description         string `json:"description,omitempty"`
	Domain              string `json:"domain,omitempty"`
	Enabled             bool   `json:"enabled"`
	Hidden              bool   `json:"hidden"`
	ID                  string `json:"id,omitempty"`
	LocalDomain         string `json:"localDomain,omitempty"`
	ModifiedBy          string `json:"modifiedBy,omitempty"`
	ModifiedTime        string `json:"modifiedTime,omitempty"`
	Name                string `json:"name,omitempty"`
	Path                string `json:"path,omitempty"`
	MicroTenantID       string `json:"microtenantId,omitempty"`
	MicroTenantName     string `json:"microtenantName,omitempty"`
	TrustUntrustedCert  bool   `json:"trustUntrustedCert"`
}

type AppServerGroups struct {
	ID string `json:"id"`
}

func (service *Service) Get(id string) (*BrowserAccess, *http.Response, error) {
	v := new(BrowserAccess)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+browserAccessEndpoint, id)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) GetByName(BaName string) (*BrowserAccess, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + browserAccessEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[BrowserAccess](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, BaName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no browser access application named '%s' was found", BaName)
}

func (service *Service) GetByApplicationType(applicationType string, expandAll bool) ([]BrowserAccess, *http.Response, error) {
	if applicationType != "BROWSER_ACCESS" && applicationType != "SECURE_REMOTE_ACCESS" && applicationType != "INSPECT" {
		return nil, nil, fmt.Errorf("invalid applicationType '%s'. Valid types are 'BROWSER_ACCESS', 'SECURE_REMOTE_ACCESS', 'INSPECT'", applicationType)
	}
	// Constructing the query parameters as part of the URL
	relativeURL := fmt.Sprintf("%s%s%s?applicationType=%s&expandAll=%t&page=1&pagesize=20",
		mgmtConfig, service.Client.Config.CustomerID, applicationTypeEndpoint, applicationType, expandAll)
	filter := common.Filter{} // Initialize an empty filter or with minimal required fields

	list, resp, err := common.GetAllPagesGenericWithCustomFilters[BrowserAccess](service.Client, relativeURL, filter)
	if err != nil {
		return nil, nil, err
	}

	return list, resp, nil
}

func (service *Service) Create(browserAccess BrowserAccess) (*BrowserAccess, *http.Response, error) {
	v := new(BrowserAccess)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+browserAccessEndpoint, common.Filter{MicroTenantID: service.microTenantID}, browserAccess, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(id string, browserAccess *BrowserAccess) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+browserAccessEndpoint, id)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.microTenantID}, browserAccess, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(id string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+browserAccessEndpoint, id)
	resp, err := service.Client.NewRequestDo("DELETE", path, common.DeleteApplicationQueryParams{ForceDelete: true, MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) GetAll() ([]BrowserAccess, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + browserAccessEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[BrowserAccess](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	result := []BrowserAccess{}
	// filter browser access apps
	for _, item := range list {
		if len(item.ClientlessApps) > 0 {
			result = append(result, item)
		}
	}
	return result, resp, nil
}
