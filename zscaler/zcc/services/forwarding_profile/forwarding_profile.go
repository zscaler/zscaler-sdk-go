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
	ID                          IntOrString                  `json:"id,omitempty"`
	Active                      string                       `json:"active"`
	Name                        string                       `json:"name"`
	ConditionType               int                          `json:"conditionType"`
	DnsServers                  string                       `json:"dnsServers"`
	DnsSearchDomains            string                       `json:"dnsSearchDomains"`
	EnableLWFDriver             string                       `json:"enableLWFDriver"`
	Hostname                    string                       `json:"hostname"`
	ResolvedIpsForHostname      string                       `json:"resolvedIpsForHostname"`
	TrustedSubnets              string                       `json:"trustedSubnets"`
	TrustedGateways             string                       `json:"trustedGateways"`
	TrustedDhcpServers          string                       `json:"trustedDhcpServers"`
	TrustedEgressIps            string                       `json:"trustedEgressIps"`
	PredefinedTrustedNetworks   bool                         `json:"predefinedTrustedNetworks"`
	PredefinedTnAll             bool                         `json:"predefinedTnAll"`
	ForwardingProfileActions    []ForwardingProfileAction    `json:"forwardingProfileActions"`
	ForwardingProfileZpaActions []ForwardingProfileZpaAction `json:"forwardingProfileZpaActions"`
	EnableUnifiedTunnel         int                          `json:"enableUnifiedTunnel"`
	UnifiedTunnel               []UnifiedTunnel              `json:"unifiedTunnel"`
	EnableAllDefaultAdaptersTN  int                          `json:"enableAllDefaultAdaptersTN"`
	EnableSplitVpnTN            int                          `json:"enableSplitVpnTN"`
	EvaluateTrustedNetwork      int                          `json:"evaluateTrustedNetwork"`
	SkipTrustedCriteriaMatch    int                          `json:"skipTrustedCriteriaMatch"`
	TrustedNetworkIds           []int                        `json:"trustedNetworkIds"`
	TrustedNetworks             []string                     `json:"trustedNetworks"`
	TrustedNetworkIdsSelected   []int                        `json:"trustedNetworkIdsSelected"`
}

type ForwardingProfileAction struct {
	NetworkType                        int             `json:"networkType"`
	ActionType                         int             `json:"actionType"`
	SystemProxy                        int             `json:"systemProxy"`
	CustomPac                          string          `json:"customPac"`
	EnablePacketTunnel                 int             `json:"enablePacketTunnel"`
	SystemProxyData                    SystemProxyData `json:"systemProxyData"`
	PrimaryTransport                   int             `json:"primaryTransport"`
	DTLSTimeout                        int             `json:"DTLSTimeout"`
	UDPTimeout                         int             `json:"UDPTimeout"`
	TLSTimeout                         int             `json:"TLSTimeout"`
	MtuForZadapter                     IntOrString     `json:"mtuForZadapter"`
	BlockUnreachableDomainsTraffic     IntOrString     `json:"blockUnreachableDomainsTraffic"`
	AllowTLSFallback                   int             `json:"allowTLSFallback"`
	Tunnel2FallbackType                int             `json:"tunnel2FallbackType"`
	SendAllDNSToTrustedServer          int             `json:"sendAllDNSToTrustedServer"`
	DropIpv6Traffic                    IntOrString     `json:"dropIpv6Traffic"`
	RedirectWebTraffic                 IntOrString     `json:"redirectWebTraffic"`
	DropIpv6IncludeTrafficInT2         IntOrString     `json:"dropIpv6IncludeTrafficInT2"`
	UseTunnel2ForProxiedWebTraffic     int             `json:"useTunnel2ForProxiedWebTraffic"`
	UseTunnel2ForUnencryptedWebTraffic int             `json:"useTunnel2ForUnencryptedWebTraffic"`
	PathMtuDiscovery                   int             `json:"pathMtuDiscovery"`
	LatencyBasedZenEnablement          IntOrString     `json:"latencyBasedZenEnablement"`
	ZenProbeInterval                   int             `json:"zenProbeInterval"`
	ZenProbeSampleSize                 int             `json:"zenProbeSampleSize"`
	ZenThresholdLimit                  int             `json:"zenThresholdLimit"`
	DropIpv6TrafficInIpv6Network       IntOrString     `json:"dropIpv6TrafficInIpv6Network"`
	OptimiseForUnstableConnections     int             `json:"optimiseForUnstableConnections"`
	LatencyBasedServerEnablement       int             `json:"latencyBasedServerEnablement,omitempty"`
	LbsProbeInterval                   int             `json:"lbsProbeInterval,omitempty"`
	LbsProbeSampleSize                 int             `json:"lbsProbeSampleSize,omitempty"`
	LbsThresholdLimit                  int             `json:"lbsThresholdLimit,omitempty"`
	LatencyBasedServerMTEnablement     int             `json:"latencyBasedServerMTEnablement,omitempty"`
	IsSameAsOnTrustedNetwork           bool            `json:"isSameAsOnTrustedNetwork,omitempty"`
}

