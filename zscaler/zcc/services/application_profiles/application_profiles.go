package application_profiles

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	appProfilesEndpoint = "/zcc/papi/public/v1/application-profiles"
)

// ApplicationProfilesResponse is the wrapper returned by the GET list endpoint.
type ApplicationProfilesResponse struct {
	TotalCount int                  `json:"totalCount"`
	Policies   []ApplicationProfile `json:"policies"`
}

type ApplicationProfile struct {
	DeviceType                   string           `json:"deviceType,omitempty"`
	ID                           int              `json:"id,omitempty"`
	Name                         string           `json:"name,omitempty"`
	Description                  string           `json:"description,omitempty"`
	PacURL                       string           `json:"pac_url,omitempty"`
	Active                       int              `json:"active"`
	RuleOrder                    int              `json:"ruleOrder,omitempty"`
	LogMode                      int              `json:"logMode,omitempty"`
	LogLevel                     int              `json:"logLevel,omitempty"`
	LogFileSize                  int              `json:"logFileSize,omitempty"`
	ReauthPeriod                 string           `json:"reauth_period,omitempty"`
	ReactivateWebSecurityMinutes string           `json:"reactivateWebSecurityMinutes,omitempty"`
	HighlightActiveControl       int              `json:"highlightActiveControl,omitempty"`
	SendDisableServiceReason     int              `json:"sendDisableServiceReason,omitempty"`
	RefreshKerberosToken         int              `json:"refreshKerberosToken,omitempty"`
	EnableDeviceGroups           int              `json:"enableDeviceGroups,omitempty"`
	Groups                       []int            `json:"groups,omitempty"`
	DeviceGroups                 []int            `json:"deviceGroups,omitempty"`
	OnNetPolicy                  interface{}      `json:"onNetPolicy,omitempty"`
	NotificationTemplateContract interface{}      `json:"notificationTemplateContract,omitempty"`
	NotificationTemplateId       int              `json:"notificationTemplateId,omitempty"`
	ForwardingProfileId          int              `json:"forwardingProfileId,omitempty"`
	ZiaPostureConfigId           int              `json:"ziaPostureConfigId,omitempty"`
	PolicyToken                  string           `json:"policyToken,omitempty"`
	TunnelZappTraffic            int              `json:"tunnelZappTraffic,omitempty"`
	GroupAll                     int              `json:"groupAll,omitempty"`
	Users                        []int            `json:"users,omitempty"`
	PolicyExtension              PolicyExtension  `json:"policyExtension,omitempty"`
	DisasterRecovery             DisasterRecovery `json:"disasterRecovery,omitempty"`
	ZiaPostureConfig             interface{}      `json:"ziaPostureConfig,omitempty"`
	GroupIds                     []int            `json:"groupIds,omitempty"`
	DeviceGroupIds               []int            `json:"deviceGroupIds,omitempty"`
	UserIds                      []int            `json:"userIds,omitempty"`
	BypassAppIds                 []int            `json:"bypassAppIds,omitempty"`
	AppServiceIds                []string         `json:"appServiceIds,omitempty"`
	BypassCustomAppIds           []int            `json:"bypassCustomAppIds,omitempty"`
	BypassApps                   interface{}      `json:"bypassApps,omitempty"`
	BypassCustomApps             interface{}      `json:"bypassCustomApps,omitempty"`
	AppServices                  []AppService     `json:"appServices,omitempty"`
	Passcode                     string           `json:"passcode,omitempty"`
	LogoutPassword               string           `json:"logout_password,omitempty"`
	DisablePassword              string           `json:"disable_password,omitempty"`
	UninstallPassword            string           `json:"uninstall_password,omitempty"`
	ShowVPNTunNotification       int              `json:"showVPNTunNotification,omitempty"`
	UseTunnelSDK4_3              int              `json:"useTunnelSDK4_3,omitempty"`
	Ipv6Mode                     int              `json:"ipv6Mode,omitempty"`
}

