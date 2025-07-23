package cloudnss

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	nssFeedsEndpoint            = "/zia/api/v1/nssFeeds"
	nssTestConnectivityEndpoint = "/zia/api/v1/nssFeeds/testConnectivity"
	nssFeedOutputEndpoint       = "/zia/api/v1/nssFeeds/feedOutputDefaults"
)

var (
	supportedTypes = map[string]bool{
		"ADMIN_AUDIT": true, "WEBLOG": true, "ALERT": true, "FWLOG": true,
		"DNSLOG": true, "MULTIFEEDLOG": true, "CASB_FILELOG": true, "CASB_MAILLOG": true,
		"ECLOG": true, "EC_DNSLOG": true, "CASB_ITSM": true, "CASB_CRM": true,
		"CASB_CODE_REPO": true, "CASB_COLLAB": true, "CASB_PCS": true,
		"USER_ACT_REP": true, "USER_COUNT_ALERT": true, "USER_IMP_TRAVEL_ALERT": true,
		"ENDPOINT_DLP": true, "EC_EVENTLOG": true, "EMAIL_DLP": true,
	}

	supportedMultiFeedTypes = map[string]bool{
		"ANY": true, "NONE": true, "WEB": true, "FW": true, "DNS": true, "EMAIL": true,
		"NSSALERTS": true, "BW": true, "CSTAT": true, "TUNNEL": true, "CASB": true,
		"EC_SESS": true, "EC_DNS": true, "USER_ACT_REP": true, "USER_ACT_ALERT": true,
		"USER_APP_STATE": true, "NSS_ADMIN_AUDIT": true, "ENDPOINT_DLP": true,
		"EC_METRICS": true, "EC_EVENT": true, "EMAIL_DLP": true, "ZPA_USER_ACT_LOG": true,
		"ZPA_EVENT_LOG": true, "ZPA_APP_PROTECTION_LOG": true, "ZPA_APP_CONNECTOR_STATUS_LOG": true,
		"ZPA_ZEN_STATUS_LOG": true, "ZPA_USER_STATUS_LOG": true, "EXTRANET_SESS": true,
	}

	supportedFieldFormats = map[string]bool{
		"QRADAR": true, "CSV": true, "TAB_SEPARATED": true, "SPLUNK_CIM": true,
		"ARCSIGHT_CEF": true, "SYMANTEC_MSS": true, "LOGRHYTHM": true,
		"NAME_VALUE_PAIRS": true, "RSA_SECURITY": true, "JSON": true,
	}
)

