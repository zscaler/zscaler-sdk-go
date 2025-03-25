package web_policy

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	baseWebPolicyEndpoint     = "/zcc/papi/public/v1/web/policy"
	baseWebAppServiceEndpoint = "/zcc/papi/public/v1/webAppService"
)

type WebPolicy struct {
	Active                    string           `json:"active"`
	AllowUnreachablePac       bool             `json:"allowUnreachablePac"`
	AndroidPolicy             AndroidPolicy    `json:"androidPolicy"`
	AppIdentityNames          []string         `json:"appIdentityNames"`
	AppServiceIds             []int            `json:"appServiceIds"`
	AppServiceNames           []string         `json:"appServiceNames"`
	BypassAppIds              []int            `json:"bypassAppIds"`
	BypassCustomAppIds        []int            `json:"bypassCustomAppIds"`
	Description               string           `json:"description"`
	DeviceGroupIds            []int            `json:"deviceGroupIds"`
	DeviceGroupNames          []string         `json:"deviceGroupNames"`
	DeviceType                string           `json:"device_type"`
	DisasterRecovery          DisasterRecovery `json:"disasterRecovery"`
	EnableDeviceGroups        string           `json:"enableDeviceGroups"`
	ForwardingProfileId       int              `json:"forwardingProfileId"`
	GroupAll                  string           `json:"groupAll"`
	GroupIds                  []int            `json:"groupIds"`
	GroupNames                []string         `json:"groupNames"`
	HighlightActiveControl    string           `json:"highlightActiveControl"`
	ID                        string           `json:"id"`
	IosPolicy                 IosPolicy        `json:"iosPolicy"`
	LinuxPolicy               LinuxPolicy      `json:"linuxPolicy"`
	LogFileSize               string           `json:"logFileSize"`
	LogLevel                  string           `json:"logLevel"`
	LogMode                   string           `json:"logMode"`
	MacPolicy                 MacPolicy        `json:"macPolicy"`
	Name                      string           `json:"name"`
	PacURL                    string           `json:"pac_url"`
	PolicyExtension           PolicyExtension  `json:"policyExtension"`
	ReactivateWebSecurityMins string           `json:"reactivateWebSecurityMinutes"`
	ReauthPeriod              string           `json:"reauth_period"`
	RuleOrder                 string           `json:"ruleOrder"`
	SendDisableServiceReason  string           `json:"sendDisableServiceReason"`
	TunnelZappTraffic         string           `json:"tunnelZappTraffic"`
	UserIds                   []int            `json:"userIds"`
	UserNames                 []string         `json:"userNames"`
	WindowsPolicy             WindowsPolicy    `json:"windowsPolicy"`
	ZiaPostureConfigId        int              `json:"ziaPostureConfigId"`
}

type AndroidPolicy struct {
	AllowedApps   string `json:"allowed_apps"`
	BillingDay    string `json:"billing_day"`
	BypassApps    string `json:"bypass_android_apps"`
	BypassMmsApps string `json:"bypass_mms_apps"`
	CustomText    string `json:"custom_text"`
	DisablePass   string `json:"disable_password"`
	EnableVerbose string `json:"enableVerboseLog"`
	Enforced      string `json:"enforced"`
	InstallSsl    string `json:"install_ssl_certs"`
	Limit         string `json:"limit"`
	LogoutPass    string `json:"logout_password"`
	QuotaRoaming  string `json:"quota_in_roaming"`
	UninstallPass string `json:"uninstall_password"`
	WifiSsid      string `json:"wifi_ssid"`
}

type DisasterRecovery struct {
	AllowZiaTest        bool   `json:"allowZiaTest"`
	AllowZpaTest        bool   `json:"allowZpaTest"`
	EnableZiaDR         bool   `json:"enableZiaDR"`
	EnableZpaDR         bool   `json:"enableZpaDR"`
	PolicyId            string `json:"policyId"`
	UseZiaGlobalDb      bool   `json:"useZiaGlobalDb"`
	ZiaDRRecoveryMethod int    `json:"ziaDRRecoveryMethod"`
	ZiaDomainName       string `json:"ziaDomainName"`
	ZiaGlobalDbURL      string `json:"ziaGlobalDbUrl"`
	ZiaGlobalDbURLV2    string `json:"ziaGlobalDbUrlv2"`
	ZiaPacURL           string `json:"ziaPacUrl"`
	ZiaSecretKeyData    string `json:"ziaSecretKeyData"`
	ZiaSecretKeyName    string `json:"ziaSecretKeyName"`
	ZpaDomainName       string `json:"zpaDomainName"`
	ZpaSecretKeyData    string `json:"zpaSecretKeyData"`
	ZpaSecretKeyName    string `json:"zpaSecretKeyName"`
}

