package shadowitreport

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	cloudApplicationsEndpoint    = "/cloudApplications/lite"
	customTagsEndpoint           = "/customTags"
	appExportEndpoint            = "/shadowIT/applications/export"
	appExportCsvEndpointTemplate = "/shadowIT/applications/%s/exportCsv"
	appBulkUpdateEndpoint        = "/cloudApplications/bulkUpdate"
)

type CloudApplicationsAndCustomTags struct {
	// Unique identifier of the cloud application
	ID int `json:"id"`

	// The name of the cloud application
	Name string `json:"name,omitempty"`
}

type CloudApplicationsExport struct {
	Duration                         string                   `json:"duration,omitempty"`
	Application                      []string                 `json:"application,omitempty"`
	AppName                          string                   `json:"appName,omitempty"`
	ApplicationCategory              []string                 `json:"applicationCategory,omitempty"`
	DataConsumed                     []common.DataConsumed    `json:"dataConsumed,omitempty"`
	RiskIndex                        []int                    `json:"riskIndex,omitempty"`
	Order                            *common.Order            `json:"order,omitempty"`
	SanctionedState                  []string                 `json:"sanctionedState,omitempty"`
	Employees                        []string                 `json:"employees,omitempty"`
	SupportedCertifications          *SupportedCertifications `json:"supportedCertifications,omitempty"`
	SourceIpRestriction              []string                 `json:"sourceIpRestriction,omitempty"`
	MfaSupport                       []string                 `json:"mfaSupport,omitempty"`
	AdminAuditLogs                   []string                 `json:"adminAuditLogs,omitempty"`
	HadBreachInLast3Years            []string                 `json:"hadBreachInLast3Years,omitempty"`
	HavePoorItemsOfService           []string                 `json:"havePoorItemsOfService,omitempty"`
	PasswordStrength                 []string                 `json:"passwordStrength,omitempty"`
	SslPinned                        []string                 `json:"sslPinned,omitempty"`
	Evasive                          []string                 `json:"evasive,omitempty"`
	HaveHTTPSecurityHeaderSupport    []string                 `json:"haveHTTPSecurityHeaderSupport,omitempty"`
	DnsCAAPolicy                     []string                 `json:"dnsCAAPolicy,omitempty"`
	HaveWeakCipherSupport            []string                 `json:"haveWeakCipherSupport,omitempty"`
	SslCertificationValidity         []string                 `json:"sslCertificationValidity,omitempty"`
	MalwareScanningContent           []string                 `json:"malwareScanningContent,omitempty"`
	FileSharing                      []string                 `json:"fileSharing,omitempty"`
	RemoteAccessScreenSharing        []string                 `json:"remoteAccessScreenSharing,omitempty"`
	SenderPolicyFramework            []string                 `json:"senderPolicyFramework,omitempty"`
	DomainKeysIdentifiedMail         []string                 `json:"domainKeysIdentifiedMail,omitempty"`
	DomainBasedMessageAuthentication []string                 `json:"domainBasedMessageAuthentication,omitempty"`
	VulnerableDisclosureProgram      []string                 `json:"vulnerableDisclosureProgram,omitempty"`
	WafSupport                       []string                 `json:"wafSupport,omitempty"`
	Vulnerability                    []string                 `json:"vulnerability,omitempty"`
	ValidSSLCertificate              []string                 `json:"validSSLCertificate,omitempty"`
	DataEncryptionInTransit          []string                 `json:"dataEncryptionInTransit,omitempty"`
	VulnerableToHeartBleed           []string                 `json:"vulnerableToHeartBleed,omitempty"`
	VulnerableToPoodle               []string                 `json:"vulnerableToPoodle,omitempty"`
	VulnerableToLogJam               []string                 `json:"vulnerableToLogJam,omitempty"`
	CertKeySize                      *CertKeySize             `json:"certKeySize,omitempty"`
}

type SupportedCertifications struct {
	Operation string   `json:"operation,omitempty"`
	Value     []string `json:"value,omitempty"`
}

type CertKeySize struct {
	Operation string   `json:"operation,omitempty"`
	Value     []string `json:"value,omitempty"`
}

type CloudApplicationsExportCSV struct {
	Duration      string                `json:"duration,omitempty"`
	Application   []string              `json:"application,omitempty"`
	Order         *common.Order         `json:"order,omitempty"`
	DownloadBytes []common.DataConsumed `json:"downloadBytes,omitempty"`
	UploadBytes   []common.DataConsumed `json:"uploadBytes,omitempty"`
	DataConsumed  []common.DataConsumed `json:"dataConsumed,omitempty"`
	Users         []User                `json:"users,omitempty"`
	Locations     []Location            `json:"locations,omitempty"`
	Departments   []Department          `json:"departments,omitempty"`
}

