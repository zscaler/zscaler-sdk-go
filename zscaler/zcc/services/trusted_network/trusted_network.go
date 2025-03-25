package trusted_network

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	trustedNetworkEndpoint = "/zcc/papi/public/v1/webTrustedNetwork"
)

type TrustedNetwork struct {
	ID                     string `json:"id"`
	Active                 bool   `json:"active"`
	CompanyID              string `json:"companyId"`
	ConditionType          int    `json:"conditionType"`
	CreatedBy              string `json:"createdBy"`
	DnsSearchDomains       string `json:"dnsSearchDomains"`
	DnsServers             string `json:"dnsServers"`
	EditedBy               string `json:"editedBy"`
	Guid                   string `json:"guid"`
	Hostnames              string `json:"hostnames"`
	NetworkName            string `json:"networkName"`
	ResolvedIpsForHostname string `json:"resolvedIpsForHostname"`
	Ssids                  string `json:"ssids"`
	TrustedDhcpServers     string `json:"trustedDhcpServers"`
	TrustedEgressIps       string `json:"trustedEgressIps"`
	TrustedGateways        string `json:"trustedGateways"`
	TrustedSubnets         string `json:"trustedSubnets"`
}

type TrustedNetworksResponse struct {
	TotalCount              int              `json:"totalCount"`
	TrustedNetworkContracts []TrustedNetwork `json:"trustedNetworkContracts"`
}

func GetMultipleTrustedNetworks(ctx context.Context, service *zscaler.Service, search, searchType string, page, pageSize *int) (*TrustedNetworksResponse, *http.Response, error) {
	endpoint := fmt.Sprintf("%s/listByCompany", trustedNetworkEndpoint)

	queryParams := common.QueryParams{
		Search: search,
		// SearchType: searchType,
	}
	if page != nil {
		queryParams.Page = *page
	}
	if pageSize != nil {
		queryParams.PageSize = *pageSize
	}

	var response TrustedNetworksResponse

	resp, err := service.Client.NewZccRequestDo(ctx, "GET", endpoint, queryParams, nil, &response)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to fetch trusted networks: %w", err)
	}

	if response.TotalCount > 0 && len(response.TrustedNetworkContracts) == 0 {
		return nil, resp, fmt.Errorf("totalCount is %d, but no trusted networks are returned", response.TotalCount)
	}

	return &response, resp, nil
}

func GetTrustedNetworkByName(ctx context.Context, service *zscaler.Service, name string) (*TrustedNetwork, *http.Response, error) {
	pageSize := 1000
	page := 1

	for {
		res, resp, err := GetMultipleTrustedNetworks(ctx, service, "", "", &page, &pageSize)
		if err != nil {
			return nil, resp, err
		}

		for _, tn := range res.TrustedNetworkContracts {
			if tn.NetworkName == name {
				return &tn, resp, nil
			}
		}

		if len(res.TrustedNetworkContracts) < pageSize {
			break
		}
		page++
	}

	return nil, nil, fmt.Errorf("trusted network with name '%s' not found", name)
}

func GetTrustedNetworkByID(ctx context.Context, service *zscaler.Service, id string) (*TrustedNetwork, *http.Response, error) {
	pageSize := 1000
	page := 1

	for {
		res, resp, err := GetMultipleTrustedNetworks(ctx, service, "", "", &page, &pageSize)
		if err != nil {
			return nil, resp, err
		}

		for _, tn := range res.TrustedNetworkContracts {
			if tn.ID == id {
				return &tn, resp, nil
			}
		}

		// If we got less than the page size, we’re done
		if len(res.TrustedNetworkContracts) < pageSize {
			break
		}
		page++
	}

	return nil, nil, fmt.Errorf("trusted network with ID %s not found", id)
}

func CreateTrustedNetwork(ctx context.Context, service *zscaler.Service, network *TrustedNetwork) (*TrustedNetwork, *http.Response, error) {
	if network == nil {
		return nil, nil, errors.New("network is required")
	}

	url := fmt.Sprintf("%s/create", trustedNetworkEndpoint)

	// Send creation request
	resp, err := service.Client.NewZccRequestDo(ctx, "POST", url, nil, network, nil)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create trusted network: %w", err)
	}

	// Lookup the resource by networkName since the response doesn't include ID
	createdNetwork, _, err := GetTrustedNetworkByName(ctx, service, network.NetworkName)
	if err != nil {
		return nil, resp, fmt.Errorf("trusted network created, but failed to retrieve by name '%s': %w", network.NetworkName, err)
	}

	service.Client.GetLogger().Printf("[DEBUG] returning trusted network from create (via lookup): ID=%s, Name=%s", createdNetwork.ID, createdNetwork.NetworkName)
	return createdNetwork, resp, nil
}

func UpdateTrustedNetwork(ctx context.Context, service *zscaler.Service, network *TrustedNetwork) (*TrustedNetwork, *http.Response, error) {
	if network == nil {
		return nil, nil, errors.New("network is required")
	}

	url := fmt.Sprintf("%s/edit", trustedNetworkEndpoint)

	var updatedNetwork TrustedNetwork

	resp, err := service.Client.NewZccRequestDo(ctx, "PUT", url, nil, network, &updatedNetwork)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to update trusted network: %w", err)
	}

	service.Client.GetLogger().Printf("[DEBUG] returning trusted network from update: %s", updatedNetwork.ID)
	return &updatedNetwork, resp, nil
}

func DeleteTrustedNetwork(ctx context.Context, service *zscaler.Service, networkID string) (*http.Response, error) {
	if networkID == "" {
		return nil, fmt.Errorf("network ID is required for deletion")
	}

	// Construct the delete endpoint URL
	endpoint := fmt.Sprintf("%s/%s/delete", trustedNetworkEndpoint, networkID)

	// Make the DELETE request using NewZccRequestDo
	resp, err := service.Client.NewZccRequestDo(ctx, "DELETE", endpoint, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to delete trusted network: %w", err)
	}

	service.Client.GetLogger().Printf("[DEBUG] Deleted trusted network ID: %s", networkID)
	return resp, nil
}
