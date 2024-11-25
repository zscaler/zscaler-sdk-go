package applicationsegmentpra

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"
)

const (
	mgmtConfig              = "/zpa/mgmtconfig/v1/admin/customers/"
	appSegmentPraEndpoint   = "/application"
	applicationTypeEndpoint = "/application/getAppsByType"
)

type AppSegmentPRA struct {
	ID                        string                    `json:"id,omitempty"`
	DomainNames               []string                  `json:"domainNames,omitempty"`
	Name                      string                    `json:"name,omitempty"`
	Description               string                    `json:"description,omitempty"`
	Enabled                   bool                      `json:"enabled"`
	PassiveHealthEnabled      bool                      `json:"passiveHealthEnabled"`
	SelectConnectorCloseToApp bool                      `json:"selectConnectorCloseToApp"`
	DoubleEncrypt             bool                      `json:"doubleEncrypt"`
	AppRecommendationId       string                    `json:"appRecommendationId,omitempty"`
	ConfigSpace               string                    `json:"configSpace,omitempty"`
	Applications              string                    `json:"applications,omitempty"`
	BypassType                string                    `json:"bypassType,omitempty"`
	MatchStyle                string                    `json:"matchStyle,omitempty"`
	BypassOnReauth            bool                      `json:"bypassOnReauth,omitempty"`
	FQDNDnsCheck              bool                      `json:"fqdnDnsCheck"`
	HealthCheckType           string                    `json:"healthCheckType,omitempty"`
	IsCnameEnabled            bool                      `json:"isCnameEnabled"`
	IpAnchored                bool                      `json:"ipAnchored"`
	HealthReporting           string                    `json:"healthReporting,omitempty"`
	IcmpAccessType            string                    `json:"icmpAccessType,omitempty"`
	SegmentGroupID            string                    `json:"segmentGroupId"`
	SegmentGroupName          string                    `json:"segmentGroupName,omitempty"`
	CreationTime              string                    `json:"creationTime,omitempty"`
	ModifiedBy                string                    `json:"modifiedBy,omitempty"`
	ModifiedTime              string                    `json:"modifiedTime,omitempty"`
	TCPKeepAlive              string                    `json:"tcpKeepAlive,omitempty"`
	IsIncompleteDRConfig      bool                      `json:"isIncompleteDRConfig"`
	UseInDrMode               bool                      `json:"useInDrMode"`
	MicroTenantID             string                    `json:"microtenantId,omitempty"`
	MicroTenantName           string                    `json:"microtenantName,omitempty"`
	TCPAppPortRange           []common.NetworkPorts     `json:"tcpPortRange,omitempty"`
	UDPAppPortRange           []common.NetworkPorts     `json:"udpPortRange,omitempty"`
	TCPPortRanges             []string                  `json:"tcpPortRanges,omitempty"`
	UDPPortRanges             []string                  `json:"udpPortRanges,omitempty"`
	ServerGroups              []servergroup.ServerGroup `json:"serverGroups,omitempty"`
	DefaultIdleTimeout        string                    `json:"defaultIdleTimeout,omitempty"`
	DefaultMaxAge             string                    `json:"defaultMaxAge,omitempty"`
	PRAApps                   []PRAApps                 `json:"praApps"`
	CommonAppsDto             CommonAppsDto             `json:"commonAppsDto"`
	SharedMicrotenantDetails  SharedMicrotenantDetails  `json:"sharedMicrotenantDetails,omitempty"`
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
	DeletedPraApps []string     `json:"deletedPraApps,omitempty"`
}

type AppsConfig struct {
	ID                  string   `json:"id,omitempty"`
	AppID               string   `json:"appId"`
	PRAAppID            string   `json:"praAppId"`
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

type PRAApps struct {
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	AppID               string `json:"appId"`
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

func Get(ctx context.Context, service *zscaler.Service, id string) (*AppSegmentPRA, *http.Response, error) {
	v := new(AppSegmentPRA)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+appSegmentPraEndpoint, id)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, praName string) (*AppSegmentPRA, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + appSegmentPraEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AppSegmentPRA](ctx, service.Client, relativeURL, common.Filter{Search: praName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, praName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no pra application named '%s' was found", praName)
}

func Create(ctx context.Context, service *zscaler.Service, appSegmentPra AppSegmentPRA) (*AppSegmentPRA, *http.Response, error) {
	v := new(AppSegmentPRA)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+appSegmentPraEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, appSegmentPra, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, id string, appSegmentPra *AppSegmentPRA) (*http.Response, error) {
	// Step 1: Retrieve the existing resource to get current `appId` and `praAppId`
	existingResource, _, err := Get(ctx, service, id)
	if err != nil {
		return nil, err
	}

	// Set the primary app ID
	appSegmentPra.ID = existingResource.ID

	// Step 2: Map existing `praApps` entries by `Name` to get `praAppId` for each sub-application
	existingPraApps := make(map[string]PRAApps)
	for _, praApp := range existingResource.PRAApps {
		existingPraApps[praApp.Name] = praApp
	}

	// Step 3: Inject `appId` and `praAppId` into each entry in `appsConfig`
	for i, appConfig := range appSegmentPra.CommonAppsDto.AppsConfig {
		if existingApp, ok := existingPraApps[appConfig.Name]; ok {
			appSegmentPra.CommonAppsDto.AppsConfig[i].AppID = existingResource.ID // main app ID
			appSegmentPra.CommonAppsDto.AppsConfig[i].PRAAppID = existingApp.ID   // praAppId for sub-app
		}
	}

	// Check if `commonAppsDto` actually has entries, set to nil if empty
	if len(appSegmentPra.CommonAppsDto.AppsConfig) == 0 {
		appSegmentPra.CommonAppsDto = CommonAppsDto{}
	}
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+appSegmentPraEndpoint, id)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, appSegmentPra, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, id string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+appSegmentPraEndpoint, id)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.DeleteApplicationQueryParams{ForceDelete: true, MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]AppSegmentPRA, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + appSegmentPraEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AppSegmentPRA](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	result := []AppSegmentPRA{}
	// filter pra apps
	for _, item := range list {
		if len(item.PRAApps) > 0 {
			result = append(result, item)
		}
	}
	return result, resp, nil
}