type SystemProxyData struct {
	ProxyAction             int    `json:"proxyAction"`
	EnableAutoDetect        int    `json:"enableAutoDetect"`
	EnablePAC               int    `json:"enablePAC"`
	PacURL                  string `json:"pacURL"`
	EnableProxyServer       int    `json:"enableProxyServer"`
	ProxyServerAddress      string `json:"proxyServerAddress"`
	ProxyServerPort         string `json:"proxyServerPort"`
	BypassProxyForPrivateIP int    `json:"bypassProxyForPrivateIP"`
	PerformGPUpdate         int    `json:"performGPUpdate"`
	PacDataPath             string `json:"pacDataPath"`
}

type ForwardingProfileZpaAction struct {
	NetworkType                    int         `json:"networkType"`
	ActionType                     int         `json:"actionType"`
	PrimaryTransport               int         `json:"primaryTransport"`
	DTLSTimeout                    int         `json:"DTLSTimeout"`
	TLSTimeout                     int         `json:"TLSTimeout"`
	MtuForZadapter                 int         `json:"mtuForZadapter"`
	SendTrustedNetworkResultToZpa  int         `json:"sendTrustedNetworkResultToZpa"`
	PartnerInfo                    PartnerInfo `json:"partnerInfo"`
	LatencyBasedServerEnablement   int         `json:"latencyBasedZpaServerEnablement"`
	LbsProbeInterval               int         `json:"lbsZpaProbeInterval"`
	LbsProbeSampleSize             int         `json:"lbsZpaProbeSampleSize"`
	LbsThresholdLimit              int         `json:"lbsZpaThresholdLimit"`
	LatencyBasedServerMTEnablement int         `json:"latencyBasedServerMTEnablement"`
	IsSameAsOnTrustedNetwork       bool        `json:"isSameAsOnTrustedNetwork"`
}

type PartnerInfo struct {
	PrimaryTransport int `json:"primaryTransport"`
	AllowTlsFallback int `json:"allowTlsFallback"`
	MtuForZadapter   int `json:"mtuForZadapter"`
}

type UnifiedTunnel struct {
	NetworkType                    int             `json:"networkType"`
	ActionTypeZIA                  int             `json:"actionTypeZIA"`
	ActionTypeZPA                  int             `json:"actionTypeZPA"`
	PrimaryTransport               int             `json:"primaryTransport"`
	DTLSTimeout                    int             `json:"DTLSTimeout"`
	TLSTimeout                     int             `json:"TLSTimeout"`
	MtuForZadapter                 int             `json:"mtuForZadapter"`
	AllowTLSFallback               int             `json:"allowTLSFallback"`
	PathMtuDiscovery               int             `json:"pathMtuDiscovery"`
	OptimiseForUnstableConnections int             `json:"optimiseForUnstableConnections"`
	Tunnel2FallbackType            int             `json:"tunnel2FallbackType"`
	RedirectWebTraffic             int             `json:"redirectWebTraffic"`
	DropIpv6Traffic                int             `json:"dropIpv6Traffic"`
	DropIpv6TrafficInIpv6Network   int             `json:"dropIpv6TrafficInIpv6Network"`
	BlockUnreachableDomainsTraffic int             `json:"blockUnreachableDomainsTraffic"`
	DropIpv6IncludeTrafficInT2     int             `json:"dropIpv6IncludeTrafficInT2"`
	SendAllDNSToTrustedServer      int             `json:"sendAllDNSToTrustedServer"`
	SystemProxyData                SystemProxyData `json:"systemProxyData"`
	SameAsOnTrusted                int             `json:"sameAsOnTrusted"`
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
	endpoint := fmt.Sprintf("%s/listByCompany", webForwardingProfileEndpoint)

	queryParams := struct {
		Search   string `url:"search,omitempty"`
		Page     int    `url:"page,omitempty"`
		PageSize int    `url:"pageSize,omitempty"`
	}{
		Search: search,
	}

	if page != nil {
		queryParams.Page = *page
	}
	if pageSize != nil {
		queryParams.PageSize = *pageSize
	}

	var profiles []ForwardingProfile
	_, err := service.Client.NewZccRequestDo(ctx, "GET", endpoint, queryParams, nil, &profiles)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forwarding profiles: %w", err)
	}

	return profiles, nil
}

func CreateForwardingProfile(ctx context.Context, service *zscaler.Service, request *ForwardingProfileRequest) (*CreateUpdateResponse, error) {
	if request == nil {
		return nil, errors.New("request is required")
	}

	url := fmt.Sprintf("%s/edit", webForwardingProfileEndpoint)

	var response CreateUpdateResponse
	_, err := service.Client.NewZccRequestDo(ctx, "POST", url, nil, request, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create/update forwarding profile: %w", err)
	}

	service.Client.GetLogger().Printf("[DEBUG] forwarding profile create/update response: %+v", response)
	return &response, nil
}

func DeleteForwardingProfile(ctx context.Context, service *zscaler.Service, profileID int) (*http.Response, error) {
	endpoint := fmt.Sprintf("%s/%d/delete", webForwardingProfileEndpoint, profileID)

	err := service.Client.Delete(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
