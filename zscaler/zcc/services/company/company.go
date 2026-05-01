package company

import (
	"context"
	"errors"
	"fmt"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	getCompanyInfoEndpoint = "/zcc/papi/public/v1/getCompanyInfo"
	setCompanyInfoEndpoint = "/zcc/papi/public/v1/setCompanyInfo"
)

// CompanyInfo represents the full company configuration returned by the
// GET /getCompanyInfo endpoint and accepted by the PUT /setCompanyInfo
// endpoint. Field types intentionally mirror the API payload (some flags
// are returned as numeric strings like "0"/"1" while others are integers).
type CompanyInfo struct {
	OrgID                                    string                   `json:"orgId,omitempty"`
	MasterCustomerID                         string                   `json:"masterCustomerId,omitempty"`
	Name                                     string                   `json:"name,omitempty"`
	BusinessName                             string                   `json:"businessName"`
	BusinessContactNumber                    string                   `json:"businessContactNumber"`
	ActivationRecipient                      string                   `json:"activationRecipient"`
	ActivationCopy                           string                   `json:"activationCopy"`
	MdmStatus                                string                   `json:"mdmStatus"`
	SendEmail                                string                   `json:"sendEmail"`
	ProxyEnabled                             string                   `json:"proxyEnabled"`
	ZpnEnabled                               string                   `json:"zpnEnabled"`
	UpmEnabled                               string                   `json:"upmEnabled"`
	ZadEnabled                               string                   `json:"zadEnabled"`
	EnableDeceptionForAll                    string                   `json:"enableDeceptionForAll"`
	DlpEnabled                               string                   `json:"dlpEnabled"`
	TunnelProtocolType                       string                   `json:"tunnelProtocolType"`
	SecureAgentBasic                         string                   `json:"secureAgentBasic"`
	SecureAgentAdvanced                      string                   `json:"secureAgentAdvanced"`
	SupportAdminEmail                        string                   `json:"supportAdminEmail"`
	SupportEnabled                           int                      `json:"supportEnabled"`
	FetchLogsForAdminsEnabled                int                      `json:"fetchLogsForAdminsEnabled"`
	EnableRectifyUtils                       int                      `json:"enableRectifyUtils"`
	SupportTicketEnabled                     int                      `json:"supportTicketEnabled"`
	DisableLoggingControls                   int                      `json:"disableLoggingControls"`
	DefaultAuthType                          int                      `json:"defaultAuthType"`
	Version                                  string                   `json:"version,omitempty"`
	PolicyActivationRequired                 int                      `json:"policyActivationRequired"`
	EnableAutofillUsername                   int                      `json:"enableAutofillUsername"`
	AutoFillUsingLoginHint                   int                      `json:"autoFillUsingLoginHint"`
	DcServiceReadOnly                        int                      `json:"dcServiceReadOnly"`
	EnableTunnelZappTrafficToggle            string                   `json:"enableTunnelZappTrafficToggle"`
	MachineIdpAuth                           string                   `json:"machineIdpAuth"`
	LinuxVisibility                          string                   `json:"linuxVisibility"`
	RegistryPathForPac                       string                   `json:"registryPathForPac"`
	UsePollsetForSocketReactor               string                   `json:"usePollsetForSocketReactor"`
	EnableDtlsForZpa                         string                   `json:"enableDtlsForZpa"`
	UseV8JsEngine                            string                   `json:"useV8JsEngine"`
	DisableParallelIpv4AndIPv6               string                   `json:"disableParallelIpv4AndIPv6"`
	Send64BitBuild                           string                   `json:"send64BitBuild"`
	UseAddIfscopeRoute                       string                   `json:"useAddIfscopeRoute"`
	UseClearArpCache                         string                   `json:"useClearArpCache"`
	UseDnsPriorityOrdering                   string                   `json:"useDnsPriorityOrdering"`
	EnableBrowserAuth                        string                   `json:"enableBrowserAuth"`
	EnablePublicAPI                          string                   `json:"enablePublicAPI"`
	DisableReasonVisibility                  string                   `json:"disableReasonVisibility"`
	FollowRoutingTable                       string                   `json:"followRoutingTable"`
	UseDefaultAdapterForDNS                  string                   `json:"useDefaultAdapterForDNS"`
	EnableMinimumDeviceCleanupAsOne          string                   `json:"enableMinimumDeviceCleanupAsOne"`
	DnsPriorityOrderingForTrustedDnsCriteria string                   `json:"dnsPriorityOrderingForTrustedDnsCriteria"`
	MachineTunnelPosture                     string                   `json:"machineTunnelPosture"`
	ZpaPartnerLogin                          string                   `json:"zpaPartnerLogin"`
	ProxyPort                                int                      `json:"proxyPort"`
	DnsCacheTtlWindows                       int                      `json:"dnsCacheTtlWindows"`
	DnsCacheTtlMac                           int                      `json:"dnsCacheTtlMac"`
	DnsCacheTtlAndroid                       int                      `json:"dnsCacheTtlAndroid"`
	DnsCacheTtlIos                           int                      `json:"dnsCacheTtlIos"`
	DnsCacheTtlLinux                         int                      `json:"dnsCacheTtlLinux"`
	ZpaClientCertExpInDays                   int                      `json:"zpaClientCertExpInDays"`
	EnableFlowLogger                         string                   `json:"enableFlowLogger"`
	FlowLoggingBufferLimit                   int                      `json:"flowLoggingBufferLimit"`
	FlowLoggingTimeInterval                  int                      `json:"flowLoggingTimeInterval"`
	PostureBasedService                      string                   `json:"postureBasedService"`
	EnablePostureBasedProfile                string                   `json:"enablePostureBasedProfile"`
	DisasterRecovery                         string                   `json:"disasterRecovery"`
	ZiaGlobalDbUrlForDR                      string                   `json:"ziaGlobalDbUrlForDR,omitempty"`
	EnableReactUI                            string                   `json:"enableReactUI"`
	LaunchReactUIbyDefault                   string                   `json:"launchReactUIbyDefault"`
	DlpNotification                          string                   `json:"dlpNotification"`
	VpnGatewayCharLimit                      int                      `json:"vpnGatewayCharLimit"`
	DeviceGroupsCount                        int                      `json:"deviceGroupsCount"`
	VpnBypassRefreshInterval                 int                      `json:"vpnBypassRefreshInterval"`
	DestIncludeExcludeCharLimit              int                      `json:"destIncludeExcludeCharLimit"`
	IpV6SupportForTunnel2                    string                   `json:"ipV6SupportForTunnel2"`
	DestIncludeExcludeCharLimitForIpv6       int                      `json:"destIncludeExcludeCharLimitForIpv6"`
	EnableSetProxyOnVPNAdapters              string                   `json:"enableSetProxyOnVPNAdapters"`
	DisableDNSRouteExclusion                 string                   `json:"disableDNSRouteExclusion"`
	ShowVPNTunNotification                   string                   `json:"showVPNTunNotification"`
	AddAppBypassToVPNGateway                 string                   `json:"addAppBypassToVPNGateway"`
	EnableZscalerFirewall                    string                   `json:"enableZscalerFirewall"`
	PersistentZscalerFirewall                string                   `json:"persistentZscalerFirewall"`
	ClearMupCache                            string                   `json:"clearMupCache"`
	ExecuteGpoUpdate                         string                   `json:"executeGpoUpdate"`
	EnablePortBasedZPAFilter                 string                   `json:"enablePortBasedZPAFilter"`
	EnableAntiTampering                      string                   `json:"enableAntiTampering"`
	ZpaReauthEnabled                         int                      `json:"zpaReauthEnabled"`
	ZpaAutoReauthTimeout                     int                      `json:"zpaAutoReauthTimeout"`
	EnableZpaAuthUserName                    int                      `json:"enableZpaAuthUserName"`
	EnableGlobalZCCTelemetry                 int                      `json:"enableGlobalZCCTelemetry"`
	ConfigureTunnel2fallbackForZia           string                   `json:"configureTunnel2fallbackForZia"`
	WebAppConfig                             *WebAppConfig            `json:"webAppConfig,omitempty"`
	EnableInstallWebView2                    string                   `json:"enableInstallWebView2"`
	EnableCustomProxyPorts                   string                   `json:"enableCustomProxyPorts"`
	InterceptZIATrafficAllAdapters           string                   `json:"interceptZIATrafficAllAdapters"`
	SwaggerLink                              string                   `json:"swaggerLink,omitempty"`
	EnableOneIdAdmin                         string                   `json:"enableOneIdAdmin"`
	EnableOneIdUser                          string                   `json:"enableOneIdUser"`
	RestrictAdminAccess                      string                   `json:"restrictAdminAccess"`
	EnableZiaUserDepartmentSync              string                   `json:"enableZiaUserDepartmentSync"`
	EnableUDPTransportSelection              string                   `json:"enableUDPTransportSelection"`
	ComputeDeviceGroupsForZIA                string                   `json:"computeDeviceGroupsForZIA"`
	ComputeDeviceGroupsForZPA                string                   `json:"computeDeviceGroupsForZPA"`
	ComputeDeviceGroupsForZDX                string                   `json:"computeDeviceGroupsForZDX"`
	ComputeDeviceGroupsForZAD                string                   `json:"computeDeviceGroupsForZAD"`
	UseTunnel2SmeForTunnel1                  string                   `json:"useTunnel2SmeForTunnel1"`
	MaCloudName                              string                   `json:"maCloudName,omitempty"`
	ZiaCloudName                             string                   `json:"ziaCloudName,omitempty"`
	Zt2HealthProbeInterval                   int                      `json:"zt2HealthProbeInterval"`
	DevicePostureFrequency                   []DevicePostureFrequency `json:"devicePostureFrequency,omitempty"`
	ZdxManualRollout                         string                   `json:"zdxManualRollout"`
	WinZdxLiteEnabled                        string                   `json:"winZdxLiteEnabled"`
	TelemetryDefault                         int                      `json:"telemetryDefault"`
}

