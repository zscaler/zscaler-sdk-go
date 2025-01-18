package apptotal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	appEndpoint              = "/zia/api/v1/apps/app"
	appSearchEndpoint        = "/zia/api/v1/apps/search"
	appViewsListEndpoint     = "/zia/api/v1/app_views/list"
	appViewsListBaseEndpoint = "/zia/api/v1/app_views"
)

// ######################### 3RD-PARTY APP GOVERNANCE API - SEARCH APP CATALOG #################
type ApplicationCatalog struct {
	Name              string          `json:"name"`
	Publisher         Publisher       `json:"publisher"`
	Platform          string          `json:"platform"`
	Description       string          `json:"description"`
	RedirectUrls      []string        `json:"redirectUrls"`
	WebsiteUrls       []string        `json:"websiteUrls"`
	Categories        []string        `json:"categories"`
	Tags              []string        `json:"tags"`
	PermissionLevel   float64         `json:"permissionLevel"`
	RiskScore         float64         `json:"riskScore"`
	Risk              string          `json:"risk"`
	ExternalIds       []ExternalID    `json:"externalIds"`
	ClientId          string          `json:"clientId"`
	Permissions       []Permission    `json:"permissions"`
	Compliance        []string        `json:"compliance"`
	DataRetention     string          `json:"dataRetention"`
	ClientType        string          `json:"clientType"`
	LogoUrl           string          `json:"logoUrl"`
	PrivacyPolicyUrl  string          `json:"privacyPolicyUrl"`
	TermsOfServiceUrl string          `json:"termsOfServiceUrl"`
	MarketplaceUrl    string          `json:"marketplaceUrl"`
	MarketplaceData   MarketplaceData `json:"marketplaceData"`
	PlatformVerified  bool            `json:"platformVerified"`
	CanonicVerified   bool            `json:"canonicVerified"`
	DeveloperEmail    string          `json:"developerEmail"`
	ConsentScreenshot string          `json:"consentScreenshot"`
	IPAddresses       []IPAddress     `json:"ipAddresses"`
	ExtractedUrls     []string        `json:"extractedUrls"`
	ExtractedApiCalls []string        `json:"extractedApiCalls"`
	Vulnerabilities   []Vulnerability `json:"vulnerabilities"`
	ApiActivities     []ApiActivity   `json:"apiActivities"`
	Risks             []Risk          `json:"risks"`
	Insights          []Insight       `json:"insights"`
	Instances         []Instance      `json:"instances"`
}

type Publisher struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SiteUrl     string `json:"siteUrl"`
	LogoUrl     string `json:"logoUrl"`
}

type ExternalID struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Permission struct {
	Scope       string `json:"scope"`
	Service     string `json:"service"`
	Description string `json:"description"`
	AccessType  string `json:"accessType"`
	Level       string `json:"level"`
}

type MarketplaceData struct {
	Stars     int `json:"stars"`
	Downloads int `json:"downloads"`
	Reviews   int `json:"reviews"`
}

type IPAddress struct {
	ISPName     string `json:"ispName"`
	IPAddress   string `json:"ipAddress"`
	ProxyType   string `json:"proxyType"`
	UsageType   string `json:"usageType"`
	DomainName  string `json:"domainName"`
	CountryCode string `json:"countryCode"`
}

type Vulnerability struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	CVEId    string `json:"cveId"`
	Summary  string `json:"summary"`
	Severity string `json:"severity"`
}

type ApiActivity struct {
	OperationType string  `json:"operationType"`
	Percentage    float64 `json:"percentage"`
}

type Risk struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Severity    string `json:"severity"`
}

type Insight struct {
	Description string            `json:"description"`
	Timestamp   int64             `json:"timestamp"`
	Urls        map[string]string `json:"urls"`
}

type Instance struct {
	ID             string `json:"id"`
	IntegrationId  string `json:"integrationId"`
	Status         string `json:"status"`
	Classification string `json:"classification"`
}

// ######################### 3RD-PARTY APP GOVERNANCE API - SUBMIT APP FOR ANALYSIS #################
type AppPayload struct {
	AppID string `json:"appId"`
}

// ######################### 3RD-PARTY APP GOVERNANCE API - SEARCH APP BY NAME #################################
type PaginationResponse struct {
	Count       int      `json:"count"`
	CurrentPage int      `json:"currentPage"`
	Data        []Result `json:"data"`
}

type Result struct {
	Result AppDetails `json:"result"`
}

type AppDetails struct {
	AppID     string `json:"appId"`
	Name      string `json:"name"`
	Provider  string `json:"provider"`
	Publisher string `json:"publisher"`
}

// ######################### 3RD-PARTY APP GOVERNANCE API - RETRIEVES LIST OF CUSTOM VIEWS ###################
type AppViewsList struct {
	ID        string      `json:"id"`
	AppIDs    []string    `json:"appIds"` // Used to retrieve all assets
	Name      string      `json:"name"`
	CreatedBy string      `json:"createdBy"`
	CreatedAt int64       `json:"createdAt"`
	Spec      SpecDetails `json:"spec"`
}

