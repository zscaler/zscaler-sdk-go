package custom_file_types

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
	customFilefileType          = "/zia/api/v1/customFileTypes"
	customFileTypeCountEndpoint = "/zia/api/v1/customFileTypes/count"
)

type CustomFileTypes struct {
	// System generated identifier for a file-type policy
	ID int `json:"id,omitempty"`

	// The name of the File Type rule
	Name string `json:"name,omitempty"`

	// Additional information about the custom file type, if any.
	Description string `json:"description,omitempty"`

	// Specifies the file type extension. The maximum extension length is 10 characters. Existing Zscaler extensions cannot be added to custom file types.
	Extension string `json:"extension,omitempty"`

	// File type ID. This ID is assigned and maintained for all file types including predefined and custom file types, and this value is different from the custom file type ID.
	FileTypeID int `json:"fileTypeId,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, fileID int) (*CustomFileTypes, error) {
	var fileTypes CustomFileTypes
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", customFilefileType, fileID), &fileTypes)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning custom file types from Get: %d", fileTypes.ID)
	return &fileTypes, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, ruleName string) (*CustomFileTypes, error) {
	var customFileTypes []CustomFileTypes
	err := common.ReadAllPages(ctx, service.Client, customFilefileType, &customFileTypes)
	if err != nil {
		return nil, err
	}
	for _, fileTypeControlRule := range customFileTypes {
		if strings.EqualFold(fileTypeControlRule.Name, ruleName) {
			return &fileTypeControlRule, nil
		}
	}
	return nil, fmt.Errorf("no custom file types found with name: %s", ruleName)
}

func Create(ctx context.Context, service *zscaler.Service, ruleID *CustomFileTypes) (*CustomFileTypes, error) {
	resp, err := service.Client.Create(ctx, customFilefileType, *ruleID)
	if err != nil {
		return nil, err
	}

	createdCustomFileType, ok := resp.(*CustomFileTypes)
	if !ok {
		return nil, errors.New("object returned from api was not a custom file types pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning custom file types from create: %d", createdCustomFileType.ID)
	return createdCustomFileType, nil
}

func Update(ctx context.Context, service *zscaler.Service, fileID int, custotmFileTypes *CustomFileTypes) (*CustomFileTypes, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", customFilefileType, fileID), *custotmFileTypes)
	if err != nil {
		return nil, err
	}
	updatedCustomFileType, _ := resp.(*CustomFileTypes)

	service.Client.GetLogger().Printf("[DEBUG]returning updates from custom file types from update: %d", updatedCustomFileType.ID)
	return updatedCustomFileType, nil
}

func Delete(ctx context.Context, service *zscaler.Service, fileID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", customFilefileType, fileID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetCustomFileTypes(ctx context.Context, service *zscaler.Service) ([]CustomFileTypes, error) {
	var customFileTypes []CustomFileTypes
	err := common.ReadAllPages(ctx, service.Client, customFilefileType, &customFileTypes)
	return customFileTypes, err
}

// GetCustomFileTypeCountFilterOptions represents optional filter parameters for GetCustomFileTypeCount
type GetCustomFileTypeCountFilterOptions struct {
	// Query string to filter custom file types by name or other attributes
	Query *string
}

// GetCustomFileTypeCount retrieves the count of custom file types available.
// The API returns a simple integer count.
func GetCustomFileTypeCount(ctx context.Context, service *zscaler.Service, opts *GetCustomFileTypeCountFilterOptions) (int, error) {
	var count int
	endpoint := customFileTypeCountEndpoint

	// Build query parameters from filter options
	queryParams := url.Values{}
	if opts != nil && opts.Query != nil && *opts.Query != "" {
		queryParams.Add("query", *opts.Query)
	}

	// Build endpoint with query parameters
	baseQuery := queryParams.Encode()
	if baseQuery != "" {
		endpoint += "?" + baseQuery
	}

	err := service.Client.Read(ctx, endpoint, &count)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve custom file type count: %w", err)
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning custom file type count: %d", count)
	return count, nil
}
