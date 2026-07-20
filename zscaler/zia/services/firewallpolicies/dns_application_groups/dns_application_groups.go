package dns_application_groups

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	dnsApplicationGroupsEndpoint = "/zia/api/v1/dnsApplicationGroups"
)

type DnsApplicationGroup struct {
	ID              int      `json:"id,omitempty"`
	Name            string   `json:"name,omitempty"`
	DnsApplications []string `json:"dnsApplications,omitempty"`
	Description     string   `json:"description,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, groupID int) (*DnsApplicationGroup, error) {
	var dnsApplicationGroups DnsApplicationGroup
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", dnsApplicationGroupsEndpoint, groupID), &dnsApplicationGroups)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning dns application group from Get: %d", dnsApplicationGroups.ID)
	return &dnsApplicationGroups, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, groupName string) (*DnsApplicationGroup, error) {
	// Use GetAll to leverage API and verify exact match
	dnsGroups, err := GetAll(ctx, service)
	if err != nil {
		return nil, err
	}
	for _, dnsGroup := range dnsGroups {
		if strings.EqualFold(dnsGroup.Name, groupName) {
			return &dnsGroup, nil
		}
	}
	return nil, fmt.Errorf("no dns application group found with name: %s", groupName)
}

func Create(ctx context.Context, service *zscaler.Service, groupID *DnsApplicationGroup) (*DnsApplicationGroup, error) {
	resp, err := service.Client.Create(ctx, dnsApplicationGroupsEndpoint, *groupID)
	if err != nil {
		return nil, err
	}

	createdDnsApplicationGroup, ok := resp.(*DnsApplicationGroup)
	if !ok {
		return nil, errors.New("object returned from api was not an dns application group pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning dns application group from create: %d", createdDnsApplicationGroup.ID)
	return createdDnsApplicationGroup, nil
}

func Update(ctx context.Context, service *zscaler.Service, groupID int, dnsGroup *DnsApplicationGroup) (*DnsApplicationGroup, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", dnsApplicationGroupsEndpoint, groupID), *dnsGroup)
	if err != nil {
		return nil, err
	}
	updatedDNSApplicationGroup, _ := resp.(*DnsApplicationGroup)

	service.Client.GetLogger().Printf("[DEBUG]returning dns application group from update: %d", updatedDNSApplicationGroup.ID)
	return updatedDNSApplicationGroup, nil
}

func Delete(ctx context.Context, service *zscaler.Service, groupID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", dnsApplicationGroupsEndpoint, groupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]DnsApplicationGroup, error) {
	var dnsGroups []DnsApplicationGroup
	err := service.Client.Read(ctx, dnsApplicationGroupsEndpoint, &dnsGroups)
	return dnsGroups, err
}
