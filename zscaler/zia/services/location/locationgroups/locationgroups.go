package locationgroups

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	locationGroupEndpoint = "/zia/api/v1/locations/groups"
)

type LocationGroup struct {
	// Unique identifier for the location group
	ID int `json:"id,omitempty"`

	// Location group name
	Name string `json:"name,omitempty"`

	// Indicates the location group was deleted
	Deleted bool `json:"deleted,omitempty"`

	// The location group's type (i.e., Static or Dynamic)
	GroupType string `json:"groupType,omitempty"`

	// A dynamic location group's criteria. This is ignored if the groupType is Static.
	DynamicLocationGroupCriteria *DynamicLocationGroupCriteria `json:"dynamicLocationGroupCriteria,omitempty"`

	// Additional information about the location group
	Comments string `json:"comments"`

	// The Name-ID pairs of the locations that are assigned to the static location group. This is ignored if the groupType is Dynamic.
	Locations []common.IDNameExtensions `json:"locations"`

	// Automatically populated with the current ZIA admin user, after a successful POST or PUT request.
	LastModUser *LastModUser `json:"lastModUser"`

	// Automatically populated with the current time, after a successful POST or PUT request.
	LastModTime int  `json:"lastModTime"`
	Predefined  bool `json:"predefined"`
}

type DynamicLocationGroupCriteria struct {
	// A sub-string to match location name. Valid operators are contains, starts with, and ends with",
	Name *Name `json:"name,omitempty"`

	// One or more countries from a predefined set
	Countries []string `json:"countries,omitempty"`

	// A sub-string to match city. Valid operators are starts with, ends with, contains, and exact match operators.
	City *City `json:"city,omitempty"`

	// One or more values from a predefined set of SD-WAN partner list to display partner names.
	ManagedBy []ManagedBy `json:"managedBy,omitempty"`

	// Enforce Authentication. Required when ports are enabled, IP Surrogate is enabled, or Kerberos Authentication is enabled.
	EnforceAuthentication bool `json:"enforceAuthentication"`

	// Enable AUP. When set to true, AUP is enabled for the location.
	EnforceAup bool `json:"enforceAup"`

	// Enable Firewall. When set to true, Firewall is enabled for the location.
	EnforceFirewallControl bool `json:"enforceFirewallControl"`

	// Enable XFF Forwarding. When set to true, traffic is passed to Zscaler Cloud via the X-Forwarded-For (XFF) header.
	EnableXffForwarding bool `json:"enableXffForwarding"`

	// Enable Caution. When set to true, a caution notifcation is enabled for the location.
	EnableCaution bool `json:"enableCaution"`

	// Enable Bandwidth Control. When set to true, Bandwidth Control is enabled for the location.
	EnableBandwidthControl bool `json:"enableBandwidthControl"`

	// One or more location profiles from a predefined set
	Profiles []string `json:"profiles"`
}

type Name struct {
	// String value to be matched or partially matched
	MatchString string `json:"matchString,omitempty"`

	// Operator that performs match action
	MatchType string `json:"matchType,omitempty"`
}

type City struct {
	// String value to be matched or partially matched
	MatchString string `json:"matchString,omitempty"`

	// Operator that performs match action
	MatchType string `json:"matchType,omitempty"`
}

type LastModUser struct {
	// Identifier that uniquely identifies an entity
	ID int `json:"id,omitempty"`

	// The configured name of the entity
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type ManagedBy struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

func GetLocationGroup(ctx context.Context, service *zscaler.Service, groupID int) (*LocationGroup, error) {
	var locationGroup LocationGroup
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", locationGroupEndpoint, groupID), &locationGroup)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]returning location group from Get: %d", locationGroup.ID)
	return &locationGroup, nil
}

func GetLocationGroupByName(ctx context.Context, service *zscaler.Service, locationGroupName string) (*LocationGroup, error) {
	var locationGroups []LocationGroup
	err := common.ReadAllPages(ctx, service.Client, fmt.Sprintf("%s?name=%s", locationGroupEndpoint, url.QueryEscape(locationGroupName)), &locationGroups)
	if err != nil {
		return nil, err
	}
	for _, locationGroup := range locationGroups {
		if strings.EqualFold(locationGroup.Name, locationGroupName) {
			return &locationGroup, nil
		}
	}
	return nil, fmt.Errorf("no location group found with name: %s", locationGroupName)
}

// GetGroupType queries the location group by its type
func GetGroupType(ctx context.Context, service *zscaler.Service, gType string) (*LocationGroup, error) {
	var groupTypes []LocationGroup
	err := service.Client.Read(ctx, fmt.Sprintf("%s?groupType=%s", locationGroupEndpoint, url.QueryEscape(gType)), &groupTypes)
	if err != nil {
		return nil, err
	}
	for _, locationGroup := range groupTypes {
		if strings.EqualFold(locationGroup.GroupType, gType) {
			return &locationGroup, nil
		}
	}
	return nil, fmt.Errorf("no group type found with name: %s", gType)
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]LocationGroup, error) {
	var locationGroups []LocationGroup
	err := common.ReadAllPages(ctx, service.Client, locationGroupEndpoint, &locationGroups)
	return locationGroups, err
}
