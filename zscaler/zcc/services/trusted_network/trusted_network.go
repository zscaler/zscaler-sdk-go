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

func GetMultipleTrustedNetworks(ctx context.Context, service *zscaler.Service, search, searchType string, page, pageSize *int) (*TrustedNetworksResponse, error) {
	// Construct the endpoint URL
	endpoint := fmt.Sprintf("%s/listByCompany", trustedNetworkEndpoint)

	// Construct query parameters
	queryParams := common.QueryParams{
		Search:     search,
		SearchType: searchType,
	}
	if page != nil {
		queryParams.Page = *page
	}
	if pageSize != nil {
		queryParams.PageSize = *pageSize
	}

	// Fetch the API response
	var response TrustedNetworksResponse
	_, err := service.Client.NewRequestDo(ctx, "GET", endpoint, queryParams, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trusted networks: %w", err)
	}

	// Handle cases where totalCount > 0 but no results are returned
	if response.TotalCount > 0 && len(response.TrustedNetworkContracts) == 0 {
		return nil, fmt.Errorf("totalCount is %d, but no trusted networks are returned", response.TotalCount)
	}

	return &response, nil
}

func CreateTrustedNetwork(ctx context.Context, service *zscaler.Service, network *TrustedNetwork) (*TrustedNetwork, error) {
	if network == nil {
		return nil, errors.New("rule is required")
	}

	// Construct the URL for the create endpoint
	url := fmt.Sprintf("%s/create", trustedNetworkEndpoint)

	// Initialize a variable to hold the response
	var createdNetwork TrustedNetwork

	// Make the POST request to create the trusted network
	_, err := service.Client.NewRequestDo(ctx, "POST", url, nil, network, &createdNetwork)
	if err != nil {
		return nil, fmt.Errorf("failed to create trusted network: %w", err)
	}

	service.Client.GetLogger().Printf("[DEBUG] returning trusted network from create: %s", createdNetwork.ID)
	return &createdNetwork, nil
}

func UpdateTrustedNetwork(ctx context.Context, service *zscaler.Service, network *TrustedNetwork) (*TrustedNetwork, error) {
	if network == nil {
		return nil, errors.New("network is required")
	}

	// Construct the URL for the update endpoint
	url := fmt.Sprintf("%s/edit", trustedNetworkEndpoint)

	// Initialize a variable to hold the response
	var updatedNetwork TrustedNetwork

	// Make the PUT request to update the trusted network
	_, err := service.Client.NewRequestDo(ctx, "PUT", url, nil, network, &updatedNetwork)
	if err != nil {
		return nil, fmt.Errorf("failed to update trusted network: %w", err)
	}

	service.Client.GetLogger().Printf("[DEBUG] returning trusted network from update: %s", updatedNetwork)
	return &updatedNetwork, nil
}

func DeleteTrustedNetwork(ctx context.Context, service *zscaler.Service, networkID int) (*http.Response, error) {
	// Construct the complete endpoint with /delete
	endpoint := fmt.Sprintf("%s/%d/delete", trustedNetworkEndpoint, networkID)

	// Make the DELETE request
	err := service.Client.Delete(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
