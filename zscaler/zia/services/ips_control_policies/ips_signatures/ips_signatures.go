package ips_signatures

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	ipsSignaturesEndpoint         = "/zia/api/v1/ipsSignatureRules"
	ipsSignaturesExportEndpoint   = "/zia/api/v1/ipsSignatureRules/export"
	ipsSignaturesImportEndpoint   = "/zia/api/v1/ipsSignatureRules/import"
	ipsSignaturesValidateEndpoint = "/zia/api/v1/ipsSignatureRules/validateRuleText"

	contentTypeCSV = "text/csv"
)

type IPSSignatureRules struct {
	// System-generated identifier for the custom IPS signature rule
	ID int `json:"id"`

	// Custom IPS signature rule name
	Name string `json:"name,omitempty"`

	// The rule text in Suricata/Snort syntax that defines the custom IPS signature
	RuleText string `json:"ruleText,omitempty"`

	// Additional information about the custom signature rule
	Description string `json:"description,omitempty"`

	// The threat category that is assigned to the custom signature rule
	Category *IPSSignatureCategory `json:"category,omitempty"`

	// A Boolean value that indicates whether the custom signature rule is enabled
	// and is ready to be used in IPS Control rules via the assigned threat category
	Enabled bool `json:"enabled"`

	// A Boolean value that indicates whether the custom signature rule is deleted
	Deleted bool `json:"deleted,omitempty"`

	// Unix timestamp (in seconds) when the rule was promoted; 0 if not yet promoted
	PromoteTime int `json:"promoteTime,omitempty"`

	// Unix timestamp (in seconds) when the rule text was last modified
	RuleTextModTime int `json:"ruleTextModTime,omitempty"`

	// A Boolean value that indicates whether the rule was submitted for dynamic validation
	DynamicValidationSubmitted bool `json:"dynamicValidationSubmitted,omitempty"`

	// A Boolean value that indicates whether dynamic validation rejected the rule
	DynamicValidationRejected bool `json:"dynamicValidationRejected,omitempty"`

	// A Boolean value that indicates whether dynamic validation succeeded for the rule
	DynamicValidationSucceeded bool `json:"dynamicValidationSucceeded,omitempty"`

	// A Boolean value that indicates whether the rule was disabled from Zscaler Cloud Management
	DisabledFromZSCM bool `json:"disabledFromZSCM,omitempty"`

	// Reject code returned by dynamic validation
	DynamicValRejectCode int `json:"dynamicValRejectCode,omitempty"`
}

