package pacfiles

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	pacfileEndpoint = "/zia/api/v1/pacFiles"
)

type PACFileConfig struct {
	// The unique identifier for the PAC file
	ID int `json:"id,omitempty"`

	//The name of the PAC file
	Name string `json:"name,omitempty"`

	// The description of the PAC file
	Description string `json:"description,omitempty"`

	// The domain of your organization to which the PAC file applies
	Domain string `json:"domain,omitempty"`

	// The URL location of the PAC file that is auto-generated when the PAC file is added for the first time.
	// Note: This value is not required in POST and PUT requests and the value is ignored if present.
	PACUrl string `json:"pacUrl,omitempty"`

	// The content of the PAC file.
	// To learn more, see Writing a PAC File. https://help.zscaler.com/zia/writing-pac-file
	PACContent string `json:"pacContent,omitempty"`

	// Indicates whether the PAC file is editable
	Editable bool `json:"editable,omitempty"`

	// Obfuscated domain name of the organization to which this PAC file applies
	PACSubURL string `json:"pacSubURL,omitempty"`

	// A Boolean value that indicates whether the PAC file URL is obfuscated.
	// If this value is true, the obfuscated URL is returned in the pacSubURL field.
	PACUrlObfuscated bool `json:"pacUrlObfuscated,omitempty"`

	// Indicates the verification status of the PAC file and if any errors are identified in the syntax
	// Supported values are: VERIFY_NOERR, VERIFY_ERR, NOVERIFY
	PACVerificationStatus string `json:"pacVerificationStatus,omitempty"`

	// Indicates the status of a specific version of a PAC file
	// as whether it is deployed, staged for deployment, or is marked as the last known good version.
	// Supported values are: DEPLOYED, STAGE, LKG
	PACVersionStatus string `json:"pacVersionStatus,omitempty"`

	// The version number of the PAC file
	PACVersion int `json:"pacVersion,omitempty"`

	// The commit message entered while saving the PAC file version as indicated by the pacVersion field
	PACCommitMessage string `json:"pacCommitMessage,omitempty"`

	// The number of times the PAC file was used during the last 30 days
	TotalHits int `json:"totalHits,omitempty"`

	// The timestamp when the PAC file was last modified. This value is represented in Unix time.
	LastModificationTime int64 `json:"lastModificationTime,omitempty"`

	// The username of the admin who last modified the PAC file
	LastModifiedBy LastModifiedBy `json:"lastModifiedBy,omitempty"`

	// The timestamp when the PAC file was created. This value is represented in Unix time.
	CreateTime int64 `json:"createTime,omitempty"`
}

type LastModifiedBy struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	ExternalID string                 `json:"externalId,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type PacContentRequest struct {
	PacContent string `json:"pacContent"`
}

// PacResult represents the result of PAC file validation
type PacResult struct {
	Success      bool                   `json:"success"`
	Messages     []PacValidationMessage `json:"messages"`
	WarningCount int                    `json:"warningCount"`
	ErrorCount   int                    `json:"errorCount"`
}

// PacValidationMessage represents each validation message returned by the API
type PacValidationMessage struct {
	Severity  int    `json:"severity"`
	EndLine   int    `json:"endLine"`
	EndColumn int    `json:"endColumn"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	Message   string `json:"message"`
	Fatal     bool   `json:"fatal"`
}

func GetPacFiles(ctx context.Context, service *zscaler.Service, filter string) ([]PACFileConfig, error) {
	var pacFiles []PACFileConfig
	endpoint := pacfileEndpoint
	// If the filter is provided, append it to the endpoint URL
	if filter != "" {
		endpoint += fmt.Sprintf("?filter=%s", url.QueryEscape(filter))
	}
	// Call the common ReadAllPages function
	err := common.ReadAllPages(ctx, service.Client, endpoint, &pacFiles)
	return pacFiles, err
}

func GetPacFileByName(ctx context.Context, service *zscaler.Service, pacFileName string) (*PACFileConfig, error) {
	var pacFiles []PACFileConfig
	err := common.ReadAllPages(ctx, service.Client, pacfileEndpoint, &pacFiles)
	if err != nil {
		return nil, err
	}
	for _, pacFile := range pacFiles {
		if strings.EqualFold(pacFile.Name, pacFileName) {
			return &pacFile, nil
		}
	}
	return nil, fmt.Errorf("no pac file found with name: %s", pacFileName)
}

func GetPacFileVersion(ctx context.Context, service *zscaler.Service, pacID int, filter string) ([]PACFileConfig, error) {
	// Initialize an empty slice to hold the PAC file versions
	var pacFiles []PACFileConfig

	// Construct the endpoint URL with the optional filter parameter
	endpoint := fmt.Sprintf("%s/%d/version", pacfileEndpoint, pacID)
	if filter != "" {
		// Add the filter query parameter if provided
		endpoint = fmt.Sprintf("%s?filter=%s", endpoint, filter)
	}

	// Make the API request to retrieve the versions
	err := service.Client.Read(ctx, endpoint, &pacFiles)
	if err != nil {
		return nil, err
	}

	// Log the PAC file version IDs for debugging
	for _, pacFile := range pacFiles {
		service.Client.GetLogger().Printf("[DEBUG] Returning PAC file version ID: %d", pacFile.ID)
	}

	return pacFiles, nil
}