type NSSFeed struct {
	ID                          int                       `json:"id"`
	Name                        string                    `json:"name,omitempty"`
	FeedStatus                  string                    `json:"feedStatus,omitempty"`
	NssLogType                  string                    `json:"nssLogType,omitempty"`
	NssFeedType                 string                    `json:"nssFeedType,omitempty"`
	FeedOutputFormat            string                    `json:"feedOutputFormat"`
	UserObfuscation             string                    `json:"userObfuscation"`
	TimeZone                    string                    `json:"timeZone,omitempty"`
	CustomEscapedCharacter      []string                  `json:"customEscapedCharacter,omitempty"`
	EpsRateLimit                int                       `json:"epsRateLimit,omitempty"`
	JsonArrayToggle             bool                      `json:"jsonArrayToggle,omitempty"`
	SiemType                    string                    `json:"siemType,omitempty"`
	MaxBatchSize                int                       `json:"maxBatchSize,omitempty"`
	ConnectionURL               string                    `json:"connectionURL,omitempty"`
	AuthenticationToken         string                    `json:"authenticationToken,omitempty"`
	ConnectionHeaders           []string                  `json:"connectionHeaders,omitempty"`
	LastSuccessFullTest         int                       `json:"lastSuccessFullTest,omitempty"`
	TestConnectivityCode        int                       `json:"testConnectivityCode,omitempty"`
	Base64EncodedCertificate    string                    `json:"base64EncodedCertificate,omitempty"`
	NssType                     string                    `json:"nssType,omitempty"`
	ClientID                    string                    `json:"clientId,omitempty"`
	ClientSecret                string                    `json:"clientSecret,omitempty"`
	AuthenticationUrl           string                    `json:"authenticationUrl,omitempty"`
	GrantType                   string                    `json:"grantType,omitempty"`
	Scope                       string                    `json:"scope,omitempty"`
	CloudNSS                    bool                      `json:"cloudNss,omitempty"`
	OauthAuthentication         bool                      `json:"oauthAuthentication,omitempty"`
	ServerIps                   []string                  `json:"serverIps,omitempty"`
	ClientIps                   []string                  `json:"clientIps,omitempty"`
	Domains                     []string                  `json:"domains,omitempty"`
	DNSRequestTypes             []string                  `json:"dnsRequestTypes,omitempty"`
	DNSResponseTypes            []string                  `json:"dnsResponseTypes,omitempty"`
	DNSResponses                []string                  `json:"dnsResponses,omitempty"`
	Durations                   []string                  `json:"durations,omitempty"`
	DNSActions                  []string                  `json:"dnsActions,omitempty"`
	FirewallLoggingMode         string                    `json:"firewallLoggingMode,omitempty"`
	ClientSourceIps             []string                  `json:"clientSourceIps,omitempty"`
	FirewallActions             []string                  `json:"firewallActions,omitempty"`
	Countries                   []string                  `json:"countries,omitempty"`
	ServerSourcePorts           []string                  `json:"serverSourcePorts,omitempty"`
	ClientSourcePorts           []string                  `json:"clientSourcePorts,omitempty"`
	ActionFilter                string                    `json:"actionFilter,omitempty"`
	EmailDlpPolicyAction        string                    `json:"emailDlpPolicyAction,omitempty"`
	Direction                   string                    `json:"direction,omitempty"`
	Event                       string                    `json:"event,omitempty"`
	PolicyReasons               []string                  `json:"policyReasons,omitempty"`
	ProtocolTypes               []string                  `json:"protocolTypes,omitempty"`
	UserAgents                  []string                  `json:"userAgents,omitempty"`
	RequestMethods              []string                  `json:"requestMethods,omitempty"`
	CasbSeverity                []string                  `json:"casbSeverity,omitempty"`
	CasbPolicyTypes             []string                  `json:"casbPolicyTypes,omitempty"`
	CasbApplications            []string                  `json:"casbApplications,omitempty"`
	CasbAction                  []string                  `json:"casbAction,omitempty"`
	URLSuperCategories          []string                  `json:"urlSuperCategories,omitempty"`
	WebApplications             []string                  `json:"webApplications,omitempty"`
	WebApplicationClasses       []string                  `json:"webApplicationClasses,omitempty"`
	MalwareNames                []string                  `json:"malwareNames,omitempty"`
	URLClasses                  []string                  `json:"urlClasses,omitempty"`
	MalwareClasses              []string                  `json:"malwareClasses,omitempty"`
	AdvancedThreats             []string                  `json:"advancedThreats,omitempty"`
	ResponseCodes               []string                  `json:"responseCodes,omitempty"`
	NwApplications              []string                  `json:"nwApplications,omitempty"`
	NatActions                  []string                  `json:"natActions,omitempty"`
	TrafficForwards             []string                  `json:"trafficForwards,omitempty"`
	WebTrafficForwards          []string                  `json:"webTrafficForwards,omitempty"`
	TunnelTypes                 []string                  `json:"tunnelTypes,omitempty"`
	Alerts                      []string                  `json:"alerts,omitempty"`
	ObjectType                  []string                  `json:"objectType,omitempty"`
	Activity                    []string                  `json:"activity,omitempty"`
	ObjectType1                 []string                  `json:"objectType1,omitempty"`
	ObjectType2                 []string                  `json:"objectType2,omitempty"`
	EndPointDLPLogType          []string                  `json:"endPointDLPLogType,omitempty"`
	EmailDLPLogType             []string                  `json:"emailDLPLogType,omitempty"`
	FileTypeSuperCategories     []string                  `json:"fileTypeSuperCategories,omitempty"`
	FileTypeCategories          []string                  `json:"fileTypeCategories,omitempty"`
	CasbFileType                []string                  `json:"casbFileType,omitempty"`
	CasbFileTypeSuperCategories []string                  `json:"casbFileTypeSuperCategories,omitempty"`
	MessageSize                 []string                  `json:"messageSize,omitempty"`
	FileSizes                   []string                  `json:"fileSizes,omitempty"`
	RequestSizes                []string                  `json:"requestSizes,omitempty"`
	ResponseSizes               []string                  `json:"responseSizes,omitempty"`
	TransactionSizes            []string                  `json:"transactionSizes,omitempty"`
	InBoundBytes                []string                  `json:"inBoundBytes,omitempty"`
	OutBoundBytes               []string                  `json:"outBoundBytes,omitempty"`
	DownloadTime                []string                  `json:"downloadTime,omitempty"`
	ScanTime                    []string                  `json:"scanTime,omitempty"`
	ServerSourceIps             []string                  `json:"serverSourceIps,omitempty"`
	ServerDestinationIps        []string                  `json:"serverDestinationIps,omitempty"`
	TunnelIps                   []string                  `json:"tunnelIps,omitempty"`
	InternalIps                 []string                  `json:"internalIps,omitempty"`
	TunnelSourceIps             []string                  `json:"tunnelSourceIps,omitempty"`
	TunnelDestIps               []string                  `json:"tunnelDestIps,omitempty"`
	ClientDestinationIps        []string                  `json:"clientDestinationIps,omitempty"`
	AuditLogType                []string                  `json:"auditLogType,omitempty"`
	ProjectName                 []string                  `json:"projectName,omitempty"`
	RepoName                    []string                  `json:"repoName,omitempty"`
	ObjectName                  []string                  `json:"objectName,omitempty"`
	ChannelName                 []string                  `json:"channelName,omitempty"`
	FileSource                  []string                  `json:"fileSource,omitempty"`
	FileName                    []string                  `json:"fileName,omitempty"`
	SessionCounts               []string                  `json:"sessionCounts,omitempty"`
	AdvUserAgents               []string                  `json:"advUserAgents,omitempty"`
	RefererUrls                 []string                  `json:"refererUrls,omitempty"`
	HostNames                   []string                  `json:"hostNames,omitempty"`
	FullUrls                    []string                  `json:"fullUrls,omitempty"`
	ThreatNames                 []string                  `json:"threatNames,omitempty"`
	PageRiskIndexes             []string                  `json:"pageRiskIndexes,omitempty"`
	ClientDestinationPorts      []string                  `json:"clientDestinationPorts,omitempty"`
	TunnelSourcePort            []string                  `json:"tunnelSourcePort,omitempty"`
	CasbTenant                  []common.CommonNSS        `json:"casbTenant,omitempty"`
	Locations                   []common.CommonNSS        `json:"locations,omitempty"`
	LocationGroups              []common.CommonNSS        `json:"locationGroups,omitempty"`
	Users                       []common.CommonNSS        `json:"users,omitempty"`
	Departments                 []common.CommonNSS        `json:"departments,omitempty"`
	SenderName                  []common.CommonNSS        `json:"senderName,omitempty"`
	Buckets                     []common.CommonNSS        `json:"buckets,omitempty"`
	VPNCredentials              []common.CommonNSS        `json:"vpnCredentials,omitempty"`
	ExternalOwners              []common.IDNameExtensions `json:"externalOwners,omitempty"`
	ExternalCollaborators       []common.IDNameExtensions `json:"externalCollaborators,omitempty"`
	InternalCollaborators       []common.IDNameExtensions `json:"internalCollaborators,omitempty"`
	ItsmObjectType              []common.IDNameExtensions `json:"itsmObjectType,omitempty"`
	URLCategories               []common.IDNameExtensions `json:"urlCategories,omitempty"`
	DLPEngines                  []common.IDNameExtensions `json:"dlpEngines,omitempty"`
	DLPDictionaries             []common.IDNameExtensions `json:"dlpDictionaries,omitempty"`
	Rules                       []common.IDNameExtensions `json:"rules,omitempty"`
	NwServices                  []common.IDNameExtensions `json:"nwServices,omitempty"`
}

