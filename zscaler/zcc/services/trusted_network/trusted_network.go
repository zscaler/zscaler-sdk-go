package trusted_network

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	trustedNetworkEndpoint = "/zcc/papi/public/v1/webTrustedNetwork"
)

type TrustedNetwork struct {
	ID                     string `json:"id,omitempty"`
	Active                 bool   `json:"active"`
	CompanyID              string `json:"companyId,omitempty"`
	ConditionType          int    `json:"conditionType"`
	CreatedBy              string `json:"createdBy,omitempty"`
	DnsSearchDomains       string `json:"dnsSearchDomains"`
	DnsServers             string `json:"dnsServers"`
	EditedBy               string `json:"editedBy,omitempty"`
	Guid                   string `json:"guid,omitempty"`
	Hostnames              string `json:"hostnames"`
	NetworkName            string `json:"networkName"`
	ResolvedIpsForHostname string `json:"resolvedIpsForHostname"`
	Ssids                  string `json:"ssids,omitempty"`
	TrustedDhcpServers     string `json:"trustedDhcpServers"`
	TrustedEgressIps       string `json:"trustedEgressIps,omitempty"`
	TrustedGateways        string `json:"trustedGateways"`
	TrustedSubnets         string `json:"trustedSubnets"`
}

type TrustedNetworksResponse struct {
	TotalCount              int              `json:"totalCount"`
	TrustedNetworkContracts []TrustedNetwork `json:"trustedNetworkContracts"`
}

// TrustedNetworkMutationResponse is the JSON body returned by POST /webTrustedNetwork/create and
// PUT /webTrustedNetwork/edit on 200 OK (e.g. {"success":"true","errorCode":"0"}).
// The contract (id, networkName, etc.) is not returned; resolve the resource via GetTrustedNetworkByName
// after create or GetTrustedNetworkByID after update.
type TrustedNetworkMutationResponse struct {
	Success   string `json:"success"`
	ErrorCode string `json:"errorCode"`
}

func validateTrustedNetworkMutationResponse(r TrustedNetworkMutationResponse) error {
	if r.ErrorCode != "" && r.ErrorCode != "0" {
		return fmt.Errorf("trusted network mutation failed: errorCode=%q success=%q", r.ErrorCode, r.Success)
	}
	return nil
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

	var apiResp TrustedNetworkMutationResponse
	resp, err := service.Client.NewZccRequestDo(ctx, "POST", url, nil, network, &apiResp)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create trusted network: %w", err)
	}
	if err := validateTrustedNetworkMutationResponse(apiResp); err != nil {
		return nil, resp, err
	}

	// POST does not return the new id; resolve by networkName via listByCompany.
	// The API has eventual consistency — the newly created resource may not
	// appear in the list endpoint immediately. Retry a few times with a delay.
	var createdNetwork *TrustedNetwork
	var lookupErr error
	for attempt := 1; attempt <= 6; attempt++ {
		createdNetwork, _, lookupErr = GetTrustedNetworkByName(ctx, service, network.NetworkName)
		if lookupErr == nil && createdNetwork != nil {
			break
		}
		service.Client.GetLogger().Printf("[DEBUG] trusted network name lookup attempt %d failed, retrying in 2s: %v", attempt, lookupErr)
		time.Sleep(2 * time.Second)
	}
	if createdNetwork == nil {
		return nil, resp, fmt.Errorf("trusted network created, but failed to retrieve by name '%s': %w", network.NetworkName, lookupErr)
	}

	service.Client.GetLogger().Printf("[DEBUG] returning trusted network from create (via lookup): ID=%s, Name=%s", createdNetwork.ID, createdNetwork.NetworkName)
	return createdNetwork, resp, nil
}

func UpdateTrustedNetwork(ctx context.Context, service *zscaler.Service, network *TrustedNetwork) (*TrustedNetwork, *http.Response, error) {
	if network == nil {
		return nil, nil, errors.New("network is required")
	}
	if network.ID == "" {
		return nil, nil, errors.New("network ID is required for update")
	}

	url := fmt.Sprintf("%s/edit", trustedNetworkEndpoint)

	var apiResp TrustedNetworkMutationResponse
	resp, err := service.Client.NewZccRequestDo(ctx, "PUT", url, nil, network, &apiResp)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to update trusted network: %w", err)
	}
	if err := validateTrustedNetworkMutationResponse(apiResp); err != nil {
		return nil, resp, err
	}

	// PUT returns only success/errorCode, not the contract; re-fetch by id (id must be set on the request body).
	refreshed, _, refreshErr := GetTrustedNetworkByID(ctx, service, network.ID)
	if refreshErr != nil {
		return nil, resp, fmt.Errorf("trusted network updated, but failed to re-fetch by id %q: %w", network.ID, refreshErr)
	}

	service.Client.GetLogger().Printf("[DEBUG] returning trusted network from update (via refresh): ID=%s", refreshed.ID)
	return refreshed, resp, nil
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
