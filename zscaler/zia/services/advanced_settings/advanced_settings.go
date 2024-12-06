package advanced_settings

import (
	"context"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	advSettingsEndpoint = "/zia/api/v1/advancedSettings"
)

type AdvancedSettings struct {
	AuthBypassUrls                                         []string                `json:"authBypassUrls,omitempty"`
	KerberosBypassUrls                                     []string                `json:"kerberosBypassUrls,omitempty"`
	DigestAuthBypassUrls                                   []string                `json:"digestAuthBypassUrls,omitempty"`
	DnsResolutionOnTransparentProxyExemptUrls              []string                `json:"dnsResolutionOnTransparentProxyExemptUrls,omitempty"`
	DnsResolutionOnTransparentProxyUrls                    []string                `json:"dnsResolutionOnTransparentProxyUrls,omitempty"`
	EnableDnsResolutionOnTransparentProxy                  bool                    `json:"enableDnsResolutionOnTransparentProxy,omitempty"`
	EnableIPv6DnsResolutionOnTransparentProxy              bool                    `json:"enableIPv6DnsResolutionOnTransparentProxy,omitempty"`
	EnableIPv6DnsOptimizationOnAllTransparentProxy         bool                    `json:"enableIPv6DnsOptimizationOnAllTransparentProxy,omitempty"`
	EnableEvaluatePolicyOnGlobalSSLBypass                  bool                    `json:"enableEvaluatePolicyOnGlobalSSLBypass,omitempty"`
	EnableOffice365                                        bool                    `json:"enableOffice365,omitempty"`
	LogInternalIp                                          bool                    `json:"logInternalIp,omitempty"`
	EnforceSurrogateIpForWindowsApp                        bool                    `json:"enforceSurrogateIpForWindowsApp,omitempty"`
	TrackHttpTunnelOnHttpPorts                             bool                    `json:"trackHttpTunnelOnHttpPorts,omitempty"`
	BlockHttpTunnelOnNonHttpPorts                          bool                    `json:"blockHttpTunnelOnNonHttpPorts,omitempty"`
	BlockDomainFrontingOnHostHeader                        bool                    `json:"blockDomainFrontingOnHostHeader,omitempty"`
	ZscalerClientConnector1AndPacRoadWarriorInFirewall     bool                    `json:"zscalerClientConnector1AndPacRoadWarriorInFirewall,omitempty"`
	CascadeUrlFiltering                                    bool                    `json:"cascadeUrlFiltering,omitempty"`
	EnablePolicyForUnauthenticatedTraffic                  bool                    `json:"enablePolicyForUnauthenticatedTraffic,omitempty"`
	BlockNonCompliantHttpRequestOnHttpPorts                bool                    `json:"blockNonCompliantHttpRequestOnHttpPorts,omitempty"`
	EnableAdminRankAccess                                  bool                    `json:"enableAdminRankAccess,omitempty"`
	Http2NonbrowserTrafficEnabled                          bool                    `json:"http2NonbrowserTrafficEnabled,omitempty"`
	EcsForAllEnabled                                       bool                    `json:"ecsForAllEnabled,omitempty"`
	DynamicUserRiskEnabled                                 bool                    `json:"dynamicUserRiskEnabled,omitempty"`
	BlockConnectHostSniMismatch                            bool                    `json:"blockConnectHostSniMismatch,omitempty"`
	PreferSniOverConnHost                                  bool                    `json:"preferSniOverConnHost,omitempty"`
	SipaXffHeaderEnabled                                   bool                    `json:"sipaXffHeaderEnabled,omitempty"`
	BlockNonHttpOnHttpPortEnabled                          bool                    `json:"blockNonHttpOnHttpPortEnabled,omitempty"`
	UISessionTimeout                                       int                     `json:"uiSessionTimeout,omitempty"`
	EcsObject                                              common.IDNameExternalID `json:"ecsObject,omitempty"`
	AuthBypassApps                                         []string                `json:"authBypassApps,omitempty"`
	KerberosBypassApps                                     []string                `json:"kerberosBypassApps,omitempty"`
	BasicBypassApps                                        []string                `json:"basicBypassApps,omitempty"`
	DigestAuthBypassApps                                   []string                `json:"digestAuthBypassApps,omitempty"`
	DnsResolutionOnTransparentProxyExemptApps              []string                `json:"dnsResolutionOnTransparentProxyExemptApps,omitempty"`
	DnsResolutionOnTransparentProxyIPv6ExemptApps          []string                `json:"dnsResolutionOnTransparentProxyIPv6ExemptApps,omitempty"`
	DnsResolutionOnTransparentProxyApps                    []string                `json:"dnsResolutionOnTransparentProxyApps,omitempty"`
	DnsResolutionOnTransparentProxyIPv6Apps                []string                `json:"dnsResolutionOnTransparentProxyIPv6Apps,omitempty"`
	BlockDomainFrontingApps                                []string                `json:"blockDomainFrontingApps,omitempty"`
	PreferSniOverConnHostApps                              []string                `json:"preferSniOverConnHostApps,omitempty"`
	DnsResolutionOnTransparentProxyExemptUrlCategories     []string                `json:"dnsResolutionOnTransparentProxyExemptUrlCategories,omitempty"`
	DnsResolutionOnTransparentProxyIPv6ExemptUrlCategories []string                `json:"dnsResolutionOnTransparentProxyIPv6ExemptUrlCategories,omitempty"`
	DnsResolutionOnTransparentProxyUrlCategories           []string                `json:"dnsResolutionOnTransparentProxyUrlCategories,omitempty"`
	DnsResolutionOnTransparentProxyIPv6UrlCategories       []string                `json:"dnsResolutionOnTransparentProxyIPv6UrlCategories,omitempty"`
	AuthBypassUrlCategories                                []string                `json:"authBypassUrlCategories,omitempty"`
	DomainFrontingBypassUrlCategories                      []string                `json:"domainFrontingBypassUrlCategories,omitempty"`
	KerberosBypassUrlCategories                            []string                `json:"kerberosBypassUrlCategories,omitempty"`
	BasicBypassUrlCategories                               []string                `json:"basicBypassUrlCategories,omitempty"`
	HttpRangeHeaderRemoveUrlCategories                     []string                `json:"httpRangeHeaderRemoveUrlCategories,omitempty"`
	DigestAuthBypassUrlCategories                          []string                `json:"digestAuthBypassUrlCategories,omitempty"`
	SniDnsOptimizationBypassUrlCategories                  []string                `json:"sniDnsOptimizationBypassUrlCategories,omitempty"`
}

func GetAdvancedSettings(ctx context.Context, service *zscaler.Service) (*AdvancedSettings, error) {
	var advSettings AdvancedSettings
	err := service.Client.Read(ctx, advSettingsEndpoint, &advSettings)
	if err != nil {
		return nil, err
	}
	return &advSettings, nil
}

func UpdateAdvancedSettings(ctx context.Context, service *zscaler.Service, advancedSettings *AdvancedSettings) (*AdvancedSettings, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, (advSettingsEndpoint), *advancedSettings)
	if err != nil {
		return nil, nil, err
	}
	updatedAdvancedSettings, _ := resp.(*AdvancedSettings)

	service.Client.GetLogger().Printf("[DEBUG]returning updates rule label from update: %d", updatedAdvancedSettings)
	return updatedAdvancedSettings, nil, nil
}
