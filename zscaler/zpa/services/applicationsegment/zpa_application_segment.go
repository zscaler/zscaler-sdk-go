package applicationsegment

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentbrowseraccess"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"
)

const (
	mgmtConfig         = "/zpa/mgmtconfig/v1/admin/customers/"
	appSegmentEndpoint = "/application"
)

type ApplicationSegmentResource struct {
	ID                        string                    `json:"id,omitempty"`
	DomainNames               []string                  `json:"domainNames,omitempty"`
	Name                      string                    `json:"name,omitempty"`
	Description               string                    `json:"description,omitempty"`
	Enabled                   bool                      `json:"enabled"`
	ExtranetEnabled           bool                      `json:"extranetEnabled"`
	APIProtectionEnabled      bool                      `json:"apiProtectionEnabled"`
	AutoAppProtectEnabled     bool                      `json:"autoAppProtectEnabled"`
	ADPEnabled                bool                      `json:"adpEnabled"`
	PassiveHealthEnabled      bool                      `json:"passiveHealthEnabled"`
	DoubleEncrypt             bool                      `json:"doubleEncrypt"`
	ConfigSpace               string                    `json:"configSpace,omitempty"`
	Applications              string                    `json:"applications,omitempty"`
	BypassType                string                    `json:"bypassType,omitempty"`
	BypassOnReauth            bool                      `json:"bypassOnReauth,omitempty"`
	HealthCheckType           string                    `json:"healthCheckType,omitempty"`
	IsCnameEnabled            bool                      `json:"isCnameEnabled"`
	IpAnchored                bool                      `json:"ipAnchored"`
	FQDNDnsCheck              bool                      `json:"fqdnDnsCheck"`
	HealthReporting           string                    `json:"healthReporting,omitempty"`
	SelectConnectorCloseToApp bool                      `json:"selectConnectorCloseToApp"`
	IcmpAccessType            string                    `json:"icmpAccessType,omitempty"`
	AppRecommendationId       string                    `json:"appRecommendationId,omitempty"`
	SegmentGroupID            string                    `json:"segmentGroupId"`
	SegmentGroupName          string                    `json:"segmentGroupName,omitempty"`
	CreationTime              string                    `json:"creationTime,omitempty"`
	ModifiedBy                string                    `json:"modifiedBy,omitempty"`
	ModifiedTime              string                    `json:"modifiedTime,omitempty"`
	TCPKeepAlive              string                    `json:"tcpKeepAlive,omitempty"`
	IsIncompleteDRConfig      bool                      `json:"isIncompleteDRConfig"`
	UseInDrMode               bool                      `json:"useInDrMode"`
	InspectTrafficWithZia     bool                      `json:"inspectTrafficWithZia"`
	WeightedLoadBalancing     bool                      `json:"weightedLoadBalancing"`
	MicroTenantID             string                    `json:"microtenantId,omitempty"`
	MicroTenantName           string                    `json:"microtenantName,omitempty"`
	MatchStyle                string                    `json:"matchStyle,omitempty"`
	ReadOnly                  bool                      `json:"readOnly,omitempty"`
	RestrictionType           string                    `json:"restrictionType,omitempty"`
	ZscalerManaged            bool                      `json:"zscalerManaged,omitempty"`
	TCPPortRanges             []string                  `json:"tcpPortRanges"`
	UDPPortRanges             []string                  `json:"udpPortRanges"`
	TCPAppPortRange           []common.NetworkPorts     `json:"tcpPortRange,omitempty"`
	UDPAppPortRange           []common.NetworkPorts     `json:"udpPortRange,omitempty"`
	ServerGroups              []servergroup.ServerGroup `json:"serverGroups"`
	DefaultIdleTimeout        string                    `json:"defaultIdleTimeout,omitempty"`
	DefaultMaxAge             string                    `json:"defaultMaxAge,omitempty"`
	// CommonAppsDto             applicationsegmentpra.CommonAppsDto              `json:"commonAppsDto,omitempty"`
	ClientlessApps           []applicationsegmentbrowseraccess.ClientlessApps `json:"clientlessApps,omitempty"`
	ShareToMicrotenants      []string                                         `json:"shareToMicrotenants"`
	SharedMicrotenantDetails SharedMicrotenantDetails                         `json:"sharedMicrotenantDetails,omitempty"`
	ZPNERID                  *common.ZPNERID                                  `json:"zpnErId"`
	Tags                     []Tag                                            `json:"tags,omitempty"`
	PolicyStyle              string                                           `json:"policyStyle,omitempty"`
}