// WebAppConfig holds the nested webAppConfig object returned by the
// company info endpoint. All values are strings as the API serializes
// them.
type WebAppConfig struct {
	EnableFipsMode                                  string `json:"enableFipsMode"`
	DeviceCleanup                                   string `json:"deviceCleanup"`
	SyncTimeHours                                   string `json:"syncTimeHours"`
	HideNonFedSettings                              string `json:"hideNonFedSettings"`
	HideAuditLogs                                   string `json:"hideAuditLogs"`
	ActivatePolicy                                  string `json:"activatePolicy"`
	TrustedNetwork                                  string `json:"trustedNetwork"`
	ProcessPostures                                 string `json:"processPostures"`
	ZpaReauth                                       string `json:"zpaReauth"`
	InactiveDeviceCleanup                           string `json:"inactiveDeviceCleanup"`
	ZpaAuthUsername                                 string `json:"zpaAuthUsername"`
	MachineTunnel                                   string `json:"machineTunnel"`
	CacheSystemProxy                                string `json:"cacheSystemProxy"`
	HideDTLSSupportSettings                         string `json:"hideDTLSSupportSettings"`
	MachineToken                                    string `json:"machineToken"`
	ApplicationBypassInfo                           string `json:"applicationBypassInfo"`
	TunnelTwoForAndroidDevices                      string `json:"tunnelTwoForAndroidDevices"`
	TunnelTwoForiOSDevices                          string `json:"tunnelTwoForiOSDevices"`
	OwnershipVariablePosture                        string `json:"ownershipVariablePosture"`
	BlockUnreachableDomainsTrafficFlag              string `json:"blockUnreachableDomainsTrafficFlag"`
	PrioritizeIPv4OverIpv6                          string `json:"prioritizeIPv4OverIpv6"`
	CrowdStrikeZTAScoreVisibility                   string `json:"crowdStrikeZTAScoreVisibility"`
	NotificationForZPAReauthVisibility              string `json:"notificationForZPAReauthVisibility"`
	CrlCheckVisibilityFlag                          string `json:"crlCheckVisibilityFlag"`
	DedicatedProxyPortsVisibility                   string `json:"dedicatedProxyPortsVisibility"`
	RemoteFetchLogs                                 string `json:"remoteFetchLogs"`
	MsDefenderPostureVisibility                     string `json:"msDefenderPostureVisibility"`
	ExitPasswordVisibility                          string `json:"exitPasswordVisibility"`
	CollectZdxLocationVisibility                    string `json:"collectZdxLocationVisibility"`
	UseV8JsEngineVisibility                         string `json:"useV8JsEngineVisibility"`
	ZdxDisablePasswordVisibility                    string `json:"zdxDisablePasswordVisibility"`
	ZadDisablePasswordVisibility                    string `json:"zadDisablePasswordVisibility"`
	ZpaDisablePasswordVisibility                    string `json:"zpaDisablePasswordVisibility"`
	DefaultProtocolForZPA                           string `json:"defaultProtocolForZPA"`
	DropIpv6TrafficVisibility                       string `json:"dropIpv6TrafficVisibility"`
	MacCacheSystemProxyVisibility                   string `json:"macCacheSystemProxyVisibility"`
	UseWsaPollForZpa                                string `json:"useWsaPollForZpa"`
	Enable64BitFeature                              string `json:"enable64BitFeature"`
	AntivirusPostureVisibility                      string `json:"antivirusPostureVisibility"`
	SystemProxyOnAnyNetworkChangeVisibility         string `json:"systemProxyOnAnyNetworkChangeVisibility"`
	DevicePostureOsVersionVisibility                string `json:"devicePostureOsVersionVisibility"`
	SccmConfigVisibility                            string `json:"sccmConfigVisibility"`
	BrowserAuthFlagVisibility                       string `json:"browserAuthFlagVisibility"`
	InstallWebView2FlagVisibility                   string `json:"installWebView2FlagVisibility"`
	AllowWebView2ToFollowSPVisibility               string `json:"allowWebView2ToFollowSPVisibility"`
	EnableIpv6ResolutionForZscalerDomainsVisibility string `json:"enableIpv6ResolutionForZscalerDomainsVisibility"`
	DisableReasonVisibility                         string `json:"disableReasonVisibility"`
	FollowRoutingTableVisibility                    string `json:"followRoutingTableVisibility"`
	ZiaDevicePostureVisibility                      string `json:"ziaDevicePostureVisibility"`
	UseCustomDNS                                    string `json:"useCustomDNS"`
	UseDefaultAdapterForDNSVisibility               string `json:"useDefaultAdapterForDNSVisibility"`
	T2FallbackBlockAllTrafficAndTlsFallback         string `json:"t2FallbackBlockAllTrafficAndTlsFallback"`
	OverrideT2ProtocolSetting                       string `json:"overrideT2ProtocolSetting"`
	GrantAccessToZscalerLogFolderVisibility         string `json:"grantAccessToZscalerLogFolderVisibility"`
	AdminManagementVisibility                       string `json:"adminManagementVisibility"`
	RedirectWebTrafficToZccListeningProxyVisibility string `json:"redirectWebTrafficToZccListeningProxyVisibility"`
	UseZtunnel2_0ForProxiedWebTrafficVisibility     string `json:"useZtunnel2_0ForProxiedWebTrafficVisibility"`
	SplitVpnVisibility                              string `json:"splitVpnVisibility"`
	EvaluateTrustedNetworkVisibility                string `json:"evaluateTrustedNetworkVisibility"`
	VpnAdaptersConfigurationVisibility              string `json:"vpnAdaptersConfigurationVisibility"`
	VpnServicesVisibility                           string `json:"vpnServicesVisibility"`
	SkipTrustedCriteriaMatchVisibility              string `json:"skipTrustedCriteriaMatchVisibility"`
	ExternalDeviceIdVisibility                      string `json:"externalDeviceIdVisibility"`
	FlowLoggerLoopbackTypeVisibility                string `json:"flowLoggerLoopbackTypeVisibility"`
	FlowLoggerZPATypeVisibility                     string `json:"flowLoggerZPATypeVisibility"`
	FlowLoggerVPNTypeVisibility                     string `json:"flowLoggerVPNTypeVisibility"`
	FlowLoggerVPNTunnelTypeVisibility               string `json:"flowLoggerVPNTunnelTypeVisibility"`
	FlowLoggerDirectTypeVisibility                  string `json:"flowLoggerDirectTypeVisibility"`
	UseZscalerNotificationFramework                 string `json:"useZscalerNotificationFramework"`
	FallbackToGatewayDomain                         string `json:"fallbackToGatewayDomain"`
	ZccRevertVisibility                             string `json:"zccRevertVisibility"`
	ForceZccRevertVisibility                        string `json:"forceZccRevertVisibility"`
	DisasterRecoveryVisibility                      string `json:"disasterRecoveryVisibility"`
	DeviceGroupVisibility                           string `json:"deviceGroupVisibility"`
	IpV6SupportForTunnel2                           string `json:"ipV6SupportForTunnel2"`
	PathMtuDiscovery                                string `json:"pathMtuDiscovery"`
	PostureDiscEncryptionVisibilityForLinux         string `json:"postureDiscEncryptionVisibilityForLinux"`
	PostureMsDefenderVisibilityForLinux             string `json:"postureMsDefenderVisibilityForLinux"`
	PostureOsVersionVisibilityForLinux              string `json:"postureOsVersionVisibilityForLinux"`
	PostureCrowdStrikeZTAScoreVisibilityForLinux    string `json:"postureCrowdStrikeZTAScoreVisibilityForLinux"`
	FlowLoggerZCCBlockedTrafficVisibility           string `json:"flowLoggerZCCBlockedTrafficVisibility"`
	FlowLoggerIntranetTrafficVisibility             string `json:"flowLoggerIntranetTrafficVisibility"`
	CustomMTUForZpaVisibility                       string `json:"customMTUForZpaVisibility"`
	ZpaAutoReauthTimeoutVisibility                  string `json:"zpaAutoReauthTimeoutVisibility"`
	ForceZpaAuthExpireVisibility                    string `json:"forceZpaAuthExpireVisibility"`
	EnableSetProxyOnVPNAdaptersVisibility           string `json:"enableSetProxyOnVPNAdaptersVisibility"`
	DnsServerRouteExclusionVisibility               string `json:"dnsServerRouteExclusionVisibility"`
	EnableSeparateOtpForDevice                      string `json:"enableSeparateOtpForDevice"`
	UninstallPasswordForProfileVisibility           string `json:"uninstallPasswordForProfileVisibility"`
	ZpaAdvanceReauthVisibility                      string `json:"zpaAdvanceReauthVisibility"`
	LatencyBasedZenEnablementVisibility             string `json:"latencyBasedZenEnablementVisibility"`
	DynamicZPAServiceEdgeAssignmenttVisibility      string `json:"dynamicZPAServiceEdgeAssignmenttVisibility"`
	CustomProxyPortsVisibility                      string `json:"customProxyPortsVisibility"`
	DomainInclusionExclusionForDNSRequestVisibility string `json:"domainInclusionExclusionForDNSRequestVisibility"`
	AppNotificationConfigVisibility                 string `json:"appNotificationConfigVisibility"`
	EnableAntiTamperingVisibility                   string `json:"enableAntiTamperingVisibility"`
	StrictEnforcementStatusVisibility               string `json:"strictEnforcementStatusVisibility"`
	AntiTamperingOtpSupportVisibility               string `json:"antiTamperingOtpSupportVisibility"`
	OverrideATCmdByPolicyVisibility                 string `json:"overrideATCmdByPolicyVisibility"`
	DeviceTrustLevelVisibility                      string `json:"deviceTrustLevelVisibility"`
	SourcePortBasedBypassesVisibility               string `json:"sourcePortBasedBypassesVisibility"`
	ProcessBasedApplicationBypassVisibility         string `json:"processBasedApplicationBypassVisibility"`
	CustomBasedApplicationBypassVisibility          string `json:"customBasedApplicationBypassVisibility"`
	ClientCertificateTemplateVisibility             string `json:"clientCertificateTemplateVisibility"`
	SupportedZccVersionChartVisibility              string `json:"supportedZccVersionChartVisibility"`
	IosIpv6ModeVisibility                           string `json:"iosIpv6ModeVisibility"`
	DeviceGroupMultiplePosturesVisibility           string `json:"deviceGroupMultiplePosturesVisibility"`
	DropNonZscalerPacketsVisibility                 string `json:"dropNonZscalerPacketsVisibility"`
	ZccSyntheticIPRangeVisibility                   string `json:"zccSyntheticIPRangeVisibility"`
	DevicePostureFrequencyVisibility                string `json:"devicePostureFrequencyVisibility"`
	EnforceSplitDNSVisibility                       string `json:"enforceSplitDNSVisibility"`
	DataProtectionVisibility                        string `json:"dataProtectionVisibility"`
	DropQuicTrafficVisibility                       string `json:"dropQuicTrafficVisibility"`
	TruncateLargeUDPDNSResponseVisibility           string `json:"truncateLargeUDPDNSResponseVisibility"`
	PrioritizeDnsExclusionsVisibility               string `json:"prioritizeDnsExclusionsVisibility"`
	FetchLogConfigurationOptionVisibility           string `json:"fetchLogConfigurationOptionVisibility"`
	EnableSerialNumberVisibility                    string `json:"enableSerialNumberVisibility"`
	SupportMultiplePWLPostures                      string `json:"supportMultiplePWLPostures"`
	RestrictRemotePacketCaptureVisibility           string `json:"restrictRemotePacketCaptureVisibility"`
	EnableApplicationBasedBypassForMacVisibility    string `json:"enableApplicationBasedBypassForMacVisibility"`
	RemoveExemptedContainersVisibility              string `json:"removeExemptedContainersVisibility"`
	CaptivePortalDetectionVisibility                string `json:"captivePortalDetectionVisibility"`
	DeviceGroupInProfileVisibility                  string `json:"deviceGroupInProfileVisibility"`
	UpdateDnsSearchOrder                            string `json:"updateDnsSearchOrder"`
	InstallActivityBasedMonitoringDriverVisibility  string `json:"installActivityBasedMonitoringDriverVisibility"`
	SlowRolloutZCC                                  string `json:"slowRolloutZCC"`
	ZccTunnelVersionVisibility                      string `json:"zccTunnelVersionVisibility"`
	AntiTamperingStatusVisibility                   string `json:"antiTamperingStatusVisibility"`
	LbbThresholdRankToPercentMapping                string `json:"lbbThresholdRankToPercentMapping"`
	RemoveZscalerSslCertUrl                         string `json:"removeZscalerSslCertUrl"`
	LbzThresholdRankToPercentMapping                string `json:"lbzThresholdRankToPercentMapping"`
	SplashScreenUrl                                 string `json:"splashScreenUrl"`
	SplashScreenVisibility                          string `json:"splashScreenVisibility"`
	TrustedNetworkRangeCriteriaVisibility           string `json:"trustedNetworkRangeCriteriaVisibility"`
	TrustedEgressIpsVisibility                      string `json:"trustedEgressIpsVisibility"`
	DomainProfileDetectionVisibility                string `json:"domainProfileDetectionVisibility"`
	AllInboundTrafficVisibility                     string `json:"allInboundTrafficVisibility"`
	ExportLogsForNonAdminVisibility                 string `json:"exportLogsForNonAdminVisibility"`
	EnableAutoLogSnippetVisibility                  string `json:"enableAutoLogSnippetVisibility"`
	EnableCliVisibility                             string `json:"enableCliVisibility"`
	ZccUserTypeVisibility                           string `json:"zccUserTypeVisibility"`
	InstallWindowsFirewallInboundRule               string `json:"installWindowsFirewallInboundRule"`
	RetryAfterInSeconds                             string `json:"retryAfterInSeconds"`
	AzureADPostureVisibility                        string `json:"azureADPostureVisibility"`
	ServerCertPostureVisibility                     string `json:"serverCertPostureVisibility"`
	PerformCRLCheckServerPostureVisibility          string `json:"performCRLCheckServerPostureVisibility"`
	AutoFillUsingLoginHintVisibility                string `json:"autoFillUsingLoginHintVisibility"`
	SendDefaultPolicyForInvalidPolicyToken          string `json:"sendDefaultPolicyForInvalidPolicyToken"`
	EnableZccPasswordSettings                       string `json:"enableZccPasswordSettings"`
	CliPasswordExpiryMinutes                        string `json:"cliPasswordExpiryMinutes"`
	SsoUsingWindowsPrimaryAccount                   string `json:"ssoUsingWindowsPrimaryAccount"`
	EnableVerboseLog                                string `json:"enableVerboseLog"`
	ZpaAuthExpOnWinLogonSession                     string `json:"zpaAuthExpOnWinLogonSession"`
	ZpaAuthExpOnWinSessionLockVisibility            string `json:"zpaAuthExpOnWinSessionLockVisibility"`
	EnableZccSlowRolloutByDefault                   string `json:"enableZccSlowRolloutByDefault"`
	PurgeKerberosPreferredDCCacheVisibility         string `json:"purgeKerberosPreferredDCCacheVisibility"`
	PostureJamfDetectionVisibility                  string `json:"postureJamfDetectionVisibility"`
	PostureJamfDeviceRiskVisibility                 string `json:"postureJamfDeviceRiskVisibility"`
	WindowsAPCaptivePortalDetectionVisibility       string `json:"windowsAPCaptivePortalDetectionVisibility"`
	WindowsAPEnableFailOpenVisibility               string `json:"windowsAPEnableFailOpenVisibility"`
	AutomaticCaptureDuration                        string `json:"automaticCaptureDuration"`
	ForceLocationRefreshSccm                        string `json:"forceLocationRefreshSccm"`
	EnablePostureFailureDashboard                   string `json:"enablePostureFailureDashboard"`
	EnableOneIDPhase2Changes                        string `json:"enableOneIDPhase2Changes"`
	DropIpv6TrafficInIpv6NetworkVisibility          string `json:"dropIpv6TrafficInIpv6NetworkVisibility"`
	EnablePosturesForPartner                        string `json:"enablePosturesForPartner"`
	EnablePartnerConfigInPrimaryPolicy              string `json:"enablePartnerConfigInPrimaryPolicy"`
	EnableOneIDAdminMigrationChanges                string `json:"enableOneIDAdminMigrationChanges"`
	DdilConfigVisibility                            string `json:"ddilConfigVisibility"`
	AddZDXServiceEntitlement                        string `json:"addZDXServiceEntitlement"`
	UseZcdn                                         string `json:"useZcdn"`
	DeleteDHCPOption121RoutesVisibility             string `json:"deleteDHCPOption121RoutesVisibility"`
	ZdxRolloutControlVisibility                     string `json:"zdxRolloutControlVisibility"`
	ShowM365ServicesInAppBypasses                   string `json:"showM365ServicesInAppBypasses"`
	AllowWebView2IgnoreClientCertErrors             string `json:"allowWebView2IgnoreClientCertErrors"`
	LinuxRPMBuildVisibility                         string `json:"linuxRPMBuildVisibility"`
	HelpBannerDataVisibility                        string `json:"helpBannerDataVisibility"`
	ZpaOnlyDeviceCleanupVisibility                  string `json:"zpaOnlyDeviceCleanupVisibility"`
	AppProfileFailOpenPolicyVisibility              string `json:"appProfileFailOpenPolicyVisibility"`
	ShowRegistryOptionInEnforceAndNone              string `json:"showRegistryOptionInEnforceAndNone"`
	StrictEnforcementNotificationVisibility         string `json:"strictEnforcementNotificationVisibility"`
	CrowdStrikeZTAOsScoreVisibility                 string `json:"crowdStrikeZTAOsScoreVisibility"`
	CrowdStrikeZTASensorConfigScoreVisibility       string `json:"crowdStrikeZTASensorConfigScoreVisibility"`
	ResizeWindowToFitToPageVisibility               string `json:"resizeWindowToFitToPageVisibility"`
	EnableZCCFailCloseSettingsForSEMode             string `json:"enableZCCFailCloseSettingsForSEMode"`
}