type CommonNSS struct {
	ID          int    `json:"id,omitempty"`
	PID         int    `json:"pid,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Deleted     bool   `json:"deleted,omitempty"`
	GetlID      int    `json:"getlId,omitempty"`
}

type WebApplication struct {
	Val                 int    `json:"val,omitempty"`
	WebApplicationClass string `json:"webApplicationClass,omitempty"`
	BackendName         string `json:"backendName,omitempty"`
	OriginalName        string `json:"originalName,omitempty"`
	Extended            bool   `json:"extended,omitempty"`
	Misc                bool   `json:"misc,omitempty"`
	Name                string `json:"name,omitempty"`
	Deprecated          bool   `json:"deprecated,omitempty"`
}

type IDNameDescription struct {
	ID          int    `json:"id,omitempty"`
	PID         int    `json:"pid,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Deleted     bool   `json:"deleted,omitempty"`
	GetlID      int    `json:"getlId,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, feedID int) (*NSSFeed, error) {
	var rule NSSFeed
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", nssFeedsEndpoint, feedID), &rule)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning nss feed from Get: %d", rule.ID)
	return &rule, nil
}

func GetTestConnectivity(ctx context.Context, service *zscaler.Service, feedID int) (*NSSFeed, error) {
	var feed NSSFeed
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", nssTestConnectivityEndpoint, feedID), &feed)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning nss feed from Get: %d", feed.ID)
	return &feed, nil
}

func GetFeedOutputDefaults(ctx context.Context, service *zscaler.Service, params map[string]string) (map[string]string, error) {
	// Validate parameters
	if err := validateFeedOutputParams(params); err != nil {
		return nil, err
	}

	// Build query string if parameters are provided
	queryParams := ""
	if len(params) > 0 {
		q := url.Values{}
		for key, value := range params {
			q.Set(key, value)
		}
		queryParams = "?" + q.Encode()
	}

	// Make the API call
	endpoint := nssFeedOutputEndpoint + queryParams
	var response map[string]string
	err := service.Client.Read(ctx, endpoint, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, feedName string) (*NSSFeed, error) {
	var feeds []NSSFeed
	err := common.ReadAllPages(ctx, service.Client, nssFeedsEndpoint, &feeds)
	if err != nil {
		return nil, err
	}
	for _, feed := range feeds {
		if strings.EqualFold(feed.Name, feedName) {
			return &feed, nil
		}
	}
	return nil, fmt.Errorf("no nss feed found with name: %s", feedName)
}

func Create(ctx context.Context, service *zscaler.Service, feed *NSSFeed) (*NSSFeed, error) {

	// Proceed with creating the feed
	resp, err := service.Client.Create(ctx, nssFeedsEndpoint, *feed)
	if err != nil {
		return nil, err
	}

	createdFeeds, ok := resp.(*NSSFeed)
	if !ok {
		return nil, errors.New("object returned from api was not a feed Pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning rule from create: %d", createdFeeds.ID)
	return createdFeeds, nil
}

func Update(ctx context.Context, service *zscaler.Service, feedID int, rules *NSSFeed) (*NSSFeed, error) {

	// Proceed with updating the rule
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", nssFeedsEndpoint, feedID), *rules)
	if err != nil {
		return nil, err
	}

	updatedFeeds, ok := resp.(*NSSFeed)
	if !ok {
		return nil, errors.New("object returned from api was not a feed Pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG] returning nss feed from update: %d", updatedFeeds.ID)
	return updatedFeeds, nil
}

func Delete(ctx context.Context, service *zscaler.Service, feedID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", nssFeedsEndpoint, feedID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]NSSFeed, error) {
	var rules []NSSFeed
	err := common.ReadAllPages(ctx, service.Client, nssFeedsEndpoint, &rules)
	return rules, err
}

func validateFeedOutputParams(params map[string]string) error {
	if t, ok := params["type"]; ok {
		if !supportedTypes[t] {
			return fmt.Errorf("invalid value for type: %s", t)
		}
	}

	if mf, ok := params["multiFeedType"]; ok {
		if !supportedMultiFeedTypes[mf] {
			return fmt.Errorf("invalid value for multiFeedType: %s", mf)
		}
	}

	if ff, ok := params["fieldFormat"]; ok {
		if !supportedFieldFormats[ff] {
			return fmt.Errorf("invalid value for fieldFormat: %s", ff)
		}
	}

	return nil
}