func GetPacVersionID(ctx context.Context, service *zscaler.Service, pacID, pacVersion int, filter string) (*PACFileConfig, error) {
	// Initialize an empty PACFileConfig object to hold the response
	var pacFile PACFileConfig

	// Construct the endpoint URL with pacID and pacVersion
	endpoint := fmt.Sprintf("%s/%d/version/%d", pacfileEndpoint, pacID, pacVersion)

	// If the optional filter is provided, add it as a query parameter
	if filter != "" {
		// Add the filter query parameter if provided
		endpoint = fmt.Sprintf("%s?filter=%s", endpoint, url.QueryEscape(filter))
	}

	// Make the API request to retrieve the PAC file version
	err := service.Client.Read(ctx, endpoint, &pacFile)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve PAC file version: %w", err)
	}

	// Log the PAC file version ID for debugging
	service.Client.GetLogger().Printf("[DEBUG] Returning PAC file version ID: %d", pacFile.ID)

	return &pacFile, nil
}

func CreatePacFile(ctx context.Context, service *zscaler.Service, file *PACFileConfig) (*PACFileConfig, error) {
	resp, err := service.Client.Create(ctx, pacfileEndpoint, *file)
	if err != nil {
		return nil, fmt.Errorf("failed to create PAC file: %w", err)
	}

	createdPacFiles, ok := resp.(*PACFileConfig)
	if !ok {
		return nil, errors.New("object returned from api was not a pac file Pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning pac file from create: %d", createdPacFiles.ID)
	return createdPacFiles, nil
}

// UpdatePacFile updates a PAC file version with a specified action.
// pacID: the PAC file ID
// pacVersion: the PAC file version
// pacVersionAction: the action to be performed on the PAC file version (DEPLOY, STAGE, LKG, UNSTAGE, REMOVE_LKG)
// newLKGVer: optional, specify if removing LKG and want to assign a new version
func UpdatePacFile(ctx context.Context, service *zscaler.Service, pacID, pacVersion int, pacVersionAction string, file *PACFileConfig, newLKGVer *int) (*PACFileConfig, error) {
	// Construct the endpoint URL with pacID, pacVersion, and pacVersionAction
	endpoint := fmt.Sprintf("%s/%d/version/%d/action/%s", pacfileEndpoint, pacID, pacVersion, pacVersionAction)

	// If newLKGVer is provided, add it as a query parameter
	if newLKGVer != nil {
		endpoint = fmt.Sprintf("%s?newLKGVer=%d", endpoint, *newLKGVer)
	}

	// Send the request with the PACFileConfig payload
	resp, err := service.Client.UpdateWithPut(ctx, endpoint, *file)
	if err != nil {
		return nil, fmt.Errorf("failed to update PAC file version: %w", err)
	}

	// Parse the response into PACFileConfig
	updatedPacFile, ok := resp.(*PACFileConfig)
	if !ok {
		return nil, errors.New("object returned from API was not a PAC file pointer")
	}

	// Log the result for debugging
	service.Client.GetLogger().Printf("[DEBUG] Returning updated PAC file: %d", updatedPacFile.ID)

	return updatedPacFile, nil
}

func DeletePacFile(ctx context.Context, service *zscaler.Service, fileID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", pacfileEndpoint, fileID))
	if err != nil {
		return nil, fmt.Errorf("failed to delete PAC file: %w", err)
	}

	return nil, nil
}

func ValidatePacFile(ctx context.Context, service *zscaler.Service, pacContent string) (*PacResult, error) {
	// The URL for PAC validation
	url := fmt.Sprintf("%s/validate", pacfileEndpoint)

	// Use the CreateWithRawPayload function to send the raw PAC content
	resp, err := service.Client.CreateWithRawPayload(ctx, url, pacContent)
	if err != nil {
		return nil, fmt.Errorf("failed to validate PAC file: %w", err)
	}

	// Check if the response is empty
	if len(resp) == 0 {
		return nil, fmt.Errorf("no response received from the API")
	}

	// Unmarshal the response into PacResult
	var validationResult PacResult
	err = json.Unmarshal(resp, &validationResult)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal PAC validation response: %w", err)
	}

	// Log the result for debugging
	service.Client.GetLogger().Printf("[DEBUG] PAC validation result: %+v", validationResult)

	return &validationResult, nil
}

func CreateClonedPacFileVersion(ctx context.Context, service *zscaler.Service, pacID int, clonedPacVersion int, deleteVersion *int, file *PACFileConfig) (*PACFileConfig, error) {
	// Construct the endpoint URL with the provided pacID and clonedPacVersion
	endpoint := fmt.Sprintf("%s/%d/version/%d", pacfileEndpoint, pacID, clonedPacVersion)

	// If deleteVersion is provided (not nil), add it as a query parameter
	if deleteVersion != nil {
		endpoint = fmt.Sprintf("%s?deleteVersion=%d", endpoint, *deleteVersion)
	}

	// Send the request to create a new PAC file version by cloning the specified version
	resp, err := service.Client.Create(ctx, endpoint, *file)
	if err != nil {
		return nil, fmt.Errorf("failed to create PAC file version: %w", err)
	}

	// Parse the response
	createdPacFile, ok := resp.(*PACFileConfig)
	if !ok {
		return nil, errors.New("object returned from API was not a PAC file pointer")
	}

	// Log the result for debugging
	service.Client.GetLogger().Printf("[DEBUG] Returning PAC file version from create: %d", createdPacFile.ID)

	return createdPacFile, nil
}
