package forwarding_profile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	webForwardingProfileEndpoint = "/zcc/papi/public/v1/webForwardingProfile"
)

type ForwardingProfile struct {
	Active                      string                       `json:"active"`
	ConditionType               int                          `json:"conditionType"`
	DnsSearchDomains            string                       `json:"dnsSearchDomains"`
	DnsServers                  string                       `json:"dnsServers"`
	EnableLWFDriver             string                       `json:"enableLWFDriver"`
	EnableSplitVpnTN            int                          `json:"enableSplitVpnTN"`
	EvaluateTrustedNetwork      int                          `json:"evaluateTrustedNetwork"`
	ForwardingProfileActions    []ForwardingProfileAction    `json:"forwardingProfileActions"`
	ForwardingProfileZpaActions []ForwardingProfileZpaAction `json:"forwardingProfileZpaActions"`
	Hostname                    string                       `json:"hostname"`
	ID                          IntOrString                  `json:"id"`
	Name                        string                       `json:"name"`
	PredefinedTnAll             bool                         `json:"predefinedTnAll"`
	PredefinedTrustedNetworks   bool                         `json:"predefinedTrustedNetworks"`
	ResolvedIpsForHostname      string                       `json:"resolvedIpsForHostname"`
	SkipTrustedCriteriaMatch    int                          `json:"skipTrustedCriteriaMatch"`
	TrustedDhcpServers          string                       `json:"trustedDhcpServers"`
	TrustedEgressIps            string                       `json:"trustedEgressIps"`
	TrustedGateways             string                       `json:"trustedGateways"`
	TrustedNetworkIds           []int                        `json:"trustedNetworkIds"`
	TrustedNetworks             []string                     `json:"trustedNetworks"`
	TrustedSubnets              string                       `json:"trustedSubnets"`
}

type ForwardingProfileAction struct {
	DTLSTimeout                    int             `json:"DTLSTimeout"`
	TLSTimeout                     int             `json:"TLSTimeout"`
	UDPTimeout                     int             `json:"UDPTimeout"`
	ActionType                     int             `json:"actionType"`
	AllowTLSFallback               int             `json:"allowTLSFallback"`
	BlockUnreachableDomainsTraffic int             `json:"blockUnreachableDomainsTraffic"`
	CustomPac                      string          `json:"customPac"`
	DropIpv6IncludeTrafficInT2     int             `json:"dropIpv6IncludeTrafficInT2"`
	DropIpv6Traffic                int             `json:"dropIpv6Traffic"`
	DropIpv6TrafficInIpv6Network   int             `json:"dropIpv6TrafficInIpv6Network"`
	EnablePacketTunnel             int             `json:"enablePacketTunnel"`
	LatencyBasedZenEnablement      int             `json:"latencyBasedZenEnablement"`
	MtuForZadapter                 int             `json:"mtuForZadapter"`
	NetworkType                    int             `json:"networkType"`
	PathMtuDiscovery               int             `json:"pathMtuDiscovery"`
	PrimaryTransport               int             `json:"primaryTransport"`
	RedirectWebTraffic             int             `json:"redirectWebTraffic"`
	SystemProxy                    int             `json:"systemProxy"`
	SystemProxyData                SystemProxyData `json:"systemProxyData"`
	Tunnel2FallbackType            int             `json:"tunnel2FallbackType"`
	UseTunnel2ForProxiedWebTraffic int             `json:"useTunnel2ForProxiedWebTraffic"`
	ZenProbeInterval               int             `json:"zenProbeInterval"`
	ZenProbeSampleSize             int             `json:"zenProbeSampleSize"`
	ZenThresholdLimit              int             `json:"zenThresholdLimit"`
}

type SystemProxyData struct {
	BypassProxyForPrivateIP int    `json:"bypassProxyForPrivateIP"`
	EnableAutoDetect        int    `json:"enableAutoDetect"`
	EnablePAC               int    `json:"enablePAC"`
	EnableProxyServer       int    `json:"enableProxyServer"`
	PacDataPath             string `json:"pacDataPath"`
	PacURL                  string `json:"pacURL"`
	PerformGPUpdate         int    `json:"performGPUpdate"`
	ProxyAction             int    `json:"proxyAction"`
	ProxyServerAddress      string `json:"proxyServerAddress"`
	ProxyServerPort         string `json:"proxyServerPort"`
}

