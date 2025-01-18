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
	// Custom URLs that are exempted from cookie authentication for users
	AuthBypassUrls []string `json:"authBypassUrls,omitempty"`

	// Custom URLs that are exempted from Kerberos authentication
	KerberosBypassUrls []string `json:"kerberosBypassUrls,omitempty"`

	// Custom URLs that are exempted from Digest authentication
	DigestAuthBypassUrls []string `json:"digestAuthBypassUrls,omitempty"`

	// URLs that are excluded from DNS optimization on transparent proxy mode
	DnsResolutionOnTransparentProxyExemptUrls []string `json:"dnsResolutionOnTransparentProxyExemptUrls,omitempty"`

	// URLs to which DNS optimization on transparent proxy mode applies
	DnsResolutionOnTransparentProxyUrls []string `json:"dnsResolutionOnTransparentProxyUrls,omitempty"`

	// A Boolean value indicating whether DNS optimization is enabled or disabled for Z-Tunnel 2.0 and transparent proxy mode traffic (e.g., traffic via GRE or IPSec tunnels without a PAC file).
	EnableDnsResolutionOnTransparentProxy bool `json:"enableDnsResolutionOnTransparentProxy,omitempty"`

	// A Boolean value indicating whether DNS optimization is enabled or disabled for IPv6 traffic sent via Z-Tunnel 2.0 and transparent proxy mode traffic (e.g., traffic via GRE or IPSec tunnels without a PAC file).
	EnableIPv6DnsResolutionOnTransparentProxy bool `json:"enableIPv6DnsResolutionOnTransparentProxy,omitempty"`

	// A Boolean value indicating whether DNS optimization is enabled or disabled for all IPv6 transparent proxy traffic
	EnableIPv6DnsOptimizationOnAllTransparentProxy bool `json:"enableIPv6DnsOptimizationOnAllTransparentProxy,omitempty"`

	// A Boolean value indicating whether policy evaluation for global SSL bypass traffic is enabled or not
	EnableEvaluatePolicyOnGlobalSSLBypass bool `json:"enableEvaluatePolicyOnGlobalSSLBypass,omitempty"`

	// A Boolean value indicating whether Microsoft Office 365 One Click Configuration is enabled or not
	EnableOffice365 bool `json:"enableOffice365,omitempty"`

	// A Boolean value indicating whether to log internal IP address present in X-Forwarded-For (XFF) proxy header or not
	LogInternalIp bool `json:"logInternalIp,omitempty"`

	// Enforce Surrogate IP authentication for Windows app traffic
	EnforceSurrogateIpForWindowsApp bool `json:"enforceSurrogateIpForWindowsApp,omitempty"`

	// A Boolean value indicating whether to apply configured policies on tunneled HTTP traffic sent via a CONNECT method request on port 80
	TrackHttpTunnelOnHttpPorts bool `json:"trackHttpTunnelOnHttpPorts,omitempty"`

	// A Boolean value indicating whether HTTP CONNECT method requests to non-standard ports are allowed or not (i.e., requests directed to ports other than the standard HTTP and HTTPS ports, 80 and 443)
	BlockHttpTunnelOnNonHttpPorts bool `json:"blockHttpTunnelOnNonHttpPorts,omitempty"`

	// A Boolean value indicating whether to block HTTP and HTTPS transactions that have an FQDN mismatch
	BlockDomainFrontingOnHostHeader bool `json:"blockDomainFrontingOnHostHeader,omitempty"`

	// A Boolean value indicating whether to apply the Firewall rules configured without a specified location criteria (or with the Road Warrior location) to remote user traffic forwarded via Z-Tunnel 1.0 or PAC files
	ZscalerClientConnector1AndPacRoadWarriorInFirewall bool `json:"zscalerClientConnector1AndPacRoadWarriorInFirewall,omitempty"`

	// A Boolean value indicating whether to apply the URL Filtering policy even when the Cloud App Control policy already allows a transaction explicitly
	CascadeUrlFiltering bool `json:"cascadeUrlFiltering,omitempty"`

	// A Boolean value indicating whether policies that include user and department criteria can be configured and applied for unauthenticated traffic
	EnablePolicyForUnauthenticatedTraffic bool `json:"enablePolicyForUnauthenticatedTraffic,omitempty"`

	// A Boolean value indicating whether to allow or block traffic that is not compliant with RFC HTTP protocol standards
	BlockNonCompliantHttpRequestOnHttpPorts bool `json:"blockNonCompliantHttpRequestOnHttpPorts,omitempty"`

	// A Boolean value indicating whether ranks are enabled for admins to allow admin ranks in policy configuration and management
	EnableAdminRankAccess bool `json:"enableAdminRankAccess,omitempty"`

	// A Boolean value indicating whether or not HTTP/2 should be the default web protocol for accessing various applications at your organizational level
	Http2NonbrowserTrafficEnabled bool `json:"http2NonbrowserTrafficEnabled,omitempty"`

	// A Boolean value indicating whether or not to include the ECS option in all DNS queries, originating from all locations and remote users.
	EcsForAllEnabled bool `json:"ecsForAllEnabled,omitempty"`

	// A Boolean value indicating whether to dynamically update user risk score by tracking risky user activities in real time
	DynamicUserRiskEnabled bool `json:"dynamicUserRiskEnabled,omitempty"`

	// A Boolean value indicating whether CONNECT host and SNI mismatch (i.e., CONNECT host doesn't match the SSL/TLS client hello SNI) is blocked or not
	BlockConnectHostSniMismatch bool `json:"blockConnectHostSniMismatch,omitempty"`

	// A Boolean value indicating whether or not to use the SSL/TLS client hello Server Name Indication (SNI) for DNS resolution instead of the CONNECT host for forward proxy connections
	PreferSniOverConnHost bool `json:"preferSniOverConnHost,omitempty"`

	// A Boolean value indicating whether or not to insert XFF header to all traffic forwarded from ZIA to ZPA, including source IP-anchored and ZIA-inspected ZPA application traffic.
	SipaXffHeaderEnabled bool `json:"sipaXffHeaderEnabled,omitempty"`

	// A Boolean value indicating whether non-HTTP Traffic on HTTP and HTTPS ports are allowed or blocked
	BlockNonHttpOnHttpPortEnabled bool `json:"blockNonHttpOnHttpPortEnabled,omitempty"`

	// Specifies the login session timeout for admins accessing the ZIA Admin Portal
	UISessionTimeout int `json:"uiSessionTimeout,omitempty"`

	// The ECS prefix that must be used in DNS queries when the ECS option is enabled
	EcsObject common.IDNameExternalID `json:"ecsObject,omitempty"`

	// Cloud applications that are exempted from cookie authenticatio
	AuthBypassApps []string `json:"authBypassApps,omitempty"`

	// Cloud applications that are exempted from Kerberos authentication
	KerberosBypassApps []string `json:"kerberosBypassApps,omitempty"`

	// Cloud applications that are exempted from Basic authentication
	BasicBypassApps []string `json:"basicBypassApps,omitempty"`

	// Cloud applications that are exempted from Digest authentication
	DigestAuthBypassApps []string `json:"digestAuthBypassApps,omitempty"`

	// Cloud applications that are excluded from DNS optimization on transparent proxy mode
	DnsResolutionOnTransparentProxyExemptApps []string `json:"dnsResolutionOnTransparentProxyExemptApps,omitempty"`

	// Cloud applications that are excluded from DNS optimization for IPv6 addresses on transparent proxy mode
	DnsResolutionOnTransparentProxyIPv6ExemptApps []string `json:"dnsResolutionOnTransparentProxyIPv6ExemptApps,omitempty"`

	// Cloud applications to which DNS optimization on transparent proxy mode applies
	DnsResolutionOnTransparentProxyApps []string `json:"dnsResolutionOnTransparentProxyApps,omitempty"`

	// Cloud applications to which DNS optimization for IPv6 addresses on transparent proxy mode applies
	DnsResolutionOnTransparentProxyIPv6Apps []string `json:"dnsResolutionOnTransparentProxyIPv6Apps,omitempty"`

	// Applications that are exempted from domain fronting
	BlockDomainFrontingApps []string `json:"blockDomainFrontingApps,omitempty"`

	// Applications that are exempted from the preferSniOverConnHost setting
	PreferSniOverConnHostApps []string `json:"preferSniOverConnHostApps,omitempty"`

	// URL categories that are excluded from DNS optimization on transparent proxy mode
	DnsResolutionOnTransparentProxyExemptUrlCategories []string `json:"dnsResolutionOnTransparentProxyExemptUrlCategories,omitempty"`

	// IPv6 URL categories that are excluded from DNS optimization on transparent proxy mode
	DnsResolutionOnTransparentProxyIPv6ExemptUrlCategories []string `json:"dnsResolutionOnTransparentProxyIPv6ExemptUrlCategories,omitempty"`

	// URL categories to which DNS optimization on transparent proxy mode applies
	DnsResolutionOnTransparentProxyUrlCategories []string `json:"dnsResolutionOnTransparentProxyUrlCategories,omitempty"`

	// IPv6 URL categories to which DNS optimization on transparent proxy mode applies
	DnsResolutionOnTransparentProxyIPv6UrlCategories []string `json:"dnsResolutionOnTransparentProxyIPv6UrlCategories,omitempty"`
	// URL categories that are exempted from cookie authentication
	AuthBypassUrlCategories []string `json:"authBypassUrlCategories,omitempty"`

	// URL categories that are exempted from domain fronting
	DomainFrontingBypassUrlCategories []string `json:"domainFrontingBypassUrlCategories,omitempty"`

	// URL categories that are exempted from Kerberos authentication
	KerberosBypassUrlCategories []string `json:"kerberosBypassUrlCategories,omitempty"`

	// URL categories that are exempted from Basic authentication
	BasicBypassUrlCategories []string `json:"basicBypassUrlCategories,omitempty"`

	// URL categories for which HTTP range headers must be removed
	HttpRangeHeaderRemoveUrlCategories []string `json:"httpRangeHeaderRemoveUrlCategories,omitempty"`

	// URL categories that are exempted from Digest authentication
	DigestAuthBypassUrlCategories []string `json:"digestAuthBypassUrlCategories,omitempty"`

	// URL categories that are excluded from the preferSniOverConnHost setting (i.e., prefer SSL/TLS client hello SNI for DNS resolution instead of the CONNECT host for forward proxy connections)
	SniDnsOptimizationBypassUrlCategories []string `json:"sniDnsOptimizationBypassUrlCategories,omitempty"`
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
