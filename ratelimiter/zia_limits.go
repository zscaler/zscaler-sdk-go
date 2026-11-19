package ratelimiter

import (
	"net/http"
	"time"
)

// ZIA rate limits, from the ZIA API Rate Limit Summary
// (https://help.zscaler.com/legacy-apis/api-rate-limit-summary). The API
// enforces them per method and resource URI. The summary notes that rate
// limits are subject to change.

// ZIAPathPrefix is the prefix of ZIA paths on the OneAPI gateway.
const ZIAPathPrefix = "/zia/api/v1"

// ZIADefaultLimits are the limits of most ZIA resources: GET 2/sec and
// 1000/hr, POST, PUT and DELETE 1/sec and 400/hr.
func ZIADefaultLimits() map[string]EndpointLimit {
	read := EndpointLimit{Rate{2, time.Second}, 1000}
	write := EndpointLimit{Rate{1, time.Second}, 400}
	return map[string]EndpointLimit{
		http.MethodGet:    read,
		http.MethodPost:   write,
		http.MethodPut:    write,
		http.MethodPatch:  write,
		http.MethodDelete: write,
	}
}

// ZIAEndpointRules are the resources whose limits are stricter than
// ZIADefaultLimits. Resources with looser limits keep the defaults, and daily
// quotas are left to the API.
func ZIAEndpointRules() []EndpointRule {
	return []EndpointRule{
		{http.MethodGet, "/alertSubscriptions", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/alertSubscriptions/{alertSubscriptionId}", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodPost, "/apps-services", EndpointLimit{Rate{1, time.Minute}, 4}},
		{http.MethodPut, "/apps-services/{appServiceId}", EndpointLimit{Rate{1, time.Minute}, 4}},
		{http.MethodPost, "/auditlogEntryReport", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodPost, "/azure/refreshConnectionStatus", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodPost, "/azure/tunnelConfiguration", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodDelete, "/azure/unConfigureTunnel/{hubId}", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodPost, "/azure/virtualHubSync", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodGet, "/bandwidthControlRules", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/bandwidthControlRules/lite", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/bandwidthControlRules/{ruleId}", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/casbDlpRules", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/casbDlpRules/all", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/casbDlpRules/{ruleId}", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/casbMalwareRules", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/casbMalwareRules/all", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/casbMalwareRules/{ruleId}", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodPut, "/cloudApplications/bulkUpdate", EndpointLimit{Rate{1, time.Minute}, 4}},
		{http.MethodGet, "/configAudit", EndpointLimit{Rate{1, time.Hour}, 1}},
		{http.MethodGet, "/configAudit/ipVisibility", EndpointLimit{Rate{1, time.Hour}, 1}},
		{http.MethodGet, "/configAudit/pacFile", EndpointLimit{Rate{1, time.Hour}, 1}},
		{http.MethodGet, "/customFileTypes", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/dcExclusions", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/dnatRules", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/dnatRules/{ruleId}", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/dnsGateways", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/dnsGateways/{id}", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/endPointApplicationGroups", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/endPointApplicationGroups/policies", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/endPointApplicationGroups/{id}", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/endPointApplications/cloudApps", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/endPointApplications/cloudApps/count", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/endPointApplications/count", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/endPointApplications/customApps", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/endPointApplications/customApps/{id}", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/endPointApplications/policies", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodPost, "/endPointDlpRules", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodPut, "/endPointDlpRules/{id}", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodPost, "/endPointDlpRules/{id}/subRule", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodPut, "/endPointDlpRules/{id}/subRule/{subRuleId}", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodGet, "/extranet", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodPost, "/extranet", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodGet, "/extranet/lite", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/extranet/{id}", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodPut, "/extranet/{id}", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodDelete, "/extranet/{id}", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodGet, "/fileTypeCategories", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/groups", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodPost, "/intermediateCaCertificate/finalizeCert/{certId}", EndpointLimit{Rate{1, time.Minute}, 4}},
		{http.MethodPost, "/intermediateCaCertificate/generateCsr/{certId}", EndpointLimit{Rate{1, time.Minute}, 4}},
		{http.MethodPost, "/intermediateCaCertificate/KeyPair/{certId}", EndpointLimit{Rate{1, time.Minute}, 4}},
		{http.MethodPost, "/intermediateCaCertificate/uploadcert/{certId}", EndpointLimit{Rate{1, time.Minute}, 4}},
		{http.MethodPost, "/intermediateCaCertificate/uploadcertchain/{certId}", EndpointLimit{Rate{1, time.Minute}, 4}},
		{http.MethodPost, "/intermediateCaCertificate/verifyKeyAttestation/{certId}", EndpointLimit{Rate{1, time.Minute}, 4}},
		{http.MethodDelete, "/intermediateCaCertificate/{certId}", EndpointLimit{Rate{1, time.Minute}, 4}},
		{http.MethodGet, "/iotDiscovery/deviceList", EndpointLimit{Rate{1, time.Minute}, 4}},
		{http.MethodGet, "/isolationVotiroCdr/policies", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/isolationVotiroCdr/tokenConfig", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodPost, "/locations/bulkDelete", EndpointLimit{Rate{1, time.Minute}, 4}},
		{http.MethodGet, "/nssCollectors", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodDelete, "/pacFiles/{pacId}", EndpointLimit{Rate{1, time.Minute}, 4}},
		{http.MethodGet, "/passwordExpiry/settings", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodPost, "/riskProfiles", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodPut, "/riskProfiles/{profileId}", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodGet, "/security/advanced", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodPost, "/shadowIT/applications/export", EndpointLimit{Rate{1, time.Minute}, 4}},
		{http.MethodPost, "/shadowIT/applications/{entity}/exportCsv", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodPost, "/sslSettings/generatecsr", EndpointLimit{Rate{1, time.Hour}, 1}},
		{http.MethodPost, "/sslSettings/uploadcert/text", EndpointLimit{Rate{1, time.Minute}, 4}},
		{http.MethodPost, "/sslSettings/uploadcertchain/text", EndpointLimit{Rate{1, time.Minute}, 4}},
		{http.MethodPost, "/status/activate", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodGet, "/subclouds", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/subclouds/isLastDcInCountry/{id}", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/timeIntervals", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/timeIntervals/{Id}", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodPost, "/urlCategories/review/domains", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodPut, "/urlCategories/review/domains", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodGet, "/users", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodPost, "/users/bulkDelete", EndpointLimit{Rate{1, time.Minute}, 4}},
		{http.MethodGet, "/virtualZenClusters", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/virtualZenClusters/{virtualZenClusterId}", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/virtualZenNodes", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/virtualZenNodes/{virtualZenNodeId}", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodPost, "/vpnCredentials/bulkDelete", EndpointLimit{Rate{1, time.Minute}, 4}},
		{http.MethodGet, "/webApplicationRules/ruleTypeMapping", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodGet, "/webApplicationRules/{rule_type}", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodPost, "/webApplicationRules/{rule_type}", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodPost, "/webApplicationRules/{rule_type}/availableActions", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodPost, "/webApplicationRules/{rule_type}/duplicate/{ruleId}", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodGet, "/webApplicationRules/{rule_type}/{ruleId}", EndpointLimit{Rate{1, time.Second}, 400}},
		{http.MethodPut, "/webApplicationRules/{rule_type}/{ruleId}", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodDelete, "/webApplicationRules/{rule_type}/{ruleId}", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodDelete, "/workloadGroups", EndpointLimit{Rate{10, time.Minute}, 40}},
		{http.MethodDelete, "/workloadGroups/{workloadGroupId}", EndpointLimit{Rate{10, time.Minute}, 40}},
	}
}

// NewZIAEndpointRateLimiter returns the ZIA OneAPI rate limiter.
func NewZIAEndpointRateLimiter() *EndpointRateLimiter {
	return NewEndpointRateLimiter(ZIAPathPrefix, ZIADefaultLimits(), ZIAEndpointRules())
}
