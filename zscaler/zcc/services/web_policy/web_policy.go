package web_policy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	baseWebPolicyEndpoint     = "/zcc/papi/public/v1/web/policy"
	baseWebAppServiceEndpoint = "/zcc/papi/public/v1/webAppService"
)

// LabelValuePair models the `{label, value}` form-state objects the ZCC
// UI generates for every dropdown / autocomplete it renders, and which
// the /web/policy/edit endpoint then expects to see echoed back in the
// PUT body (ruleOrderSelectedOption, logModeSelected, ipv6ModeSelected,
// billingDaySelectedOption, browserAuthType, ziaDRMethod, etc.).
//
// `Value` is an `any` because the upstream payload is inconsistent —
// some pickers use integer values ({"label":"Debug","value":3}) and
// others use strings ({"label":"1","value":"1"}). Modelling it as a
// concrete type would force callers to pick one form and lose data on
// the other.
type LabelValuePair struct {
	Label string `json:"label"`
	Value any    `json:"value"`
}

// EndToEndDiagnostics is the per-network-context end-to-end diagnostics
// toggle the macOS web policy ships at the top level (and embedded as
// JSON inside policyExtension.zdxLiteConfigObj).
type EndToEndDiagnostics struct {
	Trusted         int `json:"trusted"`
	VpnTrusted      int `json:"vpnTrusted"`
	OffTrusted      int `json:"offTrusted"`
	SplitVpnTrusted int `json:"splitVpnTrusted"`
}

// LocationRulesetEntry / LocationRulesetPolicies model the small
// locationRulesetPolicies block the API expects inside policyExtension.
// Both nested entries are present even when no ruleset is bound (their
// `id` is then 0), so the wire shape is fixed.
type LocationRulesetEntry struct {
	ID int `json:"id"`
}

type LocationRulesetPolicies struct {
	SplitVpnTrusted LocationRulesetEntry `json:"splitVpnTrusted"`
	VpnTrusted      LocationRulesetEntry `json:"vpnTrusted"`
}