type PolicyExtension struct {
	SourcePortBasedBypasses                      string                      `json:"sourcePortBasedBypasses,omitempty"`
	PacketTunnelExcludeList                      string                      `json:"packetTunnelExcludeList,omitempty"`
	PacketTunnelIncludeList                      string                      `json:"packetTunnelIncludeList,omitempty"`
	CustomDNS                                    string                      `json:"customDNS,omitempty"`
	ExitPassword                                 string                      `json:"exitPassword,omitempty"`
	UseV8JsEngine                                string                      `json:"useV8JsEngine,omitempty"`
	ZdxDisablePassword                           string                      `json:"zdxDisablePassword,omitempty"`
	ZdDisablePassword                            string                      `json:"zdDisablePassword,omitempty"`
	ZpaDisablePassword                           string                      `json:"zpaDisablePassword,omitempty"`
	ZdpDisablePassword                           string                      `json:"zdpDisablePassword,omitempty"`
	FollowRoutingTable                           string                      `json:"followRoutingTable,omitempty"`
	UseWsaPollForZpa                             string                      `json:"useWsaPollForZpa,omitempty"`
	UseDefaultAdapterForDNS                      string                      `json:"useDefaultAdapterForDNS,omitempty"`
	UseZscalerNotificationFramework              string                      `json:"useZscalerNotificationFramework,omitempty"`
	SwitchFocusToNotification                    string                      `json:"switchFocusToNotification,omitempty"`
	FallbackToGatewayDomain                      string                      `json:"fallbackToGatewayDomain,omitempty"`
	EnableZCCRevert                              string                      `json:"enableZCCRevert,omitempty"`
	ZccRevertPassword                            string                      `json:"zccRevertPassword,omitempty"`
	ZpaAuthExpOnSleep                            int                         `json:"zpaAuthExpOnSleep,omitempty"`
	ZpaAuthExpOnSysRestart                       int                         `json:"zpaAuthExpOnSysRestart,omitempty"`
	ZpaAuthExpOnNetIpChange                      int                         `json:"zpaAuthExpOnNetIpChange,omitempty"`
	InstantForceZPAReauthStateUpdate             int                         `json:"instantForceZPAReauthStateUpdate,omitempty"`
	ZpaAuthExpOnWinLogonSession                  int                         `json:"zpaAuthExpOnWinLogonSession,omitempty"`
	ZpaAuthExpOnWinSessionLock                   int                         `json:"zpaAuthExpOnWinSessionLock,omitempty"`
	ZpaAuthExpSessionLockStateMinTimeInSecond    int                         `json:"zpaAuthExpSessionLockStateMinTimeInSecond,omitempty"`
	PacketTunnelExcludeListForIPv6               string                      `json:"packetTunnelExcludeListForIPv6,omitempty"`
	PacketTunnelIncludeListForIPv6               string                      `json:"packetTunnelIncludeListForIPv6,omitempty"`
	EnableSetProxyOnVPNAdapters                  int                         `json:"enableSetProxyOnVPNAdapters,omitempty"`
	DisableDNSRouteExclusion                     int                         `json:"disableDNSRouteExclusion,omitempty"`
	AdvanceZpaReauth                             bool                        `json:"advanceZpaReauth"`
	UseProxyPortForT1                            string                      `json:"useProxyPortForT1,omitempty"`
	UseProxyPortForT2                            string                      `json:"useProxyPortForT2,omitempty"`
	AllowPacExclusionsOnly                       string                      `json:"allowPacExclusionsOnly,omitempty"`
	InterceptZIATrafficAllAdapters               string                      `json:"interceptZIATrafficAllAdapters,omitempty"`
	EnableAntiTampering                          string                      `json:"enableAntiTampering,omitempty"`
	OverrideATCmdByPolicy                        string                      `json:"overrideATCmdByPolicy,omitempty"`
	ReactivateAntiTamperingTime                  int                         `json:"reactivateAntiTamperingTime,omitempty"`
	EnforceSplitDNS                              int                         `json:"enforceSplitDNS,omitempty"`
	DropQuicTraffic                              int                         `json:"dropQuicTraffic,omitempty"`
	EnableZdpService                             string                      `json:"enableZdpService,omitempty"`
	UpdateDnsSearchOrder                         int                         `json:"updateDnsSearchOrder,omitempty"`
	TruncateLargeUDPDNSResponse                  int                         `json:"truncateLargeUDPDNSResponse,omitempty"`
	PrioritizeDnsExclusions                      int                         `json:"prioritizeDnsExclusions,omitempty"`
	PurgeKerberosPreferredDCCache                string                      `json:"purgeKerberosPreferredDCCache,omitempty"`
	DeleteDHCPOption121Routes                    string                      `json:"deleteDHCPOption121Routes,omitempty"`
	EnableLocationPolicyOverride                 int                         `json:"enableLocationPolicyOverride,omitempty"`
	EnableCustomTheme                            int                         `json:"enableCustomTheme,omitempty"`
	LocationRulesetPolicies                      LocationRulesetPolicies     `json:"locationRulesetPolicies,omitempty"`
	GenerateCliPasswordContract                  GenerateCliPasswordContract `json:"generateCliPasswordContract,omitempty"`
	ZdxLiteConfigObj                             string                      `json:"zdxLiteConfigObj,omitempty"`
	DdilConfig                                   string                      `json:"ddilConfig,omitempty"`
	ZccFailCloseSettingsIpBypasses               string                      `json:"zccFailCloseSettingsIpBypasses,omitempty"`
	ZccFailCloseSettingsExitUninstallPassword    string                      `json:"zccFailCloseSettingsExitUninstallPassword,omitempty"`
	ZccFailCloseSettingsLockdownOnTunnelProcExit int                         `json:"zccFailCloseSettingsLockdownOnTunnelProcessExit,omitempty"`
	ZccFailCloseSettingsLockdownOnFirewallError  int                         `json:"zccFailCloseSettingsLockdownOnFirewallError,omitempty"`
	ZccFailCloseSettingsLockdownOnDriverError    int                         `json:"zccFailCloseSettingsLockdownOnDriverError,omitempty"`
	ZccFailCloseSettingsThumbPrint               string                      `json:"zccFailCloseSettingsThumbPrint,omitempty"`
	ZccAppFailOpenPolicy                         int                         `json:"zccAppFailOpenPolicy,omitempty"`
	ZccTunnelFailPolicy                          int                         `json:"zccTunnelFailPolicy,omitempty"`
	FollowGlobalForPartnerLogin                  string                      `json:"followGlobalForPartnerLogin,omitempty"`
	UserAllowedToAddPartner                      string                      `json:"userAllowedToAddPartner,omitempty"`
	AllowClientCertCachingForWebView2            string                      `json:"allowClientCertCachingForWebView2,omitempty"`
	ShowConfirmationDialogForCachedCert          string                      `json:"showConfirmationDialogForCachedCert,omitempty"`
	EnableFlowBasedTunnel                        int                         `json:"enableFlowBasedTunnel,omitempty"`
	EnableNetworkTrafficProcessMapping           int                         `json:"enableNetworkTrafficProcessMapping,omitempty"`
	EnableLocalPacketCapture                     string                      `json:"enableLocalPacketCapture,omitempty"`
	OneIdMTDeviceAuthEnabled                     string                      `json:"oneIdMTDeviceAuthEnabled,omitempty"`
	EnableCustomProxyDetection                   string                      `json:"enableCustomProxyDetection,omitempty"`
	PreventAutoReauthDuringDeviceLock            string                      `json:"preventAutoReauthDuringDeviceLock,omitempty"`
	UseEndPointLocationForDCSelection            string                      `json:"useEndPointLocationForDCSelection,omitempty"`
	EnableCrashReporting                         int                         `json:"enableCrashReporting,omitempty"`
	RecacheSystemProxy                           string                      `json:"recacheSystemProxy,omitempty"`
	EnableAutomaticPacketCapture                 int                         `json:"enableAutomaticPacketCapture,omitempty"`
	EnableAPCforCriticalSections                 int                         `json:"enableAPCforCriticalSections,omitempty"`
	EnableAPCforOtherSections                    int                         `json:"enableAPCforOtherSections,omitempty"`
	EnablePCAdditionalSpace                      int                         `json:"enablePCAdditionalSpace,omitempty"`
	PcAdditionalSpace                            int                         `json:"pcAdditionalSpace,omitempty"`
	ClientConnectorUiLanguage                    int                         `json:"clientConnectorUiLanguage,omitempty"`
	BlockPrivateRelay                            string                      `json:"blockPrivateRelay,omitempty"`
	BypassDNSTrafficUsingUDPProxy                int                         `json:"bypassDNSTrafficUsingUDPProxy,omitempty"`
	ReconnectTunOnWakeup                         int                         `json:"reconnectTunOnWakeup,omitempty"`
	BrowserAuthType                              string                      `json:"browserAuthType,omitempty"`
	UseDefaultBrowser                            string                      `json:"useDefaultBrowser,omitempty"`
}

