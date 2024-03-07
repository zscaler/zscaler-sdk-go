package applicationsegmentpra

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig              = "/mgmtconfig/v1/admin/customers/"
	appSegmentPraEndpoint   = "/application"
	applicationTypeEndpoint = "/application/getAppsByType"
)

type AppSegmentPRA struct {
	ID                        string                   `json:"id,omitempty"`
	DomainNames               []string                 `json:"domainNames,omitempty"`
	Name                      string                   `json:"name,omitempty"`
	Description               string                   `json:"description,omitempty"`
	Enabled                   bool                     `json:"enabled"`
	PassiveHealthEnabled      bool                     `json:"passiveHealthEnabled"`
	SelectConnectorCloseToApp bool                     `json:"selectConnectorCloseToApp"`
	DoubleEncrypt             bool                     `json:"doubleEncrypt"`
	AppRecommendationId       string                   `json:"appRecommendationId,omitempty"`
	ConfigSpace               string                   `json:"configSpace,omitempty"`
	Applications              string                   `json:"applications,omitempty"`
	BypassType                string                   `json:"bypassType,omitempty"`
	BypassOnReauth            bool                     `json:"bypassOnReauth,omitempty"`
	FQDNDnsCheck              bool                     `json:"fqdnDnsCheck"`
	HealthCheckType           string                   `json:"healthCheckType,omitempty"`
	IsCnameEnabled            bool                     `json:"isCnameEnabled"`
	IpAnchored                bool                     `json:"ipAnchored"`
	HealthReporting           string                   `json:"healthReporting,omitempty"`
	IcmpAccessType            string                   `json:"icmpAccessType,omitempty"`
	SegmentGroupID            string                   `json:"segmentGroupId"`
	SegmentGroupName          string                   `json:"segmentGroupName,omitempty"`
	CreationTime              string                   `json:"creationTime,omitempty"`
	ModifiedBy                string                   `json:"modifiedBy,omitempty"`
	ModifiedTime              string                   `json:"modifiedTime,omitempty"`
	TCPKeepAlive              string                   `json:"tcpKeepAlive,omitempty"`
	IsIncompleteDRConfig      bool                     `json:"isIncompleteDRConfig"`
	UseInDrMode               bool                     `json:"useInDrMode"`
	MicroTenantID             string                   `json:"microtenantId,omitempty"`
	MicroTenantName           string                   `json:"microtenantName,omitempty"`
	TCPAppPortRange           []common.NetworkPorts    `json:"tcpPortRange,omitempty"`
	UDPAppPortRange           []common.NetworkPorts    `json:"udpPortRange,omitempty"`
	TCPPortRanges             []string                 `json:"tcpPortRanges,omitempty"`
	UDPPortRanges             []string                 `json:"udpPortRanges,omitempty"`
	ServerGroups              []AppServerGroups        `json:"serverGroups,omitempty"`
	DefaultIdleTimeout        string                   `json:"defaultIdleTimeout,omitempty"`
	DefaultMaxAge             string                   `json:"defaultMaxAge,omitempty"`
	SRAAppsDto                []SRAAppsDto             `json:"praApps"`
	CommonAppsDto             CommonAppsDto            `json:"commonAppsDto,omitempty"`
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
type CommonAppsDto struct {
	AppsConfig     []AppsConfig `json:"appsConfig,omitempty"`
	DeletedSraApps []string     `json:"deletedSraApps,omitempty"`
}

type AppsConfig struct {
	ID                  string   `json:"id,omitempty"`
	AppID               string   `json:"appId,omitempty"`
	Name                string   `json:"name,omitempty"`
	Description         string   `json:"description,omitempty"`
	Enabled             bool     `json:"enabled,omitempty"`
	AppTypes            []string `json:"appTypes,omitempty"`
	ApplicationPort     string   `json:"applicationPort,omitempty"`
	ApplicationProtocol string   `json:"applicationProtocol,omitempty"`
	Cname               string   `json:"cname,omitempty"`
	ConnectionSecurity  string   `json:"connectionSecurity,omitempty"`
	Domain              string   `json:"domain,omitempty"`
	Hidden              bool     `json:"hidden,omitempty"`
	LocalDomain         string   `json:"localDomain,omitempty"`
	Portal              bool     `json:"portal,omitempty"`
}

type SRAAppsDto struct {
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	AppID               string `json:"appId,omitempty"`
	ApplicationPort     string `json:"applicationPort,omitempty"`
	ApplicationProtocol string `json:"applicationProtocol,omitempty"`
	CertificateID       string `json:"certificateId,omitempty"`
	CertificateName     string `json:"certificateName,omitempty"`
	ConnectionSecurity  string `json:"connectionSecurity,omitempty"`
	Hidden              bool   `json:"hidden"`
	Portal              bool   `json:"portal"`
	Description         string `json:"description,omitempty"`
	Domain              string `json:"domain,omitempty"`
	Enabled             bool   `json:"enabled"`
	MicroTenantID       string `json:"microtenantId,omitempty"`
	MicroTenantName     string `json:"microtenantName,omitempty"`
}

type AppServerGroups struct {
	ID string `json:"id"`
}

func (service *Service) Get(id string) (*AppSegmentPRA, *http.Response, error) {
	v := new(AppSegmentPRA)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentPraEndpoint, id)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) GetByName(BaName string) (*AppSegmentPRA, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appSegmentPraEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AppSegmentPRA](service.Client, relativeURL, common.Filter{Search: BaName, MicroTenantID: service.microTenantID})
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