// IPSSignatureCategory represents the threat category object embedded in an
// IPSSignatureRules response. It is a single object (not an array) and uses
// the localization flag isNameL10nTag rather than the common ExternalID field,
// so we declare a dedicated type instead of reusing common.IDNameExternalID.
type IPSSignatureCategory struct {
	// Unique identifier of the threat category
	ID int `json:"id,omitempty"`

	// Name of the threat category (e.g., ADVANCED_SECURITY)
	Name string `json:"name,omitempty"`

	// Indicates whether the name is a localization tag rather than a literal label
	IsNameL10nTag bool `json:"isNameL10nTag,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, signatureID int) (*IPSSignatureRules, error) {
	var ipsSignature IPSSignatureRules
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", ipsSignaturesEndpoint, signatureID), &ipsSignature)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning IPS signature from Get: %d", ipsSignature.ID)
	return &ipsSignature, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, signatureName string) (*IPSSignatureRules, error) {
	// Use GetAll to leverage the single API call and built-in pagination
	ipsSignatures, err := GetAll(ctx, service)
	if err != nil {
		return nil, err
	}
	// Search for exact match (case-insensitive)
	for _, ipsSignature := range ipsSignatures {
		if strings.EqualFold(ipsSignature.Name, signatureName) {
			return &ipsSignature, nil
		}
	}
	return nil, fmt.Errorf("no IPS signature found with name: %s", signatureName)
}

func Create(ctx context.Context, service *zscaler.Service, ipsSignature *IPSSignatureRules) (*IPSSignatureRules, *http.Response, error) {
	resp, err := service.Client.Create(ctx, ipsSignaturesEndpoint, *ipsSignature)
	if err != nil {
		return nil, nil, err
	}

	createdIPSSignature, ok := resp.(*IPSSignatureRules)
	if !ok {
		return nil, nil, errors.New("object returned from api was not an IPS signature pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new IPS signature from create: %d", createdIPSSignature.ID)
	return createdIPSSignature, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, signatureID int, ipsSignature *IPSSignatureRules) (*IPSSignatureRules, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", ipsSignaturesEndpoint, signatureID), *ipsSignature)
	if err != nil {
		return nil, nil, err
	}
	updatedIPSSignature, _ := resp.(*IPSSignatureRules)

	service.Client.GetLogger().Printf("[DEBUG]returning updated IPS signature from update: %d", updatedIPSSignature.ID)
	return updatedIPSSignature, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, signatureID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", ipsSignaturesEndpoint, signatureID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]IPSSignatureRules, error) {
	var ipsSignatures []IPSSignatureRules
	err := common.ReadAllPages(ctx, service.Client, ipsSignaturesEndpoint, &ipsSignatures)
	return ipsSignatures, err
}

// IPSSignatureRulesImportStatus represents the status of a custom IPS signature
// rules CSV import. It is returned by GET /ipsSignatureRules/import to report
// progress and outcome of the most recent import.
type IPSSignatureRulesImportStatus struct {
	// Overall status of the import (e.g., INIT, IN_PROGRESS, COMPLETED, FAILED)
	Status string `json:"status,omitempty"`

	// Number of records successfully added during the import
	TotalRecordsAdded int `json:"totalRecordsAdded"`

	// Number of records successfully deleted during the import
	TotalRecordsDeleted int `json:"totalRecordsDeleted"`

	// Number of records successfully updated during the import
	TotalRecordsUpdated int `json:"totalRecordsUpdated"`

	// List of records that could not be processed during the import
	FailedRecords []IPSSignatureRulesFailedRecord `json:"failedRecords,omitempty"`

	// Number of records processed so far
	ProcessedRecords int `json:"processedRecords"`

	// Total number of records in the import file
	TotalRecordsInImport int `json:"totalRecordsInImport"`

	// List of errors encountered during the import (not tied to a specific record)
	Errors []IPSSignatureRulesImportError `json:"errors,omitempty"`

	// Percentage of the import that has been completed (0-100)
	PercentComplete int `json:"percentComplete"`

	// Top-level error code returned when the import fails
	ErrorCode string `json:"errorCode,omitempty"`
}

// IPSSignatureRulesFailedRecord represents a single record from a custom IPS
// signature rules CSV import that the API was unable to process.
type IPSSignatureRulesFailedRecord struct {
	// Error code categorizing the failure (e.g., CONFIGURATION)
	ErrorCode string `json:"errorCode,omitempty"`

	// Name of the custom IPS signature rule that failed to import
	Name string `json:"name,omitempty"`

	// Action that was attempted on the record (e.g., ADD, UPDATE, DELETE)
	Action string `json:"action,omitempty"`

	// Human-readable description of the failure
	Description string `json:"description,omitempty"`
}

// IPSSignatureRulesImportError represents a global error encountered during a
// custom IPS signature rules CSV import (i.e., not tied to a specific record).
type IPSSignatureRulesImportError struct {
	// Error code categorizing the failure (e.g., CONFIGURATION)
	ErrorCode string `json:"errorCode,omitempty"`

	// Human-readable description of the error
	Description string `json:"description,omitempty"`
}

// IPSSignatureRuleTextValidationRequest is the JSON body accepted by the
// POST /ipsSignatureRules/validateRuleText endpoint. The API rejects raw
// string payloads; the rule text must be wrapped in this object.
type IPSSignatureRuleTextValidationRequest struct {
	// Rule text in Suricata/Snort syntax to validate
	RuleText string `json:"ruleText"`
}

// IPSSignatureRulesValidation represents the result of validating a custom IPS
// signature rule text via POST /ipsSignatureRules/validateRuleText. It maps to
// the SMMsgStatusInfo model and reports syntax errors, duplicate signatures,
// and similar conditions.
//
// Note: The Zscaler API uses HTTP-level signaling for validation. A well-formed
// rule returns HTTP 200 with status=0 and empty error fields. An invalid rule
// returns HTTP 400 with a standard error envelope (INVALID_INPUT_ARGUMENT and
// the diagnostic in "message"); the SDK surfaces that as a Go error, not as a
// populated struct.
type IPSSignatureRulesValidation struct {
	// Numeric validation status code returned by the API (0 indicates success)
	Status int `json:"status"`

	// Position in the rule text where the error occurred
	ErrPosition int `json:"errPosition,omitempty"`

	// Error message for the rule text corresponding to the error position
	ErrMsg string `json:"errMsg,omitempty"`

	// Rule text parameter that caused the error
	ErrParameter string `json:"errParameter,omitempty"`

	// Suggestion for error correction
	ErrSuggestion string `json:"errSuggestion,omitempty"`

	// Optional list of IDs used by certain validations
	IDList []int `json:"idList,omitempty"`

	// Optional map of sub-identifiers used by certain validations
	SubIdsMap map[string]interface{} `json:"subIdsMap,omitempty"`
}

// NOTE: The /ipsSignatureRules/{export,import} endpoints are currently broken
// upstream (see internal bug tracker). The three functions below are kept as
// reference for when the backend is fixed, but are commented out so callers
// cannot accidentally hit the broken endpoints. Re-enable in lockstep with
// the dev test script at local_test/OneAPI/zia_dev_tests/ips_signature/main.go
// and the corresponding integration tests.

/*
// GetImportIPSSignatureRulesStatus retrieves the status of the most recent
// custom IPS signature rules CSV import (GET /ipsSignatureRules/import).
func GetImportIPSSignatureRulesStatus(ctx context.Context, service *zscaler.Service) (*IPSSignatureRulesImportStatus, error) {
	var status IPSSignatureRulesImportStatus
	if err := service.Client.Read(ctx, ipsSignaturesImportEndpoint, &status); err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning custom IPS signature rules import status: %s (%d%%)", status.Status, status.PercentComplete)
	return &status, nil
}

// ExportIPSSignatureRules exports custom IPS signature rules to a CSV file
// (GET /ipsSignatureRules/export). The raw CSV bytes are returned to the
// caller, which can then write them to disk or stream them onward.
func ExportIPSSignatureRules(ctx context.Context, service *zscaler.Service) ([]byte, error) {
	csvBytes, err := service.Client.ReadRaw(ctx, ipsSignaturesExportEndpoint, contentTypeCSV)
	if err != nil {
		return nil, fmt.Errorf("failed to export custom IPS signature rules: %w", err)
	}

	service.Client.GetLogger().Printf("[DEBUG] Exported custom IPS signature rules CSV: %d bytes", len(csvBytes))
	return csvBytes, nil
}

// ImportIPSSignatureRules uploads a CSV file containing custom IPS signature
// rules (POST /ipsSignatureRules/import). The CSV must use the same format as
// the Sample Import CSV file provided by Zscaler. Per the API contract, this
// POST does not return a structured response body; callers should follow up
// with GetImportIPSSignatureRulesStatus to track progress.
func ImportIPSSignatureRules(ctx context.Context, service *zscaler.Service, csvData []byte) (*http.Response, error) {
	if len(csvData) == 0 {
		return nil, errors.New("csv data is required to import custom IPS signature rules")
	}

	_, resp, err := service.Client.CreateWithRawPayloadAndContentType(ctx, ipsSignaturesImportEndpoint, csvData, contentTypeCSV)
	if err != nil {
		return nil, fmt.Errorf("failed to import custom IPS signature rules: %w", err)
	}

	service.Client.GetLogger().Printf("[DEBUG] Imported custom IPS signature rules CSV (%d bytes), status: %d", len(csvData), resp.StatusCode)
	return resp, nil
}
*/

// ValidateIPSSignatureRuleText validates a custom IPS signature rule text
// (POST /ipsSignatureRules/validateRuleText). The rule text is wrapped in a
// JSON object ({"ruleText": "..."}) as required by the API and the response
// is decoded into an SMMsgStatusInfo describing any syntax errors, duplicate
// signatures, or similar conditions detected.
func ValidateIPSSignatureRuleText(ctx context.Context, service *zscaler.Service, ruleText string) (*IPSSignatureRulesValidation, error) {
	if ruleText == "" {
		return nil, errors.New("rule text is required to validate a custom IPS signature rule")
	}

	var validation IPSSignatureRulesValidation
	err := service.Client.CreateWithJSONResponse(ctx, ipsSignaturesValidateEndpoint, IPSSignatureRuleTextValidationRequest{RuleText: ruleText}, &validation)
	if err != nil {
		return nil, fmt.Errorf("failed to validate custom IPS signature rule text: %w", err)
	}

	service.Client.GetLogger().Printf("[DEBUG] Custom IPS signature rule validation status: %d", validation.Status)
	return &validation, nil
}