// WebPolicy models a ZCC web policy / app profile. The JSON wire shape
// matches what the ZCC OneAPI /web/policy/edit endpoint accepts; the
// response /listByCompany endpoint returns a superset that we read back
// into the same struct.
//
// device_type is serialized as a JSON number on the wire (1=iOS,
// 2=Android, 3=Windows, 4=macOS, 5=Linux). The companion deviceType
// string field that the API also returns on reads (e.g.
// "DEVICE_TYPE_MAC") is intentionally not modelled because callers
// already know which device type they are working with.
//
// The five per-OS policy blocks are pointer-typed with `omitempty` so
// that a payload created for one device type does not also include
// empty blocks for the other four — the API rejects extra OS blocks
// with a 400.
//
// The integer-on-the-wire scalars (ruleOrder, logMode, logLevel, etc.)
// are typed as common.IntOrString because the ZCC API is inconsistent:
// every POST/PUT endpoint requires a JSON number, but the GET response
// quotes some of them as strings ("0", "12", ...). IntOrString marshals
// as a number and unmarshals from either form (or from null / empty
// string, which become 0).
//
// Many fields appear twice on the wire — once at the top level and once
// inside policyExtension — with subtly different types (top-level int
// vs string inside the extension, etc.). The ZCC UI generates both, and
// the API expects both. The Go struct mirrors this faithfully: a Go
// field name suffixed with `Top` lives at the WebPolicy level, while
// the unsuffixed one lives inside PolicyExtension. The DefaultMacosWebPolicy
// constructor populates both with sensible defaults captured from a
// known-working UI-generated request body.
type WebPolicy struct {
	// Core identity / lifecycle. `id` is omitempty so a fresh create body
	// matches the UI shape (no id key at all when the resource is new);
	// the API uses the absence of `id` as the "create" signal on /edit.
	ID          string             `json:"id,omitempty"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Active      string             `json:"active"`
	DeviceType  int                `json:"device_type"`
	RuleOrder   common.IntOrString `json:"ruleOrder"`

	// AllowUnreachablePac is the legacy flag the provider has surfaced
	// since the first SDK v2 days; the API echoes it on reads even when
	// the UI no longer renders it. omitempty keeps it off the wire when
	// the default Go zero value applies, matching the UI capture.
	AllowUnreachablePac bool `json:"allowUnreachablePac,omitempty"`

	// Targeting (groups / users / device groups / app services). The
	// legacy *Ids / *Names slices are omitempty because the UI replaced
	// them with the new `groups` / `users` / `deviceGroups` collections;
	// emitting them as empty arrays when unused diverges from the UI's
	// wire shape and trips some API validations. They are still
	// populated on read responses, so we keep them in the struct.
	Groups                                []any              `json:"groups"`
	Users                                 []any              `json:"users"`
	GroupAll                              common.IntOrString `json:"groupAll"`
	GroupIds                              []int              `json:"groupIds,omitempty"`
	GroupNames                            []string           `json:"groupNames,omitempty"`
	UserIds                               []int              `json:"userIds,omitempty"`
	UserNames                             []string           `json:"userNames,omitempty"`
	AppIdentityNames                      []string           `json:"appIdentityNames,omitempty"`
	AppServiceIds                         []int              `json:"appServiceIds"`
	AppServiceNames                       []string           `json:"appServiceNames,omitempty"`
	AppServiceCustomIdsSelected           []any              `json:"appServiceCustomIdsSelected"`
	BypassAppIds                          []int              `json:"bypassAppIds"`
	BypassCustomAppIds                    []int              `json:"bypassCustomAppIds"`
	BypassMacAppIds                       []any              `json:"bypassMacAppIds"`
	DeviceGroupIds                        []int              `json:"deviceGroupIds,omitempty"`
	DeviceGroupNames                      []string           `json:"deviceGroupNames,omitempty"`
	DeviceGroups                          []any              `json:"deviceGroups"`
	DeviceGroupsOption                    int                `json:"deviceGroupsOption"`
	DeviceGroupsSelected                  []any              `json:"deviceGroupsSelected"`
	UsersOption                           int                `json:"usersOption"`
	UsersSelected                         []any              `json:"usersSelected"`
	ZccFailCloseSettingsAppByPassIdsTop   []int              `json:"zccFailCloseSettingsAppByPassIds"`
	ZccFailCloseSettingsAppByPassSelected []any              `json:"zccFailCloseSettingsAppByPassSelected"`

	// Forwarding / posture profiles. ZiaPostureConfigId is omitempty
	// because test.json does not surface it on a fresh create; we still
	// flatten it on read.
	ForwardingProfileId int   `json:"forwardingProfileId"`
	ZiaPostureProfile   []any `json:"ziaPostureProfile"`
	ZiaPostureConfigId  int   `json:"ziaPostureConfigId,omitempty"`

	// Logging / log mode picker
	LogMode         common.IntOrString `json:"logMode"`
	LogLevel        common.IntOrString `json:"logLevel"`
	LogFileSize     common.IntOrString `json:"logFileSize"`
	LogModeSelected *LabelValuePair    `json:"logModeSelected,omitempty"`

	// Captive portal + diagnostics
	EnableCaptivePortalDetection      int                 `json:"enableCaptivePortalDetection"`
	EnableFailOpen                    int                 `json:"enableFailOpen"`
	CaptivePortalWebSecDisableMinutes int                 `json:"captivePortalWebSecDisableMinutes"`
	CaptivePortalUrlId                []LabelValuePair    `json:"captivePortalUrlId"`
	EndToEndDiagnostics               EndToEndDiagnostics `json:"endToEndDiagnostics"`
	EndToEndDiagnosticsSelected       []any               `json:"endToEndDiagnosticsSelected"`
	LocalMetrics                      int                 `json:"localMetrics"`
	FlowLoggingSelected               []any               `json:"flowLoggingSelected"`
	BlockDomainSelected               []any               `json:"blockDomainSelected"`
	BlockInboundTrafficSelected       []any               `json:"blockInboundTrafficSelected"`
	NotificationTemplateSelected      []any               `json:"notificationTemplateSelected"`

	// PAC config
	PacURL      string `json:"pac_url"`
	PacType     int    `json:"pacType"`
	PacDataPath string `json:"pacDataPath"`

	// MDM / billing / mobile
	Mdm               int    `json:"mdm"`
	Passcode          string `json:"passcode"`
	ExitPassword      string `json:"exit_password"`
	Limit             string `json:"limit"`
	BillingDay        string `json:"billing_day"`
	AllowedApps       string `json:"allowed_apps"`
	CustomText        string `json:"custom_text"`
	BypassMmsApps     int    `json:"bypass_mms_apps"`
	QuotaInRoaming    int    `json:"quota_in_roaming"`
	WifiSSID          string `json:"wifi_ssid"`
	BypassAndroidApps []int  `json:"bypass_android_apps"`
	Enforced          int    `json:"enforced"`

	// Registry / Windows-ish defaults (still required for macOS bodies — the API echoes them)
	RegistryPath                      string `json:"registryPath"`
	RegistryName                      string `json:"registryName"`
	InstallSslCertsTop                common.IntOrString `json:"install_ssl_certs"`
	DisableLoopBackRestriction        int    `json:"disableLoopBackRestriction"`
	RemoveExemptedContainers          int    `json:"removeExemptedContainers"`
	OverrideWPAD                      int    `json:"overrideWPAD"`
	RestartWinHttpSvc                 int    `json:"restartWinHttpSvc"`
	InstallWindowsFirewallInboundRule string `json:"installWindowsFirewallInboundRule"`
	ForceLocationRefreshSccm          int    `json:"forceLocationRefreshSccm"`
	WfpMtr                            int    `json:"wfpMtr"`
	EnableLocalPacketCaptureTabValue  int    `json:"enableLocalPacketCaptureTabValue"`
	RefreshKerberosToken              int    `json:"refreshKerberosToken"`

	// Nullable nested configs the UI sends (defaults are JSON null)
	FlowLoggerConfig             any `json:"flowLoggerConfig"`
	DomainProfileDetectionConfig any `json:"domainProfileDetectionConfig"`
	AllInboundTrafficConfig      any `json:"allInboundTrafficConfig"`

	// Cosmetic / runtime knobs at the top level
	HighlightActiveControl     common.IntOrString `json:"highlightActiveControl"`
	SendDisableServiceReason   common.IntOrString `json:"sendDisableServiceReason"`
	TunnelZappTraffic          common.IntOrString `json:"tunnelZappTraffic"`
	EnableDeviceGroups         common.IntOrString `json:"enableDeviceGroups,omitempty"`
	ReactivateWebSecurityMins  common.IntOrString `json:"reactivateWebSecurityMinutes"`
	ReauthPeriod               common.IntOrString `json:"reauth_period"`
	ClearArpCacheTop           int                `json:"clearArpCache"`
	EnableZscalerFirewallTop   string             `json:"enableZscalerFirewall"`
	PersistentZscalerFirewallTop int              `json:"persistentZscalerFirewall"`
	CacheSystemProxyTop        int                `json:"cacheSystemProxy"`
	DnsPriorityOrderingTop     []string           `json:"dnsPriorityOrdering"`
	EnableZdpServiceTop        int                `json:"enableZdpService"`
	DisableParallelIpv4AndIPv6 int                `json:"disableParallelIpv4AndIPv6"`
	DisableParallelIpv4andIpv6 string             `json:"disableParallelIpv4andIpv6"`

	// Top-level "selected" pickers (UI form-state mirrors)
	RuleOrderSelectedOption    *LabelValuePair  `json:"ruleOrderSelectedOption,omitempty"`
	BillingDaySelectedOption   *LabelValuePair  `json:"billingDaySelectedOption,omitempty"`
	Ipv6ModeSelected           *LabelValuePair  `json:"ipv6ModeSelected,omitempty"`
	ZpaAutoReauthTimeoutTop    []LabelValuePair `json:"zpaAutoReauthTimeout"`
	PcAdditionalSpaceTop       []LabelValuePair `json:"pcAdditionalSpace"`
	BrowserAuthTypeTop         *LabelValuePair  `json:"browserAuthType,omitempty"`
	ClientConnectorUiLanguageSelected []LabelValuePair `json:"clientConnectorUiLanguageSelected"`

	// Machine token / ZPA reauth scheduling
	MachineTokenOption                            int    `json:"machineTokenOption"`
	MachineTokenSelectedOption                    int    `json:"machineTokenSelectedOption"`
	ZpaAuthExpSessionLockStateMinTimeInSecondTop  string `json:"zpaAuthExpSessionLockStateMinTimeInSecond"`
	ForceZpaAuthenticationToExpire                []any  `json:"forceZpaAuthenticationToExpire"`
	ZpaReauthConfigTop                            []any  `json:"zpaReauthConfig"`

	// DR mirrors of the disasterRecovery block (top-level form-state)
	ZiaDRMethodTop *LabelValuePair `json:"ziaDRMethod,omitempty"`

	// Top-level CLI / disable-without-password trio
	AllowZpaDisableWithoutPasswordTop bool `json:"allowZpaDisableWithoutPassword"`
	AllowZiaDisableWithoutPasswordTop bool `json:"allowZiaDisableWithoutPassword"`
	AllowZdxDisableWithoutPasswordTop bool `json:"allowZdxDisableWithoutPassword"`

	// Top-level DNS / split-tunnel flags (these duplicate PolicyExtension entries with different types)
	UseDefaultAdapterForDNSTop  string `json:"useDefaultAdapterForDNS"`
	UpdateDnsSearchOrderTop     string `json:"updateDnsSearchOrder"`
	EnforceSplitDNSTop          string `json:"enforceSplitDNS"`
	DisableDNSRouteExclusionTop string `json:"disableDNSRouteExclusion"`
	EnableSetProxyOnVPNAdaptersTop int `json:"enableSetProxyOnVPNAdapters"`
	DropQuicTrafficTop          string `json:"dropQuicTraffic"`
	FollowRoutingTableTop       string `json:"followRoutingTable"`

	// Top-level partner / fail-close / packet capture / packet tunnel mirrors
	VpnGatewaysTop                                  []any    `json:"vpnGateways"`
	PartnerDomainsTop                               []any    `json:"partnerDomains"`
	ZccFailCloseSettingsIpBypassesTop               []any    `json:"zccFailCloseSettingsIpBypasses"`
	ZccFailCloseSettingsLockdownOnTunnelProcessExitTop int   `json:"zccFailCloseSettingsLockdownOnTunnelProcessExit"`
	ZccFailCloseSettingsExitUninstallPasswordTop    string   `json:"zccFailCloseSettingsExitUninstallPassword"`
	UserAllowedToAddPartnerTop                      int      `json:"userAllowedToAddPartner"`
	FollowGlobalForPartnerLoginTop                  string   `json:"followGlobalForPartnerLogin"`
	FollowGlobalForZpaReauthTop                     string   `json:"followGlobalForZpaReauth"`
	FollowGlobalForPacketCaptureTop                 string   `json:"followGlobalForPacketCapture"`
	EnableLocalPacketCaptureTop                     string   `json:"enableLocalPacketCapture"`
	EnableLocalPacketCaptureV2Top                   []any    `json:"enableLocalPacketCaptureV2"`
	PacketTunnelIncludeListTop                      []string `json:"packetTunnelIncludeList"`
	PacketTunnelExcludeListTop                      []string `json:"packetTunnelExcludeList"`
	PacketTunnelIncludeListForIPv6Top               []string `json:"packetTunnelIncludeListForIPv6"`
	PacketTunnelExcludeListForIPv6Top               []string `json:"packetTunnelExcludeListForIPv6"`
	PacketTunnelDnsIncludeListTop                   []string `json:"packetTunnelDnsIncludeList"`
	PacketTunnelDnsExcludeListTop                   []string `json:"packetTunnelDnsExcludeList"`
	SourcePortBasedBypassesTop                      []string `json:"sourcePortBasedBypasses"`
	UseV8JsEngineTop                                string   `json:"useV8JsEngine"`
	PrioritizeDnsExclusionsTop                      string   `json:"prioritizeDnsExclusions"`

	// Trusted-network buckets the UI mirrors at the top level (empty lists by default)
	VpnTrusted      []any `json:"vpnTrusted"`
	SplitVpnTrusted []any `json:"splitVpnTrusted"`
	Trusted         []any `json:"trusted"`
	OffTrusted      []any `json:"offTrusted"`
	CustomDNSTop    []any `json:"customDNS"`

	// Top-level zCC revert / proxy detection / language / crash reporting
	EnableZCCRevertTop                    bool   `json:"enableZCCRevert"`
	EnableCustomProxyDetectionTop         string `json:"enableCustomProxyDetection"`
	ClientConnectorUiLanguageTop          int    `json:"clientConnectorUiLanguage"`
	OneIdMTDeviceAuthEnabledTop           string `json:"oneIdMTDeviceAuthEnabled"`
	PreventAutoReauthDuringDeviceLockTop  string `json:"preventAutoReauthDuringDeviceLock"`
	InstantForceZPAReauthStateUpdateTop   int    `json:"instantForceZPAReauthStateUpdate"`
	EnableNetworkTrafficProcessMappingTop int    `json:"enableNetworkTrafficProcessMapping"`
	UseEndPointLocationForDCSelectionTop  string `json:"useEndPointLocationForDCSelection"`
	RecacheSystemProxyTop                 string `json:"recacheSystemProxy"`
	EnableLocationPolicyOverrideTop       int    `json:"enableLocationPolicyOverride"`
	BlockPrivateRelayTop                  string `json:"blockPrivateRelay"`
	EnableCrashReportingTop               string `json:"enableCrashReporting"`
	EnableAutomaticPacketCaptureTop       string `json:"enableAutomaticPacketCapture"`
	EnableAPCforCriticalSectionsTop       string `json:"enableAPCforCriticalSections"`
	EnableAPCforOtherSectionsTop          string `json:"enableAPCforOtherSections"`
	EnablePCAdditionalSpaceTop            string `json:"enablePCAdditionalSpace"`

	// Anti-tampering top-level mirror
	ReactivateAntiTamperingTimeTop int `json:"reactivateAntiTamperingTime"`

	// Browser auth defaults (int form, sits next to the BrowserAuthTypeTop picker object)
	UseDefaultBrowserTop int `json:"useDefaultBrowser"`

	// Per-OS embedded policy blocks (only one is non-nil at a time)
	AndroidPolicy *AndroidPolicy `json:"androidPolicy,omitempty"`
	IosPolicy     *IosPolicy     `json:"iosPolicy,omitempty"`
	LinuxPolicy   *LinuxPolicy   `json:"linuxPolicy,omitempty"`
	MacPolicy     *MacPolicy     `json:"macPolicy,omitempty"`
	WindowsPolicy *WindowsPolicy `json:"windowsPolicy,omitempty"`

	// Nested policy/recovery blocks
	PolicyExtension  PolicyExtension  `json:"policyExtension"`
	DisasterRecovery DisasterRecovery `json:"disasterRecovery"`
}

type AndroidPolicy struct {
	AllowedApps      string `json:"allowedApps"`
	BillingDay       string `json:"billingDay"`
	BypassAndroidApp string `json:"bypassAndroidApps"`
	BypassMmsApps    string `json:"bypassMmsApps"`
	CustomText       string `json:"customText"`
	DisablePassword  string `json:"disablePassword"`
	EnableVerboseLog string `json:"enableVerboseLog"`
	Enforced         string `json:"enforced"`
	InstallCerts     string `json:"installCerts"`
	Limit            string `json:"limit"`
	LogoutPassword   string `json:"logoutPassword"`
	QuotaRoaming     string `json:"quotaRoaming"`
	UninstallPass    string `json:"uninstallPassword"`
	WifiSsid         string `json:"wifissid"`
}

type IosPolicy struct {
	DisablePassword        string `json:"disablePassword"`
	Ipv6Mode               string `json:"ipv6Mode"`
	LogoutPassword         string `json:"logoutPassword"`
	Passcode               string `json:"passcode"`
	ShowVPNTunNotification string `json:"showVPNTunNotification"`
	UninstallPassword      string `json:"uninstallPassword"`
}

type LinuxPolicy struct {
	DisablePassword   string `json:"disablePassword"`
	InstallCerts      string `json:"installCerts"`
	LogoutPassword    string `json:"logoutPassword"`
	UninstallPassword string `json:"uninstallPassword"`
}

// MacPolicy models the macPolicy block of a ZCC web policy. JSON tags
// match the wire shape the API actually uses, captured from a working
// UI-generated request body:
//
//   - the password/cert fields are serialized in snake_case
//     (disable_password / install_ssl_certs / logout_password /
//     uninstall_password); earlier versions of this struct used
//     camelCase tags which the API silently ignored, leading to
//     {"success":"false","id":0} on /edit.
//   - install_ssl_certs is a JSON number, not a string.
//   - browser_auth_type / use_default_browser / captive_portal_config
//     are required by the API for macOS policy creates.
type MacPolicy struct {
	AddIfscopeRoute                      string             `json:"addIfscopeRoute"`
	BrowserAuthType                      int                `json:"browserAuthType"`
	CacheSystemProxy                     string             `json:"cacheSystemProxy"`
	CaptivePortalConfig                  string             `json:"captivePortalConfig"`
	ClearArpCache                        string             `json:"clearArpCache"`
	DisablePassword                      string             `json:"disable_password"`
	DnsPriorityOrdering                  string             `json:"dnsPriorityOrdering"`
	DnsPriorityOrderingForTrustedDnsCrit string             `json:"dnsPriorityOrderingForTrustedDnsCriteria"`
	EnableAppBasedBypass                 string             `json:"enableApplicationBasedBypass"`
	EnableZscalerFirewall                string             `json:"enableZscalerFirewall"`
	InstallSslCerts                      common.IntOrString `json:"install_ssl_certs"`
	LogoutPassword                       string             `json:"logout_password"`
	PersistentZscalerFirewall            string             `json:"persistentZscalerFirewall"`
	UninstallPassword                    string             `json:"uninstall_password"`
	UseDefaultBrowser                    int                `json:"useDefaultBrowser"`
}

type WindowsPolicy struct {
	CacheSystemProxy              int    `json:"cacheSystemProxy"`
	CaptivePortalConfig           string `json:"captivePortalConfig"`
	DisableLoopBackRestriction    int    `json:"disableLoopBackRestriction"`
	DisableParallelIpv4andIpv6    string `json:"disableParallelIpv4andIpv6"`
	DisablePassword               string `json:"disablePassword"`
	FlowLoggerConfig              string `json:"flowLoggerConfig"`
	ForceLocationRefreshSccm      int    `json:"forceLocationRefreshSccm"`
	InstallWindowsFirewallInbound int    `json:"installWindowsFirewallInboundRule"`
	InstallCerts                  string `json:"installCerts"`
	LogoutPassword                string `json:"logoutPassword"`
	OverrideWPAD                  int    `json:"overrideWPAD"`
	PacDataPath                   string `json:"pacDataPath"`
	PacType                       int    `json:"pacType"`
	PrioritizeIPv4                int    `json:"prioritizeIPv4"`
	RemoveExemptedContainers      int    `json:"removeExemptedContainers"`
	RestartWinHttpSvc             int    `json:"restartWinHttpSvc"`
	TriggerDomainProfleDetection  int    `json:"triggerDomainProfleDetection"`
	UninstallPassword             string `json:"uninstallPassword"`
	WfpDriver                     int    `json:"wfpDriver"`
}

// DisasterRecovery models the disasterRecovery block inside a WebPolicy.
// JSON tags match the wire shape the ZCC API uses (verified against a
// real /listByCompany response): ziaDRMethod (not ziaDRRecoveryMethod),
// ziaRSAPubKeyName / ziaRSAPubKey, zpaRSAPubKeyName / zpaRSAPubKey, and
// the supplemental ziaCustomDbUrl field that custom DR setups use.
//
// PolicyId / ZiaGlobalDbURL / ZiaGlobalDbURLV2 / ZiaPacURL are
// omitempty because the UI capture never emits them in a fresh create
// body — they're populated on read responses but absent on writes.
type DisasterRecovery struct {
	AllowZiaTest     bool   `json:"allowZiaTest"`
	AllowZpaTest     bool   `json:"allowZpaTest"`
	EnableZiaDR      bool   `json:"enableZiaDR"`
	EnableZpaDR      bool   `json:"enableZpaDR"`
	PolicyId         string `json:"policyId,omitempty"`
	UseZiaGlobalDb   bool   `json:"useZiaGlobalDb"`
	ZiaDRMethod      int    `json:"ziaDRMethod"`
	ZiaCustomDbUrl   string `json:"ziaCustomDbUrl"`
	ZiaDomainName    string `json:"ziaDomainName"`
	ZiaGlobalDbURL   string `json:"ziaGlobalDbUrl,omitempty"`
	ZiaGlobalDbURLV2 string `json:"ziaGlobalDbUrlv2,omitempty"`
	ZiaPacURL        string `json:"ziaPacUrl,omitempty"`
	ZiaRSAPubKey     string `json:"ziaRSAPubKey"`
	ZiaRSAPubKeyName string `json:"ziaRSAPubKeyName"`
	ZpaDomainName    string `json:"zpaDomainName"`
	ZpaRSAPubKey     string `json:"zpaRSAPubKey"`
	ZpaRSAPubKeyName string `json:"zpaRSAPubKeyName"`
}

// PolicyExtension models the policyExtension nested block. The field
// set mirrors the wire shape of the UI-generated request body — every
// flag the API expects is present, with the wire type the API uses
// (numbers for some, quoted strings for others; the API is inconsistent
// and the SDK matches it field-by-field).
type PolicyExtension struct {
	// CLI / disable-without-password contract
	GenerateCliPasswordContract GenerateCliPasswordContract `json:"generateCliPasswordContract"`

	// VPN / partner / packet-capture
	VpnGateways                                     string `json:"vpnGateways"`
	PartnerDomains                                  string `json:"partnerDomains"`
	ZccFailCloseSettingsIpBypasses                  string `json:"zccFailCloseSettingsIpBypasses"`
	ZccFailCloseSettingsLockdownOnTunnelProcessExit string `json:"zccFailCloseSettingsLockdownOnTunnelProcessExit"`
	ZccFailCloseSettingsExitUninstallPassword       string `json:"zccFailCloseSettingsExitUninstallPassword"`
	ZccFailCloseSettingsAppByPassIds                []int    `json:"zccFailCloseSettingsAppByPassIds"`
	ZccFailCloseSettingsAppByPassNames              []string `json:"zccFailCloseSettingsAppByPassNames,omitempty"`
	ZccFailCloseSettingsThumbPrint                  string   `json:"zccFailCloseSettingsThumbPrint,omitempty"`
	ZccFailCloseSettingsLockdownOnFirewallError     string `json:"zccFailCloseSettingsLockdownOnFirewallError"`
	ZccFailCloseSettingsLockdownOnDriverError       string `json:"zccFailCloseSettingsLockdownOnDriverError"`
	UserAllowedToAddPartner                         string `json:"userAllowedToAddPartner"`
	FollowGlobalForPartnerLogin                     string `json:"followGlobalForPartnerLogin"`
	FollowGlobalForZpaReauth                        string `json:"followGlobalForZpaReauth"`
	FollowGlobalForPacketCapture                    string `json:"followGlobalForPacketCapture"`
	EnableLocalPacketCapture                        string `json:"enableLocalPacketCapture"`
	EnableLocalPacketCaptureV2                      int    `json:"enableLocalPacketCaptureV2"`
	EnableFlowBasedTunnel                           string `json:"enableFlowBasedTunnel"`

	// ZPA reauth scheduling. The fields are common.IntOrString because the
	// upstream API returns some as quoted strings and others as numbers
	// depending on which read endpoint you hit. PUT/PUT requires numbers,
	// which IntOrString's MarshalJSON guarantees.
	ZpaReauthConfig                   any                `json:"zpaReauthConfig"`
	ZpaAutoReauthTimeout              common.IntOrString `json:"zpaAutoReauthTimeout"`
	ZpaAuthExpOnSleep                 common.IntOrString `json:"zpaAuthExpOnSleep"`
	ZpaAuthExpOnSysRestart            common.IntOrString `json:"zpaAuthExpOnSysRestart"`
	ZpaAuthExpOnNetIpChange           common.IntOrString `json:"zpaAuthExpOnNetIpChange"`
	InstantForceZPAReauthStateUpdate  common.IntOrString `json:"instantForceZPAReauthStateUpdate"`
	ZpaAuthExpOnWinLogonSession       common.IntOrString `json:"zpaAuthExpOnWinLogonSession"`
	ZpaAuthExpOnWinSessionLock        common.IntOrString `json:"zpaAuthExpOnWinSessionLock"`
	ZpaAuthExpSessionLockStateMinTime string             `json:"zpaAuthExpSessionLockStateMinTimeInSecond"`
	AdvanceZpaReauth                  bool               `json:"advanceZpaReauth"`
	AdvanceZpaReauthTime              int                `json:"advanceZpaReauthTime,omitempty"`

	// DNS / split-tunnel / packet tunnel CSV strings
	ExitPassword                    string `json:"exitPassword"`
	FollowRoutingTable              string `json:"followRoutingTable"`
	UseDefaultAdapterForDNS         string `json:"useDefaultAdapterForDNS"`
	UpdateDnsSearchOrder            string `json:"updateDnsSearchOrder"`
	UseZscalerNotificationFramework string `json:"useZscalerNotificationFramework"`
	SwitchFocusToNotification       string `json:"switchFocusToNotification"`
	FallbackToGatewayDomain         string `json:"fallbackToGatewayDomain"`
	UseProxyPortForT1               string `json:"useProxyPortForT1"`
	UseProxyPortForT2               string `json:"useProxyPortForT2"`
	AllowPacExclusionsOnly          string `json:"allowPacExclusionsOnly"`
	UseWsaPollForZpa                string `json:"useWsaPollForZpa"`
	EnableZCCRevert                 string `json:"enableZCCRevert"`
	ZccRevertPassword               string `json:"zccRevertPassword"`
	EnableSetProxyOnVPNAdapters     string             `json:"enableSetProxyOnVPNAdapters"`
	DisableDNSRouteExclusion        common.IntOrString `json:"disableDNSRouteExclusion"`
	PacketTunnelIncludeListForIPv6  string             `json:"packetTunnelIncludeListForIPv6"`
	InterceptZIATrafficAllAdapters  common.IntOrString `json:"interceptZIATrafficAllAdapters"`
	EnableAntiTampering             common.IntOrString `json:"enableAntiTampering"`
	ReactivateAntiTamperingTime     int                `json:"reactivateAntiTamperingTime"`
	SourcePortBasedBypasses         string             `json:"sourcePortBasedBypasses"`
	EnforceSplitDNS                 common.IntOrString `json:"enforceSplitDNS"`
	DropQuicTraffic                 common.IntOrString `json:"dropQuicTraffic"`
	ZdpDisablePassword              string `json:"zdpDisablePassword"`
	UseV8JsEngine                   string `json:"useV8JsEngine"`
	ZdDisablePassword               string `json:"zdDisablePassword"`
	ZdxDisablePassword              string `json:"zdxDisablePassword"`
	ZpaDisablePassword              string `json:"zpaDisablePassword"`
	BypassDNSTrafficUsingUDPProxy   string `json:"bypassDNSTrafficUsingUDPProxy"`
	ReconnectTunOnWakeup            string `json:"reconnectTunOnWakeup"`
	EnableCustomTheme               int    `json:"enableCustomTheme"`
	DeleteDHCPOption121Routes       string `json:"deleteDHCPOption121Routes"`
	MachineIdpAuth                  bool   `json:"machineIdpAuth"`
	Nonce                           string `json:"nonce"`
	PacketTunnelDnsExcludeList      string `json:"packetTunnelDnsExcludeList"`
	PacketTunnelDnsIncludeList      string `json:"packetTunnelDnsIncludeList"`
	PacketTunnelExcludeList         string `json:"packetTunnelExcludeList"`
	PacketTunnelExcludeListForIPv6  string `json:"packetTunnelExcludeListForIPv6"`
	PacketTunnelIncludeList         string `json:"packetTunnelIncludeList"`
	TruncateLargeUDPDNSResponse     common.IntOrString      `json:"truncateLargeUDPDNSResponse"`
	OverrideATCmdByPolicy           common.IntOrString      `json:"overrideATCmdByPolicy"`
	PurgeKerberosPreferredDCCache   common.IntOrString      `json:"purgeKerberosPreferredDCCache"`
	RscModeOnAllAdapters            common.IntOrString      `json:"rscModeOnAllAdapters"`
	EnableAdapterHardwareOffloading common.IntOrString      `json:"enableAdapterHardwareOffloading"`
	SupportZPASearchDomainsInTRP    common.IntOrString      `json:"supportZPASearchDomainsInTRP"`
	PrioritizeDnsExclusions         common.IntOrString      `json:"prioritizeDnsExclusions"`
	LocationRulesetPolicies         LocationRulesetPolicies `json:"locationRulesetPolicies"`
	DdilConfig                      string                  `json:"ddilConfig"`
	ZccAppFailOpenPolicy            common.IntOrString      `json:"zccAppFailOpenPolicy"`
	ZccTunnelFailPolicy             common.IntOrString      `json:"zccTunnelFailPolicy"`

	// Certificate caching / device auth / process mapping
	AllowClientCertCachingForWebView2   string `json:"allowClientCertCachingForWebView2"`
	ShowConfirmationDialogForCachedCert string `json:"showConfirmationDialogForCachedCert"`
	OneIdMTDeviceAuthEnabled            string `json:"oneIdMTDeviceAuthEnabled"`
	PreventAutoReauthDuringDeviceLock   string             `json:"preventAutoReauthDuringDeviceLock"`
	ClientConnectorUiLanguage           common.IntOrString `json:"clientConnectorUiLanguage"`
	EnableNetworkTrafficProcessMapping  common.IntOrString `json:"enableNetworkTrafficProcessMapping"`
	UseEndPointLocationForDCSelection   string             `json:"useEndPointLocationForDCSelection"`
	RecacheSystemProxy                  string             `json:"recacheSystemProxy"`
	EnableLocationPolicyOverride        common.IntOrString `json:"enableLocationPolicyOverride"`
	BlockPrivateRelay                   string             `json:"blockPrivateRelay"`
	EnableAutomaticPacketCapture        string             `json:"enableAutomaticPacketCapture"`
	EnableAPCforCriticalSections        string             `json:"enableAPCforCriticalSections"`
	EnableAPCforOtherSections            string            `json:"enableAPCforOtherSections"`
	EnablePCAdditionalSpace             string             `json:"enablePCAdditionalSpace"`
	PcAdditionalSpace                   string             `json:"pcAdditionalSpace"`
	EnableCustomProxyDetection          string             `json:"enableCustomProxyDetection"`
	EnableCrashReporting                string             `json:"enableCrashReporting"`
	EnableZdpService                    common.IntOrString `json:"enableZdpService"`

	// Custom DNS / ZDX-lite trailing fields
	CustomDNS        string `json:"customDNS,omitempty"`
	ZdxLiteConfigObj string `json:"zdxLiteConfigObj"`
}

// GenerateCliPasswordContract carries the per-product disable-without-
// password switches inside policyExtension.generateCliPasswordContract.
// The API expects all four toggles to be present even when CLI disable
// is off; the UI capture always emits them. PolicyId is omitempty
// because the UI does not send it on create — it's only populated on
// read responses for already-saved policies.
type GenerateCliPasswordContract struct {
	EnableCli                      bool `json:"enableCli"`
	AllowZpaDisableWithoutPassword bool `json:"allowZpaDisableWithoutPassword"`
	AllowZiaDisableWithoutPassword bool `json:"allowZiaDisableWithoutPassword"`
	AllowZdxDisableWithoutPassword bool `json:"allowZdxDisableWithoutPassword"`
	PolicyId                       int  `json:"policyId,omitempty"`
}
type WebPolicyActivation struct {
	DeviceType int `json:"deviceType"`
	PolicyId   int `json:"policyId"`
}

// UnmarshalJSON allows WebPolicy.ID to be decoded from either a JSON string
// ("123") or a JSON number (123). The /edit endpoint returns the id as an
// unquoted integer ({"success":"true","id":205241}), while the listByCompany
// endpoint typically returns it as a string. The custom unmarshaler bridges
// both shapes so the typed struct is always populated correctly.
func (w *WebPolicy) UnmarshalJSON(data []byte) error {
	type alias WebPolicy
	aux := &struct {
		ID json.Number `json:"id"`
		*alias
	}{alias: (*alias)(w)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	w.ID = aux.ID.String()
	return nil
}

// EditResponse models the bare response returned by PUT /edit, which is of
// the form {"success":"true","id":205241}. The id is captured as a
// json.Number so it survives both quoted and unquoted JSON.
type EditResponse struct {
	Success string      `json:"success"`
	ID      json.Number `json:"id"`
}

// ErrWebPolicyPartialDecode is returned alongside a minimally-hydrated
// *WebPolicy when GetWebPolicyByID found the requested record in the
// listByCompany response but could not strict-decode it into the typed
// struct (typically because one or more fields in the response use a
// JSON number where the struct declares a Go string, e.g. `"active":1.0`).
//
// Callers should treat the returned *WebPolicy as best-effort: the id,
// name, description, device_type and active fields are populated from
// the raw map but every other field will be at its Go zero value. The
// recommended pattern is for the caller (Create/Update) to fall back to
// the request body it just wrote, since the API has already confirmed
// success and the PUT body is authoritative state.
var ErrWebPolicyPartialDecode = errors.New("web policy: strict decode of listByCompany entry failed")

// GetWebPolicyByID fetches a single web policy filtered by its id within a
// device type. The underlying listByCompany endpoint returns untyped maps
// where many fields come back as JSON numbers even though the strict
// WebPolicy struct types them as strings (e.g. `"active":1.0`), so the
// lookup is done in two phases:
//
//  1. Match the requested id against the raw map directly, normalizing
//     across the float64 / json.Number / string shapes that the JSON
//     decoder may produce for the `id` field.
//  2. Once matched, attempt a strict json.Unmarshal of that entry into
//     the typed WebPolicy. If the strict decode succeeds, return the
//     full struct with nil error. If it fails because of a field-level
//     type mismatch, return a minimal hydration from the raw map (id,
//     name, description, device_type, active) together with the
//     sentinel ErrWebPolicyPartialDecode so the caller can choose to
//     use its own request body as the authoritative source of state.
//
// deviceType uses the integer convention defined in common (1=iOS,
// 2=Android, 3=Windows, 4=macOS, 5=Linux).
func GetWebPolicyByID(ctx context.Context, service *zscaler.Service, id string, deviceType int) (*WebPolicy, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	dt := deviceType
	raw, err := GetPolicyListByCompanyID(ctx, service, nil, nil, nil, nil, &dt)
	if err != nil {
		return nil, fmt.Errorf("failed to list web policies (deviceType=%d): %w", deviceType, err)
	}
	for _, m := range raw {
		if normalizeWebPolicyID(m["id"]) != id {
			continue
		}
		bytes, marshalErr := json.Marshal(m)
		if marshalErr != nil {
			return nil, fmt.Errorf("web policy with id %q (deviceType=%d): re-marshal raw entry: %w", id, deviceType, marshalErr)
		}
		var policy WebPolicy
		if unmarshalErr := json.Unmarshal(bytes, &policy); unmarshalErr != nil {
			service.Client.GetLogger().Printf("[WARN] web policy %s (deviceType=%d) strict decode failed: %v — returning minimal record from raw map", id, deviceType, unmarshalErr)
			return minimalWebPolicyFromRaw(m, id, deviceType), fmt.Errorf("%w: %v", ErrWebPolicyPartialDecode, unmarshalErr)
		}
		// Ensure id+device_type are populated even if the response left
		// either blank (the listByCompany endpoint sometimes drops the
		// numeric id when the entry came from a fresh /edit write).
		if policy.ID == "" {
			policy.ID = id
		}
		if policy.DeviceType == 0 {
			policy.DeviceType = deviceType
		}
		return &policy, nil
	}
	return nil, fmt.Errorf("web policy with id %q (deviceType=%d) not found", id, deviceType)
}

// normalizeWebPolicyID renders the raw `id` value into the canonical
// string form the rest of the SDK (and the /edit response) uses. The
// listByCompany payload encodes ids as JSON numbers (e.g. `507745.0`),
// which json.Unmarshal into `interface{}` materializes as float64; the
// PUT /edit response uses a json.Number. Both must compare equal to the
// stringified id the caller already has in hand.
func normalizeWebPolicyID(v interface{}) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case json.Number:
		return x.String()
	case float64:
		if x == math.Trunc(x) {
			return strconv.FormatInt(int64(x), 10)
		}
		return strconv.FormatFloat(x, 'f', -1, 64)
	case float32:
		return normalizeWebPolicyID(float64(x))
	case int:
		return strconv.Itoa(x)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	}
	return fmt.Sprintf("%v", v)
}

// minimalWebPolicyFromRaw hydrates the subset of WebPolicy fields the
// provider's flatten path reads from when the full strict decode of the
// listByCompany entry has failed. Anything not surfaced here will keep
// its zero value, which is fine because RunUpsert overlays the just-sent
// payload on top of this struct before persisting state.
func minimalWebPolicyFromRaw(m map[string]interface{}, id string, deviceType int) *WebPolicy {
	p := &WebPolicy{ID: id, DeviceType: deviceType}
	if v, ok := m["name"]; ok {
		if s, ok := v.(string); ok {
			p.Name = s
		}
	}
	if v, ok := m["description"]; ok {
		if s, ok := v.(string); ok {
			p.Description = s
		}
	}
	if v, ok := m["active"]; ok {
		p.Active = normalizeWebPolicyID(v)
	}
	return p
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

// UpdateWebPolicy creates or updates a ZCC web policy. The /edit endpoint
// is used for both operations — supplying an empty/zero id creates a new
// policy, while a populated id updates the existing record. The API
// response is of the form {"success":"true","id":205241}; callers should
// use the returned EditResponse.ID combined with the policy's deviceType
// to refetch the full record via GetWebPolicyByID.
func UpdateWebPolicy(ctx context.Context, service *zscaler.Service, webPolicy *WebPolicy) (*EditResponse, error) {
	if webPolicy == nil {
		return nil, errors.New("web policy is required")
	}

	url := fmt.Sprintf("%s/edit", baseWebPolicyEndpoint)

	var resp EditResponse
	_, err := service.Client.NewZccRequestDo(ctx, "PUT", url, nil, webPolicy, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to update web policy: %w", err)
	}

	if resp.Success != "true" {
		return nil, fmt.Errorf("API rejected web policy update (success=%q)", resp.Success)
	}

	service.Client.GetLogger().Printf("[DEBUG] web policy /edit response: success=%s id=%s", resp.Success, resp.ID.String())
	return &resp, nil
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
