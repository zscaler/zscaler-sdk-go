package groups

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"encoding/json"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zid/services/common"
)

const (
	groupsEndpoint = "/admin/api/v1/groups"
)

// Group represents a single group record
type Groups struct {
	Name                      string                    `json:"name,omitempty"`
	Description               string                    `json:"description,omitempty"`
	ID                        string                    `json:"id,omitempty"`
	Source                    string                    `json:"source,omitempty"`
	IsDynamicGroup            bool                      `json:"isDynamicGroup,omitempty"`
	DynamicGroup              bool                      `json:"dynamicGroup,omitempty"`
	AdminEntitlementEnabled   bool                      `json:"adminEntitlementEnabled,omitempty"`
	ServiceEntitlementEnabled bool                      `json:"serviceEntitlementEnabled,omitempty"`
	IDP                       *common.IDNameDisplayName `json:"idp,omitempty"`
}

type UserID struct {
	ID string `json:"id"`
}

type GroupsResponse = common.PaginationResponse[Groups]

func Get(ctx context.Context, service *zscaler.Service, groupID string) (*Groups, error) {
	var group Groups
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%s", groupsEndpoint, groupID), &group)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning group from Get: %s", group.ID)
	return &group, nil
}

// GetAll retrieves all groups with optional pagination and filtering parameters
func GetAll(ctx context.Context, service *zscaler.Service, queryParams *common.PaginationQueryParams) ([]Groups, error) {
	return common.ReadAllPagesWithPagination[Groups](ctx, service.Client, groupsEndpoint, queryParams)
}

// GetByName retrieves groups by searching through paginated data for the specified name
func GetByName(ctx context.Context, service *zscaler.Service, name string) ([]Groups, error) {
	var allGroups []Groups
	var currentOffset int
	pageSize := 100 // Use a reasonable page size for searching

	for {
		// Create query params for current page
		queryParams := common.NewPaginationQueryParams(pageSize)
		queryParams.WithOffset(currentOffset)

		// Get current page
		pageResponse, err := common.ReadPageWithPagination[Groups](ctx, service.Client, groupsEndpoint, &queryParams)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page at offset %d: %w", currentOffset, err)
		}

		// Search through records in this page
		for _, group := range pageResponse.Records {
			if strings.Contains(strings.ToLower(group.Name), strings.ToLower(name)) {
				allGroups = append(allGroups, group)
			}
		}

		// Check if we've reached the end
		if len(pageResponse.Records) < pageSize || pageResponse.NextLink == "" {
			break
		}

		// Move to next page
		currentOffset += len(pageResponse.Records)
	}

	return allGroups, nil
}

// GetUsers retrieves users within a specific group
func GetUsers(ctx context.Context, service *zscaler.Service, groupID string, queryParams *common.PaginationQueryParams) ([]interface{}, error) {
	usersEndpoint := fmt.Sprintf("%s/%s/users", groupsEndpoint, groupID)
	return common.ReadAllPagesWithPagination[interface{}](ctx, service.Client, usersEndpoint, queryParams)
}

func Create(ctx context.Context, service *zscaler.Service, groups *Groups) (*Groups, *http.Response, error) {
	resp, err := service.Client.Create(ctx, groupsEndpoint, *groups)
	if err != nil {
		return nil, nil, err
	}

	createdGroup, ok := resp.(*Groups)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a group pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new group from create: %d", createdGroup.ID)
	return createdGroup, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, groupID int, groups *Groups) (*Groups, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", groupsEndpoint, groupID), *groups)
	if err != nil {
		return nil, nil, err
	}
	updatedGroup, _ := resp.(*Groups)

	service.Client.GetLogger().Printf("[DEBUG]returning updates group from update: %d", updatedGroup.ID)
	return updatedGroup, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, groupID string) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%s", groupsEndpoint, groupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func AddUserListToGroup(ctx context.Context, service *zscaler.Service, groupID string, userIDs []string) (*Groups, *http.Response, error) {
	var payload []UserID
	for _, userID := range userIDs {
		payload = append(payload, UserID{ID: userID})
	}

	resp, err := service.Client.CreateWithSlicePayload(ctx, fmt.Sprintf("%s/%s/users", groupsEndpoint, groupID), payload)
	if err != nil {
		return nil, nil, err
	}

	var createdGroup Groups
	if len(resp) > 0 {
		err = json.Unmarshal(resp, &createdGroup)
		if err != nil {
			return nil, nil, err
		}
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new group from create: %s", createdGroup.ID)
	return &createdGroup, nil, nil
}

func ReplaceUserListInGroup(ctx context.Context, service *zscaler.Service, groupID string, userIDs []string) (*Groups, *http.Response, error) {
	var payload []UserID
	for _, userID := range userIDs {
		payload = append(payload, UserID{ID: userID})
	}

	resp, err := service.Client.UpdateWithSlicePayload(ctx, fmt.Sprintf("%s/%s/users", groupsEndpoint, groupID), payload)
	if err != nil {
		return nil, nil, err
	}
	var updatedGroup Groups
	if len(resp) > 0 {
		err = json.Unmarshal(resp, &updatedGroup)
		if err != nil {
			return nil, nil, err
		}
	}

	service.Client.GetLogger().Printf("[DEBUG]returning updates group from update: %s", updatedGroup.ID)
	return &updatedGroup, nil, nil
}

func AddUserToGroup(ctx context.Context, service *zscaler.Service, groupID, userID string) (*http.Response, error) {
	// Pass empty struct since Create requires a struct, not nil
	emptyPayload := struct{}{}
	resp, err := service.Client.Create(ctx, fmt.Sprintf("%s/%s/users/%s", groupsEndpoint, groupID, userID), emptyPayload)
	if err != nil {
		return nil, err
	}

	httpResp, _ := resp.(*http.Response)
	return httpResp, nil
}

func DeleteUserFromGroup(ctx context.Context, service *zscaler.Service, groupID, userID string) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%s/users/%s", groupsEndpoint, groupID, userID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