func (service *Service) GetByApplicationType(applicationType string, expandAll bool) ([]AppSegmentPRA, *http.Response, error) {
	if applicationType != "BROWSER_ACCESS" && applicationType != "SECURE_REMOTE_ACCESS" && applicationType != "INSPECT" {
		return nil, nil, fmt.Errorf("invalid applicationType '%s'. Valid types are 'BROWSER_ACCESS', 'SECURE_REMOTE_ACCESS', 'INSPECT'", applicationType)
	}
	// Constructing the query parameters as part of the URL
	relativeURL := fmt.Sprintf("%s%s%s?applicationType=%s&expandAll=%t&page=1&pagesize=20",
		mgmtConfig, service.Client.Config.CustomerID, applicationTypeEndpoint, applicationType, expandAll)
	filter := common.Filter{} // Initialize an empty filter or with minimal required fields

	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AppSegmentPRA](service.Client, relativeURL, filter)
	if err != nil {
		return nil, nil, err
	}

	return list, resp, nil
}

func (service *Service) Create(appSegmentPra AppSegmentPRA) (*AppSegmentPRA, *http.Response, error) {
	v := new(AppSegmentPRA)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+appSegmentPraEndpoint, common.Filter{MicroTenantID: service.microTenantID}, appSegmentPra, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// return the new items that were added to slice1
func difference(slice1 []AppsConfig, slice2 []AppsConfig) []AppsConfig {
	var diff []AppsConfig
	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1.Domain == s2.Domain || s1.Name == s2.Name {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, s1)
		}
	}
	return diff
}

func mapSraApp(SRAAppsDto []SRAAppsDto) []AppsConfig {
	result := []AppsConfig{}
	for _, app := range SRAAppsDto {
		result = append(result, AppsConfig{
			Name:   app.Name,
			Domain: app.Domain,
			ID:     app.ID,
			AppID:  app.AppID,
		})
	}
	return result
}

func appToListStringIDs(apps []AppsConfig) []string {
	result := []string{}
	for _, app := range apps {
		result = append(result, app.ID)
	}
	return result
}

func (service *Service) Update(id string, appSegmentPra *AppSegmentPRA) (*http.Response, error) {
	existingResource, _, err := service.Get(id)
	if err != nil {
		return nil, err
	}
	existingApps := mapSraApp(existingResource.SRAAppsDto)
	newApps := difference(appSegmentPra.CommonAppsDto.AppsConfig, existingApps)
	removedApps := difference(existingApps, appSegmentPra.CommonAppsDto.AppsConfig)
	appSegmentPra.CommonAppsDto.AppsConfig = newApps
	appSegmentPra.CommonAppsDto.DeletedSraApps = appToListStringIDs(removedApps)
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentPraEndpoint, id)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.microTenantID}, appSegmentPra, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(id string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appSegmentPraEndpoint, id)
	resp, err := service.Client.NewRequestDo("DELETE", path, common.DeleteApplicationQueryParams{ForceDelete: true, MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) GetAll() ([]AppSegmentPRA, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appSegmentPraEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AppSegmentPRA](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	result := []AppSegmentPRA{}
	// filter pra apps
	for _, item := range list {
		if len(item.CommonAppsDto.AppsConfig) == 0 || (!common.InList(item.CommonAppsDto.AppsConfig[0].AppTypes, "SECURE_REMOTE_ACCESS") && !common.InList(item.CommonAppsDto.AppsConfig[0].AppTypes, "INSPECT")) {
			result = append(result, item)
		}
	}
	return result, resp, nil
}
