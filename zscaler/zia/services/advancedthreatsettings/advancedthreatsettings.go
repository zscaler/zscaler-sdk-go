package advancedthreatsettings

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	advThreatSettingsEndpoint  = "/zia/api/v1/cyberThreatProtection/advancedThreatSettings"
	maliciousUrlsEndpoint      = "/zia/api/v1/cyberThreatProtection/maliciousUrls"
	securityExceptionsEndpoint = "/zia/api/v1/cyberThreatProtection/securityExceptions"
)

type AdvancedThreatSettings struct {
	// The Page Risk tolerance index set between 0 and 100 (100 being the highest risk).
	// Users are blocked from accessing web pages with higher Page Risk than the specified value.
	RiskTolerance int `json:"riskTolerance"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for suspicious web pages
	RiskToleranceCapture bool `json:"riskToleranceCapture"`

	// A Boolean value specifying whether connections to known Command & Control (C2) Servers are allowed or blocked
	CmdCtlServerBlocked bool `json:"cmdCtlServerBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for connections to known C2 servers
	CmdCtlServerCapture bool `json:"cmdCtlServerCapture"`

	// A Boolean value specifying whether botnets are allowed or blocked from sending or receiving commands to unknown servers
	CmdCtlTrafficBlocked bool `json:"cmdCtlTrafficBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for botnets
	CmdCtlTrafficCapture bool `json:"cmdCtlTrafficCapture"`

	// A Boolean value specifying whether known malicious sites and content are allowed or blocked
	MalwareSitesBlocked bool `json:"malwareSitesBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for malicious sites
	MalwareSitesCapture bool `json:"malwareSitesCapture"`

	// A Boolean value specifying whether sites are allowed or blocked from accessing vulnerable ActiveX controls that are known to have been exploited.
	ActiveXBlocked bool `json:"activeXBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for ActiveX controls
	ActiveXCapture bool `json:"activeXCapture"`

	// A Boolean value specifying whether known web browser vulnerabilities prone to exploitation are allowed or blocked.
	BrowserExploitsBlocked bool `json:"browserExploitsBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for browser exploits
	BrowserExploitsCapture bool `json:"browserExploitsCapture"`

	// A Boolean value specifying whether known file format vulnerabilities and suspicious or malicious content in Microsoft Office or PDF documents are allowed or blocked
	FileFormatVulnerabilitiesBlocked bool `json:"fileFormatVunerabilitesBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for file format vulnerabilities
	FileFormatVulnerabilitiesCapture bool `json:"fileFormatVunerabilitesCapture"`

	// A Boolean value specifying whether known phishing sites are allowed or blocked
	KnownPhishingSitesBlocked bool `json:"knownPhishingSitesBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for known phishing sites
	KnownPhishingSitesCapture bool `json:"knownPhishingSitesCapture"`

	// A Boolean value specifying whether to allow or block suspected phishing sites identified through heuristic detection.
	// The Zscaler service can inspect the content of a website for indications that it might be a phishing site.
	SuspectedPhishingSitesBlocked bool `json:"suspectedPhishingSitesBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for suspected phishing sites
	SuspectedPhishingSitesCapture bool `json:"suspectedPhishingSitesCapture"`

	// A Boolean value specifying whether to allow or block any detections of communication and callback traffic associated with spyware agents and data transmission
	SuspectAdwareSpywareSitesBlocked bool `json:"suspectAdwareSpywareSitesBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for suspected adware and spyware sites
	SuspectAdwareSpywareSitesCapture bool `json:"suspectAdwareSpywareSitesCapture"`

	// Boolean value specifying whether to allow or block web pages that pretend to contain useful information, to get higher ranking in search engine results or drive traffic to phishing, adware, or spyware distribution sites.
	WebspamBlocked bool `json:"webspamBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for web spam
	WebspamCapture bool `json:"webspamCapture"`

	// A Boolean value specifying whether to allow or block IRC traffic being tunneled over HTTP/S
	IrcTunnellingBlocked bool `json:"ircTunnellingBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for IRC tunnels
	IrcTunnellingCapture bool `json:"ircTunnellingCapture"`

	// A Boolean value specifying whether to allow or block applications and methods used to obscure the destination and the content accessed by the user, therefore blocking traffic to anonymizing web proxies.
	AnonymizerBlocked bool `json:"anonymizerBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for anonymizers
	AnonymizerCapture bool `json:"anonymizerCapture"`

	// A Boolean value specifying whether to allow or block third-party websites that gather cookie information, which can be used to personally identify users, track internet activity, or steal a user's session or sensitive information.
	CookieStealingBlocked bool `json:"cookieStealingBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for cookie stealing
	CookieStealingPCAPEnabled bool `json:"cookieStealingPCAPEnabled"`

	// A Boolean value specifying whether to allow or block this type of cross-site scripting (XSS)
	PotentialMaliciousRequestsBlocked bool `json:"potentialMaliciousRequestsBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for (XSS) attacks
	PotentialMaliciousRequestsCapture bool `json:"potentialMaliciousRequestsCapture"`

	// A Boolean value specifying whether to allow or block requests to websites located in specific countries based on the ISO3166 mapping of countries to their IP address space
	BlockedCountries []string `json:"blockedCountries"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for blocked countries
	BlockCountriesCapture bool `json:"blockCountriesCapture"`

	// A Boolean value specifying whether to allow or block the usage of BitTorrent, a popular P2P file sharing application that supports content download with encryption.
	BitTorrentBlocked bool `json:"bitTorrentBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for BitTorrent
	BitTorrentCapture bool `json:"bitTorrentCapture"`

	// A Boolean value specifying whether to allow or block the usage of Tor, a popular P2P anonymizer protocol with support for encryption.
	TorBlocked bool `json:"torBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for Tor
	TorCapture bool `json:"torCapture"`

	// A Boolean value specifying whether to allow or block access to Google Hangouts, a popular P2P VoIP application.
	GoogleTalkBlocked bool `json:"googleTalkBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for Google Hangouts
	GoogleTalkCapture bool `json:"googleTalkCapture"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for SSH tunnels
	SshTunnellingBlocked bool `json:"sshTunnellingBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for SSH tunnels
	SshTunnellingCapture bool `json:"sshTunnellingCapture"`

	// A Boolean value specifying whether to allow or block cryptocurrency mining network traffic and scripts
	// which can negatively impact endpoint device performance and potentially lead to a misuse of company resources.
	CryptoMiningBlocked bool `json:"cryptoMiningBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for cryptomining
	CryptoMiningCapture bool `json:"cryptoMiningCapture"`

	// A Boolean value specifying whether to allow or block websites known to contain adware or
	// spyware that displays malicious advertisements that can collect users' information without their knowledge
	AdSpywareSitesBlocked bool `json:"adSpywareSitesBlocked"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for adware and spyware sites
	AdSpywareSitesCapture bool `json:"adSpywareSitesCapture"`

	// A Boolean value specifying whether to allow or block domains that are suspected to be generated using domain generation algorithms (DGA)
	DgaDomainsBlocked bool `json:"dgaDomainsBlocked"`

	// A Boolean value specifying whether to send alerts upon detecting unknown or suspicious C2 traffic
	AlertForUnknownOrSuspiciousC2Traffic bool `json:"alertForUnknownOrSuspiciousC2Traffic"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for DGA domains
	DgaDomainsCapture bool `json:"dgaDomainsCapture"`

	// A Boolean value specifying whether packet capture (PCAP) is enabled or not for malicious URLs
	MaliciousUrlsCapture bool `json:"maliciousUrlsCapture"`
}

