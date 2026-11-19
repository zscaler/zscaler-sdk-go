package networkservices

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	networkServicesEndpoint = "/zia/api/v1/networkServices"

	// networkServicesPageSize is the number of services requested per page.
	// The API rejects a pageSize above its limit with INVALID_INPUT_ARGUMENT
	// "Invalid Page Size" (1,000 is rejected on some clouds), and pages with a
	// default size when pageSize is omitted, so a single unpaged read can miss
	// services. 127 is accepted on every tenant.
	networkServicesPageSize = 127
)

type NetworkServices struct {
	ID            int            `json:"id"`
	Name          string         `json:"name,omitempty"`
	Tag           string         `json:"tag,omitempty"`
	SrcTCPPorts   []NetworkPorts `json:"srcTcpPorts,omitempty"`
	DestTCPPorts  []NetworkPorts `json:"destTcpPorts,omitempty"`
	SrcUDPPorts   []NetworkPorts `json:"srcUdpPorts,omitempty"`
	DestUDPPorts  []NetworkPorts `json:"destUdpPorts,omitempty"`
	Type          string         `json:"type,omitempty"`
	Description   string         `json:"description,omitempty"`
	Protocol      string         `json:"protocol,omitempty"`
	IsNameL10nTag bool           `json:"isNameL10nTag,omitempty"`
}

type NetworkPorts struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, serviceID int) (*NetworkServices, error) {
	var networkServices NetworkServices
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", networkServicesEndpoint, serviceID), &networkServices)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning network services from Get: %d", networkServices.ID)
	return &networkServices, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, networkServiceName string, protocol, locale *string) (*NetworkServices, error) {
	networkServices, err := GetAllNetworkServices(ctx, service, protocol, locale)
	if err != nil {
		return nil, err
	}

	// If name is empty and protocol or locale is provided, return the first matching service
	if networkServiceName == "" {
		if protocol != nil || locale != nil {
			if len(networkServices) > 0 {
				return &networkServices[0], nil
			}
			return nil, fmt.Errorf("no network services found with the provided filters (protocol: %v, locale: %v)", protocol, locale)
		}
		return nil, fmt.Errorf("name parameter is required when protocol and locale are not provided")
	}

	// Search for a service matching the name
	for _, networkService := range networkServices {
		if strings.EqualFold(networkService.Name, networkServiceName) {
			return &networkService, nil
		}
	}
	return nil, fmt.Errorf("no network services found with name: %s", networkServiceName)
}

func Create(ctx context.Context, service *zscaler.Service, networkService *NetworkServices) (*NetworkServices, error) {
	resp, err := service.Client.Create(ctx, networkServicesEndpoint, *networkService)
	if err != nil {
		return nil, err
	}

	createdNetworkServices, ok := resp.(*NetworkServices)
	if !ok {
		return nil, errors.New("object returned from api was not a network service pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning network service from create: %d", createdNetworkServices.ID)
	return createdNetworkServices, nil
}

func Update(ctx context.Context, service *zscaler.Service, serviceID int, networkService *NetworkServices) (*NetworkServices, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", networkServicesEndpoint, serviceID), *networkService)
	if err != nil {
		return nil, nil, err
	}
	updatedNetworkServices, _ := resp.(*NetworkServices)

	service.Client.GetLogger().Printf("[DEBUG]returning network service from Update: %d", updatedNetworkServices.ID)
	return updatedNetworkServices, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, serviceID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", networkServicesEndpoint, serviceID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetAllNetworkServices retrieves all network services, optionally filtered by
// protocol and localized by locale. The filters are sent on every page. The
// optional JMESPath filter from the context is applied to the combined result.
//
// It pages with its own loop rather than common.ReadAllPages, whose default
// page size of 1,000 this endpoint rejects on some clouds.
func GetAllNetworkServices(ctx context.Context, service *zscaler.Service, protocol, locale *string) ([]NetworkServices, error) {
	var networkServices []NetworkServices
	seen := make(map[int]struct{})

	for page := 1; ; page++ {
		queryParams := url.Values{}
		if protocol != nil && *protocol != "" {
			queryParams.Set("protocol", *protocol)
		}
		if locale != nil && *locale != "" {
			queryParams.Set("locale", *locale)
		}
		queryParams.Set("page", strconv.Itoa(page))
		queryParams.Set("pageSize", strconv.Itoa(networkServicesPageSize))
		endpoint := fmt.Sprintf("%s?%s", networkServicesEndpoint, queryParams.Encode())

		var pageItems []NetworkServices
		if err := service.Client.Read(ctx, endpoint, &pageItems); err != nil {
			return nil, err
		}

		added := 0
		for _, svc := range pageItems {
			if _, ok := seen[svc.ID]; ok {
				continue
			}
			seen[svc.ID] = struct{}{}
			networkServices = append(networkServices, svc)
			added++
		}

		// A short page is the last one. Also stop if the API returned more
		// services than requested (it ignored pageSize and sent the whole list)
		// or only services already read (it ignored page): requesting the next
		// page would return the same services forever.
		if len(pageItems) != networkServicesPageSize || added == 0 {
			break
		}
	}

	return zscaler.ApplyJMESPathFromContext(ctx, networkServices)
}
