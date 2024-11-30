package advanced_settings

import (
	"context"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	advSettingsEndpoint = "/zia/api/v1/advancedSettings"
)

type AdvancedSettings struct {
	AuthBypassUrls                                         []string                `json:"authBypassUrls"`
	KerberosBypassUrls                                     []string                `json:"kerberosBypassUrls"`
	DigestAuthBypassUrls                                   []string                `json:"digestAuthBypassUrls"`
	DnsResolutionOnTransparentProxyExemptUrls              []string                `json:"dnsResolutionOnTransparentProxyExemptUrls"`
	DnsResolutionOnTransparentProxyUrls                    []string                `json:"dnsResolutionOnTransparentProxyUrls"`
	EnableDnsResolutionOnTransparentProxy                  bool                    `json:"enableDnsResolutionOnTransparentProxy"`
	EnableIPv6DnsResolutionOnTransparentProxy              bool                    `json:"enableIPv6DnsResolutionOnTransparentProxy"`
	EnableIPv6DnsOptimizationOnAllTransparentProxy         bool                    `json:"enableIPv6DnsOptimizationOnAllTransparentProxy"`
	EnableEvaluatePolicyOnGlobalSSLBypass                  bool                    `json:"enableEvaluatePolicyOnGlobalSSLBypass"`
	EnableOffice365                                        bool                    `json:"enableOffice365"`
	LogInternalIp                                          bool                    `json:"logInternalIp"`
	EnforceSurrogateIpForWindowsApp                        bool                    `json:"enforceSurrogateIpForWindowsApp"`
	TrackHttpTunnelOnHttpPorts                             bool                    `json:"trackHttpTunnelOnHttpPorts"`
	BlockHttpTunnelOnNonHttpPorts                          bool                    `json:"blockHttpTunnelOnNonHttpPorts"`
	BlockDomainFrontingOnHostHeader                        bool                    `json:"blockDomainFrontingOnHostHeader"`
	ZscalerClientConnector1AndPacRoadWarriorInFirewall     bool                    `json:"zscalerClientConnector1AndPacRoadWarriorInFirewall"`
	CascadeUrlFiltering                                    bool                    `json:"cascadeUrlFiltering"`
	EnablePolicyForUnauthenticatedTraffic                  bool                    `json:"enablePolicyForUnauthenticatedTraffic"`
	BlockNonCompliantHttpRequestOnHttpPorts                bool                    `json:"blockNonCompliantHttpRequestOnHttpPorts"`
	EnableAdminRankAccess                                  bool                    `json:"enableAdminRankAccess"`
	Http2NonbrowserTrafficEnabled                          bool                    `json:"http2NonbrowserTrafficEnabled"`
	EcsForAllEnabled                                       bool                    `json:"ecsForAllEnabled"`
	DynamicUserRiskEnabled                                 bool                    `json:"dynamicUserRiskEnabled"`
	BlockConnectHostSniMismatch                            bool                    `json:"blockConnectHostSniMismatch"`
	PreferSniOverConnHost                                  bool                    `json:"preferSniOverConnHost"`
	SipaXffHeaderEnabled                                   bool                    `json:"sipaXffHeaderEnabled"`
	BlockNonHttpOnHttpPortEnabled                          bool                    `json:"blockNonHttpOnHttpPortEnabled"`
	UISessionTimeout                                       int                     `json:"uiSessionTimeout"`
	EcsObject                                              common.IDNameExternalID `json:"ecsObject"`
	AuthBypassApps                                         []string                `json:"authBypassApps"`
	KerberosBypassApps                                     []string                `json:"kerberosBypassApps"`
	BasicBypassApps                                        []string                `json:"basicBypassApps"`
	DigestAuthBypassApps                                   []string                `json:"digestAuthBypassApps"`
	DnsResolutionOnTransparentProxyExemptApps              []string                `json:"dnsResolutionOnTransparentProxyExemptApps"`
	DnsResolutionOnTransparentProxyIPv6ExemptApps          []string                `json:"dnsResolutionOnTransparentProxyIPv6ExemptApps"`
	DnsResolutionOnTransparentProxyApps                    []string                `json:"dnsResolutionOnTransparentProxyApps"`
	DnsResolutionOnTransparentProxyIPv6Apps                []string                `json:"dnsResolutionOnTransparentProxyIPv6Apps"`
	BlockDomainFrontingApps                                []string                `json:"blockDomainFrontingApps"`
	PreferSniOverConnHostApps                              []string                `json:"preferSniOverConnHostApps"`
	DnsResolutionOnTransparentProxyExemptUrlCategories     []string                `json:"dnsResolutionOnTransparentProxyExemptUrlCategories"`
	DnsResolutionOnTransparentProxyIPv6ExemptUrlCategories []string                `json:"dnsResolutionOnTransparentProxyIPv6ExemptUrlCategories"`
	DnsResolutionOnTransparentProxyUrlCategories           []string                `json:"dnsResolutionOnTransparentProxyUrlCategories"`
	DnsResolutionOnTransparentProxyIPv6UrlCategories       []string                `json:"dnsResolutionOnTransparentProxyIPv6UrlCategories"`
	AuthBypassUrlCategories                                []string                `json:"authBypassUrlCategories"`
	DomainFrontingBypassUrlCategories                      []string                `json:"domainFrontingBypassUrlCategories"`
	KerberosBypassUrlCategories                            []string                `json:"kerberosBypassUrlCategories"`
	BasicBypassUrlCategories                               []string                `json:"basicBypassUrlCategories"`
	HttpRangeHeaderRemoveUrlCategories                     []string                `json:"httpRangeHeaderRemoveUrlCategories"`
	DigestAuthBypassUrlCategories                          []string                `json:"digestAuthBypassUrlCategories"`
	SniDnsOptimizationBypassUrlCategories                  []string                `json:"sniDnsOptimizationBypassUrlCategories"`
}

func GetAdvancedSettings(ctx context.Context, service *zscaler.Service) (*AdvancedSettings, error) {
	var advSettings AdvancedSettings
	err := service.Client.Read(ctx, advSettingsEndpoint, &advSettings)
	if err != nil {
		return nil, err
	}
	return &advSettings, nil
}