// Denylist URLs for ATP policy
type MaliciousURLs struct {
	// Allowlist URLs that are not inspected by the ATP policy
	MaliciousUrls []string `json:"maliciousUrls"`
}

// Security exceptions for ATP policy
type SecurityExceptions struct {
	// Allowlist URLs that are not inspected by the ATP policy
	BypassUrls []string `json:"bypassUrls"`
}

func GetAdvancedThreatSettings(ctx context.Context, service *zscaler.Service) (*AdvancedThreatSettings, error) {
	var advThreatSettings AdvancedThreatSettings
	err := service.Client.Read(ctx, advThreatSettingsEndpoint, &advThreatSettings)
	if err != nil {
		return nil, err
	}
	return &advThreatSettings, nil
}

func UpdateAdvancedThreatSettings(ctx context.Context, service *zscaler.Service, settings AdvancedThreatSettings) (*AdvancedThreatSettings, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, advThreatSettingsEndpoint, settings)
	if err != nil {
		return nil, nil, err
	}

	advThreatSettings, ok := resp.(*AdvancedThreatSettings)
	if !ok {
		return nil, nil, fmt.Errorf("unexpected response type")
	}
	service.Client.GetLogger().Printf("[DEBUG] Updated Advanced Threat Settings : %+v", advThreatSettings)
	return advThreatSettings, nil, nil
}

func GetMaliciousURLs(ctx context.Context, service *zscaler.Service) (*MaliciousURLs, error) {
	var urls MaliciousURLs
	err := service.Client.Read(ctx, maliciousUrlsEndpoint, &urls)
	if err != nil {
		return nil, err
	}
	service.Client.GetLogger().Printf("[DEBUG] Returning malicious urls from Get: %v", urls)
	return &urls, nil
}

func UpdateMaliciousURLs(ctx context.Context, service *zscaler.Service, urls MaliciousURLs) (*MaliciousURLs, error) {
	currentUrls, err := GetMaliciousURLs(ctx, service)
	if err != nil {
		return nil, err
	}
	newUrls := zscaler.Difference(urls.MaliciousUrls, currentUrls.MaliciousUrls)
	removedUrls := zscaler.Difference(currentUrls.MaliciousUrls, urls.MaliciousUrls)
	if len(newUrls) > 0 {
		_, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s?action=ADD_TO_LIST", maliciousUrlsEndpoint), MaliciousURLs{newUrls})
		if err != nil {
			return nil, err
		}
	}
	if len(removedUrls) > 0 {
		_, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s?action=REMOVE_FROM_LIST", maliciousUrlsEndpoint), MaliciousURLs{removedUrls})
		if err != nil {
			return nil, err
		}
	}
	return &urls, nil
}

func GetSecurityExceptions(ctx context.Context, service *zscaler.Service) (*SecurityExceptions, error) {
	var bypassUrls SecurityExceptions
	err := service.Client.Read(ctx, securityExceptionsEndpoint, &bypassUrls)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning bypass urls from Get: %v", bypassUrls)
	return &bypassUrls, nil
}

func UpdateSecurityExceptions(ctx context.Context, service *zscaler.Service, urls SecurityExceptions) (*SecurityExceptions, error) {
	// Overwrite the bypass URLs with the provided list
	_, err := service.Client.UpdateWithPut(ctx, securityExceptionsEndpoint, urls)
	if err != nil {
		return nil, err
	}

	return &urls, nil
}
