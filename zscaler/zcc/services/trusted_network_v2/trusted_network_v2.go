package trusted_network_v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	trustedNetworkEndpointV2 = "/zcc/papi/public/v2/trusted-networks"
)

// Predefined values for GetAllFilterOptions.Type. The trusted-networks v2
// list endpoint scopes the keyword search to the named field.
const (
	FilterTypeName               = "NAME"
	FilterTypeDNSServers         = "DNS_SERVERS"
	FilterTypeDNSSearchDomains   = "DNS_SEARCH_DOMAINS"
	FilterTypeHostnameIP         = "HOST_NAME_IP"
	FilterTypeTrustedSubnets     = "TRUSTED_SUBNETS"
	FilterTypeTrustedGateways    = "TRUSTED_GATEWAYS"
	FilterTypeTrustedDHCPServers = "TRUSTED_DHCP_SERVERS"
	FilterTypeTrustedEgressIPs   = "TRUSTED_EGRESS_IPS"
	FilterTypeSSID               = "SSID"
)

// GetAllFilterOptions models the documented optional query parameters for
// GET /zcc/papi/public/v2/trusted-networks. Pagination (skip/perPage) is
// handled by the pagination helper; callers only supply filters here.
type GetAllFilterOptions struct {
	// Keyword filters records by name (substring match). When combined
	// with Type, the keyword is matched against that specific field.
	Keyword string
	// Type narrows the keyword search to a specific field. Use one of the
	// FilterType* constants; empty leaves the API default (all fields).
	Type string
}

type TrustedNetworkV2 struct {
	ID                     int      `json:"id,omitempty"`
	CompanyID              int      `json:"companyId,omitempty"`
	ZPAID                  string   `json:"zpaId,omitempty"`
	Active                 bool     `json:"active,omitempty"`
	ConditionType          string   `json:"conditionType,omitempty"`
	Name                   string   `json:"name,omitempty"`
	CreatedBy              string   `json:"createdBy,omitempty"`
	DNSSearchDomains       []string `json:"dnsSearchDomains,omitempty"`
	DNSServerIPs           []string `json:"dnsServerIps,omitempty"`
	EditedBy               string   `json:"editedBy,omitempty"`
	Guid                   string   `json:"guid,omitempty"`
	Hostname               string   `json:"hostname,omitempty"`
	NetworkName            string   `json:"networkName,omitempty"`
	ResolvedIPsForHostname []string `json:"resolvedIpsForHostname,omitempty"`
	SSID                   string   `json:"ssid,omitempty"`
	TrustedDhcpServersIPs  []string `json:"trustedDhcpServersIps,omitempty"`
	TrustedEgressIPs       []string `json:"trustedEgressIps,omitempty"`
	TrustedGatewayIPs      []string `json:"trustedGatewayIps,omitempty"`
	TrustedSubnetIPs       []string `json:"trustedSubnetIps,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, networkID int) (*TrustedNetworkV2, error) {
	var trustedNetwork TrustedNetworkV2
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", trustedNetworkEndpointV2, networkID), &trustedNetwork)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning trusted network from Get: %d", trustedNetwork.ID)
	return &trustedNetwork, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, networkName string) (*TrustedNetworkV2, error) {
	// Narrow the server-side search to the NAME field. Final exact
	// match is still done client-side because keyword is a substring
	// filter, not an equality match.
	networks, err := GetAll(ctx, service, &GetAllFilterOptions{
		Keyword: networkName,
		Type:    FilterTypeName,
	})
	if err != nil {
		return nil, err
	}
	for _, network := range networks {
		if strings.EqualFold(network.Name, networkName) {
			return &network, nil
		}
	}
	return nil, fmt.Errorf("no trusted network found with name: %s", networkName)
}

func Create(ctx context.Context, service *zscaler.Service, network *TrustedNetworkV2) (*TrustedNetworkV2, *http.Response, error) {
	if network == nil {
		return nil, nil, errors.New("trusted network is required")
	}

	var created TrustedNetworkV2
	resp, err := service.Client.NewZccRequestDo(ctx, "POST", trustedNetworkEndpointV2, nil, network, &created)
	if err != nil {
		return nil, resp, err
	}

	service.Client.GetLogger().Printf("[DEBUG] returning new trusted network from create: %d", created.ID)
	return &created, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, networkID int, network *TrustedNetworkV2) (*TrustedNetworkV2, *http.Response, error) {
	if network == nil {
		return nil, nil, errors.New("trusted network is required")
	}

	endpoint := fmt.Sprintf("%s/%d", trustedNetworkEndpointV2, networkID)
	var updated TrustedNetworkV2
	resp, err := service.Client.NewZccRequestDo(ctx, "PUT", endpoint, nil, network, &updated)
	if err != nil {
		return nil, resp, err
	}

	service.Client.GetLogger().Printf("[DEBUG] returning updated trusted network from update: %d", updated.ID)
	return &updated, resp, nil
}

func PartialUpdate(ctx context.Context, service *zscaler.Service, networkID int, network *TrustedNetworkV2) (*TrustedNetworkV2, *http.Response, error) {
	if network == nil {
		return nil, nil, errors.New("trusted network is required")
	}

	endpoint := fmt.Sprintf("%s/%d", trustedNetworkEndpointV2, networkID)
	var updated TrustedNetworkV2
	resp, err := service.Client.NewZccRequestDo(ctx, "PATCH", endpoint, nil, network, &updated)
	if err != nil {
		return nil, resp, err
	}

	service.Client.GetLogger().Printf("[DEBUG] returning trusted network from partial update: %d", updated.ID)
	return &updated, resp, nil
}

func Delete(ctx context.Context, service *zscaler.Service, networkID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", trustedNetworkEndpointV2, networkID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service, opts *GetAllFilterOptions) ([]TrustedNetworkV2, error) {
	params := common.QueryParamsV2{}
	if opts != nil {
		params.Keyword = opts.Keyword
		params.Type = opts.Type
	}
	return common.ReadAllPagesV2[TrustedNetworkV2](ctx, service.Client, trustedNetworkEndpointV2, params, common.DefaultPageSize)
}