type IosPolicy struct {
	DisablePassword        string `json:"disable_password"`
	Ipv6Mode               string `json:"ipv6Mode"`
	LogoutPassword         string `json:"logout_password"`
	Passcode               string `json:"passcode"`
	ShowVPNTunNotification string `json:"showVPNTunNotification"`
	UninstallPassword      string `json:"uninstall_password"`
}

type LinuxPolicy struct {
	DisablePassword string `json:"disable_password"`
	InstallSslCerts string `json:"install_ssl_certs"`
	LogoutPassword  string `json:"logout_password"`
	UninstallPass   string `json:"uninstall_password"`
}

type MacPolicy struct {
	AddIfscopeRoute                      string `json:"addIfscopeRoute"`
	CacheSystemProxy                     string `json:"cacheSystemProxy"`
	ClearArpCache                        string `json:"clearArpCache"`
	DisablePassword                      string `json:"disable_password"`
	DnsPriorityOrdering                  string `json:"dnsPriorityOrdering"`
	DnsPriorityOrderingForTrustedDnsCrit string `json:"dnsPriorityOrderingForTrustedDnsCriteria"`
	EnableAppBasedBypass                 string `json:"enableApplicationBasedBypass"`
	EnableZscalerFirewall                string `json:"enableZscalerFirewall"`
	InstallSslCerts                      string `json:"install_ssl_certs"`
	LogoutPassword                       string `json:"logout_password"`
	PersistentZscalerFirewall            string `json:"persistentZscalerFirewall"`
	UninstallPassword                    string `json:"uninstall_password"`
}