type SpecDetails struct {
	Map string `json:"map"`
}

// ######################### 3RD-PARTY APP GOVERNANCE API - RETRIEVES ALL ASSETS #################################
type AppViewAppsResponse struct {
	ID     string      `json:"id"`
	AppIDs []string    `json:"appIds"` // Used to retrieve all assets
	Spec   SpecDetails `json:"spec"`
}

// ######################### 3RD-PARTY APP GOVERNANCE API - SEARCH APP CATALOG #################
func GetAllApps(ctx context.Context, service *zscaler.Service, appID, appURL string, verbose bool) (*ApplicationCatalog, error) {
	// Validate that only one of appID or appURL is provided
	if appID != "" && appURL != "" {
		return nil, fmt.Errorf("only one of app_id or url can be provided")
	}

	// Use url.Values to build the query parameters
	queryParams := url.Values{}
	if appID != "" {
		queryParams.Add("app_id", appID)
	}
	if appURL != "" {
		queryParams.Add("url", appURL)
	}
	if verbose {
		queryParams.Add("verbose", "true")
	}

	// Append query parameters to the endpoint
	endpoint := appEndpoint
	if len(queryParams) > 0 {
		endpoint = fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())
	}

	// Fetch application details from the API
	var appDetails ApplicationCatalog
	err := service.Client.Read(ctx, endpoint, &appDetails)
	if err != nil {
		return nil, err
	}
	return &appDetails, nil
}

// ######################### 3RD-PARTY APP GOVERNANCE API - SUBMIT APP FOR ANALYSIS #################
// CreateApp submits an app for analysis in the 3rd-Party App Governance Sandbox.
// The API accepts a single attribute in the payload: "appId".
func CreateApp(ctx context.Context, service *zscaler.Service, appID string) (*AppPayload, *http.Response, error) {
	// Construct the payload
	payload := &AppPayload{
		AppID: appID,
	}

	// Send the POST request with the payload
	resp, err := service.Client.Create(ctx, appEndpoint, payload)
	if err != nil {
		return nil, nil, err
	}

	// Parse the response into the AppPayload struct
	createdApp, ok := resp.(*AppPayload)
	if !ok {
		return nil, nil, errors.New("object returned from API was not an AppPayload pointer")
	}

	// Log the result and return
	service.Client.GetLogger().Printf("[DEBUG] Returning new app from create: %s", createdApp.AppID)
	return createdApp, nil, nil
}

// ######################### 3RD-PARTY APP GOVERNANCE API - SEARCH APP BY NAME #################################
// AppsSearch searches for apps by name with optional limit and pagination.
func AppsSearch(ctx context.Context, service *zscaler.Service, appName string, page, limit *int) ([]PaginationResponse, error) {
	if appName == "" {
		return nil, errors.New("appName is required")
	}

	// Construct query parameters
	queryParams := url.Values{}
	encodedAppName := url.QueryEscape(appName) // Encode appName for RFC 3986 compliance
	queryParams.Set("appName", encodedAppName)

	// Add optional parameters
	if page != nil {
		queryParams.Set("page", strconv.Itoa(*page))
	}
	if limit != nil {
		if *limit > 200 { // Enforce maximum limit of 200
			queryParams.Set("limit", "200")
		} else {
			queryParams.Set("limit", strconv.Itoa(*limit))
		}
	}

	// Build the full endpoint with query parameters
	queryString := queryParams.Encode()
	queryString = strings.ReplaceAll(queryString, "+", "%20") // Ensure spaces are encoded as %20
	fullEndpoint := fmt.Sprintf("%s?%s", appSearchEndpoint, queryString)

	// Make the request
	var appsSearch []PaginationResponse
	err := service.Client.Read(ctx, fullEndpoint, &appsSearch)
	if err != nil {
		return nil, err
	}

	return appsSearch, nil
}

// ######################### 3RD-PARTY APP GOVERNANCE API - RETRIEVES LIST OF CUSTOM VIEWS #################################
func GetAppViewsResponse(ctx context.Context, service *zscaler.Service) (*AppViewsList, error) {
	var appViewsList AppViewsList
	err := service.Client.Read(ctx, appViewsListEndpoint, &appViewsList)
	if err != nil {
		return nil, err
	}
	return &appViewsList, nil
}

func GetAppViewAppsResponse(ctx context.Context, service *zscaler.Service, appViewID string) (*AppViewAppsResponse, error) {
	if appViewID == "" {
		return nil, errors.New("appViewID is required")
	}

	// Construct the endpoint with the appViewID
	endpoint := fmt.Sprintf("%s/%s/apps", appViewsListBaseEndpoint, appViewID)

	var appViewsList AppViewAppsResponse
	err := service.Client.Read(ctx, endpoint, &appViewsList)
	if err != nil {
		return nil, err
	}
	return &appViewsList, nil
}
