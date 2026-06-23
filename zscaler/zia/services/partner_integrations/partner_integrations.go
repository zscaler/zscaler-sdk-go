package partner_integrations

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	partnerIntegrationEndpoint = "/zia/api/v1/integrationPartners"
)

type PartnerIntegrations struct {
	Id       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

func GetPartnerIntegrations(ctx context.Context, service *zscaler.Service) (*PartnerIntegrations, error) {
	var partnerIntegrations PartnerIntegrations
	err := service.Client.Read(ctx, partnerIntegrationEndpoint, &partnerIntegrations)
	if err != nil {
		return nil, err
	}
	return &partnerIntegrations, nil
}

// ////////////////////////////////////
// CROWDSTRIKE PARTNER INTEGRATIONS //
// ////////////////////////////////////
type Crowdstrike struct {
	CrowdStrikeResponse   []CrowdStrikeResponse `json:"crowdStrikeResponse,omitempty"`
	CrowdStrikePagination CrowdStrikePagination `json:"crowdStrikePagination,omitempty"`
	CrowdStrikeErrors     []CrowdStrikeErrors   `json:"crowdStrikeErrors,omitempty"`
}

type CrowdStrikeResponse struct {
	EndPointLink      string `json:"endPointLink,omitempty"`
	DeviceId          string `json:"device_id,omitempty"`
	SystemProductName string `json:"system_product_name,omitempty"`
	Hostname          string `json:"hostname,omitempty"`
	LocalIp           string `json:"local_ip,omitempty"`
	ExternalIp        string `json:"external_ip,omitempty"`
	MacAddress        string `json:"mac_address,omitempty"`
	OsVersion         string `json:"os_version,omitempty"`
	Status            string `json:"status,omitempty"`
	FileStatus        string `json:"file_status,omitempty"`
	PlatformName      string `json:"platform_name,omitempty"`
	FirstSeen         string `json:"first_seen,omitempty"`
	LastSeen          string `json:"last_seen,omitempty"`
}

type CrowdStrikePagination struct {
	Offset   string `json:"offset,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Total    int    `json:"total,omitempty"`
	NextPage string `json:"next_page,omitempty"`
}

type CrowdStrikeErrors struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

type GetAllFilterOptions struct {
	// Filters based on the IOC type. Supported value: MD5
	Type *string

	// Filters based on the IOC value
	Value *string

	// Specifies the page size
	Limit *int

	// Specifies the page offset
	Offset *string

	// Filters based on the partner JSON type
	// Supported Values: CROWDSTRIKE_CREDENTIALS, CARBON_BLACK_CREDENTIALS, ATP_DEFENDER_CREDENTIALS, UNIT_TESTING_CS, UNIT_TESTING_CB
	PartnerJsonType *string
}

// GetCrowdstrikeEndpoints retrieves the list of CrowdStrike endpoints based on
// the indicator of compromise (IOC) query.
//
// The endpoint returns a single envelope object containing the matching
// endpoints (crowdStrikeResponse), cursor-based pagination metadata
// (crowdStrikePagination), and any partner errors (crowdStrikeErrors). It is not
// a bare paginated array, so the standard page/pageSize pagination helper does
// not apply. Pagination is driven by the optional limit/offset parameters and
// the offset/next_page values returned in crowdStrikePagination.
func GetCrowdstrikeEndpoints(ctx context.Context, service *zscaler.Service, opts *GetAllFilterOptions) (*Crowdstrike, error) {
	var crowdstrike Crowdstrike
	endpoint := partnerIntegrationEndpoint + "/crowdStrike/endpoints"

	queryParams := url.Values{}
	if opts != nil {
		if opts.Type != nil && *opts.Type != "" {
			queryParams.Set("type", *opts.Type)
		}
		if opts.Value != nil && *opts.Value != "" {
			queryParams.Set("value", *opts.Value)
		}
		if opts.Limit != nil {
			queryParams.Set("limit", strconv.Itoa(*opts.Limit))
		}
		if opts.Offset != nil && *opts.Offset != "" {
			queryParams.Set("offset", *opts.Offset)
		}
		if opts.PartnerJsonType != nil && *opts.PartnerJsonType != "" {
			queryParams.Set("partnerJsonType", *opts.PartnerJsonType)
		}
	}
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	if err := service.Client.Read(ctx, endpoint, &crowdstrike); err != nil {
		return nil, err
	}
	return &crowdstrike, nil
}

func AcceptCrowdstrikeEndpointList(ctx context.Context, service *zscaler.Service, endpointList *Crowdstrike) (*Crowdstrike, *http.Response, error) {
	resp, err := service.Client.Create(ctx, partnerIntegrationEndpoint+"/crowdStrike/endpoints", *endpointList)
	if err != nil {
		return nil, nil, err
	}

	endpointList, ok := resp.(*Crowdstrike)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a crowdstrike endpoint list pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new crowdstrike endpoint list from create")
	return endpointList, nil, nil
}

func GetCrowdstrikeWhiteListedBaseURL(ctx context.Context, service *zscaler.Service, opts *GetAllFilterOptions) (*Crowdstrike, error) {
	var crowdstrike Crowdstrike
	endpoint := partnerIntegrationEndpoint + "/crowdStrike/whitelistedBaseUrls"

	queryParams := url.Values{}
	if opts != nil {

		if opts.PartnerJsonType != nil && *opts.PartnerJsonType != "" {
			queryParams.Set("partnerJsonType", *opts.PartnerJsonType)
		}
	}
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	if err := service.Client.Read(ctx, endpoint, &crowdstrike); err != nil {
		return nil, err
	}
	return &crowdstrike, nil
}

// //////////////////////////////////////////
// Microsoft Defender Partner Integration///
// //////////////////////////////////////////
type MicrosoftDefender struct {
	SHA256     string `json:"sha256,omitempty"`
	SHA1       string `json:"sha1,omitempty"`
	PageNumber int    `json:"pageNumber,omitempty"`
	PageSize   string `json:"pageSize,omitempty"`
	Offset     string `json:"offset,omitempty"`
	MachineIDs string `json:"machineIds,omitempty"`
}

func ConfigureMicrosoftDefender(ctx context.Context, service *zscaler.Service, endpointList *MicrosoftDefender) (*MicrosoftDefender, *http.Response, error) {
	resp, err := service.Client.Create(ctx, partnerIntegrationEndpoint+"/microsoftDefender/endpoints", *endpointList)
	if err != nil {
		return nil, nil, err
	}

	endpointList, ok := resp.(*MicrosoftDefender)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a microsoft defender endpoint list pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new microsoft defender endpoint list from create")
	return endpointList, nil, nil
}

// //////////////////////////////////////////////
// Sandbox Report MD5 Hash Partner Integration///
// //////////////////////////////////////////////
type SandboxReportMD5 struct {
	ThreatName      string `json:"threatName,omitempty"`
	SandboxCategory string `json:"sandboxCategory,omitempty"`
	SandboxScore    int    `json:"sandboxScore,omitempty"`
	FileType        string `json:"fileType,omitempty"`
	FileSize        int    `json:"fileSize,omitempty"`
	MD5             string `json:"md5,omitempty"`
	SHA1            string `json:"sha1,omitempty"`
	SHA256          string `json:"sha256,omitempty"`
	Ssdeep          string `json:"ssdeep,omitempty"`
	ThreatLink      string `json:"threatLink,omitempty"`
	Message         string `json:"message,omitempty"`
	OriginLanguage  string `json:"originLanguage,omitempty"`
	OriginCountry   string `json:"originCountry,omitempty"`
}

// GetReportMD5Hash retrieves the Sandbox Detail Report based on the MD5 hash of
// the file. The md5 path parameter is the only required parameter.
func GetReportMD5Hash(ctx context.Context, service *zscaler.Service, md5Hash string) (*SandboxReportMD5, error) {
	var sandboxReport SandboxReportMD5
	endpoint := fmt.Sprintf("%s/sandbox/report/%s", partnerIntegrationEndpoint, md5Hash)
	err := service.Client.Read(ctx, endpoint, &sandboxReport)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning sandbox report from Get for MD5 hash '%s': %+v", md5Hash, sandboxReport)
	return &sandboxReport, nil
}