type SharedMicrotenantDetails struct {
	SharedFromMicrotenant SharedFromMicrotenant `json:"sharedFromMicrotenant,omitempty"`
	SharedToMicrotenants  []SharedToMicrotenant `json:"sharedToMicrotenants,omitempty"`
}

// Tag represents a tag associated with an application segment
type Tag struct {
	Namespace common.CommonSummary `json:"namespace,omitempty"`
	TagKey    common.CommonSummary `json:"tagKey,omitempty"`
	TagValue  common.CommonIDName  `json:"tagValue,omitempty"`
	Origin    string               `json:"origin,omitempty"`
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

// MultiMatchUnsupportedReferencesPayload represents the payload for GetMultiMatchUnsupportedReferences
type MultiMatchUnsupportedReferencesPayload []string

// MultiMatchUnsupportedReferencesResponse represents the response from GetMultiMatchUnsupportedReferences
type MultiMatchUnsupportedReferencesResponse struct {
	ID              string   `json:"id"`
	AppSegmentName  string   `json:"appSegmentName"`
	Domains         []string `json:"domains"`
	TCPPorts        []string `json:"tcpPorts"`
	MatchStyle      string   `json:"matchStyle"`
	MicrotenantName string   `json:"microtenantName"`
}

// BulkUpdateMultiMatchPayload represents the payload for UpdatebulkUpdateMultiMatch
type BulkUpdateMultiMatchPayload struct {
	ApplicationIDs []int  `json:"applicationIds"`
	MatchStyle     string `json:"matchStyle"`
}

type ApplicationCountResponse struct {
	AppsConfigured               string `json:"appsConfigured"`
	ConfiguredDateInEpochSeconds string `json:"configuredDateInEpochSeconds"`
}

// ApplicationCurrentMaxLimitResponse represents the response from GetCurrentAndMaxLimit
type ApplicationCurrentMaxLimitResponse struct {
	CurrentAppsCount string `json:"currentAppsCount"`
	MaxAppsLimit     string `json:"maxAppsLimit"`
}

// ApplicationValidationError represents the validation error response
type ApplicationValidationError struct {
	Params []string `json:"params"`
	ID     string   `json:"id"`
	Reason string   `json:"reason"`
}

func Get(ctx context.Context, service *zscaler.Service, applicationID string) (*ApplicationSegmentResource, *http.Response, error) {
	v := new(ApplicationSegmentResource)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+appSegmentEndpoint, applicationID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, appName string) (*ApplicationSegmentResource, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + appSegmentEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ApplicationSegmentResource](ctx, service.Client, relativeURL, common.Filter{Search: appName, MicroTenantID: service.MicroTenantID()})
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

func Create(ctx context.Context, service *zscaler.Service, appSegment ApplicationSegmentResource) (*ApplicationSegmentResource, *http.Response, error) {
	v := new(ApplicationSegmentResource)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+appSegmentEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, appSegment, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, appID string, appSegmentRequest ApplicationSegmentResource) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+appSegmentEndpoint, appID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, appSegmentRequest, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, appID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+appSegmentEndpoint, appID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", relativeURL, common.DeleteApplicationQueryParams{ForceDelete: true, MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]ApplicationSegmentResource, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + appSegmentEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ApplicationSegmentResource](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func GetMultiMatchUnsupportedReferences(ctx context.Context, service *zscaler.Service, domainNames MultiMatchUnsupportedReferencesPayload) ([]MultiMatchUnsupportedReferencesResponse, *http.Response, error) {
	// Validate that at least one domain name is provided
	if len(domainNames) == 0 {
		return nil, nil, fmt.Errorf("at least one domain name must be provided")
	}

	var v []MultiMatchUnsupportedReferencesResponse
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+appSegmentEndpoint+"/multimatchUnsupportedReferences", common.Filter{MicroTenantID: service.MicroTenantID()}, domainNames, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func UpdatebulkUpdateMultiMatch(ctx context.Context, service *zscaler.Service, payload BulkUpdateMultiMatchPayload) (*http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + appSegmentEndpoint + "/bulkUpdateMultiMatch"
	resp, err := service.Client.NewRequestDo(ctx, "PUT", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, payload, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

// Need to review as the API is returning 400 error
func GetApplicationSummary(ctx context.Context, service *zscaler.Service) ([]common.CommonSummary, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + appSegmentEndpoint + "/summary"
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[common.CommonSummary](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

// ApplicationCountResponse represents the response from GetApplicationCount
func GetApplicationCount(ctx context.Context, service *zscaler.Service) ([]ApplicationCountResponse, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + appSegmentEndpoint + "/configured/count"

	var result []ApplicationCountResponse
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, &result)
	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}

func GetCurrentAndMaxLimit(ctx context.Context, service *zscaler.Service) (*ApplicationCurrentMaxLimitResponse, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + appSegmentEndpoint + "/count/currentAndMaxLimit"

	var result ApplicationCurrentMaxLimitResponse
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, &result)
	if err != nil {
		return nil, nil, err
	}

	return &result, resp, nil
}

type ApplicationMappings struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type WeightedLoadBalancerConfig struct {
	ApplicationID                string                            `json:"applicationId"`
	ApplicationToServerGroupMaps []ApplicationToServerGroupMapping `json:"applicationToServerGroupMappings"`
	WeightedLoadBalancing        bool                              `json:"weightedLoadBalancing"`
}

type ApplicationToServerGroupMapping struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Passive bool   `json:"passive"`
	Weight  string `json:"weight"`
}

func GetApplicationMappings(ctx context.Context, service *zscaler.Service, applicationID string) ([]ApplicationMappings, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + appSegmentEndpoint + "/" + applicationID + "/mappings"

	var result []ApplicationMappings
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, &result)
	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}

func GetApplicationExport(ctx context.Context, service *zscaler.Service, search string, single bool, outputPath string) (*http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + appSegmentEndpoint + "/export"

	// Build query parameters
	query := url.Values{}
	if search != "" {
		query.Set("search", search)
	}
	query.Set("single", fmt.Sprintf("%t", single))

	constructedURL := relativeURL
	if len(query) > 0 {
		constructedURL += "?" + query.Encode()
	}

	// For CSV downloads, we need to handle the raw response body
	resp, err := service.Client.NewRequestDo(ctx, "GET", constructedURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Create the output file
	file, err := os.Create(outputPath)
	if err != nil {
		return resp, fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return resp, fmt.Errorf("failed to write CSV content to file: %w", err)
	}

	service.Client.GetLogger().Printf("[INFO] CSV file saved to: %s", outputPath)
	return resp, nil
}

func ApplicationValidation(ctx context.Context, service *zscaler.Service, appSegment ApplicationSegmentResource) (*ApplicationValidationError, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + appSegmentEndpoint + "/validate"

	// For validation, we expect either success (200) or validation error (400)
	var validationError ApplicationValidationError
	resp, err := service.Client.NewRequestDo(ctx, "POST", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, appSegment, &validationError)
	if err != nil {
		return nil, resp, err
	}

	// If we get here, there was a validation error
	return &validationError, resp, nil
}

func GetWeightedLoadBalancerConfig(ctx context.Context, service *zscaler.Service, applicationID string) (*WeightedLoadBalancerConfig, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s/weightedLbConfig", mgmtConfig+service.Client.GetCustomerID()+appSegmentEndpoint, applicationID)

	result := new(WeightedLoadBalancerConfig)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, result)
	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}

func UpdateWeightedLoadBalancerConfig(ctx context.Context, service *zscaler.Service, applicationID string, payload WeightedLoadBalancerConfig) (*WeightedLoadBalancerConfig, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s/weightedLbConfig", mgmtConfig+service.Client.GetCustomerID()+appSegmentEndpoint, applicationID)

	resp, err := service.Client.NewRequestDo(ctx, "PUT", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, payload, nil)
	if err != nil {
		return nil, nil, err
	}

	return nil, resp, nil
}
