package forwarding_profile

// Request types for POST to /edit endpoint.
// The API uses different field names and types for POST vs GET responses.

type CreateUpdateResponse struct {
	Success string `json:"success"`
	ID      int    `json:"id"`
}

type ForwardingProfileRequest struct {
	ID                          string                              `json:"id"`
	Active                      int                                 `json:"active"`
	Name                        string                              `json:"name"`
	AddCondition                string                              `json:"addCondition"`
	ConditionType               int                                 `json:"conditionType"`
	DnsServers                  string                              `json:"dnsServers"`
	DnsSearchDomains            string                              `json:"dnsSearchDomains"`
	EnableLWFDriver             int                                 `json:"enableLWFDriver"`
	Hostname                    string                              `json:"hostname"`
	ResolvedIpsForHostname      string                              `json:"resolvedIpsForHostname"`
	TrustedSubnets              string                              `json:"trustedSubnets"`
	TrustedGateways             string                              `json:"trustedGateways"`
	TrustedDhcpServers          string                              `json:"trustedDhcpServers"`
	TrustedEgressIps            string                              `json:"trustedEgressIps"`
	EnableUnifiedTunnel         int                                 `json:"enableUnifiedTunnel"`
	EnableAllDefaultAdaptersTN  int                                 `json:"enableAllDefaultAdaptersTN"`
	EnableSplitVpnTN            int                                 `json:"enableSplitVpnTN"`
	EvaluateTrustedNetwork      int                                 `json:"evaluateTrustedNetwork"`
	SkipTrustedCriteriaMatch    int                                 `json:"skipTrustedCriteriaMatch"`
	PredefinedTrustedNetworks   bool                                `json:"predefinedTrustedNetworks"`
	PredefinedTnAll             bool                                `json:"predefinedTnAll"`
	TrustedNetworkIds           []int                               `json:"trustedNetworkIds"`
	ForwardingProfileActions    []ForwardingProfileActionRequest    `json:"forwardingProfileActions"`
	ForwardingProfileZpaActions []ForwardingProfileZpaActionRequest `json:"forwardingProfileZpaActions"`
	UnifiedTunnel               []UnifiedTunnelRequest              `json:"unifiedTunnel"`
}

type ForwardingProfileActionRequest struct {
	ActionType                         int                    `json:"actionType"`
	EnablePacketTunnel                 int                    `json:"enablePacketTunnel"`
	BlockUnreachableDomainsTraffic     string                 `json:"blockUnreachableDomainsTraffic"`
	DropIpv6Traffic                    int                    `json:"dropIpv6Traffic"`
	PrimaryTransport                   int                    `json:"primaryTransport"`
	UDPTimeout                         int                    `json:"UDPTimeout"`
	DTLSTimeout                        int                    `json:"DTLSTimeout"`
	TLSTimeout                         int                    `json:"TLSTimeout"`
	MtuForZadapter                     string                 `json:"mtuForZadapter"`
	AllowTLSFallback                   int                    `json:"allowTLSFallback"`
	PathMtuDiscovery                   int                    `json:"pathMtuDiscovery"`
	Tunnel2FallbackType                int                    `json:"tunnel2FallbackType"`
	UseTunnel2ForProxiedWebTraffic     int                    `json:"useTunnel2ForProxiedWebTraffic"`
	UseTunnel2ForUnencryptedWebTraffic int                    `json:"useTunnel2ForUnencryptedWebTraffic"`
	RedirectWebTraffic                 int                    `json:"redirectWebTraffic"`
	DropIpv6IncludeTrafficInT2         int                    `json:"dropIpv6IncludeTrafficInT2"`
	CustomPac                          string                 `json:"customPac"`
	SystemProxyData                    SystemProxyDataRequest `json:"systemProxyData"`
	LatencyBasedZenEnablement          int                    `json:"latencyBasedZenEnablement"`
	ZenProbeInterval                   int                    `json:"zenProbeInterval"`
	ZenProbeSampleSize                 int                    `json:"zenProbeSampleSize"`
	ZenThresholdLimit                  int                    `json:"zenThresholdLimit"`
	LatencyBasedServerEnablement       int                    `json:"latencyBasedServerEnablement"`
	LbsProbeInterval                   int                    `json:"lbsProbeInterval"`
	LbsProbeSampleSize                 int                    `json:"lbsProbeSampleSize"`
	LbsThresholdLimit                  int                    `json:"lbsThresholdLimit"`
	LatencyBasedServerMTEnablement     int                    `json:"latencyBasedServerMTEnablement"`
	NetworkType                        int                    `json:"networkType"`
	IsSameAsOnTrustedNetwork           *bool                  `json:"isSameAsOnTrustedNetwork,omitempty"`
}