type PolicyExtension struct {
	AdvanceZpaReauth                      bool                        `json:"advanceZpaReauth"`
	AdvanceZpaReauthTime                  int                         `json:"advanceZpaReauthTime"`
	CustomDNS                             string                      `json:"customDNS"`
	DdilConfig                            string                      `json:"ddilConfig"`
	DeleteDHCPOption121Routes             string                      `json:"deleteDHCPOption121Routes"`
	DisableDNSRouteExclusion              string                      `json:"disableDNSRouteExclusion"`
	DropQuicTraffic                       string                      `json:"dropQuicTraffic"`
	EnableAntiTampering                   string                      `json:"enableAntiTampering"`
	EnableSetProxyOnVPNAdapters           string                      `json:"enableSetProxyOnVPNAdapters"`
	EnableZCCRevert                       string                      `json:"enableZCCRevert"`
	EnableZdpService                      string                      `json:"enableZdpService"`
	EnforceSplitDNS                       string                      `json:"enforceSplitDNS"`
	ExitPassword                          string                      `json:"exitPassword"`
	FallbackToGatewayDomain               string                      `json:"fallbackToGatewayDomain"`
	FollowGlobalForPartnerLogin           string                      `json:"followGlobalForPartnerLogin"`
	FollowRoutingTable                    string                      `json:"followRoutingTable"`
	GenerateCliPasswordContract           GenerateCliPasswordContract `json:"generateCliPasswordContract"`
	InterceptZIATrafficAllAdapters        string                      `json:"interceptZIATrafficAllAdapters"`
	MachineIdpAuth                        bool                        `json:"machineIdpAuth"`
	Nonce                                 string                      `json:"nonce"`
	OverrideATCmdByPolicy                 string                      `json:"overrideATCmdByPolicy"`
	PacketTunnelDnsExcludeList            string                      `json:"packetTunnelDnsExcludeList"`
	PacketTunnelDnsIncludeList            string                      `json:"packetTunnelDnsIncludeList"`
	PacketTunnelExcludeList               string                      `json:"packetTunnelExcludeList"`
	PacketTunnelExcludeListForIPv6        string                      `json:"packetTunnelExcludeListForIPv6"`
	PacketTunnelIncludeList               string                      `json:"packetTunnelIncludeList"`
	PacketTunnelIncludeListForIPv6        string                      `json:"packetTunnelIncludeListForIPv6"`
	PartnerDomains                        string                      `json:"partnerDomains"`
	PrioritizeDnsExclusions               string                      `json:"prioritizeDnsExclusions"`
	PurgeKerberosPreferredDCCache         string                      `json:"purgeKerberosPreferredDCCache"`
	ReactivateAntiTamperingTime           int                         `json:"reactivateAntiTamperingTime"`
	SourcePortBasedBypasses               string                      `json:"sourcePortBasedBypasses"`
	TruncateLargeUDPDNSResponse           string                      `json:"truncateLargeUDPDNSResponse"`
	UpdateDnsSearchOrder                  string                      `json:"updateDnsSearchOrder"`
	UseDefaultAdapterForDNS               string                      `json:"useDefaultAdapterForDNS"`
	UseProxyPortForT1                     string                      `json:"useProxyPortForT1"`
	UseProxyPortForT2                     string                      `json:"useProxyPortForT2"`
	UseV8JsEngine                         string                      `json:"useV8JsEngine"`
	UseWsaPollForZpa                      string                      `json:"useWsaPollForZpa"`
	UseZscalerNotificationFramework       string                      `json:"useZscalerNotificationFramework"`
	UserAllowedToAddPartner               string                      `json:"userAllowedToAddPartner"`
	VpnGateways                           string                      `json:"vpnGateways"`
	ZccAppFailOpenPolicy                  string                      `json:"zccAppFailOpenPolicy"`
	ZccFailCloseSettingsAppByPassIds      []int                       `json:"zccFailCloseSettingsAppByPassIds"`
	ZccFailCloseSettingsAppByPassNames    []string                    `json:"zccFailCloseSettingsAppByPassNames"`
	ZccFailCloseSettingsExitUninstallPass string                      `json:"zccFailCloseSettingsExitUninstallPassword"`
	ZccFailCloseSettingsIpBypasses        string                      `json:"zccFailCloseSettingsIpBypasses"`
	ZccFailCloseSettingsLockdownOnTunnel  string                      `json:"zccFailCloseSettingsLockdownOnTunnelProcessExit"`
	ZccFailCloseSettingsThumbPrint        string                      `json:"zccFailCloseSettingsThumbPrint"`
	ZccRevertPassword                     string                      `json:"zccRevertPassword"`
	ZccTunnelFailPolicy                   string                      `json:"zccTunnelFailPolicy"`
	ZdDisablePassword                     string                      `json:"zdDisablePassword"`
	ZdpDisablePassword                    string                      `json:"zdpDisablePassword"`
	ZdxDisablePassword                    string                      `json:"zdxDisablePassword"`
	ZdxLiteConfigObj                      string                      `json:"zdxLiteConfigObj"`
	ZpaAuthExpOnNetIpChange               string                      `json:"zpaAuthExpOnNetIpChange"`
	ZpaAuthExpOnSleep                     string                      `json:"zpaAuthExpOnSleep"`
	ZpaAuthExpOnSysRestart                string                      `json:"zpaAuthExpOnSysRestart"`
	ZpaAuthExpOnWinLogonSession           string                      `json:"zpaAuthExpOnWinLogonSession"`
	ZpaAuthExpOnWinSessionLock            string                      `json:"zpaAuthExpOnWinSessionLock"`
	ZpaAuthExpSessionLockStateMinTime     int                         `json:"zpaAuthExpSessionLockStateMinTimeInSecond"`
	ZpaDisablePassword                    string                      `json:"zpaDisablePassword"`
}

type GenerateCliPasswordContract struct {
	AllowZpaDisableWithoutPassword bool `json:"allowZpaDisableWithoutPassword"`
	EnableCli                      bool `json:"enableCli"`
	PolicyId                       int  `json:"policyId"`
}

type WindowsPolicy struct {
	CacheSystemProxy              int    `json:"cacheSystemProxy"`
	CaptivePortalConfig           string `json:"captivePortalConfig"`
	DisableLoopBackRestriction    int    `json:"disableLoopBackRestriction"`
	DisableParallelIpv4andIpv6    string `json:"disableParallelIpv4andIpv6"`
	DisablePassword               string `json:"disable_password"`
	FlowLoggerConfig              string `json:"flowLoggerConfig"`
	ForceLocationRefreshSccm      int    `json:"forceLocationRefreshSccm"`
	InstallWindowsFirewallInbound int    `json:"installWindowsFirewallInboundRule"`
	InstallSslCerts               string `json:"install_ssl_certs"`
	LogoutPassword                string `json:"logout_password"`
	OverrideWPAD                  int    `json:"overrideWPAD"`
	PacDataPath                   string `json:"pacDataPath"`
	PacType                       int    `json:"pacType"`
	PrioritizeIPv4                int    `json:"prioritizeIPv4"`
	RemoveExemptedContainers      int    `json:"removeExemptedContainers"`
	RestartWinHttpSvc             int    `json:"restartWinHttpSvc"`
	TriggerDomainProfleDetection  int    `json:"triggerDomainProfleDetection"`
	UninstallPassword             string `json:"uninstall_password"`
	WfpDriver                     int    `json:"wfpDriver"`
}

