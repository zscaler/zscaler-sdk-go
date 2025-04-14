package ipsourcegroups

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
	ipSourceGroupsEndpoint = "/ztw/api/v1/ipSourceGroups"
)

type IPSourceGroups struct {
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

func Get(ctx context.Context, service *zscaler.Service, ipGroupID int) (*IPSourceGroups, error) {
	var ipSourceGroups IPSourceGroups
	err := service.Client.ReadResource(ctx, fmt.Sprintf("%s/%d", ipSourceGroupsEndpoint, ipGroupID), &ipSourceGroups)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning ip source group from Get: %d", ipSourceGroups.ID)
	return &ipSourceGroups, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, ipSourceGroupsName string) (*IPSourceGroups, error) {
	var ipSourceGroups []IPSourceGroups
	err := common.ReadAllPages(ctx, service.Client, ipSourceGroupsEndpoint, &ipSourceGroups)
	if err != nil {
		return nil, err
	}
	for _, ipSourceGroup := range ipSourceGroups {
		if strings.EqualFold(ipSourceGroup.Name, ipSourceGroupsName) {
			return &ipSourceGroup, nil
		}
	}
	return nil, fmt.Errorf("no ip source group found with name: %s", ipSourceGroupsName)
}

func Create(ctx context.Context, service *zscaler.Service, ipGroupID *IPSourceGroups) (*IPSourceGroups, error) {
	resp, err := service.Client.CreateResource(ctx, ipSourceGroupsEndpoint, *ipGroupID)
	if err != nil {
		return nil, err
	}

	createdIPSourceGroups, ok := resp.(*IPSourceGroups)
	if !ok {
		return nil, errors.New("object returned from api was not an ip source group pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning ip source group from create: %d", createdIPSourceGroups.ID)
	return createdIPSourceGroups, nil
}

func Update(ctx context.Context, service *zscaler.Service, ipGroupID int, ipGroup *IPSourceGroups) (*IPSourceGroups, error) {
	resp, err := service.Client.UpdateWithPutResource(ctx, fmt.Sprintf("%s/%d", ipSourceGroupsEndpoint, ipGroupID), *ipGroup)
	if err != nil {
		return nil, err
	}
	updatedIPSourceGroups, _ := resp.(*IPSourceGroups)

	service.Client.GetLogger().Printf("[DEBUG]returning ip source group from update: %d", updatedIPSourceGroups.ID)
	return updatedIPSourceGroups, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ipGroupID int) (*http.Response, error) {
	err := service.Client.DeleteResource(ctx, fmt.Sprintf("%s/%d", ipSourceGroupsEndpoint, ipGroupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAllLite(ctx context.Context, service *zscaler.Service) ([]IPSourceGroups, error) {
	var ipSourceGroups []IPSourceGroups
	err := common.ReadAllPages(ctx, service.Client, ipSourceGroupsEndpoint+"/lite", &ipSourceGroups)
	return ipSourceGroups, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]IPSourceGroups, error) {
	var ipSourceGroups []IPSourceGroups
	err := common.ReadAllPages(ctx, service.Client, ipSourceGroupsEndpoint, &ipSourceGroups)
	return ipSourceGroups, err
}