type ForwardingProfileZpaAction struct {
	DTLSTimeout                     int         `json:"DTLSTimeout"`
	TLSTimeout                      int         `json:"TLSTimeout"`
	ActionType                      int         `json:"actionType"`
	LatencyBasedServerMTEnablement  int         `json:"latencyBasedServerMTEnablement"`
	LatencyBasedZpaServerEnablement int         `json:"latencyBasedZpaServerEnablement"`
	LbsZpaProbeInterval             int         `json:"lbsZpaProbeInterval"`
	LbsZpaProbeSampleSize           int         `json:"lbsZpaProbeSampleSize"`
	LbsZpaThresholdLimit            int         `json:"lbsZpaThresholdLimit"`
	MtuForZadapter                  int         `json:"mtuForZadapter"`
	NetworkType                     int         `json:"networkType"`
	PartnerInfo                     PartnerInfo `json:"partnerInfo"`
	PrimaryTransport                int         `json:"primaryTransport"`
	SendTrustedNetworkResultToZpa   int         `json:"sendTrustedNetworkResultToZpa"`
}

type PartnerInfo struct {
	AllowTlsFallback int `json:"allowTlsFallback"`
	MtuForZadapter   int `json:"mtuForZadapter"`
	PrimaryTransport int `json:"primaryTransport"`
}

type IntOrString int

func (i *IntOrString) UnmarshalJSON(data []byte) error {
	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		*i = IntOrString(num)
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		parsedNum, err := strconv.Atoi(str)
		if err == nil {
			*i = IntOrString(parsedNum)
			return nil
		}
	}

	return fmt.Errorf("invalid value for IntOrString: %s", string(data))
}

func GetForwardingProfileByCompanyID(ctx context.Context, service *zscaler.Service, search string, page, pageSize *int) ([]ForwardingProfile, error) {
	// Construct the endpoint URL
	endpoint := fmt.Sprintf("%s/listByCompany", webForwardingProfileEndpoint)

	// Construct query parameters
	queryParams := struct {
		Search   string `url:"search,omitempty"`
		Page     int    `url:"page,omitempty"`
		PageSize int    `url:"pageSize,omitempty"`
	}{
		Search: search,
	}

	// Add optional pagination parameters if provided
	if page != nil {
		queryParams.Page = *page
	}
	if pageSize != nil {
		queryParams.PageSize = *pageSize
	}

	// Fetch the API response
	var profiles []ForwardingProfile
	_, err := service.Client.NewRequestDo(ctx, "GET", endpoint, queryParams, nil, &profiles)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forwarding profiles: %w", err)
	}

	return profiles, nil
}

func CreateForwardingProfile(ctx context.Context, service *zscaler.Service, profile *ForwardingProfile) (*ForwardingProfile, error) {
	if profile == nil {
		return nil, errors.New("profile is required")
	}

	// Construct the URL for the create endpoint
	url := fmt.Sprintf("%s/edit", webForwardingProfileEndpoint)

	// Initialize a variable to hold the response
	var createdProfile ForwardingProfile

	// Make the POST request to create the forwarding profile
	_, err := service.Client.NewRequestDo(ctx, "POST", url, nil, profile, &createdProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to create forwarding profile: %w", err)
	}

	service.Client.GetLogger().Printf("[DEBUG] returning forwarding profile from create: %+v", createdProfile)
	return &createdProfile, nil
}

func DeleteForwardingProfile(ctx context.Context, service *zscaler.Service, profileID int) (*http.Response, error) {
	// Correct the URL to include /delete
	endpoint := fmt.Sprintf("%s/%d/delete", webForwardingProfileEndpoint, profileID)

	// Make the DELETE request
	err := service.Client.Delete(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
