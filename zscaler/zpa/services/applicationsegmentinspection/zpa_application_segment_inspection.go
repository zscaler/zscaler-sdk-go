package applicationsegmentinspection

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"
)

const (
	mgmtConfig                   = "/zpa/mgmtconfig/v1/admin/customers/"
	appSegmentInspectionEndpoint = "/application"
	applicationTypeEndpoint      = "/application/getAppsByType"
)

type AppSegmentInspection struct {
	ID                        string                    `json:"id,omitempty"`
	SegmentGroupID            string                    `json:"segmentGroupId,omitempty"`
	SegmentGroupName          string                    `json:"segmentGroupName,omitempty"`
	BypassType                string                    `json:"bypassType,omitempty"`
	BypassOnReauth            bool                      `json:"bypassOnReauth,omitempty"`
	ConfigSpace               string                    `json:"configSpace,omitempty"`
	DomainNames               []string                  `json:"domainNames,omitempty"`
	Name                      string                    `json:"name,omitempty"`
	Description               string                    `json:"description,omitempty"`
	Enabled                   bool                      `json:"enabled"`
	AdpEnabled                bool                      `json:"adpEnabled,omitempty"`
	AppRecommendationId       string                    `json:"appRecommendationId,omitempty"`
	AutoAppProtectEnabled     bool                      `json:"autoAppProtectEnabled,omitempty"`
	ICMPAccessType            string                    `json:"icmpAccessType,omitempty"`
	PassiveHealthEnabled      bool                      `json:"passiveHealthEnabled,omitempty"`
	FQDNDnsCheck              bool                      `json:"fqdnDnsCheck"`
	MatchStyle                string                    `json:"matchStyle,omitempty"`
	SelectConnectorCloseToApp bool                      `json:"selectConnectorCloseToApp"`
	DoubleEncrypt             bool                      `json:"doubleEncrypt"`
	HealthCheckType           string                    `json:"healthCheckType,omitempty"`
	IsCnameEnabled            bool                      `json:"isCnameEnabled"`
	IPAnchored                bool                      `json:"ipAnchored"`
	HealthReporting           string                    `json:"healthReporting,omitempty"`
	CreationTime              string                    `json:"creationTime,omitempty"`
	ModifiedBy                string                    `json:"modifiedBy,omitempty"`
	ModifiedTime              string                    `json:"modifiedTime,omitempty"`
	TCPKeepAlive              string                    `json:"tcpKeepAlive,omitempty"`
	IsIncompleteDRConfig      bool                      `json:"isIncompleteDRConfig"`
	UseInDrMode               bool                      `json:"useInDrMode"`
	MicroTenantID             string                    `json:"microtenantId,omitempty"`
	MicroTenantName           string                    `json:"microtenantName,omitempty"`
	TCPPortRanges             []string                  `json:"tcpPortRanges,omitempty"`
	UDPPortRanges             []string                  `json:"udpPortRanges,omitempty"`
	TCPAppPortRange           []common.NetworkPorts     `json:"tcpPortRange,omitempty"`
	UDPAppPortRange           []common.NetworkPorts     `json:"udpPortRange,omitempty"`
	TCPProtocols              []string                  `json:"tcpProtocols"`
	UDPProtocols              []string                  `json:"udpProtocols,omitempty"`
	InspectionAppDto          []InspectionAppDto        `json:"inspectionApps,omitempty"`
	CommonAppsDto             CommonAppsDto             `json:"commonAppsDto,omitempty"`
	AppServerGroups           []servergroup.ServerGroup `json:"serverGroups,omitempty"`
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
	AppsConfig         []AppsConfig `json:"appsConfig,omitempty"`
	DeletedInspectApps []string     `json:"deletedInspectApps,omitempty"`
}

