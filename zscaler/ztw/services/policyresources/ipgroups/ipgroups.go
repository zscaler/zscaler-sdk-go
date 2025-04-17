package ipgroups

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
)

const (
	ipGroupsEndpoint = "/ztw/api/v1/ipGroups"
)

type IPGroups struct {
	// A unique identifier of the source IP address group.
	ID int `json:"id"`

	// The name of the source IP address group.
	Name string `json:"name,omitempty"`

	// The description of the source IP address group.
	Description string `json:"description,omitempty"`

	// Source IP addresses added to the group.
	IPAddresses []string `json:"ipAddresses,omitempty"`

	// Indicates whether the IP group or IP pool is created in Cloud & Branch Connector (EC) or ZIA.
	CreatorContext string `json:"creatorContext,omitempty"`

	// If set to true, the destination IP address group is non-editable. This field is applicable only to predefined IP address groups, which cannot be modified.
	IsNonEditable bool `json:"isNonEditable,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, ipGroupID int) (*IPGroups, error) {
	var ipGroups IPGroups
	err := service.Client.ReadResource(ctx, fmt.Sprintf("%s/%d", ipGroupsEndpoint, ipGroupID), &ipGroups)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning ip group from Get: %d", ipGroups.ID)
	return &ipGroups, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, ipGroupsName string) (*IPGroups, error) {
	var ipGroups []IPGroups
	err := common.ReadAllPages(ctx, service.Client, ipGroupsEndpoint, &ipGroups)
	if err != nil {
		return nil, err
	}
	for _, ipGroup := range ipGroups {
		if strings.EqualFold(ipGroup.Name, ipGroupsName) {
			return &ipGroup, nil
		}
	}
	return nil, fmt.Errorf("no ip group found with name: %s", ipGroupsName)
}

func Create(ctx context.Context, service *zscaler.Service, ipGroupID *IPGroups) (*IPGroups, error) {
	resp, err := service.Client.CreateResource(ctx, ipGroupsEndpoint, *ipGroupID)
	if err != nil {
		return nil, err
	}

	createdIPGroups, ok := resp.(*IPGroups)
	if !ok {
		return nil, errors.New("object returned from api was not an ip group pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning ip group from create: %d", createdIPGroups.ID)
	return createdIPGroups, nil
}

func Update(ctx context.Context, service *zscaler.Service, ipGroupID int, ipGroup *IPGroups) (*IPGroups, error) {
	resp, err := service.Client.UpdateWithPutResource(ctx, fmt.Sprintf("%s/%d", ipGroupsEndpoint, ipGroupID), *ipGroup)
	if err != nil {
		return nil, err
	}
	updatedIPGroups, _ := resp.(*IPGroups)

	service.Client.GetLogger().Printf("[DEBUG]returning ip group from update: %d", updatedIPGroups.ID)
	return updatedIPGroups, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ipGroupID int) (*http.Response, error) {
	err := service.Client.DeleteResource(ctx, fmt.Sprintf("%s/%d", ipGroupsEndpoint, ipGroupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetIPGroupLite(ctx context.Context, service *zscaler.Service) ([]IPGroups, error) {
	var ipGroups []IPGroups
	err := common.ReadAllPages(ctx, service.Client, ipGroupsEndpoint+"/lite", &ipGroups)
	return ipGroups, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]IPGroups, error) {
	var ipGroups []IPGroups
	err := common.ReadAllPages(ctx, service.Client, ipGroupsEndpoint, &ipGroups)
	return ipGroups, err
}