type User struct {
	ID          int    `json:"id,omitempty"`
	PID         int    `json:"pid,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Deleted     bool   `json:"deleted,omitempty"`
	GetlID      int    `json:"getlId,omitempty"`
}

type Location struct {
	ID          int    `json:"id,omitempty"`
	PID         int    `json:"pid,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Deleted     bool   `json:"deleted,omitempty"`
	GetlID      int    `json:"getlId,omitempty"`
}

type Department struct {
	ID          int    `json:"id,omitempty"`
	PID         int    `json:"pid,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Deleted     bool   `json:"deleted,omitempty"`
	GetlID      int    `json:"getlId,omitempty"`
}

type ApplicationBulkUpdate struct {
	SanctionedState                string                           `json:"sanctionedState,omitempty"`
	ApplicationIDs                 []int                            `json:"applicationIds,omitempty"`
	CloudApplicationsAndCustomTags []CloudApplicationsAndCustomTags `json:"customTags,omitempty"`
}

// GetAllCloudAppsLite retrieves all cloud applications in lite format with optional pagination parameters
func GetAllCloudAppsLite(service *services.Service, pageNumber, limit *int) ([]CloudApplicationsAndCustomTags, error) {
	endpoint := cloudApplicationsEndpoint

	// Build the query parameters
	queryParams := url.Values{}
	if pageNumber != nil {
		queryParams.Add("PageNumber", strconv.Itoa(*pageNumber))
	}
	if limit != nil {
		queryParams.Add("limit", strconv.Itoa(*limit))
	}

	// Append query parameters to endpoint if any
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	var cloudAppsLite []CloudApplicationsAndCustomTags
	err := service.Client.Read(endpoint, &cloudAppsLite)
	if err != nil {
		return nil, err
	}

	return cloudAppsLite, nil
}

func GetAllCustomTags(service *services.Service) ([]CloudApplicationsAndCustomTags, error) {
	var customTags []CloudApplicationsAndCustomTags
	err := service.Client.Read(customTagsEndpoint, &customTags)
	return customTags, err
}

func CreateCloudApplicationsExport(service *services.Service, exportRequest CloudApplicationsExport) (*http.Response, error) {
	// Send the request with the export payload
	resp, err := service.Client.Create(appExportEndpoint, exportRequest)
	if err != nil {
		return nil, err
	}

	// Assert that the response is of type *http.Response (not parsing as JSON)
	httpResp, ok := resp.(*http.Response)
	if !ok {
		return nil, errors.New("unexpected response type, expected *http.Response")
	}

	// Check if the response contains a valid CSV file
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to export cloud applications, unexpected response code: %d", httpResp.StatusCode)
	}

	service.Client.Logger.Printf("[DEBUG] Successfully exported cloud applications with payload: %+v", exportRequest)

	return httpResp, nil
}

// CreateCloudApplicationsExportCSV sends a POST request to create a new CloudApplicationsExportCSV
func CreateCloudApplicationsExportCSV(service *services.Service, entity string, appExport *CloudApplicationsExportCSV) (*CloudApplicationsExportCSV, *http.Response, error) {
	// Validate the entity parameter
	if entity != "USER" && entity != "LOCATION" {
		return nil, nil, errors.New("invalid entity value; must be 'USER' or 'LOCATION'")
	}

	// Create the endpoint with the entity
	endpoint := fmt.Sprintf(appExportCsvEndpointTemplate, entity)

	resp, err := service.Client.Create(endpoint, appExport)
	if err != nil {
		return nil, nil, err
	}

	createdExport, ok := resp.(*CloudApplicationsExportCSV)
	if !ok {
		return nil, nil, errors.New("object returned from API was not a CloudApplicationsExportCSV pointer")
	}

	service.Client.Logger.Printf("[DEBUG] Successfully created new application export")
	return createdExport, nil, nil
}

// Update sends a PUT request to perform bulk updates for cloud applications
func Update(service *services.Service, rules *ApplicationBulkUpdate) (*ApplicationBulkUpdate, error) {
	resp, err := service.Client.UpdateWithPut(appBulkUpdateEndpoint, *rules)
	if err != nil {
		return nil, err
	}

	// Handle the 204 No Content case by returning nil, since no content was returned
	if resp == nil {
		service.Client.Logger.Printf("[DEBUG] No content returned from API (204 No Content)")
		return nil, nil // No content, but operation was successful
	}

	// Proceed with normal processing for non-empty responses
	updatedRules, ok := resp.(*ApplicationBulkUpdate)
	if !ok {
		return nil, errors.New("object returned from API was not an ApplicationBulkUpdate pointer")
	}

	service.Client.Logger.Printf("[DEBUG] Successfully updated application bulk update")
	return updatedRules, nil
}