type AppsConfig struct {
	ID                  string   `json:"id,omitempty"`
	AppID               string   `json:"appId,omitempty"`
	InspectAppID        string   `json:"inspectAppId"`
	Name                string   `json:"name,omitempty"`
	Description         string   `json:"description,omitempty"`
	Enabled             bool     `json:"enabled"`
	AdpEnabled          bool     `json:"adpEnabled"`
	AllowOptions        bool     `json:"allowOptions"`
	AppTypes            []string `json:"appTypes,omitempty"`
	ApplicationPort     string   `json:"applicationPort,omitempty"`
	ApplicationProtocol string   `json:"applicationProtocol,omitempty"`
	Protocols           []string `json:"protocols,omitempty"`
	CertificateID       string   `json:"certificateId,omitempty"`
	CertificateName     string   `json:"certificateName,omitempty"`
	Cname               string   `json:"cname,omitempty"`
	Domain              string   `json:"domain,omitempty"`
	Hidden              bool     `json:"hidden"`
	TrustUntrustedCert  bool     `json:"trustUntrustedCert"`
	LocalDomain         string   `json:"localDomain,omitempty"`
	Portal              bool     `json:"portal"`
}

type InspectionAppDto struct {
	ID                  string   `json:"id,omitempty"`
	AppID               string   `json:"appId,omitempty"`
	Name                string   `json:"name,omitempty"`
	Description         string   `json:"description,omitempty"`
	Enabled             bool     `json:"enabled"`
	ApplicationPort     string   `json:"applicationPort,omitempty"`
	ApplicationProtocol string   `json:"applicationProtocol,omitempty"`
	CertificateID       string   `json:"certificateId,omitempty"`
	CertificateName     string   `json:"certificateName,omitempty"`
	Domain              string   `json:"domain,omitempty"`
	Protocols           []string `json:"protocols,omitempty"`
	TrustUntrustedCert  bool     `json:"trustUntrustedCert"`
	MicroTenantID       string   `json:"microtenantId,omitempty"`
	MicroTenantName     string   `json:"microtenantName,omitempty"`
}

type AppServerGroups struct {
	ID string `json:"id"`
}

func Get(ctx context.Context, service *zscaler.Service, id string) (*AppSegmentInspection, *http.Response, error) {
	v := new(AppSegmentInspection)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+appSegmentInspectionEndpoint, id)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, appSegmentName string) (*AppSegmentInspection, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + appSegmentInspectionEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AppSegmentInspection](ctx, service.Client, relativeURL, common.Filter{Search: appSegmentName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, appSegmentName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no inspection application segment named '%s' was found", appSegmentName)
}

func Create(ctx context.Context, service *zscaler.Service, appSegmentInspection AppSegmentInspection) (*AppSegmentInspection, *http.Response, error) {
	v := new(AppSegmentInspection)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+appSegmentInspectionEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, appSegmentInspection, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, id string, appSegmentInspection *AppSegmentInspection) (*http.Response, error) {
	// Step 1: Retrieve the existing resource to get current `appId` and `InspectAppID`
	existingResource, _, err := Get(ctx, service, id)
	if err != nil {
		return nil, err
	}

	// Set the primary app ID
	appSegmentInspection.ID = existingResource.ID

	// Step 2: Map existing `inspectionApp` entries by `Name` to get `InspectAppID` for each sub-application
	existingInspectionApps := make(map[string]InspectionAppDto)
	for _, inspectionApp := range existingResource.InspectionAppDto {
		existingInspectionApps[inspectionApp.Name] = inspectionApp
	}

	// Step 3: Inject `appId` and `InspectAppID` into each entry in `appsConfig`
	for i, appConfig := range appSegmentInspection.CommonAppsDto.AppsConfig {
		if existingApp, ok := existingInspectionApps[appConfig.Name]; ok {
			appSegmentInspection.CommonAppsDto.AppsConfig[i].AppID = existingResource.ID   // main app ID
			appSegmentInspection.CommonAppsDto.AppsConfig[i].InspectAppID = existingApp.ID // InspectAppID for sub-app
		}
	}

	// Check if `commonAppsDto` actually has entries, set to nil if empty
	if len(appSegmentInspection.CommonAppsDto.AppsConfig) == 0 {
		appSegmentInspection.CommonAppsDto = CommonAppsDto{}
	}

	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+appSegmentInspectionEndpoint, id)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, appSegmentInspection, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, id string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+appSegmentInspectionEndpoint, id)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.DeleteApplicationQueryParams{ForceDelete: true, MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]AppSegmentInspection, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + appSegmentInspectionEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AppSegmentInspection](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	result := []AppSegmentInspection{}
	// filter pra apps
	for _, item := range list {
		if len(item.InspectionAppDto) > 0 {
			result = append(result, item)
		}
	}
	return result, resp, nil
}
