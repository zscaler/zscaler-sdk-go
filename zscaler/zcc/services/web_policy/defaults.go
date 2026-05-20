package web_policy

import (
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

// DefaultMacosWebPolicy returns a WebPolicy pre-populated with the exact
// scalar / list / picker defaults the ZCC UI generates when an
// administrator clicks "Save" on a fresh macOS app profile. Callers
// (typically the Terraform provider) overlay their user-supplied values
// on top — anything the operator does not configure stays at the
// server-known default, and the resulting body matches the byte-for-byte
// shape of a known-working /web/policy/edit request.
//
// The companion docs/local_dev/zcc_app_profile_macos/test.json file is
// the source of truth for these defaults; if the API requirements
// change, capture a fresh payload from the UI (browser dev tools →
// Network → "edit" PUT body) and update this constructor.
func DefaultMacosWebPolicy() WebPolicy {
	return WebPolicy{
		// Core identity / lifecycle. Name/Description/Active/RuleOrder
		// are caller-supplied; the rest of the body matches test.json.
		ID:          "",
		Name:        "",
		Description: "",
		Active:      "1",
		DeviceType:  common.DeviceTypeMacOS,
		RuleOrder:   1,

		// Targeting — every collection ships as an empty array, matching
		// the UI's "no targeting yet" state.
		Groups:                                make([]any, 0),
		Users:                                 make([]any, 0),
		GroupAll:                              0,
		GroupIds:                              make([]int, 0),
		GroupNames:                            make([]string, 0),
		UserIds:                               make([]int, 0),
		UserNames:                             make([]string, 0),
		AppIdentityNames:                      make([]string, 0),
		AppServiceIds:                         make([]int, 0),
		AppServiceNames:                       make([]string, 0),
		AppServiceCustomIdsSelected:           make([]any, 0),
		BypassAppIds:                          make([]int, 0),
		BypassCustomAppIds:                    make([]int, 0),
		BypassMacAppIds:                       make([]any, 0),
		DeviceGroupIds:                        make([]int, 0),
		DeviceGroupNames:                      make([]string, 0),
		DeviceGroups:                          make([]any, 0),
		DeviceGroupsOption:                    0,
		DeviceGroupsSelected:                  make([]any, 0),
		UsersOption:                           0,
		UsersSelected:                         make([]any, 0),
		ZccFailCloseSettingsAppByPassIdsTop:   make([]int, 0),
		ZccFailCloseSettingsAppByPassSelected: make([]any, 0),

		// Forwarding / posture
		ForwardingProfileId: 0,
		ZiaPostureProfile:   make([]any, 0),
		ZiaPostureConfigId:  0,

		// Logging defaults (logFileSize=100MB, debug log_mode=3 via the
		// picker, raw logMode left at -1 which means "follow picker")
		LogMode:     -1,
		LogLevel:    0,
		LogFileSize: 100,
		LogModeSelected: &LabelValuePair{
			Label: "Debug",
			Value: 3,
		},

		// Captive portal + diagnostics
		EnableCaptivePortalDetection:      1,
		EnableFailOpen:                    1,
		CaptivePortalWebSecDisableMinutes: 10,
		CaptivePortalUrlId: []LabelValuePair{
			{Label: "Zscaler", Value: 1},
		},
		EndToEndDiagnostics:          EndToEndDiagnostics{},
		EndToEndDiagnosticsSelected:  make([]any, 0),
		LocalMetrics:                 1,
		FlowLoggingSelected:          make([]any, 0),
		BlockDomainSelected:          make([]any, 0),
		BlockInboundTrafficSelected:  make([]any, 0),
		NotificationTemplateSelected: make([]any, 0),

		// PAC
		PacURL:      "",
		PacType:     1,
		PacDataPath: "",

		// MDM / mobile / billing — present on every device_type body
		Mdm:               0,
		Passcode:          "",
		ExitPassword:      "",
		Limit:             "1",
		BillingDay:        "1",
		AllowedApps:       "",
		CustomText:        "",
		BypassMmsApps:     0,
		QuotaInRoaming:    0,
		WifiSSID:          "",
		BypassAndroidApps: make([]int, 0),
		Enforced:          0,

		// Registry / Windows-shape fields the API echoes regardless of OS
		RegistryPath:                      "",
		RegistryName:                      "",
		InstallSslCertsTop:                1,
		DisableLoopBackRestriction:        0,
		RemoveExemptedContainers:          1,
		OverrideWPAD:                      0,
		RestartWinHttpSvc:                 0,
		InstallWindowsFirewallInboundRule: "1",
		ForceLocationRefreshSccm:          0,
		WfpMtr:                            0,
		EnableLocalPacketCaptureTabValue:  0,
		RefreshKerberosToken:              0,

		// Nullable nested configs — the UI sends explicit null
		FlowLoggerConfig:             nil,
		DomainProfileDetectionConfig: nil,
		AllInboundTrafficConfig:      nil,

		// Cosmetic / runtime knobs at the top level
		HighlightActiveControl:       0,
		SendDisableServiceReason:     0,
		TunnelZappTraffic:            0,
		EnableDeviceGroups:           0,
		ReactivateWebSecurityMins:    0,
		ReauthPeriod:                 8,
		ClearArpCacheTop:             0,
		EnableZscalerFirewallTop:     "1",
		PersistentZscalerFirewallTop: 0,
		CacheSystemProxyTop:          1,
		DnsPriorityOrderingTop: []string{
			"State:/Network/Service/com.cisco.anyconnect/DNS",
		},
		EnableZdpServiceTop:        1,
		DisableParallelIpv4AndIPv6: -1,
		DisableParallelIpv4andIpv6: "-1",

		// Top-level UI form-state pickers. RuleOrderSelectedOption mirrors
		// the UI's default of `2/2` — the picker reflects the slot the UI
		// places a fresh rule into, which is distinct from `ruleOrder` at
		// the top level. Both keys travel on the wire (see test.json).
		RuleOrderSelectedOption:  &LabelValuePair{Label: "2", Value: 2},
		BillingDaySelectedOption: &LabelValuePair{Label: "1", Value: "1"},
		Ipv6ModeSelected:         &LabelValuePair{Label: "IPv6Native", Value: 4},
		ZpaAutoReauthTimeoutTop: []LabelValuePair{
			{Label: "30", Value: 30},
		},
		PcAdditionalSpaceTop: []LabelValuePair{
			{Label: "1GB", Value: "1024"},
		},
		BrowserAuthTypeTop: &LabelValuePair{
			Label: "FOLLOW_GLOBAL_CONFIG",
			Value: -1,
		},
		ClientConnectorUiLanguageSelected: []LabelValuePair{
			{Label: "Use System Language", Value: 0},
		},

		// Machine token / ZPA reauth schedule
		MachineTokenOption:                           0,
		MachineTokenSelectedOption:                   0,
		ZpaAuthExpSessionLockStateMinTimeInSecondTop: "1",
		ForceZpaAuthenticationToExpire:               make([]any, 0),
		ZpaReauthConfigTop:                           make([]any, 0),

		// DR mirror (top-level picker)
		ZiaDRMethodTop: &LabelValuePair{
			Label: "Policy Based Access (Web only)",
			Value: 2,
		},

		// Top-level disable-without-password trio (defaults false)
		AllowZpaDisableWithoutPasswordTop: false,
		AllowZiaDisableWithoutPasswordTop: false,
		AllowZdxDisableWithoutPasswordTop: false,

		// Top-level DNS / split-tunnel flags
		UseDefaultAdapterForDNSTop:     "1",
		UpdateDnsSearchOrderTop:        "1",
		EnforceSplitDNSTop:             "0",
		DisableDNSRouteExclusionTop:    "0",
		EnableSetProxyOnVPNAdaptersTop: 1,
		DropQuicTrafficTop:             "0",
		FollowRoutingTableTop:          "1",

		// Top-level partner / fail-close / packet-tunnel mirrors
		VpnGatewaysTop:                                     make([]any, 0),
		PartnerDomainsTop:                                  make([]any, 0),
		ZccFailCloseSettingsIpBypassesTop:                  make([]any, 0),
		ZccFailCloseSettingsLockdownOnTunnelProcessExitTop: 1,
		ZccFailCloseSettingsExitUninstallPasswordTop:       "",
		UserAllowedToAddPartnerTop:                         1,
		FollowGlobalForPartnerLoginTop:                     "1",
		FollowGlobalForZpaReauthTop:                        "1",
		FollowGlobalForPacketCaptureTop:                    "1",
		EnableLocalPacketCaptureTop:                        "0",
		EnableLocalPacketCaptureV2Top:                      make([]any, 0),
		PacketTunnelIncludeListTop: []string{
			"0.0.0.0/0",
		},
		PacketTunnelExcludeListTop: []string{
			"10.0.0.0/8",
			"172.16.0.0/12",
			"192.168.0.0/16",
			"224.0.0.0/4",
			"255.255.255.255",
			"169.254.0.0/16",
		},
		PacketTunnelIncludeListForIPv6Top: make([]string, 0),
		PacketTunnelExcludeListForIPv6Top: []string{
			"[FF00::/8]",
			"[FE80::/10]",
			"[FC00::/7]",
		},
		PacketTunnelDnsIncludeListTop: make([]string, 0),
		PacketTunnelDnsExcludeListTop: make([]string, 0),
		SourcePortBasedBypassesTop: []string{
			"3389:*",
		},
		UseV8JsEngineTop:           "1",
		PrioritizeDnsExclusionsTop: "1",

		// Trusted-network buckets, empty by default
		VpnTrusted:      make([]any, 0),
		SplitVpnTrusted: make([]any, 0),
		Trusted:         make([]any, 0),
		OffTrusted:      make([]any, 0),
		CustomDNSTop:    make([]any, 0),

		// Top-level UX / diagnostics flags
		EnableZCCRevertTop:                    false,
		EnableCustomProxyDetectionTop:         "0",
		ClientConnectorUiLanguageTop:          0,
		OneIdMTDeviceAuthEnabledTop:           "0",
		PreventAutoReauthDuringDeviceLockTop:  "0",
		InstantForceZPAReauthStateUpdateTop:   0,
		EnableNetworkTrafficProcessMappingTop: 0,
		UseEndPointLocationForDCSelectionTop:  "0",
		RecacheSystemProxyTop:                 "0",
		EnableLocationPolicyOverrideTop:       0,
		BlockPrivateRelayTop:                  "0",
		EnableCrashReportingTop:               "0",
		EnableAutomaticPacketCaptureTop:       "0",
		EnableAPCforCriticalSectionsTop:       "1",
		EnableAPCforOtherSectionsTop:          "1",
		EnablePCAdditionalSpaceTop:            "1",

		ReactivateAntiTamperingTimeTop: 0,

		// Top-level useDefaultBrowser; the picker object sits in BrowserAuthTypeTop
		UseDefaultBrowserTop: 0,

		// Nested blocks
		MacPolicy:        nil, // caller overrides
		PolicyExtension:  defaultMacosPolicyExtension(),
		DisasterRecovery: defaultMacosDisasterRecovery(),
	}
}

// defaultMacosPolicyExtension returns the policyExtension nested block as
// the macOS UI emits it on a fresh save — matches the policyExtension
// portion of test.json verbatim.
func defaultMacosPolicyExtension() PolicyExtension {
	return PolicyExtension{
		GenerateCliPasswordContract: GenerateCliPasswordContract{
			EnableCli:                      false,
			AllowZpaDisableWithoutPassword: true,
			AllowZiaDisableWithoutPassword: true,
			AllowZdxDisableWithoutPassword: true,
		},

		VpnGateways:                                     "",
		PartnerDomains:                                  "",
		ZccFailCloseSettingsIpBypasses:                  "",
		ZccFailCloseSettingsLockdownOnTunnelProcessExit: "1",
		ZccFailCloseSettingsExitUninstallPassword:       "",
		ZccFailCloseSettingsAppByPassIds:                make([]int, 0),
		ZccFailCloseSettingsAppByPassNames:              make([]string, 0),
		ZccFailCloseSettingsThumbPrint:                  "",
		ZccFailCloseSettingsLockdownOnFirewallError:     "0",
		ZccFailCloseSettingsLockdownOnDriverError:       "0",
		UserAllowedToAddPartner:                         "1",
		FollowGlobalForPartnerLogin:                     "1",
		FollowGlobalForZpaReauth:                        "1",
		FollowGlobalForPacketCapture:                    "1",
		EnableLocalPacketCapture:                        "0",
		EnableLocalPacketCaptureV2:                      0,
		EnableFlowBasedTunnel:                           "0",

		ZpaReauthConfig:                   nil,
		ZpaAutoReauthTimeout:              common.IntOrString(30),
		ZpaAuthExpOnSleep:                 common.IntOrString(0),
		ZpaAuthExpOnSysRestart:            common.IntOrString(0),
		ZpaAuthExpOnNetIpChange:           common.IntOrString(0),
		InstantForceZPAReauthStateUpdate:  common.IntOrString(0),
		ZpaAuthExpOnWinLogonSession:       common.IntOrString(0),
		ZpaAuthExpOnWinSessionLock:        common.IntOrString(0),
		ZpaAuthExpSessionLockStateMinTime: "0",
		AdvanceZpaReauth:                  false,

		ExitPassword:                    "",
		FollowRoutingTable:              "1",
		UseDefaultAdapterForDNS:         "1",
		UpdateDnsSearchOrder:            "1",
		UseZscalerNotificationFramework: "0",
		SwitchFocusToNotification:       "0",
		FallbackToGatewayDomain:         "1",
		UseProxyPortForT1:               "0",
		UseProxyPortForT2:               "0",
		AllowPacExclusionsOnly:          "0",
		UseWsaPollForZpa:                "0",
		EnableZCCRevert:                 "0",
		ZccRevertPassword:               "",
		EnableSetProxyOnVPNAdapters:     "1",
		DisableDNSRouteExclusion:        common.IntOrString(0),
		PacketTunnelIncludeListForIPv6:  "",
		InterceptZIATrafficAllAdapters:  common.IntOrString(0),
		EnableAntiTampering:             common.IntOrString(0),
		ReactivateAntiTamperingTime:     0,
		SourcePortBasedBypasses:         "3389:*",
		EnforceSplitDNS:                 common.IntOrString(0),
		DropQuicTraffic:                 common.IntOrString(0),
		ZdpDisablePassword:              "",
		UseV8JsEngine:                   "1",
		ZdDisablePassword:               "",
		ZdxDisablePassword:              "",
		ZpaDisablePassword:              "",
		BypassDNSTrafficUsingUDPProxy:   "0",
		ReconnectTunOnWakeup:            "1",
		EnableCustomTheme:               0,
		DeleteDHCPOption121Routes:       `{"trusted":1,"offTrusted":1,"vpnTrusted":1,"splitVpnTrusted":1}`,
		MachineIdpAuth:                  false,
		Nonce:                           "",
		PacketTunnelDnsExcludeList:      "",
		PacketTunnelDnsIncludeList:      "",
		PacketTunnelExcludeList:         "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,224.0.0.0/4,255.255.255.255,169.254.0.0/16",
		PacketTunnelExcludeListForIPv6:  "[FF00::/8],[FE80::/10],[FC00::/7]",
		PacketTunnelIncludeList:         "0.0.0.0/0",
		TruncateLargeUDPDNSResponse:     common.IntOrString(0),
		OverrideATCmdByPolicy:           common.IntOrString(0),
		PurgeKerberosPreferredDCCache:   common.IntOrString(0),
		RscModeOnAllAdapters:            common.IntOrString(0),
		EnableAdapterHardwareOffloading: common.IntOrString(0),
		SupportZPASearchDomainsInTRP:    common.IntOrString(0),
		PrioritizeDnsExclusions:         common.IntOrString(1),

		LocationRulesetPolicies: LocationRulesetPolicies{
			SplitVpnTrusted: LocationRulesetEntry{ID: 0},
			VpnTrusted:      LocationRulesetEntry{ID: 0},
		},

		DdilConfig:           `{"ddilEnabled":0,"businessContinuityActivationDomain":"","businessContinuityTestModeEnabled":0}`,
		ZccAppFailOpenPolicy: common.IntOrString(0),
		ZccTunnelFailPolicy:  common.IntOrString(0),

		AllowClientCertCachingForWebView2:   "0",
		ShowConfirmationDialogForCachedCert: "0",
		OneIdMTDeviceAuthEnabled:            "0",
		PreventAutoReauthDuringDeviceLock:   "0",
		ClientConnectorUiLanguage:           common.IntOrString(0),
		EnableNetworkTrafficProcessMapping:  common.IntOrString(0),
		UseEndPointLocationForDCSelection:   "0",
		RecacheSystemProxy:                  "0",
		EnableLocationPolicyOverride:        common.IntOrString(0),
		BlockPrivateRelay:                   "0",
		EnableAutomaticPacketCapture:        "0",
		EnableAPCforCriticalSections:        "1",
		EnableAPCforOtherSections:           "1",
		EnablePCAdditionalSpace:             "1",
		PcAdditionalSpace:                   "512",
		EnableCustomProxyDetection:          "0",
		EnableCrashReporting:                "0",
		EnableZdpService:                    common.IntOrString(1),

		ZdxLiteConfigObj: `{"localMetrics":1,"endToEndDiagnostics":{"trusted":0,"vpnTrusted":0,"offTrusted":0,"splitVpnTrusted":0}}`,
	}
}

// defaultMacosDisasterRecovery returns the disasterRecovery nested block
// in its UI default shape: DR is off, but the picker is preset to
// "Policy Based Access (Web only)" (value 2) and useZiaGlobalDb is true.
func defaultMacosDisasterRecovery() DisasterRecovery {
	return DisasterRecovery{
		AllowZiaTest:     false,
		AllowZpaTest:     false,
		EnableZiaDR:      false,
		EnableZpaDR:      false,
		ZiaDRMethod:      2,
		ZiaCustomDbUrl:   "",
		UseZiaGlobalDb:   true,
		ZiaDomainName:    "",
		ZiaRSAPubKeyName: "",
		ZiaRSAPubKey:     "",
		ZpaDomainName:    "",
		ZpaRSAPubKeyName: "",
		ZpaRSAPubKey:     "",
	}
}

// DefaultMacosMacPolicy returns the macPolicy nested block as it appears
// in a fresh UI-generated macOS payload. The caller layers user-set
// values on top.
func DefaultMacosMacPolicy() *MacPolicy {
	return &MacPolicy{
		AddIfscopeRoute:                      "0",
		ClearArpCache:                        "0",
		EnableZscalerFirewall:                "1",
		PersistentZscalerFirewall:            "0",
		DnsPriorityOrderingForTrustedDnsCrit: "0",
		DnsPriorityOrdering:                  "State:/Network/Service/com.cisco.anyconnect/DNS",
		BrowserAuthType:                      -1,
		UseDefaultBrowser:                    0,
		CacheSystemProxy:                     "1",
		DisablePassword:                      "",
		InstallSslCerts:                      1,
		LogoutPassword:                       "",
		UninstallPassword:                    "",
		CaptivePortalConfig:                  `{"automaticCapture":1,"enableCaptivePortalDetection":1,"enableFailOpen":1,"captivePortalWebSecDisableMinutes":10,"enableEmbeddedCaptivePortal":0}`,
		EnableAppBasedBypass:                 "0",
	}
}
