package resourceservers

import (
	"context"
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zid/services/common"
)

const (
	resourceServerEndpoint = "/admin/api/v1/resource-servers"
)

type ResourceServers struct {
	ID            string          `json:"id,omitempty"`
	Name          string          `json:"name,omitempty"`
	DisplayName   string          `json:"displayName,omitempty"`
	Description   string          `json:"description,omitempty"`
	PrimaryAud    string          `json:"primaryAud,omitempty"`
	DefaultApi    bool            `json:"defaultApi,omitempty"`
	ServiceScopes []ServiceScopes `json:"serviceScopes,omitempty"`
}

type ServiceScopes struct {
	Service Service  `json:"service,omitempty"`
	Scopes  []Scopes `json:"scopes,omitempty"`
}

type Service struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	CloudName   string `json:"cloudName,omitempty"`
	OrgName     string `json:"orgName,omitempty"`
}

type Scopes struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ResourceServersResponse = common.PaginationResponse[ResourceServers]

func Get(ctx context.Context, service *zscaler.Service, resourceID string) (*ResourceServers, error) {
	var resourceServer ResourceServers
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%s", resourceServerEndpoint, resourceID), &resourceServer)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning resource server from Get: %s", resourceServer.ID)
	return &resourceServer, nil
}

// GetAll retrieves all resource servers with optional pagination and filtering parameters
func GetAll(ctx context.Context, service *zscaler.Service, queryParams *common.PaginationQueryParams) ([]ResourceServers, error) {
	return common.ReadAllPagesWithPagination[ResourceServers](ctx, service.Client, resourceServerEndpoint, queryParams)
}

// GetByName retrieves resource servers by searching through paginated data for the specified name
func GetByName(ctx context.Context, service *zscaler.Service, name string) ([]ResourceServers, error) {
	var allResources []ResourceServers
	var currentOffset int
	pageSize := 100 // Use a reasonable page size for searching

	for {
		// Create query params for current page
		queryParams := common.NewPaginationQueryParams(pageSize)
		queryParams.WithOffset(currentOffset)

		// Get current page
		pageResponse, err := common.ReadPageWithPagination[ResourceServers](ctx, service.Client, resourceServerEndpoint, &queryParams)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page at offset %d: %w", currentOffset, err)
		}

		// Search through records in this page
		for _, group := range pageResponse.Records {
			if strings.Contains(strings.ToLower(group.Name), strings.ToLower(name)) {
				allResources = append(allResources, group)
			}
		}

		// Check if we've reached the end
		if len(pageResponse.Records) < pageSize || pageResponse.NextLink == "" {
			break
		}

		// Move to next page
		currentOffset += len(pageResponse.Records)
	}

	return allResources, nil
}