// DevicePostureFrequency represents an entry in the devicePostureFrequency
// list returned by the company info endpoint.
type DevicePostureFrequency struct {
	PostureID    int    `json:"postureId"`
	PostureName  string `json:"postureName,omitempty"`
	IosValue     int    `json:"iosValue"`
	AndroidValue int    `json:"androidValue"`
	WindowsValue int    `json:"windowsValue"`
	MacValue     int    `json:"macValue"`
	LinuxValue   int    `json:"linuxValue"`
	DefaultValue int    `json:"defaultValue"`
}

// GetCompanyInfo retrieves the full company configuration. The response is
// decoded into a CompanyInfo struct so callers (e.g. the Terraform provider)
// receive every field the API returns.
func GetCompanyInfo(ctx context.Context, service *zscaler.Service) (*CompanyInfo, error) {
	var info CompanyInfo
	_, err := service.Client.NewZccRequestDo(ctx, "GET", getCompanyInfoEndpoint, nil, nil, &info)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve company info: %w", err)
	}
	return &info, nil
}

// SetCompanyInfo updates the company configuration via PUT. The API returns
// only a {"success":"true","error":"0"} envelope, so on success we re-fetch
// the full configuration via GetCompanyInfo and return that, ensuring
// callers (Terraform state) always see the authoritative server state.
func SetCompanyInfo(ctx context.Context, service *zscaler.Service, companyInfo *CompanyInfo) (*CompanyInfo, error) {
	if companyInfo == nil {
		return nil, errors.New("companyInfo is required")
	}

	var resp common.ZCCResponse
	_, err := service.Client.NewZccRequestDo(ctx, "PUT", setCompanyInfoEndpoint, nil, companyInfo, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to update company info: %w", err)
	}

	if resp.Success != "true" {
		return nil, fmt.Errorf("API rejected company info update (error: %s)", resp.Error)
	}

	service.Client.GetLogger().Printf("[DEBUG] company info update success, re-reading via GET")

	updated, err := GetCompanyInfo(ctx, service)
	if err != nil {
		return nil, fmt.Errorf("company info updated, but failed to re-read: %w", err)
	}
	return updated, nil
}