type WebPolicyActivation struct {
	DeviceType int `json:"deviceType"`
	PolicyId   int `json:"policyId"`
}

func GetPolicyListByCompanyID(ctx context.Context, service *zscaler.Service, page, pageSize *int, search, searchType *string, deviceType *int) ([]map[string]interface{}, error) {
	// Construct the URL for the listByCompany endpoint
	url := fmt.Sprintf("%s/listByCompany", baseWebPolicyEndpoint)

	// Construct query parameters dynamically
	queryParams := common.QueryParams{}
	if page != nil {
		queryParams.Page = *page
	}
	if pageSize != nil {
		queryParams.PageSize = *pageSize
	}
	if search != nil && *search != "" {
		queryParams.Search = *search
	}
	if searchType != nil && *searchType != "" {
		queryParams.SearchType = *searchType
	}
	if deviceType != nil {
		queryParams.DeviceType = *deviceType
	}

	// Fetch the API response
	var policies []map[string]interface{}
	_, err := service.Client.NewZccRequestDo(ctx, "GET", url, queryParams, nil, &policies)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve policy list: %w", err)
	}

	return policies, nil
}

func ActivateWebPolicy(ctx context.Context, service *zscaler.Service, activation *WebPolicyActivation) (*WebPolicyActivation, error) {
	if activation == nil {
		return nil, errors.New("activation is required")
	}

	// Construct the URL for the activate endpoint
	url := fmt.Sprintf("%s/activate", baseWebPolicyEndpoint)

	// Initialize a variable to hold the response
	var updatedActivation WebPolicyActivation

	// Make the PUT request to activate the web policy
	_, err := service.Client.NewZccRequestDo(ctx, "PUT", url, nil, activation, &updatedActivation)
	if err != nil {
		return nil, fmt.Errorf("failed to activate web policy: %w", err)
	}

	service.Client.GetLogger().Printf("[DEBUG] returning activation from activate: %+v", updatedActivation)
	return &updatedActivation, nil
}

func UpdateWebPolicy(ctx context.Context, service *zscaler.Service, webPolicy *WebPolicy) (*WebPolicy, error) {
	if webPolicy == nil {
		return nil, errors.New("web policy is required")
	}

	// Construct the URL for the update endpoint
	url := fmt.Sprintf("%s/edit", baseWebPolicyEndpoint)

	// Initialize a variable to hold the response
	var updatedWebPolicy WebPolicy

	// Make the PUT request to update the web policy
	_, err := service.Client.NewZccRequestDo(ctx, "PUT", url, nil, webPolicy, &updatedWebPolicy)
	if err != nil {
		return nil, fmt.Errorf("failed to update web policy: %w", err)
	}

	service.Client.GetLogger().Printf("[DEBUG] returning web policy from update: %s", updatedWebPolicy)
	return &updatedWebPolicy, nil
}

func DeleteWebPolicy(ctx context.Context, service *zscaler.Service, policyID int) (*http.Response, error) {
	// Construct the complete endpoint with /delete
	endpoint := fmt.Sprintf("%s/%d/delete", baseWebPolicyEndpoint, policyID)

	// Make the DELETE request
	err := service.Client.Delete(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetWebAppServiceInfoByCompanyID(ctx context.Context, service *zscaler.Service, page, pageSize *int, search, searchType *string, deviceType *int) ([]map[string]interface{}, error) {
	// Construct the URL for the listByCompany endpoint
	url := fmt.Sprintf("%s/listByCompany", baseWebAppServiceEndpoint)

	// Construct query parameters dynamically
	queryParams := common.QueryParams{}
	if page != nil {
		queryParams.Page = *page
	}
	if pageSize != nil {
		queryParams.PageSize = *pageSize
	}
	if search != nil && *search != "" {
		queryParams.Search = *search
	}

	// Fetch the API response
	var policies []map[string]interface{}
	_, err := service.Client.NewZccRequestDo(ctx, "GET", url, queryParams, nil, &policies)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve policy list: %w", err)
	}

	return policies, nil
}
