package networkapplicationgroups

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
	networkAppGroupsEndpoint = "/zia/api/v1/networkApplicationGroups"
)

type NetworkApplicationGroups struct {
	ID                  int      `json:"id"`
	Name                string   `json:"name,omitempty"`
	NetworkApplications []string `json:"networkApplications,omitempty"`
	Description         string   `json:"description,omitempty"`
}

func GetNetworkApplicationGroups(ctx context.Context, service *zscaler.Service, groupID int) (*NetworkApplicationGroups, error) {
	var networkApplicationGroups NetworkApplicationGroups
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", networkAppGroupsEndpoint, groupID), &networkApplicationGroups)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning network application groups from Get: %d", networkApplicationGroups.ID)
	return &networkApplicationGroups, nil
}

func GetNetworkApplicationGroupsByName(ctx context.Context, service *zscaler.Service, appGroupsName string) (*NetworkApplicationGroups, error) {
	var networkApplicationGroups []NetworkApplicationGroups
	err := common.ReadAllPages(ctx, service.Client, networkAppGroupsEndpoint, &networkApplicationGroups)
	if err != nil {
		return nil, err
	}
	for _, networkAppGroup := range networkApplicationGroups {
		if strings.EqualFold(networkAppGroup.Name, appGroupsName) {
			return &networkAppGroup, nil
		}
	}
	return nil, fmt.Errorf("no network application groups found with name: %s", appGroupsName)
}

func Create(ctx context.Context, service *zscaler.Service, applicationGroup *NetworkApplicationGroups) (*NetworkApplicationGroups, error) {
	resp, err := service.Client.Create(ctx, networkAppGroupsEndpoint, *applicationGroup)
	if err != nil {
		return nil, err
	}

	createdApplicationGroups, ok := resp.(*NetworkApplicationGroups)
	if !ok {
		return nil, errors.New("object returned from api was not a network application groups pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning network application groups from create: %d", createdApplicationGroups.ID)
	return createdApplicationGroups, nil
}

func Update(ctx context.Context, service *zscaler.Service, groupID int, applicationGroup *NetworkApplicationGroups) (*NetworkApplicationGroups, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", networkAppGroupsEndpoint, groupID), *applicationGroup)
	if err != nil {
		return nil, nil, err
	}
	updatedApplicationGroups, _ := resp.(*NetworkApplicationGroups)

	service.Client.GetLogger().Printf("[DEBUG]returning network application groups from Update: %d", updatedApplicationGroups.ID)
	return updatedApplicationGroups, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, groupID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", networkAppGroupsEndpoint, groupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAllNetworkApplicationGroups(ctx context.Context, service *zscaler.Service) ([]NetworkApplicationGroups, error) {
	var networkApplicationGroups []NetworkApplicationGroups
	err := common.ReadAllPages(ctx, service.Client, networkAppGroupsEndpoint, &networkApplicationGroups)
	return networkApplicationGroups, err
}