type LocationRulesetPolicies struct {
	OffTrusted      LocationPolicy `json:"offTrusted,omitempty"`
	Trusted         LocationPolicy `json:"trusted,omitempty"`
	VpnTrusted      LocationPolicy `json:"vpnTrusted,omitempty"`
	SplitVpnTrusted LocationPolicy `json:"splitVpnTrusted,omitempty"`
}

type LocationPolicy struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type GenerateCliPasswordContract struct {
	EnableCli                      bool `json:"enableCli"`
	AllowZpaDisableWithoutPassword bool `json:"allowZpaDisableWithoutPassword"`
	AllowZiaDisableWithoutPassword bool `json:"allowZiaDisableWithoutPassword"`
	AllowZdxDisableWithoutPassword bool `json:"allowZdxDisableWithoutPassword"`
}

type DisasterRecovery struct {
	EnableZiaDR      bool   `json:"enableZiaDR"`
	EnableZpaDR      bool   `json:"enableZpaDR"`
	ZiaDRMethod      int    `json:"ziaDRMethod,omitempty"`
	ZiaCustomDbUrl   string `json:"ziaCustomDbUrl,omitempty"`
	UseZiaGlobalDb   bool   `json:"useZiaGlobalDb"`
	ZiaGlobalDbUrl   string `json:"ziaGlobalDbUrl,omitempty"`
	ZiaGlobalDbUrlv2 string `json:"ziaGlobalDbUrlv2,omitempty"`
	ZiaDomainName    string `json:"ziaDomainName,omitempty"`
	ZiaRSAPubKeyName string `json:"ziaRSAPubKeyName,omitempty"`
	ZiaRSAPubKey     string `json:"ziaRSAPubKey,omitempty"`
	ZpaDomainName    string `json:"zpaDomainName,omitempty"`
	ZpaRSAPubKeyName string `json:"zpaRSAPubKeyName,omitempty"`
	ZpaRSAPubKey     string `json:"zpaRSAPubKey,omitempty"`
	AllowZiaTest     bool   `json:"allowZiaTest"`
	AllowZpaTest     bool   `json:"allowZpaTest"`
}

