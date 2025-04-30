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
	Name                        string                    `json:"name"`
	FeedStatus                  string                    `json:"feedStatus"`
	NssLogType                  string                    `json:"nssLogType"`
	NssFeedType                 string                    `json:"nssFeedType"`
	FeedOutputFormat            string                    `json:"feedOutputFormat"`
	UserObfuscation             string                    `json:"userObfuscation"`
	TimeZone                    string                    `json:"timeZone"`
	CustomEscapedCharacter      []string                  `json:"customEscapedCharacter"`
	EpsRateLimit                int                       `json:"epsRateLimit"`
	JsonArrayToggle             bool                      `json:"jsonArrayToggle"`
	SiemType                    string                    `json:"siemType"`
	MaxBatchSize                int                       `json:"maxBatchSize"`
	ConnectionURL               string                    `json:"connectionURL"`
	AuthenticationToken         string                    `json:"authenticationToken"`
	ConnectionHeaders           []string                  `json:"connectionHeaders"`
	LastSuccessFullTest         int                       `json:"lastSuccessFullTest"`
	TestConnectivityCode        int                       `json:"testConnectivityCode"`
	Base64EncodedCertificate    string                    `json:"base64EncodedCertificate"`
	NssType                     string                    `json:"nssType"`
	ClientID                    string                    `json:"clientId"`
	ClientSecret                string                    `json:"clientSecret"`
	AuthenticationUrl           string                    `json:"authenticationUrl"`
	GrantType                   string                    `json:"grantType"`
	Scope                       string                    `json:"scope"`
	OauthAuthentication         bool                      `json:"oauthAuthentication"`
	ServerIps                   []string                  `json:"serverIps"`
	ClientIps                   []string                  `json:"clientIps"`
	Domains                     []string                  `json:"domains"`
	DNSRequestTypes             []string                  `json:"dnsRequestTypes"`
	DNSResponseTypes            []string                  `json:"dnsResponseTypes"`
	DNSResponses                []string                  `json:"dnsResponses"`
	Durations                   []string                  `json:"durations"`
	DNSActions                  []string                  `json:"dnsActions"`
	FirewallLoggingMode         string                    `json:"firewallLoggingMode"`
	ClientSourceIps             []string                  `json:"clientSourceIps"`
	FirewallActions             []string                  `json:"firewallActions"`
	Locations                   []Location                `json:"locations"`
	Countries                   []string                  `json:"countries"`
	ServerSourcePorts           []string                  `json:"serverSourcePorts"`
	ClientSourcePorts           []string                  `json:"clientSourcePorts"`
	ActionFilter                string                    `json:"actionFilter"`
	EmailDlpPolicyAction        string                    `json:"emailDlpPolicyAction"`
	Direction                   string                    `json:"direction"`
	Event                       string                    `json:"event"`
	PolicyReasons               []string                  `json:"policyReasons"`
	ProtocolTypes               []string                  `json:"protocolTypes"`
	UserAgents                  []string                  `json:"userAgents"`
	RequestMethods              []string                  `json:"requestMethods"`
	CasbSeverity                []string                  `json:"casbSeverity"`
	CasbPolicyTypes             []string                  `json:"casbPolicyTypes"`
	CasbApplications            []string                  `json:"casbApplications"`
	CasbAction                  []string                  `json:"casbAction"`
	CasbTenant                  []Location                `json:"casbTenant"`
	URLSuperCategories          []string                  `json:"urlSuperCategories"`
	WebApplications             []WebApplication          `json:"webApplications"`
	WebApplicationClasses       []string                  `json:"webApplicationClasses"`
	MalwareNames                []string                  `json:"malwareNames"`
	URLClasses                  []string                  `json:"urlClasses"`
	MalwareClasses              []string                  `json:"malwareClasses"`
	AdvancedThreats             []string                  `json:"advancedThreats"`
	ResponseCodes               []string                  `json:"responseCodes"`
	NwApplications              []string                  `json:"nwApplications"`
	NatActions                  []string                  `json:"natActions"`
	TrafficForwards             []string                  `json:"trafficForwards"`
	WebTrafficForwards          []string                  `json:"webTrafficForwards"`
	TunnelTypes                 []string                  `json:"tunnelTypes"`
	Alerts                      []string                  `json:"alerts"`
	ObjectType                  []string                  `json:"objectType"`
	Activity                    []string                  `json:"activity"`
	ObjectType1                 []string                  `json:"objectType1"`
	ObjectType2                 []string                  `json:"objectType2"`
	EndPointDLPLogType          []string                  `json:"endPointDLPLogType"`
	EmailDLPLogType             []string                  `json:"emailDLPLogType"`
	FileTypeSuperCategories     []string                  `json:"fileTypeSuperCategories"`
	FileTypeCategories          []string                  `json:"fileTypeCategories"`
	CasbFileType                []string                  `json:"casbFileType"`
	CasbFileTypeSuperCategories []string                  `json:"casbFileTypeSuperCategories"`
	Users                       []Location                `json:"users"`
	Departments                 []Location                `json:"departments"`
	SenderName                  []Location                `json:"senderName"`
	Buckets                     []Location                `json:"buckets"`
	VPNCredentials              []Location                `json:"vpnCredentials"`
	MessageSize                 []string                  `json:"messageSize"`
	FileSizes                   []string                  `json:"fileSizes"`
	RequestSizes                []string                  `json:"requestSizes"`
	ResponseSizes               []string                  `json:"responseSizes"`
	TransactionSizes            []string                  `json:"transactionSizes"`
	InBoundBytes                []string                  `json:"inBoundBytes"`
	OutBoundBytes               []string                  `json:"outBoundBytes"`
	DownloadTime                []string                  `json:"downloadTime"`
	ScanTime                    []string                  `json:"scanTime"`
	ServerSourceIps             []string                  `json:"serverSourceIps"`
	ServerDestinationIps        []string                  `json:"serverDestinationIps"`
	TunnelIps                   []string                  `json:"tunnelIps"`
	InternalIps                 []string                  `json:"internalIps"`
	TunnelSourceIps             []string                  `json:"tunnelSourceIps"`
	TunnelDestIps               []string                  `json:"tunnelDestIps"`
	ClientDestinationIps        []string                  `json:"clientDestinationIps"`
	AuditLogType                []string                  `json:"auditLogType"`
	ProjectName                 []string                  `json:"projectName"`
	RepoName                    []string                  `json:"repoName"`
	ObjectName                  []string                  `json:"objectName"`
	ChannelName                 []string                  `json:"channelName"`
	FileSource                  []string                  `json:"fileSource"`
	FileName                    []string                  `json:"fileName"`
	SessionCounts               []string                  `json:"sessionCounts"`
	AdvUserAgents               []string                  `json:"advUserAgents"`
	RefererUrls                 []string                  `json:"refererUrls"`
	HostNames                   []string                  `json:"hostNames"`
	FullUrls                    []string                  `json:"fullUrls"`
	ThreatNames                 []string                  `json:"threatNames"`
	PageRiskIndexes             []string                  `json:"pageRiskIndexes"`
	ClientDestinationPorts      []string                  `json:"clientDestinationPorts"`
	TunnelSourcePort            []string                  `json:"tunnelSourcePort"`
	ExternalOwners              []common.IDNameExternalID `json:"externalOwners"`
	ExternalCollaborators       []common.IDNameExternalID `json:"externalCollaborators"`
	InternalCollaborators       []common.IDNameExternalID `json:"internalCollaborators"`
	ItsmObjectType              []common.IDNameExternalID `json:"itsmObjectType"`
	URLCategories               []common.IDNameExternalID `json:"urlCategories"`
	DLPEngines                  []common.IDNameExternalID `json:"dlpEngines"`
	DLPDictionaries             []common.IDNameExternalID `json:"dlpDictionaries"`
	Rules                       []common.IDNameExternalID `json:"rules"`
	NwServices                  []common.IDNameExternalID `json:"nwServices"`
}

type Location struct {
	ID          int    `json:"id"`
	PID         int    `json:"pid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Deleted     bool   `json:"deleted"`
	GetlID      int    `json:"getlId"`
}

type WebApplication struct {
	Val                 int    `json:"val"`
	WebApplicationClass string `json:"webApplicationClass"`
	BackendName         string `json:"backendName"`
	OriginalName        string `json:"originalName"`
	Extended            bool   `json:"extended"`
	Misc                bool   `json:"misc"`
	Name                string `json:"name"`
	Deprecated          bool   `json:"deprecated"`
}

type IDNameDescription struct {
	ID          int    `json:"id"`
	PID         int    `json:"pid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Deleted     bool   `json:"deleted"`
	GetlID      int    `json:"getlId"`
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