type SystemProxyDataRequest struct {
	BypassProxyForPrivateIP int    `json:"bypassProxyForPrivateIP"`
	EnableAutoDetect        int    `json:"enableAutoDetect"`
	EnablePAC               int    `json:"enablePAC"`
	EnableProxyServer       int    `json:"enableProxyServer"`
	PacURL                  string `json:"pacURL"`
	PacDataPath             string `json:"pacDataPath"`
	PerformGPUpdate         int    `json:"performGPUpdate"`
	ProxyAction             int    `json:"proxyAction"`
	ProxyServerAddress      string `json:"proxyServerAddress"`
	ProxyServerPort         string `json:"proxyServerPort"`
}

type ForwardingProfileZpaActionRequest struct {
	ActionType                     int                `json:"actionType"`
	PrimaryTransport               int                `json:"primaryTransport"`
	DTLSTimeout                    int                `json:"DTLSTimeout"`
	TLSTimeout                     int                `json:"TLSTimeout"`
	MtuForZadapter                 string             `json:"mtuForZadapter"`
	PartnerInfo                    PartnerInfoRequest `json:"partnerInfo"`
	LatencyBasedServerEnablement   int                `json:"latencyBasedServerEnablement"`
	LbsProbeSampleSize             int                `json:"lbsProbeSampleSize"`
	LbsThresholdLimit              int                `json:"lbsThresholdLimit"`
	LbsProbeInterval               int                `json:"lbsProbeInterval"`
	LatencyBasedServerMTEnablement int                `json:"latencyBasedServerMTEnablement"`
	NetworkType                    int                `json:"networkType"`
	IsSameAsOnTrustedNetwork       *bool              `json:"isSameAsOnTrustedNetwork,omitempty"`
	SendTrustedNetworkResultToZpa  int                `json:"sendTrustedNetworkResultToZpa,omitempty"`
}

type PartnerInfoRequest struct {
	PrimaryTransport int `json:"primaryTransport"`
	AllowTlsFallback int `json:"allowTlsFallback"`
	MtuForZadapter   int `json:"mtuForZadapter"`
}

type UnifiedTunnelRequest struct {
	BlockUnreachableDomainsTraffic string                 `json:"blockUnreachableDomainsTraffic"`
	DropIpv6Traffic                int                    `json:"dropIpv6Traffic"`
	PrimaryTransport               int                    `json:"primaryTransport"`
	DTLSTimeout                    int                    `json:"DTLSTimeout"`
	TLSTimeout                     int                    `json:"TLSTimeout"`
	MtuForZadapter                 string                 `json:"mtuForZadapter"`
	AllowTLSFallback               int                    `json:"allowTLSFallback"`
	PathMtuDiscovery               int                    `json:"pathMtuDiscovery"`
	Tunnel2FallbackType            int                    `json:"tunnel2FallbackType"`
	RedirectWebTraffic             int                    `json:"redirectWebTraffic"`
	DropIpv6IncludeTrafficInT2     int                    `json:"dropIpv6IncludeTrafficInT2"`
	SystemProxyData                SystemProxyDataRequest `json:"systemProxyData"`
	NetworkType                    int                    `json:"networkType"`
	SameAsOnTrusted                int                    `json:"sameAsOnTrusted,omitempty"`
	ActionTypeZIA                  int                    `json:"actionTypeZIA"`
	ActionTypeZPA                  int                    `json:"actionTypeZPA"`
}