type AppService struct {
	Active      bool          `json:"active"`
	AppDataBlob []AppDataBlob `json:"appDataBlob,omitempty"`
}

type AppDataBlob struct {
	Fqdn   string `json:"fqdn,omitempty"`
	Ipaddr string `json:"ipaddr,omitempty"`
	Port   string `json:"port,omitempty"`
}

// GetApplicationProfiles retrieves the paginated list of application profiles.
// The deviceType parameter accepts a friendly name (e.g., "windows", "ios") or
// a numeric string; it is converted to the API's integer value automatically.
func GetApplicationProfiles(ctx context.Context, service *zscaler.Service, search, searchType, deviceType string, page, pageSize *int) (*ApplicationProfilesResponse, *http.Response, error) {
	params := common.QueryParams{
		Search:     search,
		SearchType: searchType,
	}
	if page != nil {
		params.Page = *page
	}
	if pageSize != nil {
		params.PageSize = *pageSize
	}
	if deviceType != "" {
		dt, err := common.GetDeviceTypeByName(deviceType)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid deviceType: %w", err)
		}
		params.DeviceType = dt
	}

	var response ApplicationProfilesResponse
	resp, err := service.Client.NewZccRequestDo(ctx, "GET", appProfilesEndpoint, params, nil, &response)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve application profiles: %w", err)
	}
	return &response, resp, nil
}

// GetByProfileID retrieves a single application profile by its ID.
func GetByProfileID(ctx context.Context, service *zscaler.Service, profileID string) (*ApplicationProfile, *http.Response, error) {
	if profileID == "" {
		return nil, nil, fmt.Errorf("profileId is required")
	}
	endpoint := fmt.Sprintf("%s/%s", appProfilesEndpoint, profileID)

	var profile ApplicationProfile
	resp, err := service.Client.NewZccRequestDo(ctx, "GET", endpoint, nil, nil, &profile)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve application profile %s: %w", profileID, err)
	}
	return &profile, resp, nil
}

// GetByName iterates all application profiles and returns the one matching
// the given name (case-insensitive). Returns an error if not found.
func GetByName(ctx context.Context, service *zscaler.Service, name string) (*ApplicationProfile, *http.Response, error) {
	pageSize := 1000
	page := 1

	for {
		res, resp, err := GetApplicationProfiles(ctx, service, "", "", "", &page, &pageSize)
		if err != nil {
			return nil, resp, err
		}
		for _, p := range res.Policies {
			if strings.EqualFold(p.Name, name) {
				return &p, resp, nil
			}
		}
		if len(res.Policies) < pageSize {
			break
		}
		page++
	}
	return nil, nil, fmt.Errorf("no application profile found with name: %s", name)
}

// PatchApplicationProfile performs a partial update on an application profile.
func PatchApplicationProfile(ctx context.Context, service *zscaler.Service, profileID string, patch *ApplicationProfile) (*ApplicationProfile, *http.Response, error) {
	if profileID == "" {
		return nil, nil, fmt.Errorf("profileId is required")
	}
	if patch == nil {
		return nil, nil, fmt.Errorf("patch body is required")
	}
	endpoint := fmt.Sprintf("%s/%s", appProfilesEndpoint, profileID)

	var updated ApplicationProfile
	resp, err := service.Client.NewZccRequestDo(ctx, "PATCH", endpoint, nil, patch, &updated)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to patch application profile %s: %w", profileID, err)
	}

	service.Client.GetLogger().Printf("[DEBUG] returning application profile from patch: %d", updated.ID)
	return &updated, resp, nil
}
