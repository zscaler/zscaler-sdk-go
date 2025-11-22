package ipdestinationgroups

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	ipDestinationGroupsEndpoint = "/zia/api/v1/ipDestinationGroups"
)

type IPDestinationGroups struct {
	// Unique identifer for the destination IP group
	ID int `json:"id"`

	// Destination IP group name
	Name string `json:"name,omitempty"`

	// Additional information about the destination IP group
	Description string `json:"description,omitempty"`

	// Destination IP group type (i.e., the group can contain destination IP addresses or FQDNs)
	Type string `json:"type,omitempty"`

	// Destination IP addresses, FQDNs, or wildcard FQDNs added to the group.
	Addresses []string `json:"addresses,omitempty"`

	// Destination IP address URL categories. You can identify destinations based on the URL category of the domain.
	IPCategories []string `json:"ipCategories,omitempty"`

	// Destination IP address countries. You can identify destinations based on the location of a server.
	Countries []string `json:"countries,omitempty"`

	// If set to true, the destination IP address group is non-editable. This field is applicable only to predefined IP address groups, which cannot be modified.
	IsNonEditable bool `json:"isNonEditable,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, ipGroupID int) (*IPDestinationGroups, error) {
	var ipDestinationGroups IPDestinationGroups
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", ipDestinationGroupsEndpoint, ipGroupID), &ipDestinationGroups)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning ip destination group from Get: %d", ipDestinationGroups.ID)
	return &ipDestinationGroups, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, ipDestinationGroupsName string) (*IPDestinationGroups, error) {
	// Use GetAll to leverage API and verify exact match
	ipDestinationGroups, err := GetAll(ctx, service, "")
	if err != nil {
		return nil, err
	}
	for _, ipDestinationGroup := range ipDestinationGroups {
		if strings.EqualFold(ipDestinationGroup.Name, ipDestinationGroupsName) {
			return &ipDestinationGroup, nil
		}
	}
	return nil, fmt.Errorf("no ip destination group found with name: %s", ipDestinationGroupsName)
}

func Create(ctx context.Context, service *zscaler.Service, ipGroupID *IPDestinationGroups) (*IPDestinationGroups, error) {
	resp, err := service.Client.Create(ctx, ipDestinationGroupsEndpoint, *ipGroupID)
	if err != nil {
		return nil, err
	}

	createdIPDestinationGroups, ok := resp.(*IPDestinationGroups)
	if !ok {
		return nil, errors.New("object returned from api was not an ip destination group pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning ip destination group from create: %d", createdIPDestinationGroups.ID)
	return createdIPDestinationGroups, nil
}

func Update(ctx context.Context, service *zscaler.Service, ipGroupID int, ipGroup *IPDestinationGroups, override *bool) (*IPDestinationGroups, *http.Response, error) {
	// Build the URL with optional override query parameter
	url := fmt.Sprintf("%s/%d", ipDestinationGroupsEndpoint, ipGroupID)
	if override != nil {
		url = fmt.Sprintf("%s?override=%t", url, *override)
	}

	resp, err := service.Client.UpdateWithPut(ctx, url, *ipGroup)
	if err != nil {
		return nil, nil, err
	}
	updatedIPDestinationGroups, _ := resp.(*IPDestinationGroups)

	service.Client.GetLogger().Printf("[DEBUG]returning ip destination group from update: %d", updatedIPDestinationGroups.ID)
	return updatedIPDestinationGroups, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ipGroupID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", ipDestinationGroupsEndpoint, ipGroupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetAll retrieves all IP destination groups with optional type filter.
// The API doesn't support pagination and returns all results in a single response.
// excludeType parameter filters based on the IP destination group's type.
// Valid values: "DSTN_IP", "DSTN_FQDN", "DSTN_DOMAIN", "DSTN_OTHER", or "" (no filter)
func GetAll(ctx context.Context, service *zscaler.Service, excludeType string) ([]IPDestinationGroups, error) {
	var ipDestinationGroups []IPDestinationGroups
	endpoint := ipDestinationGroupsEndpoint

	// Build query parameters
	queryParams := url.Values{}
	if excludeType != "" {
		queryParams.Add("excludeType", excludeType)
	}

	// Build endpoint with query parameters
	baseQuery := queryParams.Encode()
	if baseQuery != "" {
		endpoint += "?" + baseQuery
	}

	// Use service.Client.Read directly since the API doesn't support pagination
	err := service.Client.Read(ctx, endpoint, &ipDestinationGroups)
	return ipDestinationGroups, err
}
