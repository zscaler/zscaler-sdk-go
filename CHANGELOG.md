# Changelog

# 3.1.6 (February 14, 2025)

## Notes
- Golang: **v1.22**

### Bug Fixes
[PR #305](https://github.com/zscaler/zscaler-sdk-go/pull/305) - Fixed ZIA device management function `GetDevicesByName` pagination, by removing duplicated pagination parameters.

# 3.1.5 (February 12, 2025)

## Notes
- Golang: **v1.22**

### Bug Fixes
[PR #304](https://github.com/zscaler/zscaler-sdk-go/pull/304) - Fixed ZIA `ssl_inspection` validation options.

# 3.1.4 (February 10, 2025)

## Notes
- Golang: **v1.22**

### Bug Fixes
[PR #302](https://github.com/zscaler/zscaler-sdk-go/pull/302) - Implemented fix on the legacy API clients for `ZCC`, `ZIA` and `ZPA` to prevent rate limit override after client instantiation.

### Internal Changes
[PR #302](https://github.com/zscaler/zscaler-sdk-go/pull/302) - Updated SDK Header version to `v3.1.4`

# 3.1.3 (February 5, 2025)

## Notes
- Golang: **v1.22**

### ZIA Policy Export
[PR #301](https://github.com/zscaler/zscaler-sdk-go/pull/301) - Added the following new ZIA API Endpoints:
    - Added `POST /exportPolicies` Exports the rules configured for the specified policy types to JSON files.

### Bug Fixes

[PR #298](https://github.com/zscaler/zscaler-sdk-go/pull/298) - Fixed ZCC `ReadAllPages` pagination function due to panic related to incorrect method reference.

# 3.1.2 (January 28, 2025)

## Notes
- Golang: **v1.22**

### Bug Fixes

[PR #297](https://github.com/zscaler/zscaler-sdk-go/pull/297) - Fixed ZIA SSL Inspection Proxy Gateway attribute.

# 3.1.1 (January 27, 2025)

## Notes
- Golang: **v1.22**

### Bug Fixes

[PR #296](https://github.com/zscaler/zscaler-sdk-go/pull/296) - Fixed ZIA Rate limit override issue to prevent inadivertant panic.

# 3.1.0 (January 20, 2025)

## Notes
- Golang: **v1.22**

### ZIA SSL Inspection Rules
[PR #295](https://github.com/zscaler/zscaler-sdk-go/pull/295) - Added the following new ZIA API Endpoints:
    - Added `GET /sslInspectionRules` Retrieves all SSL inspection rules.
    - Added `GET /sslInspectionRules/{ruleId}` Retrieves the SSL inspection rule based on the specified ID
    - Added `POST /sslInspectionRules` Creates a new SSL inspection rule
    - Added `PUT /sslInspectionRules/{ruleId}` Updates the SSL inspection rule based on the specified ID
    - Added `DELETE /sslInspectionRules/{ruleId}` Deletes an existing SSL inspection rule based on the specified ID

# 3.0.0 (January 20, 2025) - BREAKING CHANGES

## Notes
- Golang: **v1.23**

#### Zscaler OneAPI Support
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293): Added support for [OneAPI](https://help.zscaler.com/oneapi/understanding-oneapi) Oauth2 authentication support through [Zidentity](https://help.zscaler.com/zidentity/what-zidentity).

**NOTES** 
  - Starting at v3.0.0 version this SDK provides dual API client functionality and is backwards compatible with the legacy Zscaler API framework.
  - The new OneAPI framework is compatible only with the following products `ZCC/ZIA/ZPA`.
  - The following products `ZCON` - Cloud Connector and `ZDX` and Zscaler Digital Experience, authentication methods remain unnaffected.

Refer to the [README](https://github.com/zscaler/zscaler-sdk-go/blob/master/README.md) page for details on client instantiation, and authentication requirements on each individual product.

[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293): All API clients now support Config Setter object `ZCC/ZCON/ZDX/ZIA/ZPA`
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293): Added Ability to pass `context` to each method that is sent into the request.

#### ZCC New Endpoints
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZCC API Endpoints:
  - Added `GET /downloadServiceStatus` to download service status for all devices.
  - Added `GET /getDeviceCleanupInfo` to retrieve device cleanup information.
  - Added `PUT /setDeviceCleanupInfo` to cleanup device information.
  - Added `GET /getDeviceDetails` to retrieve device detailed information.
  - Added `GET /getAdminUsers` to retrieve mobile portal admin user.
  - Added `PUT /editAdminUser` to update mobile portal admin user.
  - Added `GET /getAdminUsersSyncInfo` to retrieve mobile portal admin user sync information.
  - Added `POST /syncZiaZdxAdminUsers` to retrieve mobile portal admin users ZIA and ZDX sync information.
  - Added `POST /syncZpaAdminUsers` to retrieve mobile portal admin users ZPA sync information.
  - Added `GET /getAdminRoles` to retrieve mobile portal admin roles.
  - Added `GET /getCompanyInfo` to retrieve company information.
  - Added `GET /getZdxGroupEntitlements` to retrieve ZDX Group entitlement enablement.
  - Added `PUT /updateZdxGroupEntitlement` to retrieve ZDX Group entitlement enablement.
  - Added `GET /updateZpaGroupEntitlement` to retrieve ZPA Group entitlement enablement.
  - Added `GET /web/policy/listByCompany` to retrieve Web Policy By Company ID.
  - Added `PUT /web/policy/activate` to activate mobile portal web policy
  - Added `PUT /web/policy/edit` to update mobile portal web policy
  - Added `DELETE /web/policy/{policyId}/delete` to delete mobile portal web policy.
  - Added `GET /webAppService/listByCompany` to retrieve Web App Service information By Company ID.
  - Added `GET /webFailOpenPolicy/listByCompany` to retrieve web Fail Open Policy information By Company ID.
  - Added `PUT /webFailOpenPolicy/edit` to update mobile portal web Fail Open Policy.
  - Added `GET /webForwardingProfile/listByCompany` to retrieve Web Forwarding Profile information By Company ID.
  - Added `POST /webForwardingProfile/edit` to create a Web Forwarding Profile.
  - Added `DELETE /webForwardingProfile/{profileId}/delete` to delete Web Forwarding Profile.
  - Added `GET /webTrustedNetwork/listByCompany` to retrieve multiple Web Trusted Network information By Company ID.
  - Added `POST /webTrustedNetwork/edit` to create Web Trusted Network resource.
  - Added `PUT /webTrustedNetwork/edit` to update Web Trusted Network resource.
  - Added `DELETE /webTrustedNetwork/{networkId}/delete` to delete Web Trusted Network resource.
  - Added `GET /getWebPrivacyInfo` to retrieve Web Privacy Info.
  - Added `GET /setWebPrivacyInfo` to update Web Privacy Info.

#### ZIA Sandbox Submission - BREAKING CHANGES
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Authentication to Zscaler Sandbox now use the following attributes during client instantiation.
 - `sandboxToken` - Can also be sourced from the `ZSCALER_SANDBOX_TOKEN` environment variable.
 - `sandboxCloud` - Can also be sourced from the `ZSCALER_SANDBOX_CLOUD` environment variable.

**NOTE** The previous `ZIA_SANDBOX_TOKEN` has been deprecated.

#### ZIA Sandbox Rules
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
  - Added `GET /sandboxRules` to retrieve the list of all Sandbox policy rules.
  - Added `GET /sandboxRules/{ruleId}` to retrieve the Sandbox policy rule information based on the specified ID.
  - Added `POST /sandboxRules` to add a Sandbox policy rule. 
  - Added `PUT /sandboxRules/{ruleId}` to update the Sandbox policy rule configuration for the specified ID.
  - Added `DELETE /sandboxRules/{ruleId}` to delete the Sandbox policy rule based on the specified ID.

#### ZIA DNS Control Rules
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
  - Added `GET /firewallDnsRules` to retrieve the list of all DNS Control policy rules.
  - Added `GET /firewallDnsRules/{ruleId}` to retrieve the DNS Control policy rule information based on the specified ID.
  - Added `POST /firewallDnsRules` to add a DNS Control policy rules. 
  - Added `PUT /firewallDnsRules/{ruleId}` to update the DNS Control policy rule configuration for the specified ID.
  - Added `DELETE /firewallDnsRules/{ruleId}` to delete the DNS Control policy rule based on the specified ID.

#### ZIA IPS Control Rules
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
  - Added `GET /firewallIpsRules` to retrieve the list of all IPS Control policy rules.
  - Added `GET /firewallIpsRules/{ruleId}` to retrieve the IPS Control policy rule information based on the specified ID.
  - Added `POST /firewallIpsRules` to add a IPS Control policy rule. 
  - Added `PUT /firewallIpsRules/{ruleId}` to update the IPS Control policy rule configuration for the specified ID.
  - Added `DELETE /firewallIpsRules/{ruleId}` to delete the IPS Control policy rule based on the specified ID.

#### ZIA File Type Control Policy
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
  - Added `GET /fileTypeRules` to retrieve the list of all File Type Control policy rules.
  - Added `GET /fileTypeRules/lite` to retrieve the list of all File Type Control policy rules.
  - Added `GET /fileTypeRules/{ruleId}` to retrieve the File Type Control policy rule information based on the specified ID.
  - Added `POST /fileTypeRules` to add a File Type Control policy rule. 
  - Added `PUT /fileTypeRules/{ruleId}` to update the File Type Control policy rule configuration for the specified ID.
  - Added `DELETE /fileTypeRules/{ruleId}` to delete the File Type Control policy rule based on the specified ID.

#### ZIA Forwarding Control Policy - Proxy Gateways
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
  - Added `GET /proxyGateways` to retrieve the proxy gateway information.
  - Added `GET /proxyGateways/lite` to retrieve the name and ID of the proxy.

#### ZIA Cloud Nanolog Streaming Service (NSS)
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
  - Added `GET /nssFeeds` to retrieve the cloud NSS feeds.
  - Added `GET /nssFeeds/{feedId}` to retrieve information about cloud NSS feed based on the specified ID.
  - Added `POST /nssFeeds` to add a new cloud NSS feed.
  - Added `PUT /nssFeeds/{feedId}` to update cloud NSS feed configuration based on the specified ID.
  - Added `DELETE /nssFeeds/{feedId}` to delete cloud NSS feed configuration based on the specified ID.
  - Added `GET /nssFeeds/feedOutputDefaults` to retrieve the default cloud NSS feed output format for different log types.
  - Added `GET /nssFeeds/testConnectivity/{feedId}` to test the connectivity of cloud NSS feed based on the specified ID
  - Added `POST /nssFeeds/validateFeedFormat` to validates the cloud NSS feed format and returns the validation result

#### ZIA Advanced Threat Protection Policy
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
  - Added `GET /cyberThreatProtection/advancedThreatSettings` to retrieve the advanced threat configuration settings.
  - Added `PUT /cyberThreatProtection/advancedThreatSettings` to update the advanced threat configuration settings.
  - Added `GET /cyberThreatProtection/maliciousUrls` to retrieve the malicious URLs added to the denylist in the Advanced Threat Protection (ATP) policy
  - Added `PUT /cyberThreatProtection/maliciousUrls` to updates the malicious URLs added to the denylist in ATP policy
  - Added `GET /cyberThreatProtection/securityExceptions` to retrieves information about the security exceptions configured for the ATP policy
  - Added `PUT /cyberThreatProtection/securityExceptions` to update security exceptions for the ATP policy
  
#### ZIA Advanced Threat Protection Policy
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
  - Added `GET /cyberThreatProtection/atpMalwareInspection` to retrieve the traffic inspection configurations of Malware Protection policy
  - Added `PUT /cyberThreatProtection/atpMalwareInspection` to update the traffic inspection configurations of Malware Protection policy.
  - Added `GET /cyberThreatProtection/atpMalwareProtocols` to retrieve the protocol inspection configurations of Malware Protection policy
  - Added `PUT /cyberThreatProtection/atpMalwareProtocols` to update the protocol inspection configurations of Malware Protection policy.
  - Added `GET /cyberThreatProtection/malwareSettings` to retrieve the malware protection policy configuration details
  - Added `PUT /cyberThreatProtection/malwareSettings` to update the malware protection policy configuration details.
  - Added `GET /cyberThreatProtection/malwarePolicy` to retrieve information about the security exceptions configured for the Malware Protection policy
  - Added `PUT /cyberThreatProtection/malwarePolicy` to update security exceptions for the Malware Protection policy. 

#### ZIA URL & Cloud App Control Policy Settings
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
  - Added `GET /advancedUrlFilterAndCloudAppSettings` to retrieve information about URL and Cloud App Control advanced policy settings
  - Added `PUT /advancedUrlFilterAndCloudAppSettings` to update the URL and Cloud App Control advanced policy settings

#### ZIA Authentication Settings
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
  - Added `GET /authSettings` to retrieve the organization's default authentication settings information, including authentication profile and Kerberos authentication information.
  - Added `GET /authSettings/lite` to retrieve organization's default authentication settings information.
  - Added `PUT /authSettings` to update the organization's default authentication settings information.

#### ZIA Advanced Settings
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
  - Added `GET /advancedSettings` to retrieve information about the advanced settings.
  - Added `PUT /advancedSettings` to update the advanced settings configuration.

#### ZIA Cloud Applications
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
  - Added `GET /cloudApplications/policy` Retrieves a list of Predefined and User Defined Cloud Applications associated with the DLP rules, Cloud App Control rules, Advanced Settings, Bandwidth Classes, and File Type Control rules.
  - Added `GET /cloudApplications/sslPolicy` Retrieves a list of Predefined and User Defined Cloud Applications associated with the SSL Inspection rules.

#### ZIA Shadow IT Report
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
- Added `PUT /cloudApplications/bulkUpdate` To Update application status and tag information for predefined or custom cloud applications based on the IDs specified
- Added `GET /cloudApplications/lite` Gets the list of predefined and custom cloud applications
- Added `GET /customTags` Gets the list of custom tags available to assign to cloud applications
- Added `POST /shadowIT/applications/export` Export the Shadow IT Report (in CSV format) for the cloud applications recognized by Zscaler based on their usage in your organization.
- Added `POST /shadowIT/applications/{entity}/exportCsv` Export the Shadow IT Report (in CSV format) for the list of users or known locations identified with using the cloud applications specified in the request.

#### ZIA Remote Assistance Support
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
  - Added `GET /remoteAssistance` to retrieve information about the Remote Assistance option.
  - Added `PUT /remoteAssistance` to update information about the Remote Assistance option. Using this option, you can allow Zscaler Support to access your organization’s ZIA Admin Portal for a specified time period to troubleshoot issues.

#### ZIA Organization Details
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
  - Added `GET /orgInformation` to retrieve detailed organization information, including headquarter location, geolocation, address, and contact details.
  - Added `GET /orgInformation/lite` to retrieve minimal organization information.
  - Added `GET /subscriptions` to retrieve information about the list of subscriptions enabled for your tenant. Subscriptions define the various features and levels of functionality that are available to your organization.

#### ZIA End User Notification
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
  - Added `GET /eun` to retrieve information browser-based end user notification (EUN) configuration details.
  - Added `PUT /eun` to update the browser-based end user notification (EUN) configuration details.

#### ZIA Admin Audit Logs
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
  - Added `GET /auditlogEntryReport` to retrieve the status of a request for an audit log report.
  - Added `POST /auditlogEntryReport` to create an audit log report for the specified time period and saves it as a CSV file.
  - Added `DELETE /auditlogEntryReport` to cancel the request to create an audit log report.
  - Added `GET /auditlogEntryReport/download` to download the most recently created audit log report.

#### ZIA Extranets
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
  - Added `GET /extranet` to retrieve the list of extranets configured for the organization
  - Added `GET /extranet/lite` Retrieves the name-ID pairs of all extranets configured for an organization
  - Added `GET /extranet/{Id}` Retrieves information about an extranet based on the specified ID.
  - Added `POST /extranet` Adds a new extranet for the organization.
  - Added `PUT /extranet/{Id}` Updates an extranet based on the specified ID
  - Added `DELETE /extranet/{Id}` Deletes an extranet based on the specified ID

#### ZIA IOT Endpoint
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA IOT API Endpoints:
  - Added `GET /iotDiscovery/deviceTypes` Retrieve the mapping between device type universally unique identifier (UUID) values and the device type names for all the device types supported by the Zscaler AI/ML.
  - Added `GET /iotDiscovery/categories` Retrieve the mapping between the device category universally unique identifier (UUID) values and the category names for all the device categories supported by the Zscaler AI/ML. The parent of device category is device type.
  - Added `GET /iotDiscovery/classifications` Retrieve the mapping between the device classification universally unique identifier (UUID) values and the classification names for all the device classifications supported by Zscaler AI/ML. The parent of device classification is device category.
  - Added `GET /iotDiscovery/deviceList` Retrieve a list of discovered devices with the following key contexts, IP address, location, ML auto-label, classification, category, and type.

#### ZIA 3rd-Party App Governance
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
  - Added `GET /apps/app` to search the 3rd-Party App Governance App Catalog by either app ID or URL.
  - Added `POST /apps/app` to submis an app for analysis in the 3rd-Party App Governance Sandbox.
  - Added `GET /apps/search` to search for an app by name. Any app whose name contains the search term (appName) is returned.
  - Added `GET /app_views/list` to retrieve the list of custom views that you have configured in the 3rd-Party App Governance.
  - Added `GET /app_views/{appViewId}/apps` to retrieves all assets (i.e., apps) that are related to a specified argument (i.e., custom view).

#### ZPA SCIM API
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - The ZPA SCIM API Client now supports instantiation via configSetter mode. See [README](https://github.com/zscaler/zscaler-sdk-go/blob/master/README.md)

# 2.74.0 (November 14, 2024)

## Notes
- Golang: **v1.22**

#### ZIA PAC Files
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following new ZIA API Endpoints:
  - Added `GET /pacFiles` to Retrieves the list of all PAC files which are in deployed state.
  - Added `GET /pacFiles/{pacId}/version` to Retrieves all versions of a PAC file based on the specified ID.
  - Added `GET /pacFiles/{pacId}/version/{pacVersion}` to Retrieves a specific version of a PAC file based on the specified ID.
  - Added `POST /pacFiles` to Adds a new custom PAC file.
  - Added `DELETE /pacFiles/{pacId}` to Deletes an existing PAC file including all of its versions based on the specified ID.
  - Added `PUT /pacFiles/{pacId}/version/{pacVersion}/action/{pacVersionAction}` to Performs the specified action on the PAC file version and updates the file status.
  - Added `POST /pacFiles/validate` to send the PAC file content for validation and returns the validation result.
  - Added `POST /pacFiles/{pacId}/version/{clonedPacVersion}` to Adds a new PAC file version by branching an existing version based on the specified ID.

### ZPA Additions

The SDK now supports interaction with the dedicated SCIM API Endpoint as described in the [Zscaler Help documentation](https://help.zscaler.com/zpa/scim-api-examples). The SCIM Service Provider Endpoints and references to `scim1.private.zscaler.com`.
To authenticate to the SCIM Service Provider Endpoint you can authenticate by providing the following information:

The ZPA Cloud is identified by several cloud name prefixes, which determines which API endpoint the requests should be sent to. The following cloud environments are supported:

* `PRODUCTION`
* `ZPATWO`
* `BETA`
* `GOV`
* `GOVUS`

### Environment variables

You can provide credentials via the `ZPA_SCIM_TOKEN`, `ZPA_IDP_ID`, `ZPA_SCIM_CLOUD` environment variables, representing your ZPA `scimToken`, `idpId`, and `scimCloud` of your ZPA account, respectively.

~> **NOTE 1** `ZPA_SCIM_CLOUD` environment variable is required, and is used to identify the correct API gateway where the API requests should be forwarded to.

~> **NOTE 2** All SCIM APIs are rate limited.

| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `scimToken`       | _(String)_ The ZPA SCIM Bearer token generated from the ZPA console.| `ZPA_SCIM_TOKEN` |    
| `idpId`       | _(String)_ The ZPA IdP ID from the onboarded Identity Provider.| `ZPA_IDP_ID` |
| `scimCloud`       | _(String)_ The ZPA SCIM Cloud for your ZPA Tenant.| `ZPA_SCIM_CLOUD` |

#### ZPA SCIM API Endpoints - (NEW)
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following ZPA SCIM API Endpoints:
  - Added `GET /Groups` Fetch All Groups with pagination
  - Added `GET /Groups/{groupId}` Fetch a Group By ID
  - Added `POST /Groups` Create a new Group
  - Added `PUT /Groups/{groupId}` Update a new Group
  - Added `PATCH /Groups/{groupId}` Partially Update a Group

[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added the following ZPA SCIM API Endpoints:
  - Added `GET /Users` Fetch All Users with pagination
  - Added `GET /Users/{userID}` Fetch a User By ID
  - Added `POST /Groups` Create a new User
  - Added `PUT /Groups/{userID}` Update a new User
  - Added `PATCH /Groups/{userID}` Partially Update a User

#### ZWA - Zscaler Workflow Automation (NEW)
[PR #293](https://github.com/zscaler/zscaler-sdk-go/pull/293) - Added new ZWA endpoint:
  - Added `GET /dlp/v1/incidents/transactions/{transactionId}` Gets the list of all DLP incidents associated with the transaction ID
  - Added `GET /dlp/v1/incidents/{dlpIncidentId}` Gets the DLP incident details based on the incident ID.
  - Added `DELETE /dlp/v1/incidents/{dlpIncidentId}` Deletes the DLP incident for the specified incident ID.
  - Added `GET /dlp/v1/incidents{dlpIncidentId}/change-history` Gets the details of updates made to an incident based on the given ID and timeline.
  - Added `GET /dlp/v1/incidents/{dlpIncidentId}/tickets` Gets the information of the ticket generated for the incident. For example, ticket type, ticket ID, ticket status, etc.
  - Added `POST /dlp/v1/incidents/{dlpIncidentId}/incident-groups/search` Filters a list of DLP incident groups to which the specified incident ID belongs.
  - Added `POST /dlp/v1/incidents/{dlpIncidentId}/close` Updates the status of the incident to resolved and closes the incident with a resolution label and a resolution code.
  - Added `POST /dlp/v1/incidents/{dlpIncidentId}/notes` Adds notes to the incident during updates or status changes.
  - Added `POST /dlp/v1/incidents/{dlpIncidentId}/labels` Assign lables (a label name and it's associated value) to DLP incidents.
  - Added `POST /dlp/v1/incidents/search` Filters DLP incidents based on the given time range and the field values.
  - Added `GET /dlp/v1/incidents/{dlpIncidentId}/triggers` Downloads the actual data that triggered the incident.
  - Added `GET /dlp/v1/incidents/{dlpIncidentId}/evidence` Gets the evidence URL of the incident. 

**Notes** 
| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `key_id`       | _(String)_ The ZWA string that contains the API key ID.| `ZWA_API_KEY_ID` |    
| `key_secret`       | _(String)_ The ZWA string that contains the key secret.| `ZWA_API_SECRET` |
| `cloud`       | _(String)_ The ZWA string containing cloud provisioned for your organization.| `ZWA_CLOUD` |

# 2.74.1 (January 4, 2024)

## Notes
- Golang: **v1.22**

### Bug Fixes
[PR #292](https://github.com/zscaler/zscaler-sdk-go/pull/292) - Fixed ZPA Double Encoding of HTTP GET Request Params - Issue #291
[PR #292](https://github.com/zscaler/zscaler-sdk-go/pull/292) - Updated go mod packages

# 2.74.0 (November 14, 2024)

## Notes
- Golang: **v1.22**

#### ZIA PAC Files
[PR #286](https://github.com/zscaler/zscaler-sdk-go/pull/286) - Added the following new ZIA API Endpoints:
  - Added `GET /pacFiles` to Retrieves the list of all PAC files which are in deployed state.
  - Added `GET /pacFiles/{pacId}/version` to Retrieves all versions of a PAC file based on the specified ID.
  - Added `GET /pacFiles/{pacId}/version/{pacVersion}` to Retrieves a specific version of a PAC file based on the specified ID.
  - Added `POST /pacFiles` to Adds a new custom PAC file.
  - Added `DELETE /pacFiles/{pacId}` to Deletes an existing PAC file including all of its versions based on the specified ID.
  - Added `PUT /pacFiles/{pacId}/version/{pacVersion}/action/{pacVersionAction}` to Performs the specified action on the PAC file version and updates the file status.
  - Added `POST /pacFiles/validate` to send the PAC file content for validation and returns the validation result.
  - Added `POST /pacFiles/{pacId}/version/{clonedPacVersion}` to Adds a new PAC file version by branching an existing version based on the specified ID.

### ZPA Additions

The SDK now supports interaction with the dedicated SCIM API Endpoint as described in the [Zscaler Help documentation](https://help.zscaler.com/zpa/scim-api-examples). The SCIM Service Provider Endpoints and references to `scim1.private.zscaler.com`.
To authenticate to the SCIM Service Provider Endpoint you can authenticate by providing the following information:

The ZPA Cloud is identified by several cloud name prefixes, which determines which API endpoint the requests should be sent to. The following cloud environments are supported:

* `PRODUCTION`
* `ZPATWO`
* `BETA`
* `GOV`
* `GOVUS`

### Environment variables

You can provide credentials via the `ZPA_SCIM_TOKEN`, `ZPA_IDP_ID`, `ZPA_SCIM_CLOUD` environment variables, representing your ZPA `scimToken`, `idpId`, and `scimCloud` of your ZPA account, respectively.

~> **NOTE 1** `ZPA_SCIM_CLOUD` environment variable is required, and is used to identify the correct API gateway where the API requests should be forwarded to.

~> **NOTE 2** All SCIM APIs are rate limited.

| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `scimToken`       | _(String)_ The ZPA SCIM Bearer token generated from the ZPA console.| `ZPA_SCIM_TOKEN` |    
| `idpId`       | _(String)_ The ZPA IdP ID from the onboarded Identity Provider.| `ZPA_IDP_ID` |
| `scimCloud`       | _(String)_ The ZPA SCIM Cloud for your ZPA Tenant.| `ZPA_SCIM_CLOUD` |

#### ZPA SCIM API Endpoints
[PR #286](https://github.com/zscaler/zscaler-sdk-go/pull/286) - Added the following ZPA SCIM API Endpoints:
  - Added `GET /Groups` Fetch All Groups with pagination
  - Added `GET /Groups/{groupId}` Fetch a Group By ID
  - Added `POST /Groups` Create a new Group
  - Added `PUT /Groups/{groupId}` Update a new Group
  - Added `PATCH /Groups/{groupId}` Partially Update a Group

[PR #286](https://github.com/zscaler/zscaler-sdk-go/pull/286) - Added the following ZPA SCIM API Endpoints:
  - Added `GET /Users` Fetch All Users with pagination
  - Added `GET /Users/{userID}` Fetch a User By ID
  - Added `POST /Groups` Create a new User
  - Added `PUT /Groups/{userID}` Update a new User
  - Added `PATCH /Groups/{userID}` Partially Update a User

# 2.732.0 (October 31, 2024)

## Notes
- Golang: **v1.22**

### Internal Changes

[PR #282](https://github.com/zscaler/zscaler-sdk-go/pull/282) - Fixed update function in all specialized ZPA Application Segments
    -`applicationsegmentpra` - The fix now automatically includes the attributes `appId` and `praAppId` in the payload during updates
    - `applicationsegmentinspection` - The fix now automatically includes the attributes `appId` and `inspectAppId` in the payload during updates
    - `applicationsegmentbrowseraccess` - The fix now automatically includes the attributes `appId` and `baAppId` in the payload during updates


# 2.731.0 (October 30, 2024)

## Notes
- Golang: **v1.22**

### Enhancements

  - Zscaler Cloud Connector (ZCON)
    - Added `GET /provUrl` endpoint to list provisioning templates.
    - Added `GET /provUrl/{id}` endpoint to retrieve a specific provisioning template.
    - Added `POST /provUrl` endpoint to create provisioning template.
    - Added `PUT /provUrl/{id}` endpoint to update a specific provisioning template.
    - Added `DELETE /provUrl/{id}` endpoint to delete a specific provisioning template.

### Internal Changes

[PR #281](https://github.com/zscaler/zscaler-sdk-go/pull/281) - Added new ZPA Attributes:
  - Resource: `applicationsegment`
    * `extranetEnabled`
    * `apiProtectionEnabled`
    * `zpnErId`

  - Resource: `policysetcontrollerv1` and `policysetcontrollerv2`
    * `disabled`
    * `extranetEnabled`
    * `extranetDTO`
    * `privilegedPortalCapabilities`

# 2.73.0 (October 30, 2024)

## Notes
- Golang: **v1.22**

### Enhancements

  - Zscaler Cloud Connector (ZCON)
    - Added `GET /provUrl` endpoint to list provisioning templates.
    - Added `GET /provUrl/{id}` endpoint to retrieve a specific provisioning template.
    - Added `POST /provUrl` endpoint to create provisioning template.
    - Added `PUT /provUrl/{id}` endpoint to update a specific provisioning template.
    - Added `DELETE /provUrl/{id}` endpoint to delete a specific provisioning template.

### Internal Changes

[PR #281](https://github.com/zscaler/zscaler-sdk-go/pull/281) - Added new ZPA Attributes:
  - Resource: `applicationsegment`
    * `extranetEnabled`
    * `apiProtectionEnabled`
    * `zpnErId`

  - Resource: `policysetcontrollerv1` and `policysetcontrollerv2`
    * `disabled`
    * `extranetEnabled`
    * `extranetDTO`
    * `privilegedPortalCapabilities`

# 2.72.5 (October 8, 2024)

## Notes
- Golang: **v1.22**

### Internal Changes

[PR #280](https://github.com/zscaler/zscaler-sdk-go/pull/280) - Added missing attribute `sourceCountries` to ZIA `firewallfilteringrule`

# 2.72.4 (October 3, 2024)

## Notes
- Golang: **v1.22**

### Internal Changes

[PR #278](https://github.com/zscaler/zscaler-sdk-go/pull/278) - Consolidated several ZPA common functions for Struct simplication.

# 2.72.3 (September 30, 2024)

## Notes
- Golang: **v1.22**

### Bug Fixes

[PR #277](https://github.com/zscaler/zscaler-sdk-go/pull/277) - Added new attributes to ZPA `servicedgegroup` and `serviceedgecontroller` packages.

# 2.72.2 (September 11, 2024)

## Notes
- Golang: **v1.22**

### Bug Fixes

[PR #276](https://github.com/zscaler/zscaler-sdk-go/pull/276) - Fixed removed `omitempty` from ZPA `microtenant_id` attribute in `policysetcontrollerv2`.

# 2.72.1 (August 16, 2024)

## Notes
- Golang: **v1.22**

### Bug Fixes

[PR #274](https://github.com/zscaler/zscaler-sdk-go/pull/274) - Added new ZIA function `GetVIPRecommendedList`, which will support all optional parameters when retrieving the list of recommended Virtual IP addresses per datacenter. The following optional parameters are now supported:
  - `routable_ip` - (Boolean) The routable IP address.
  - `within_country_only` - (Boolean) Search within country only.
  - `include_private_service_edge` - (Boolean) Include ZIA Private Service Edge VIPs.
  - `include_current_vips` - (Boolean) Include currently assigned VIPs.
  - `latitude` - (Number) The latitude coordinate of the GRE tunnel source.
  - `longitude` - (Number) The longitude coordinate of the GRE tunnel source.
  - `subcloud` - (String) The longitude coordinate of the GRE tunnel source.

# 2.72.0 (August 13, 2024)

## Notes
- Golang: **v1.22**

### ZPA Additions

#### Segment Group
- Added new optimized `V2` endpoint `PUT /segmentGroup/{segmentGroupId}` to prevent "payload.size.exceeded" error when updating a segment group with large numbers or application segments attached. [PR #273](https://github.com/zscaler/zscaler-sdk-go/pull/273)
  **NOTE** The `V1` endpoint `PUT /segmentGroup/{segmentGroupId}` will eventually be deprecated; however, this change should not affect existing Segment Group configurations.

### Bug Fixes

[PR #270](https://github.com/zscaler/zscaler-sdk-go/pull/270) - Fixed `ports` attribute from `string` to `slice of intergers` in `locationmanagement`.

# 2.71.0 (August 11, 2024)

## Notes
- Golang: **v1.22**

### ZIA Additions

#### VPN Credentials and Location Management
- Added `POST /locations/bulkDelete` Bulk delete locations up to a maximum of 100 locations per request. The response returns the location IDs that were successfully deleted. [PR #270](https://github.com/zscaler/zscaler-sdk-go/pull/270)

- Added `POST /vpnCredentials/bulkDelete` Bulk delete vpn credentails up to a maximum of 100 vpn credentials per request. The response returns the vpn credential IDs that were successfully deleted. [PR #270](https://github.com/zscaler/zscaler-sdk-go/pull/270)

#### Cloud App Control Policies

- Added `GET /webApplicationRules/ruleTypeMapping` to return backend keys that match the application type string.

### Bug Fixes

[PR #270](https://github.com/zscaler/zscaler-sdk-go/pull/270) - Fixed `ports` attribute from `string` to `slice of string` in `locationmanagement`.

# 2.70.0 (July 23, 2024)

## Notes
- Golang: **v1.22**

### ZIA Additions

#### Cloud App Control Rules
- Added `GET /webApplicationRules/{rule_type}` to Get the list of Web Application Rule by type [PR #270](https://github.com/zscaler/zscaler-sdk-go/pull/270)
- Added `GET /webApplicationRules/{rule_type}/{ruleId}` to Get a Web Application Rule by type and id[PR #270](https://github.com/zscaler/zscaler-sdk-go/pull/270)
- Added `POST /webApplicationRules/{rule_type}` to Adds a new Web Application rule.[PR #270](https://github.com/zscaler/zscaler-sdk-go/pull/270)
- Added `PUT /webApplicationRules/{rule_type}/{ruleId}` to Update a new Web Application rule.[PR #270](https://github.com/zscaler/zscaler-sdk-go/pull/270)
- Added `DELETE /webApplicationRules/{rule_type}/{ruleId}` to Delete a new Web Application rule.[PR #270](https://github.com/zscaler/zscaler-sdk-go/pull/270)

#### DLP Dictionary
- Added `GET /dlpDictionaries/{dictId}/predefinedIdentifiers` to Retrieves the list of identifiers that are available for selection in the specified hierarchical DLP dictionary.[PR #270](https://github.com/zscaler/zscaler-sdk-go/pull/270)
- Added new attribute `dlpdictionary` attributes:
  * `confidenceLevelForPredefinedDict`: `String` - The DLP confidence threshold for predefined dictionaries
  * `dictionaryCloningEnabled`: `Bool` - A Boolean constant that indicates that the cloning option is supported for the DLP dictionary using the true value.
  * `customPhraseSupported`: `Bool` - A Boolean constant that indicates that custom phrases are supported for the DLP dictionary using the true value.
  * `proximityLengthEnabled`: `Bool` - A Boolean constant that indicates whether the proximity length option is supported for a DLP dictionary or not.

#### URL Categories
- Added `POST /urlLookup` to Retrieve Zscaler's default classification for a given set of URLs (e.g., ['abc.com', 'xyz.com']).[PR #270](https://github.com/zscaler/zscaler-sdk-go/pull/270)
- Added `GET /urlCategories/urlQuota` to Gets information on the number of unique URLs that are currently provisioned for your organization as well as how many URLs you can add before reaching that number.[PR #270](https://github.com/zscaler/zscaler-sdk-go/pull/270)
- Added `GET /urlCategories/lite` to Gets a lightweight key-value list of all or custom URL categories.[PR #270](https://github.com/zscaler/zscaler-sdk-go/pull/270)
- Added `GET /urlCategories/review/domains` to find matching entries present in existing custom URL categories.[PR #270](https://github.com/zscaler/zscaler-sdk-go/pull/270)
- Added `GET /urlCategories/review/domains` Adds the list of matching URLs fetched by POST /urlCategories/review/domains to the specified custom URL categories. A maximum of 100 URL categories can be updated at once using this request.[PR #270](https://github.com/zscaler/zscaler-sdk-go/pull/270)
- Added new attribute `urlCategories2` to `urlfilteringrules` package. See [Zscaler Release Notes](https://help.zscaler.com/zia/release-upgrade-summary-2024#:~:text=Filtering%20Policy.-,Update%20to%20Cloud%20Service%20API,-The%20UrlFilteringRule%20model)

# 2.61.12 (July 8, 2024)

## Notes
- Golang: **v1.21**

### Bug Fixes

- [PR #269](https://github.com/zscaler/zscaler-sdk-go/pull/269) - Fixed ZPA App Protection (Inspection) resources with missing attributes.
- [PR #269](https://github.com/zscaler/zscaler-sdk-go/pull/269) - Fixed zpa_inspection_profile missing `overrideAction` attribute.

# 2.61.11 (July 7, 2024)

## Notes
- Golang: **v1.21**

### Bug Fixes

- [PR #269](https://github.com/zscaler/zscaler-sdk-go/pull/269) - Fixed ZPA App Protection (Inspection) resources with missing attributes.

# 2.61.10 (July 5, 2024)

## Notes
- Golang: **v1.21**

### Bug Fixes

- [PR #268](https://github.com/zscaler/zscaler-sdk-go/pull/268) - Fixed ZPA Cloud Browser Isolation resources to allow search by name and ID.
  - `cbibannercontroller`
  - `cbicertificatecontroller`
  - `cbiprofilecontroller`

# 2.61.9 (July 5, 2024)

## Notes
- Golang: **v1.21**

### Bug Fixes

- [PR #268](https://github.com/zscaler/zscaler-sdk-go/pull/268) - Fixed ZPA Cloud Browser Isolation resources to allow search by name and ID.
  - `cbibannercontroller`
  - `cbicertificatecontroller`
  - `cbiprofilecontroller`

# 2.61.8 (July 4, 2024)

## Notes
- Golang: **v1.21**

### Bug Fixes

- [PR #267](https://github.com/zscaler/zscaler-sdk-go/pull/267) - Fixed ZPA Cloud Browser Isolation resources to allow search by name and ID.
  - `cbibannercontroller`
  - `cbicertificatecontroller`
  - `cbiprofilecontroller`
  - `cbizpaprofile`

# 2.61.7 (July 2, 2024)

## Notes
- Golang: **v1.21**

### Bug Fixes

- [PR #266](https://github.com/zscaler/zscaler-sdk-go/pull/266) - Added ZIA ``locationmanagement`` package missing attributes
  - `cookiesAndProxy`
  - `iotEnforcePolicySet`
  - `ecLocation`
  - `excludeFromDynamicGroups`
  - `excludeFromManualGroups`
  - `dynamiclocationGroups`
    - `id`
    - `name`
- `staticlocationGroups`
    - `id`
    - `name`

# 2.61.6 (July 2, 2024)

## Notes
- Golang: **v1.21**

### Bug Fixes

- [PR #265](https://github.com/zscaler/zscaler-sdk-go/pull/265) - Fixed ``cbiprofilecontroller`` package `message` attribute type from `bool` to `string`

### Enhancement
- [PR #265](https://github.com/zscaler/zscaler-sdk-go/pull/265) - Included  new``cbiprofilecontroller`` attribute options:
  - `watermark` - Admins can enable watermarking per isolation profile and choose to display the user ID, date and timestamp (in UTC), and a custom message
    - `enabled`
    - `showUserId`
    - `showTimestamp`
    - `showMessage`
    - `message`
  - `debugMode` - Enable to allow starting isolation sessions in debug mode to collect troubleshooting information.
    - `filePassword` - Optional password to debug files when this mode is enabled.
  - `forwardToZia` - Optional password to debug files when this mode is enabled.
    - `organizationId` - Use the organization ID from the Company Profile section.
    - `cloudName` - The cloud name on which the organization exists. i.e `zscalertwo`
    - `pacFileUrl` - Enable to have the PAC file be configured on the Isolated browser to forward traffic via ZIA.
  - `deepLink` - Deep Linking allows users to open applications from their local machine via the rendered deep link data on an isolated web page.
    - `enabled`
    - `applications` - If no specific applications are added here, then deep linking is applied to all of your applications.

# 2.61.5 (July 2, 2024)

## Notes
- Golang: **v1.21**

### Bug Fixes

- [PR #264](https://github.com/zscaler/zscaler-sdk-go/pull/264) - Fixed ZPA `policysetcontroller` and `policysetcontrollerv2` attributes `appConnectorGroups`, `serviceEdgeGroups`, `appServerGroups` by removing the `omitempty` tag

# 2.61.4 (June 28, 2024)

## Notes
- Golang: **v1.21**

### Enhancement

- [PR #263](https://github.com/zscaler/zscaler-sdk-go/pull/263) - Added new Session Time Ticker to ZIA Client to track session renewal and expiration time.
- [PR #263](https://github.com/zscaler/zscaler-sdk-go/pull/263) - Added ZDX query string parameter `q` to assist with refined search for a user name or email.

# 2.61.3 (June 24, 2024)

## Notes
- Golang: **v1.21**

### Enhancement

- [PR #262](https://github.com/zscaler/zscaler-sdk-go/pull/262) - Enhanced filtering capability for the ZPA GetAll functions in the following packages:
  - `applicationsegmentinspection`
  - `applicationsegmentpra`

# 2.61.2 (June 21, 2024)

## Notes
- Golang: **v1.21**

### Bug Fixes

- [PR #261](https://github.com/zscaler/zscaler-sdk-go/pull/261) - Fixed an issue where sessions were not properly refreshed when expired, leading to `SESSION_NOT_VALID` errors during API requests after a period of inactivity. This ensures that sessions are correctly refreshed before making API calls, improving the reliability of the SDK.

# 2.61.1 (June 17, 2024)

## Notes
- Golang: **v1.21**

### Enhancements

- [PR #259](https://github.com/zscaler/zscaler-sdk-go/pull/259) - ZIA Activator can now be compiled directly from the [Zscaler-SDK-Go](https://github.com/zscaler/zscaler-sdk-go), by executing the command `make ziaActivator` from the root of the repository directory.

- [PR #259](https://github.com/zscaler/zscaler-sdk-go/pull/259) - Combined the following two ZIA functions `GetIncludeOnlyUrlKeyWordCounts` and `GetCustomURLCategories` for simplicity. It's possible now to set the following parameters concurrently:
  - `customOnly` - The parameter is set to true by default. If set to true, it gets information on custom URL categories only.
  - `includeOnlyUrlKeywordCounts` - 

### Fixes

- [PR #259](https://github.com/zscaler/zscaler-sdk-go/pull/259) - Added missing new `city` field attribute in the ZIA package `trafficforwarding/staticips`. The attribute block returns the following information:
  - `id` - ID of the city
  - `name` - Name of the city i.e "Toronto, Ontario, Canada"
- [PR #259](https://github.com/zscaler/zscaler-sdk-go/pull/259) - The ZIA API client, now retries on `StatusPreconditionFailed` 412.

### Documentation
- [PR #259](https://github.com/zscaler/zscaler-sdk-go/pull/259) - Added several ZDX CLI based examples into the examples. [ZDX Examples](https://github.com/zscaler/zscaler-sdk-go/tree/master/examples/zdx). The results returned by the API are displayed in table format.

# 2.61.0 (June 14, 2024)

## Notes
- Golang: **v1.21**

### Enhancements

- [PR #258](https://github.com/zscaler/zscaler-sdk-go/pull/258) - Added new ZDX API methods and Endpoints
  - `GET` - `/alerts/ongoing`
  - `GET` - `/alerts/historical`
  - `GET` - `/alerts/{alert_id}`
  - `GET` - `/alerts/{alert_id}/affected_devices`
  - `GET` - `/inventory/software`
  - `GET` - `/alerts/software/{software_key}`
  - `GET` - `/active_geo`
  - `GET` - `/devices/{deviceid}/deeptraces`
  - `POST` - `/devices/{deviceid}/deeptraces`
  - `GET` - `/devices/{deviceid}/deeptraces/{trace_id}`
  - `DELETE` - `/devices/{deviceid}/deeptraces/{trace_id}`
  - `GET` - `/devices/{deviceid}/deeptraces/{trace_id}/web_probe-metrics`
  - `GET` - `/devices/{deviceid}/deeptraces/{trace_id}/cloudpath-metrics`
  - `GET` - `/devices/{deviceid}/deeptraces/{trace_id}/cloudpath`
  - `GET` - `/devices/{deviceid}/deeptraces/{trace_id}/health-metrics`
  - `GET` - `/devices/{deviceid}/deeptraces/{trace_id}/events`
  - `GET` - `/devices/{deviceid}/deeptraces/{trace_id}/top-processes`
  - `POST` - `/analysis`
  - `GET` - `/analysis/{analysis}`
  - `DELETE` - `/analysis/{analysis}`

For details on the functionality of each of the above endpoints, please see [ZDX API Guide](https://help.zscaler.com/zdx/getting-started-zdx-api)

### Internal Changes

- [PR #258](https://github.com/zscaler/zscaler-sdk-go/pull/258) - Refactored ZPA package to centralize the `service.go` client instantiation.
- [PR #258](https://github.com/zscaler/zscaler-sdk-go/pull/258) - Refactored ZIA package to centralize the `service.go` client instantiation.
- [PR #258](https://github.com/zscaler/zscaler-sdk-go/pull/258) - Refactored ZDX package to centralize the `service.go` client instantiation.
- [PR #258](https://github.com/zscaler/zscaler-sdk-go/pull/258) - Refactored ZCC package to centralize the `service.go` client instantiation.
- [PR #258](https://github.com/zscaler/zscaler-sdk-go/pull/258) - Refactored ZCON package to centralize the `service.go` client instantiation.
- [PR #258](https://github.com/zscaler/zscaler-sdk-go/pull/258) - Enhanced test coverage statements across several functions in all API packages.

### Deprecations
- [PR #258](https://github.com/zscaler/zscaler-sdk-go/pull/258) - Deprecated ZIA `urlcategories` function ``GetCustomURLCategories``. The `customOnly` parameter is now combined within the function `GetIncludeOnlyUrlKeyWordCounts` and can be optionally set to `true`.

- [PR #258](https://github.com/zscaler/zscaler-sdk-go/pull/258) - Deprecated ZPA `ConvertV1ResponseToV2Request` from package ``policysetcontrollerv2`` package. The function now lives directly in the ZPA Terraform Provider to convert ``policysetcontroller`` v1 responses into ``policysetcontrollerv2`` format.

### Documentation
- [PR #258](https://github.com/zscaler/zscaler-sdk-go/pull/258) - Expanded README with more details regarding client instantiation and options. The README also include details regarding rate limites, retries and caching parameters.

# 2.6.0 (June 6, 2024)

## Notes
- Golang: **v1.21**

### Fixes

- [PR #256](https://github.com/zscaler/zscaler-sdk-go/pull/256) - Refactored `applicationsegment` package by splitting the following endpoints into its own packages:
  - POST - `application/{applicationId}/move` - Moves application from a parent tenant to a microtenant. An application segment can only be moved from a parent to a microtenant 
  - PUT - `application/{applicationId}/share` - Share the Application Segment between microtenants. An application can only be shared between microtenants.
To learn more about microtenants see: [About Microtenants](https://help.zscaler.com/zpa/about-microtenants)

- [PR #256](https://github.com/zscaler/zscaler-sdk-go/pull/256) - Added support to Service Edge Scheduler to configure a ServiceEdge schedule frequency to delete inactive private brokers with configured frequency.
  - GET - `serviceEdgeSchedule` - Get a Configured ServiceEdge schedule frequency.
  - POST - `serviceEdgeSchedule` - Configure a ServiceEdge schedule frequency to delete the in active private broker with configured frequency.
  - PUT - `serviceEdgeSchedule/{id}` - Modifies a ServiceEdge schedule frequency to delete the in active private broker with configured frequency.

- [PR #256](https://github.com/zscaler/zscaler-sdk-go/pull/256) - Added support customer controller endpoint `authDomains` to retrieve authentication domains for the specified customer.

# 2.5.22 (May 31, 2024)

## Notes
- Golang: **v1.21**

### Fixes

- [PR #253](https://github.com/zscaler/zscaler-sdk-go/pull/253) - Implemented new ZPA error handling to retry on new `400` and `409` error format message:

```json
  "id" : "api.concurrent.access.error",
  "reason" : "Unable to modify the resource due to concurrent change requests. Try again"
```
- [PR #253](https://github.com/zscaler/zscaler-sdk-go/pull/253) - Adjusted several ZPA integration tests to cover new use cases.

# 2.5.21 (May 22, 2024)

## Notes
- Golang: **v1.21**

### Fixes

- [PR #251](https://github.com/zscaler/zscaler-sdk-go/pull/251) - Fixed attribute `subRules` within the ZIA package to support `array`[object]. The attribute is now in its dedicated Struct.

# 2.5.2 (May 16, 2024)

## Notes
- Golang: **v1.21**

### Enhancements

- [PR #248](https://github.com/zscaler/zscaler-sdk-go/pull/248) - Added getAppsByType endpoint - Get all configured `BROWSER_ACCESS`, `INSPECT`, `SECURE_REMOTE_ACCESS` application segments.

# 2.5.1 (May 16, 2024)

## Notes
- Golang: **v1.21**

### Enhancements

- [PR #246](https://github.com/zscaler/zscaler-sdk-go/pull/246) - Added ZPA Cloud Browser Isolation External Profile new attributes:
  * `forwardToZia`
    * `enabled`
    * `organizationId`
    * `cloudName`
    * `pacFileUrl`

  * `debugMode`
    * `allowed`
    * `filePassword`

- [PR #247](https://github.com/zscaler/zscaler-sdk-go/pull/247) - Added getAppsByType endpoint - Get all configured `BROWSER_ACCESS`, `INSPECT`, `SECURE_REMOTE_ACCESS` application segments.

# 2.5.0 (May 6, 2024)

## Notes
- Golang: **v1.21**

### Enhancements

- [PR #240](https://github.com/zscaler/zscaler-sdk-go/pull/240) - Added new `Retry-After` header to ZPA API Client. Please see API Developer's documentation [here](https://help.zscaler.com/zpa/understanding-rate-limiting) for details.
- [PR #241](https://github.com/zscaler/zscaler-sdk-go/pull/241) - Added new ZIA URL Filtering Rule attribute `source_ip_groups`

# 2.4.35 (April 12, 2024)

## Notes
- Golang: **v1.21**

### Fixes

- [PR #243](https://github.com/zscaler/zscaler-sdk-go/pull/243) - Added custom error handling to ZIA `MakeAuthRequestZIA` function for more clarity during authentication failures.

# 2.4.34 (April 8, 2024)

## Notes
- Golang: **v1.21**

### Fixes

- [PR #242](https://github.com/zscaler/zscaler-sdk-go/pull/242) - Fixed ZPA `bacertificate` package by adding missing attributes `publicKey`

# 2.4.33 (April 5, 2024)

## Notes
- Golang: **v1.21**

### Fixes

- [PR #241](https://github.com/zscaler/zscaler-sdk-go/pull/241) - Fixed DLP Web Rule attributes `auditor`, `icapServer`, and `notificationTemplate` to use a common struct type `IDName`

# 2.4.32 (March 27, 2024)

## Notes
- Golang: **v1.21**

### Fixes

- [PR #239](https://github.com/zscaler/zscaler-sdk-go/pull/239) - Added function `GetByIP` in the ZIA `vpncredentials` package 

# 2.4.31 (March 16, 2024)

## Notes
- Golang: **v1.21**

### Fixes

- [PR #237](https://github.com/zscaler/zscaler-sdk-go/pull/237) - Fixed ZPA `ConvertV1ResponseToV2Request` due to missing `CONSOLE` objectType
- [PR #238](https://github.com/zscaler/zscaler-sdk-go/pull/238) - Fixed ZPA `ConvertV1ResponseToV2Request` due to missing `ZpnIsolationProfileID` attribute

# 2.4.3 (March 13, 2024)

## Notes
- Golang: **v1.21**

### Fixes

- Added ZPA ``serviceEdgeGroups`` missing attribute to policy set controller v1

# 2.4.2 (March 9, 2024)

## Notes
- Golang: **v1.21**

### Fixes

- Fixed ZPA `ConvertV1ResponseToV2Request` due to missing `CLIENT_TYPE` objectType

# 2.4.1 (March 9, 2024)

## Notes
- Golang: **v1.21**

### Fixes

- Added new ZPA `matchStyle` attribute to all application segment resources

# 2.4.0 (March 7, 2024)

## Notes
- Golang: **v1.21**

### ZPA Additions

#### Emergency Access
- Added `POST /emergencyAccess/user` to create an emergency acess user for a specified customer. [PR #226](https://github.com/zscaler/zscaler-sdk-go/pull/226)
- Added `GET /emergencyAccess/user` to get all emergency acess users for a specified customer. [PR #226](https://github.com/zscaler/zscaler-sdk-go/pull/226)
- Added `GET /emergencyAccess/user/{userId}` to get the emergency access user for a specified customer. [PR #226](https://github.com/zscaler/zscaler-sdk-go/pull/226)
- Added `PUT /emergencyAccess/user/{userId}`to update the emergency access user for thae specified customer. [PR #226](https://github.com/zscaler/zscaler-sdk-go/pull/226)
- Added `PUT /emergencyAccess/user/{userId}/activate` to activate the emergency access user for the specified customer. [PR #226](https://github.com/zscaler/zscaler-sdk-go/pull/226)
- Added `PUT /emergencyAccess/user/{userId}/deactivate` to deactivate the emergency access user for the specified customer. [PR #226](https://github.com/zscaler/zscaler-sdk-go/pull/226)

#### Policy Access Controller
- Added `POST and PUT /mgmtconfig/v2/admin/customers/{customerId}/policySet/{policySetId}/rule` endpoints for access policy rule creation. This endpoint allows for larger payload submission. [PR #228](https://github.com/zscaler/zscaler-sdk-go/pull/228)

- Added `POST and PUT /mgmtconfig/v2/admin/customers/{customerId}/policySet/{policySetId}/rule` endpoints for access policy rule creation. This endpoint allows for larger payload submission. [PR #228](https://github.com/zscaler/zscaler-sdk-go/pull/228) 

#### Privileged Remote Access Approval
- Added `GET /mgmtconfig/v1/admin/customers/{customerId}/approval` endpoint to get all PRA Approval resources for a specified customer
- Added `GET /mgmtconfig/v1/admin/customers/{customerId}/approval/{id}` endpoint to get a specific PRA Approval resources for a specified customer
- Added `POST /mgmtconfig/v1/admin/customers/{customerId}/approval` endpoint to add PRA Approval resources for a specified customer. [PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/235) 
- Added `PUT /mgmtconfig/v1/admin/customers/{customerId}/approval/{id}` endpoint to update a specific PRA Approval resources for a specified customer. [PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/235) 
- Added `DELETE /mgmtconfig/v1/admin/customers/{customerId}/approval/{id}` endpoint to delete a specific PRA Approval resources for a specified customer. [PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/235) 
- Added `DELETE /mgmtconfig/v1/admin/customers/{customerId}/approval/expired` endpoint to delete all PRA Approval resources for a specified customer. [PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/235) 

#### Privileged Remote Access Console
- Added `GET /mgmtconfig/v1/admin/customers/{customerId}/console` endpoint to get all PRA Console resources for a specified customer
- Added `GET /mgmtconfig/v1/admin/customers/{customerId}/console/{id}` endpoint to get a specific PRA Console resources for a specified customer
- Added `GET /mgmtconfig/v1/admin/customers/{customerId}/console/praPortal/{portalId}` endpoint to get privileged consoles for a specified privileged portal.
- Added `POST /mgmtconfig/v1/admin/customers/{customerId}/console` endpoint to add PRA Console resources for a specified customer. [PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/235)
- Added `POST /mgmtconfig/v1/admin/customers/{customerId}/console/bulk` endpoint to create a list of PRA Console resources to a specified privileged portal and customer. [PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/235) 
- Added `PUT /mgmtconfig/v1/admin/customers/{customerId}/console/{id}` endpoint to update a specific PRA Console resources for a specified customer. [PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/235) 
- Added `DELETE /mgmtconfig/v1/admin/customers/{customerId}/console/{id}` endpoint to delete a specific PRA Console resources for a specified customer. [PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/235) 

#### Privileged Remote Access Portal
- Added `GET /mgmtconfig/v1/admin/customers/{customerId}/praPortal` endpoint to get all PRA Portal resources for a specified customer
- Added `GET /mgmtconfig/v1/admin/customers/{customerId}/praPortal/{id}` endpoint to get a specific PRA Portal resources for a specified customer
- Added `POST /mgmtconfig/v1/admin/customers/{customerId}/praPortal` endpoint to add PRA Portal resource for a specified customer. [PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/235)
- Added `PUT /mgmtconfig/v1/admin/customers/{customerId}/praPortal/{id}` endpoint to update a specific PRA Portal resources for a specified customer. [PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/235) 
- Added `DELETE /mgmtconfig/v1/admin/customers/{customerId}/praPortal/{id}` endpoint to delete a specific PRA Portal resources for a specified customer. [PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/235) 

#### Privileged Remote Access Credential
- Added `GET /mgmtconfig/v1/admin/customers/{customerId}/credential` endpoint to get all PRA Credential resources for a specified customer
- Added `GET /mgmtconfig/v1/admin/customers/{customerId}/credential/{id}` endpoint to get a specific PRA Credential resources for a specified customer
- Added `POST /mgmtconfig/v1/admin/customers/{customerId}/credential` endpoint to add PRA Credential resource for a specified customer. [PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/235)
- Added `POST /mgmtconfig/v1/admin/customers/{customerId}/credential/move` endpoint to move PRA Credentials from one microtenant to another microtenant. [PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/235)
- Added `PUT /mgmtconfig/v1/admin/customers/{customerId}/credential/{id}` endpoint to update a specific PRA Credential resources for a specified customer. [PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/235) 
- Added `DELETE /mgmtconfig/v1/admin/customers/{customerId}/credential/{id}` endpoint to delete a specific PRA Credential resources for a specified customer. [PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/235) 

#### Application Segment
- Added `POST /mgmtconfig/v1/admin/customers/{customerId}/application/move` to move application segments from one microtenant to another. [PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/233) 
- Added `PUT /mgmtconfig/v1/admin/customers/{customerId}/application/share` to share application segments between microtenants. [PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/233) 
[PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/235) Included new application segment attribute `matchStyle` to support `Exact Match` vs. `Multimatch` configuration. [Learn More Here ](https://help.zscaler.com/zpa/using-app-segment-multimatch)
ment 

### Acceptance Tests
[PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/235) - Implemented centralized ZIA and ZPA sweep facility for tenant cleanup pre and post integration tests.

### Fixes
[PR #233](https://github.com/zscaler/zscaler-sdk-go/pull/235) - Fixed ZPA API client HTTP request to prevent undesired URL encoding or special characters.

# 2.3.11 (February 28, 2024)

## Notes
- Golang: **v1.21**

### Fixes

- [PR #230](https://github.com/zscaler/zscaler-sdk-go/pull/230) - Implemented centralized sweep for ZIA and ZPA packages.
- [PR #231](https://github.com/zscaler/zscaler-sdk-go/pull/231) - Fixed ZPA Application Segment PRA changed attribute change from `sraPortal` to `praApps`.

# 2.3.10 (February 22, 2024)

## Notes
- Golang: **v1.19**

### Fixes

- [PR #225](https://github.com/zscaler/zscaler-sdk-go/pull/225) - Fixed ZIA JSession authentication logic to use `after (now)` instead of `before (now)` to prevent specific JSessionID authentication edge cases.

# 2.3.9 (February 12, 2024)

## Notes
- Golang: **v1.21**

### Fixes

- [PR #224](https://github.com/zscaler/zscaler-sdk-go/pull/224) - Added support to ZPA Policy Access Redirection resource.
  - **NOTE** This feature is in limited availability. Contact Zscaler Support to enable this feature for your organization.

# 2.3.8 (January 31, 2024)

## Notes
- Golang: **v1.19**

### Fixes

- [PR #221](https://github.com/zscaler/zscaler-sdk-go/pull/221) - Fixed new `workloadGroups` attribute for the following resources:
  - ``Cloud Firewall Rules``
  - ``DLP Web Rules``
  - ``URL Filtering Rules``

# 2.3.7 (January 26, 2024)

## Notes
- Golang: **v1.19**

### Enhacements

- [PR #215](https://github.com/zscaler/zscaler-sdk-go/pull/215) - Added new ZPA attributes for application segment.
  - matchStyle
  - inconsistentConfigDetails

- [PR #217](https://github.com/zscaler/zscaler-sdk-go/pull/217) - Added support for ZIA Workload Groups Tagging

# 2.3.6 (January 15, 2024)

## Notes
- Golang: **v1.19**

### Enhacements

- [PR #183](https://github.com/zscaler/zscaler-sdk-go/pull/183) - (feat): Implemented New ZPA Bulk Reorder Policy Rule

# 2.3.5 (December 20, 2023)

## Notes
- Golang: **v1.19**

### Enhacements

- Removed omitempty tag from enabled attribute ZPA in Assistant Schedule struct.

# 2.3.4 (December 19, 2023)

## Notes
- Golang: **v1.19**

### Enhacements

- [PR #209](https://github.com/zscaler/zscaler-sdk-go/pull/209) - Added support to ZPA Application Segment within the ZIA Firewall Filtering rule resource. Only ZPA application segments with the Source IP Anchor option enabled are supported.

# 2.3.3 (December 18, 2023)

## Notes
- Golang: **v1.19**

### Enhacements

- [PR #207](https://github.com/zscaler/zscaler-sdk-go/pull/207) - Added missing ZIA URL Filtering Rule attribute `userRiskScoreLevels`: Supported values: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`

- [PR #207](https://github.com/zscaler/zscaler-sdk-go/pull/207) - Added New ZIA URL Filtering Rule `cbiProfile` attribute to support `ISOLATE` action.

# 2.3.2 (December 16, 2023)

## Notes
- Golang: **v1.19**

### Enhacements

- [PR #206](https://github.com/zscaler/zscaler-sdk-go/pull/206) - Added missing Web DLP rule attribute `userRiskScoreLevels`: Supported values: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`

- [PR #206](https://github.com/zscaler/zscaler-sdk-go/pull/206) - Added DLP Engine Lite endpoint to `/dlpEngines/lite`

# 2.3.1 (December 15, 2023)

## Notes
- Golang: **v1.19**

### Enhacements

- [PR #205](https://github.com/zscaler/zscaler-sdk-go/pull/205) Added ZIA Web DLP Rule new attributes:
  - `severity`
  - `subRules`
  - `parentRule`

# 2.3.0 (December 13, 2023)

## Notes
- Golang: **v1.19**

### Enhacements

- [PR #202](https://github.com/zscaler/zscaler-sdk-go/pull/202) Added support to 🆕 ZIA Cloud Browser Isolation Profile endpoint ``/browserIsolation/profiles``

# 2.2.2 (December 10, 2023)

## Notes
- Golang: **v1.19**

### Fixes

- Fixed ZPA application segment PRA for missing attribute ``UDPPortRanges``

# 2.2.1 (December 8, 2023)

## Notes
- Golang: **v1.19**

### Fixes

- Removed unsupported attributes from ZIA Forwarding control rule resource

# 2.2.0 (December xx, 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #185](https://github.com/zscaler/zscaler-sdk-go/pull/185) Added ZIA Sandbox Resources:
  - **Sandbox Quota Report** - The resource access quota for retrieving Sandbox Detail Reports is restricted to 1000 requests per day, with a rate limit of 2/sec and 1000/hour. Use GET /sandbox/report/quota to retrieve details regarding your organization's daily Sandbox API resource usage (i.e., used quota, unused quota).
  - **Sandbox Quota MD5 Hash Report** - Gets a full (i.e., complete) or summary detail report for an MD5 hash of a file that was analyzed by Sandbox.
  - **Sandbox Advanced Settings** - Gets and Upddates the custom list of MD5 file hashes that are blocked by Sandbox.
  - **Sandbox Advanced Settings Hash Count** - Gets the used and unused quota for blocking MD5 file hashes with Sandbox

- [PR #185](https://github.com/zscaler/zscaler-sdk-go/pull/185)
  - **Sandbox Submission** - Submits raw or archive files (e.g., ZIP) to Sandbox for analysis. You can submit up to 100 files per day and it supports all file types that are currently supported by Sandbox.
  - **Sandbox Out-of-Band File Inspection** - Submits raw or archive files (e.g., ZIP) to the Zscaler service for out-of-band file inspection to generate real-time verdicts for known and unknown files. It leverages capabilities such as Malware Prevention, Advanced Threat Prevention, Sandbox cloud effect, AI/ML-driven file analysis, and integrated third-party threat intelligence feeds to inspect files and classify them as benign or malicious instantaneously.

- [PR #188](https://github.com/zscaler/zscaler-sdk-go/pull/188) Added support for ZIA 🆕 Forwarding Control Policy endpoint `/forwardingRules`
- [PR #188](https://github.com/zscaler/zscaler-sdk-go/pull/188) Added support for ZIA 🆕 Custom ZPA Gateway endpoint `/zpaGateways`for use with Forwarding Control policy to forward traffic to ZPA for Source IP Anchoring
- [PR #190](https://github.com/zscaler/zscaler-sdk-go/pull/190) Added support for ZIA Group, Department and UserName using ``SortOrder`` and ``SortBy`` search criteria option
- [PR #191](https://github.com/zscaler/zscaler-sdk-go/pull/191) Added support for Zscaler Cloud & Branch Connector API endpoints. The following endpoint resources are supported:
  - `/adminRoles`
  - `/adminUsers`
  - `/ecgroup`
  - `/ecgroup/lite`
  - `/location`
  - `/location/lite`
  - `/locationTemplate`
  - `/apiKeys`
  - `/apiKeys/{keyId}/regenerate`

### Fixes

- [PR #189](https://github.com/zscaler/zscaler-sdk-go/pull/189) Fixed missing `microtenantId` and `microtenantName` attributes in ZPA browser access package.

# 2.1.6 (November 17, 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #182](https://github.com/zscaler/zscaler-sdk-go/pull/182) Added support for ZPA SCIM Group SortOrder and SortBy search criteria option
- [PR #184](https://github.com/zscaler/zscaler-sdk-go/pull/184) - Added `JSESSIONID` to every ZIA API Request

# 2.1.5 (November 1, 2023)

## Notes
- Golang: **v1.19**

### Fixes

- [PR #177](https://github.com/zscaler/zscaler-sdk-go/pull/177) Temporarily disabled Cloud Browser Isolation test edge cases to prevent some errors
- [PR #178](https://github.com/zscaler/zscaler-sdk-go/pull/178) Added missing `microtenant_id` attribute to ZPA Enrollment Certificate resource.

# 2.1.4 (October 18, 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #172](https://github.com/zscaler/zscaler-sdk-go/pull/172) Added ``GetAllSubLocations`` function to ZIA package.

# 2.1.3 (October 5, 2023)

## Notes
- Golang: **v1.19**

### Fixes

- [PR #170](https://github.com/zscaler/zscaler-sdk-go/pull/170) Fixed ZPA common custom pagination function `GetAllPagesGenericWithCustomFilters` to accommodate recent API changes on searches of objects containing multiple spaces when searching by name. Issue [#169](https://github.com/zscaler/zscaler-sdk-go/issues/169)
- [PR #171](https://github.com/zscaler/zscaler-sdk-go/pull/171) Fixed ZPA application segment PRA and Inspection to include additional attributes within the ``apps_config`` menu

# 2.1.2 (October 3, 2023)

## Notes
- Golang: **v1.19**

### Fixes

- [PR #168](https://github.com/zscaler/zscaler-sdk-go/pull/168) Restructured zia user management package directory for better organization and readability.

# 2.1.1 (September 30, 2023)

## Notes
- Golang: **v1.19**

### Fixes

- [PR #167](https://github.com/zscaler/zscaler-sdk-go/pull/167) Added ZPA LSS Config Controller `ResourceLHSRHSValue` to allow for more granular SIEM policy configuration.

# 2.1.0 (September 22, 2023)

## Notes
- Golang: **v1.19**

### Fixes

- [PR #166](https://github.com/zscaler/zscaler-sdk-go/pull/166) Added new ZIA Firewall attribute ``excludeSrcCountries``
- [PR #166](https://github.com/zscaler/zscaler-sdk-go/pull/166) Added support documents and updated README page.

# 2.1.0-beta (September 14, 2023)

## Notes
- Golang: **v1.19**

### Fixes

- [PR #164](https://github.com/zscaler/zscaler-sdk-go/pull/164) Implemented caching (BigCache) for ZIA API client.
- [PR #164](https://github.com/zscaler/zscaler-sdk-go/pull/164) Implemented detailed rate limiter per method for ZPA and ZIA API Clients. The rate limiter separates limits and frequencies for GET and other (POST, PUT, DELETE) requests for further flexibility.

# 2.0.2 (September 10, 2023)

## Notes
- Golang: **v1.19**

### Fixes

- [PR #162](https://github.com/zscaler/zscaler-sdk-go/pull/162) Fixed microtenant search criteria for ``provisioning_key``

# 2.0.0 (September 6, 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #159](https://github.com/zscaler/zscaler-sdk-go/pull/159)
  1. Zscaler Private Access Microtenant feature is now supported across the following ZPA resources:
      - ``application_controller``
      - ``app_connector_group``
      - ``application_segment``
      - ``application_segment_browser_access``
      - ``application_segment_inspection``
      - ``application_segment_pra``
      - ``app_server_controller``
      - ``machine_group``
      - ``access_policy_rule``
      - ``timeout_policy_rule``
      - ``forward_policy_rule``
      - ``inspection_policy_rule``
      - ``isolation_policy_rule``
      - ``provisioning_key``
      - ``segment_group``
      - ``server_group``
      - ``service_edge_controller``
      - ``service_edge_group``

# 1.8.0-beta (August 25, 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #153](https://github.com/zscaler/zscaler-sdk-go/pull/153)
  1. Added additional rate limit optimization
  2. Improved backoff mechanism retry strategy
  3. Updated `zpa/config.go` to use `github.com/zscaler/zscaler-sdk-go/cache` new cache mechanism to decrease number of API calls being made to the ZPA API.

⚠️ **WARNING:**: This version is being released as a Beta solution pending additional performance tests.

# 1.7.0 (August 1, 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #70](https://github.com/zscaler/zscaler-sdk-go/pull/70) Added new ZPA Microtenant Controller endpoint
``/microtenants``

- [PR #126](https://github.com/zscaler/zscaler-sdk-go/pull/126) - Added New Public ZIA DLP Engine Endpoints (POST/PUT/DELETE)

- [PR #127](https://github.com/zscaler/zscaler-sdk-go/pull/127) - Added support to the following new ZPA Cloud Browser Isolation resources:
  - Cloud Browser Isolation Banner Controller
  - Cloud Browser Isolation Certificate Controller
  - Cloud Browser Isolation Profile Controller
  - Cloud Browser Isolation Regions
  - Cloud Browser Isolation ZPA Profile

- [PR #145](https://github.com/zscaler/zscaler-sdk-go/pull/145) - Added support to ZPA GOV US Cloud. [ZPA Terraform Provider Issue#333](https://github.com/zscaler/terraform-provider-zpa/issues/333)

### Fixes

- [PR #142](https://github.com/zscaler/zscaler-sdk-go/pull/142) - Fixed filtering by email on search scim attribute values

# 1.6.4 (July, 8 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #142](https://github.com/zscaler/zscaler-sdk-go/pull/142) - Fixed filtering by email on search scim attribute values

# 1.6.3 (July, 5 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #138](https://github.com/zscaler/zscaler-sdk-go/pull/138) - Added support to ZPA QA environment
- [PR #140](https://github.com/zscaler/zscaler-sdk-go/pull/140) - Added new attribute ``waf_disabled`` to resource ``zpa_app_connector_group``

# 1.6.2 (July, 5 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #138](https://github.com/zscaler/zscaler-sdk-go/pull/138) - Added support to ZPA QA environment
- [PR #140](https://github.com/zscaler/zscaler-sdk-go/pull/140) - Added new attribute ``waf_disabled`` to resource ``zpa_app_connector_group``

# 1.6.1 (June, 21 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #140](https://github.com/zscaler/zscaler-sdk-go/pull/140) - Added new attribute ``waf_disabled`` to resource ``zpa_app_connector_group``

### Fixes

- [PR #135](https://github.com/zscaler/zscaler-sdk-go/pull/133) - Fixed ZPA Inspection Predefined Control and inspection profile resources

# 1.6.0 (June, 18 2023)

## Notes
- Golang: **v1.19**

### Fixes

- [PR #133](https://github.com/zscaler/zscaler-sdk-go/pull/133) - Included initial ZPA and ZIA integration and unit tests
- [PR #134](https://github.com/zscaler/zscaler-sdk-go/pull/134) - Included additional ZPA and ZIA integration and unit tests

# 1.5.5 (June, 10 2023)

## Notes
- Golang: **v1.19**

### Fixes

- [PR #131](https://github.com/zscaler/zscaler-sdk-go/pull/131) - Improved search mechanisms for both ZIA and ZPA resources, to ensure streamline upstream GET API requests and responses using ``search`` parameter. Notice that not all current API endpoints support the search parameter, in which case, all resources will be returned.

# 1.5.4 (June, 5 2023)

## Notes
- Golang: **v1.19**

### Fixes

- [PR #129](https://github.com/zscaler/zscaler-sdk-go/pull/129) - Added additional log information for ZIA API Client. The SDK now returns the exact authentication error message, as well as includes the ``JSESSIONID`` cookie ID information.

# 1.5.3 (May, 24 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #127](https://github.com/zscaler/zscaler-sdk-go/pull/127) - Fixed ZPA resource ``Service Edge Group`` and ``Service Edge Controller`` Struct to support attribute ``publish_ips``.

# 1.5.2 (May, 23 2023)

## Notes
- Golang: **v1.19**

### Fixes

- [PR #125](https://github.com/zscaler/zscaler-sdk-go/pull/125) - Added exception handling within the ZPA API Client to deal with simultaneous DB requests, which were affecting the ZPA Policy Access rule order creation.
  - Internal References:
    - [ET-53585](https://jira.corp.zscaler.com/browse/ET-53585)
    - [ET-48860](https://confluence.corp.zscaler.com/display/ET/ET-48860+incorrect+rules+order)

# 1.5.0 (May, 15 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #124](https://github.com/zscaler/zscaler-sdk-go/pull/124) Added ZIA DLP Exact Data Match Schema endpoints

# 1.4.7 (May, 13 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #123](https://github.com/zscaler/zscaler-sdk-go/pull/123) Improve SCIM Attribute Header search function

# 1.4.6 (May, 11 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #120](https://github.com/zscaler/zscaler-sdk-go/pull/120) Added new ZPA App Inspection Security Profiles attributes. The following new attributes have been added:
  - ``checkControlDeploymentStatus`` - Bool
  - ``controlsInfo`` - String. Support values: ``WEBSOCKET_PREDEFINED``, ``WEBSOCKET_CUSTOM``, ``THREATLABZ``, ``CUSTOM``, ``PREDEFINED``
  - ``threatlabzControls`` - List
  - ``zsDefinedControlChoice`` - String. Support values: ``ALL`` and ``SPECIFIC``
- [PR #121](https://github.com/zscaler/zscaler-sdk-go/pull/121) Added new ZPA Client Type ``zpn_client_type_zapp_partner``

### Bug Fixes

- [PR #122](https://github.com/zscaler/zscaler-sdk-go/pull/122) Fixed issue with empty IDs in the resource ``zpa_service_edge_groups``

# 1.4.5 (April, 29 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #118](https://github.com/zscaler/zscaler-sdk-go/pull/118) Added new ZIA DLP Dictionary attributes. The following new attributes have been added:
  - ``ignoreExactMatchIdmDict`` - Bool: Indicates whether to exclude documents that are a 100% match to already-indexed documents from triggering an Indexed Document Match (IDM) Dictionary.
  - ``includeBinNumbers`` - Bool: A true value denotes that the specified Bank Identification Number (BIN) values are included in the Credit Cards dictionary. A false value denotes that the specified BIN values are excluded from the Credit Cards dictionary. Note: This field is applicable only to the predefined Credit Cards dictionary and its clones.
  - ``binNumbers`` - []int: The list of Bank Identification Number (BIN) values that are included or excluded from the Credit Cards dictionary. BIN values can be specified only for Diners Club, Mastercard, RuPay, and Visa cards. Up to 512 BIN values can be configured in a dictionary. Note: This field is applicable only to the predefined Credit Cards dictionary and its clones.
  - ``dictTemplateId`` - int: ID of the predefined dictionary (original source dictionary) that is used for cloning. This field is applicable only to cloned dictionaries. Only a limited set of identification-based predefined dictionaries (e.g., Credit Cards, Social Security Numbers, National Identification Numbers, etc.) can be cloned. Up to 4 clones can be created from a predefined dictionary.
  - ``predefinedClone`` - bool: This field is set to true if the dictionary is cloned from a predefined dictionary. Otherwise, it is set to false.
  - ``proximityLengthEnabled`` - bool: This value is set to true if proximity length and high confidence phrases are enabled for the DLP dictionary.

# 1.4.4 (April, 29 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #117](https://github.com/zscaler/zscaler-sdk-go/pull/117) Fix ZIA DLP dictionary attribute ``idmProfileMatchAccuracyDetails``

# 1.4.3 (April, 28 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #114](https://github.com/zscaler/zscaler-sdk-go/pull/114) Expanded ZIA search criteria to include auditor users.
- [PR #115](https://github.com/zscaler/zscaler-sdk-go/pull/115) Fixed empty ZPA body response in case of 400 Errors
- [PR #116](https://github.com/zscaler/zscaler-sdk-go/pull/116) Fixed typo in ZIA DLP Web Rule for the attribute ``zscalerIncidentReceiver``

# 1.4.2 (April, 27 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #113](https://github.com/zscaler/zscaler-sdk-go/pull/113) Fixed ZPA Empty policy conditions or operands on update due to 500 errors

# 1.4.1 (April, 17 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #109](https://github.com/zscaler/zscaler-sdk-go/pull/109) Added ZIA DLP IDM Lite endpoints to obtain summarized information about existing IDM profiles.
- [PR #110](https://github.com/zscaler/zscaler-sdk-go/pull/110) Added extra fix for ZIA API Client to prevent SESSION_INVALID error during session timeout. The client will re-authenticate automaticallyu to renew the session.

# 1.4.0 (April, 10 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #99](https://github.com/zscaler/zscaler-sdk-go/pull/99): Fixed ZIA API Client to log the user-agent information during debug
- [PR #102](https://github.com/zscaler/zscaler-sdk-go/pull/102): Log request ID and API call duration for each request
- [PR #104](https://github.com/zscaler/zscaler-sdk-go/pull/104): Removed lock client function on ZPA client package
- [PR #106](https://github.com/zscaler/zscaler-sdk-go/pull/106): Removed lock client function on all other API clients for ZCC, ZIA, and ZDX packages
- [PR #107](https://github.com/zscaler/zscaler-sdk-go/pull/107): Implementyed refresh expired session for long requests on the ZIA API client
- [PR #108](https://github.com/zscaler/zscaler-sdk-go/pull/108): Allow updating application segment access policy groups with empty list

### Bug Fixes

- [PR #105](https://github.com/zscaler/zscaler-sdk-go/pull/105): Added function to temporarily handle ZPA upstream bad request errors.

# 1.3.5 (April, 7 2023)

## Notes
- Golang: **v1.19**

### Enhancements
- [PR #99](https://github.com/zscaler/zscaler-sdk-go/pull/99): Fixed ZIA API Client to log the user-agent information during debug
- [PR #102](https://github.com/zscaler/zscaler-sdk-go/pull/102): Log request ID and API call duration for each request
- [PR #104](https://github.com/zscaler/zscaler-sdk-go/pull/104): Removed lock client function on ZPA client package

### Bug Fixes

- [PR #105](https://github.com/zscaler/zscaler-sdk-go/pull/105): Added function to temporarily handle ZPA upstream bad request errors.

# 1.3.4 (March, 29 2023)

## Notes
- Golang: **v1.19**

### Bug Fixes

- [PR #105](https://github.com/zscaler/zscaler-sdk-go/pull/105): Added function to temporarily handle ZPA upstream bad request errors.

# 1.3.3 (March, 28 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #98](https://github.com/zscaler/zscaler-sdk-go/pull/98) Added support to Get predefined DLP engines by name and set name to ``predefinedEngineName``

# 1.3.2 (March, 27 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #97](https://github.com/zscaler/zscaler-sdk-go/pull/97) Fixed ZIA GRE Tunnel attributes.
  - Make WithinCountry a pointer for GRE Tunnel response
  - City, Region, Latitude & Longitude to VIP response
  - Implement get all by source IP & get all VIPs by all existing source IPs

# 1.3.1 (March, 25 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #93](https://github.com/zscaler/zscaler-sdk-go/pull/93) The ZIA SDK now supports search of Sublocations by Name and ID.

# 1.3.0 (March, 22 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #89](https://github.com/zscaler/zscaler-sdk-go/pull/89) The ZPA Terraform Provider API Client, will now support long runs, that exceeds the 3600 seconds token validity. Terraform will automatically request a new API bearer token at that time in order to continue the resource provisioning. This enhacement will prevent long pipeline runs from being interrupted.

- [PR #92](https://github.com/zscaler/zscaler-sdk-go/pull/92) Added ZIA Location Management Lite endpoint.

# 1.2.5 (March, 20 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #86](https://github.com/zscaler/zscaler-sdk-go/pull/86) Added new ZPA IDP Controller attributes. The following new attributes have been added:
  - ``enableArbitraryAuthDomains``
  - ``forceAuth``
  - ``loginHint``

# 1.2.4 (March, 18 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #85](https://github.com/zscaler/zscaler-sdk-go/pull/85) Added new ZIA Location Management attributes. The following new attributes have been added:
  - ``basicAuthEnabled``: Enable Basic Authentication at the location
  - ``digestAuthEnabled``: Enable Digest Authentication at the location
  - ``kerberosAuth``: Enable Kerberos Authentication at the location
  - ``iotDiscoveryEnabled``: Enable IOT Discovery at the location

# 1.2.3 (March, 16 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #83](https://github.com/zscaler/zscaler-sdk-go/pull/83) Added new ZPA platform and clienttype endpoints:
  - ``/platform``
  - ``/clientTypes``

# 1.2.2 (March, 11 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #73](https://github.com/zscaler/zscaler-sdk-go/pull/73) Added support to ZIA Basic Authentication method to the following resources:
  - Location Management
    - ``basicAuthEnabled`` - (Optional) - ``Bool``

  - User Management
    - Added new endpoint ``/enroll`` which is called when the ``authMethods`` attribute is set.
    - ``authMethods`` - (Optional) - ``String``. Supported values are: ["BASIC", "DIGEST"]

# 1.2.1 (March, 7 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #82](https://github.com/zscaler/zscaler-sdk-go/pull/82) Added the new ZPA API attributes below:
  - Application Segments
    - ``tcpKeepAlive``
    - ``isIncompleteDRConfig``
    - ``useInDrMode``
    - ``selectConnectorCloseToApp``

  - App Connector Group
    - ``useInDrMode``

# 1.2.0 (March, 6 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #78](https://github.com/zscaler/zscaler-sdk-go/pull/78) AAdded support to Zscaler Digital Experience (ZDX) API.

# 1.1.3 (February, 28 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #76](https://github.com/zscaler/zscaler-sdk-go/pull/76) Added search by Source IP function to ZIA GRE Tunnel
- [PR #76](https://github.com/zscaler/zscaler-sdk-go/pull/76) Added description to all struct attributes

# 1.1.2 (February, 28 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #76](https://github.com/zscaler/zscaler-sdk-go/pull/76) Added search by Source IP function to ZIA GRE Tunnel
- [PR #76](https://github.com/zscaler/zscaler-sdk-go/pull/76) Added description to all struct attributes

# 1.1.1 (February, 24 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #74](https://github.com/zscaler/zscaler-sdk-go/pull/74) Added ZIA endpoint ``/appServices/lite`` to retrieve supported application services within an firewall filtering rule resource
- [PR #74](https://github.com/zscaler/zscaler-sdk-go/pull/74) Added ZIA endpoint ``/appServiceGroups/lite`` to retrieve supported application services groups within an firewall filtering rule resource

# 1.1.0 (February, 24 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #74](https://github.com/zscaler/zscaler-sdk-go/pull/74) Added ZIA endpoint ``/appServices/lite`` to retrieve supported application services within an firewall filtering rule resource
- [PR #74](https://github.com/zscaler/zscaler-sdk-go/pull/74) Added ZIA endpoint ``/appServiceGroups/lite`` to retrieve supported application services groups within an firewall filtering rule resource

# 1.0.0 (February, 2 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #69](https://github.com/zscaler/zscaler-sdk-go/pull/69) Added new ZPA Isolation Profile Controller endpoint ``/isolation/profiles``

# 0.7.0 (January, 31 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #68](https://github.com/zscaler/zscaler-sdk-go/pull/68) Added the following ZIA DLP endpoint resources:
  - ``dlp_icap_servers`` - /icapServers
  - ``dlp_incident_receiver_servers`` - /incidentReceiverServers
  - ``dlp_idm_profiles`` - /idmprofile

# 0.6.1 (January, 13 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #67](https://github.com/zscaler/zscaler-sdk-go/pull/67) Added ``omitempty`` bool parameters in the ZIA URL Firewall Filtering resource ``enable_full_logging``


# 0.6.0 (January, 12 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #65](https://github.com/zscaler/zscaler-sdk-go/pull/65) Fixed pagination issue with ZIA API endpoints

# 0.5.9 (January, 12 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #65](https://github.com/zscaler/zscaler-sdk-go/pull/65) Fixed pagination issue with ZIA API endpoints

# 0.5.8 (January, 11 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #64](https://github.com/zscaler/zscaler-sdk-go/pull/64) Added new ZPA Inspection control parameters

  - ZPA Inspection Profile: ``web_socket_controls``
  - ZPA Custom Inspection Control:
    - ``control_type``: The following values are supported:
      - ``WEBSOCKET_PREDEFINED``, ``WEBSOCKET_CUSTOM``, ``ZSCALER``, ``CUSTOM``, ``PREDEFINED``

    - ``protocol_type``: The following values are supported:
      - ``HTTP``, ``WEBSOCKET_CUSTOM``, ``ZSCALER``, ``CUSTOM``, ``PREDEFINED``

# 0.5.7 (January, 4 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #63](https://github.com/zscaler/zscaler-sdk-go/pull/63) Added ``omitempty`` bool parameters in the ZIA URL Firewall Filtering resource ``enable_full_logging``

# 0.5.6 (January, 4 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #61](https://github.com/zscaler/zscaler-sdk-go/pull/61) Added ``omitempty`` bool parameters in the ZIA URL Filtering Policy resource

# 0.5.5 (December, 30 2022)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #60](https://github.com/zscaler/zscaler-sdk-go/pull/60) Added new ZIA URL Filtering rule URL category parameters to Struct

# 0.5.4 (December, 30 2022)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #60](https://github.com/zscaler/zscaler-sdk-go/pull/60) Added new ZIA URL Filtering rule URL category parameters to Struct

# 0.5.3 (December, 30 2022)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #60](https://github.com/zscaler/zscaler-sdk-go/pull/60) Added new ZIA URL Filtering rule URL category parameters to Struct

# 0.5.2 (December, 27 2022)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #59](https://github.com/zscaler/zscaler-sdk-go/pull/59) Added new ZIA URL Category parameters to Struct

# 0.5.1 (December, 17 2022)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #57](https://github.com/zscaler/zscaler-sdk-go/pull/57) Added new ZPA application segment paramenter ``select_connector_close_to_app`` to Struct

# 0.5.0 (December, 16 2022)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #56](https://github.com/zscaler/zscaler-sdk-go/pull/56) Added new Intermediate CA Certificate Endpoints for ZIA
- [PR #56](https://github.com/zscaler/zscaler-sdk-go/pull/56) Added new Event Log Entry Report Endpoints for ZIA
- [PR #56](https://github.com/zscaler/zscaler-sdk-go/pull/56) Added new Location Management IPv6 Parameters

# 0.4.1 (December, 02 2022)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #53](https://github.com/zscaler/zscaler-sdk-go/pull/53) Fixed pagination issue with ZPA endpoints

# 0.4.0 (December, 01 2022)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #52](https://github.com/zscaler/zscaler-sdk-go/pull/52) Added new ZIA intermediate CA cert endpoints

# 0.3.1 (November, 30 2022)

## Notes
- Golang: **v1.18**

### Fix

- [PR #50](https://github.com/zscaler/zscaler-sdk-go/pull/50) Changed common function to allow totalPages string values

# 0.3.0 (November, 29 2022)

## Notes
- Golang: **v1.18**

### Enhancement

- [PR #49](https://github.com/zscaler/zscaler-sdk-go/pull/49) Implemented generic function to get all resources
- [PR #49](https://github.com/zscaler/zscaler-sdk-go/pull/49) Implemented generic function to get all SCIM header attribute values

# 0.2.2 (November, 24 2022)

## Notes
- Golang: **v1.18**

### Enhancement

- [PR #44](https://github.com/zscaler/zscaler-sdk-go/pull/44) Added parameter cert_blob for zpa_ba_certificate

# 0.2.1 (November, 24 2022)

## Notes
- Golang: **v1.18**

### Enhancement

- [PR #44](https://github.com/zscaler/zscaler-sdk-go/pull/44) Added parameter cert_blob for zpa_ba_certificate

# 0.2.0 (November, 24 2022)

## Notes
- Golang: **v1.18**

### Enhancement

- [PR #44](https://github.com/zscaler/zscaler-sdk-go/pull/44) Added parameter cert_blob for zpa_ba_certificate

# 0.1.9 (November, 15 2022)

## Notes
- Golang: **v1.18**

### Enhancement

- Add associationType json to prov key

# 0.1.8 (November, 15 2022)

## Notes
- Golang: **v1.18**

### Enhancement

- [PR #43](https://github.com/zscaler/zscaler-sdk-go/pull/43) Return AssociationType in provisioning key endpoints

# 0.1.7 (November, 13 2022)

## Notes
- Golang: **v1.18**

### Enhancement

- [PR #41](https://github.com/zscaler/zscaler-sdk-go/pull/41) Allow order 0 for firewall filtering rules in ZIA cloud firewall.

# 0.1.6 (October, 22 2022)

## Notes
- Golang: **v1.18**

### Enhancement

- [PR #37](https://github.com/zscaler/zscaler-sdk-go/pull/37) Implement fix on update function for App Connector Controller resource

# 0.1.5 (October, 21 2022)

## Notes
- Golang: **v1.18**

### Enhancement

- [PR #36](https://github.com/zscaler/zscaler-sdk-go/pull/36) Implement bulk delete of service-edge-controller

# 0.1.4 (October, 21 2022)

## Notes
- Golang: **v1.18**

### Enhancement

- [PR #35](https://github.com/zscaler/zscaler-sdk-go/pull/35) Implement bulk delete of app-connector-controller

# 0.1.3 (October, 20 2022)

## Notes
- Golang: **v1.18**

### Enhancement

- [PR #34](https://github.com/zscaler/zscaler-sdk-go/pull/34) Added new application segment parameter ``forceDelete``. Setting this field to true deletes the mapping between Application Segment and Segment Group

# 0.1.2 (October, 19 2022)

## Notes
- Golang: **v1.18**

### Bug Fix
- [PR #33](https://github.com/zscaler/zscaler-sdk-go/pull/33) Fix Added ZPA missing parameters

# 0.1.1 (October, 15 2022)

## Notes
- Golang: **v1.18**

### Enhancements

- [PR #30](https://github.com/zscaler/zscaler-sdk-go/pull/30) feat(ZPA Application Segments): Filters application segments apps in GetAll API calls depending on the resource type (SECURE_REMOTE_ACCESS, BROWSER_ACCESS, INSPECTION).
- [PR #32](https://github.com/zscaler/zscaler-sdk-go/pull/32) feat(Improve Logging): This PR improves logging for the SDK for all clouds (zia, zpa & zcc) and uses common logger, we now can control the logging & it verbosity using the env var:
ZSCALER_SDK_LOG=true & ZSCALER_SDK_VERBOSE=true

# 0.1.0 (October, 12 2022)

## Notes
- Golang: **v1.18**

### Enhancements

- [PR #29 ](https://github.com/zscaler/zscaler-sdk-go/pull/29) feat(New SDK Package): Added Zscaler Client Connector (ZCC) SDK Schema
- [PR #30  ](https://github.com/zscaler/zscaler-sdk-go/pull/30) feat(ZPA Application Segments): Filters application segments apps in GetAll API calls depending on the resource type (SECURE_REMOTE_ACCESS, BROWSER_ACCESS, INSPECTION).

# 0.0.13 (September, 28 2022)

## Notes
- Golang: **v1.18**

### Enhancements

- [PR #26](https://github.com/zscaler/zscaler-sdk-go/pull/26) feat(new parameters):App Connector Group TCPQuick*
- The following new App Connector Group parameters have been added to the SDK:
  - tcpQuickAckApp - Whether TCP Quick Acknowledgement is enabled or disabled for the application.
  - tcpQuickAckAssistant - Whether TCP Quick Acknowledgement is enabled or disabled for the application.
  - tcpQuickAckReadAssistant - Whether TCP Quick Acknowledgement is enabled or disabled for the application.

# 0.0.12 (September, 28 2022)

## Notes
- Golang: **v1.18**

### Enhancements

- [PR #26](https://github.com/zscaler/zscaler-sdk-go/pull/26) feat(new parameters):App Connector Group TCPQuick*
- The following new App Connector Group parameters have been added to the SDK:
  - tcpQuickAckApp - Whether TCP Quick Acknowledgement is enabled or disabled for the application.
  - tcpQuickAckAssistant - Whether TCP Quick Acknowledgement is enabled or disabled for the application.
  - tcpQuickAckReadAssistant - Whether TCP Quick Acknowledgement is enabled or disabled for the application.

# 0.0.11 (September, 26 2022)

## Notes
- Golang: **v1.18**

### Bug Fix
- [PR #25](https://github.com/zscaler/zscaler-sdk-go/pull/25) Fix zia_user_management group attribute to hold a list of group IDs as a typeList instead of typeSet.

# 0.0.10 (September, 21 2022)

## Notes
- Golang: **v1.18**

### Bug Fix
- [PR #23](https://github.com/zscaler/zscaler-sdk-go/pull/23) Fix zia_user_management group attribute to hold a list of group IDs as a typeList instead of typeSet.

# 0.0.9 (September, 10 2022)

## Notes
- Golang: **v1.18**

### Enhancement
- [PR #20](https://github.com/zscaler/zscaler-sdk-go/pull/20) Added Support to ZPA Preview CLOUD.

# 0.0.8 (September, 2 2022)

## Notes
- Golang: **v1.18**

### Bug Fix
- [PR #18](https://github.com/zscaler/zscaler-sdk-go/pull/18) Fixed ZPA_CLOUD support for production via environment variables.

# 0.0.7 (August, 30 2022)

## Notes
- Golang: **v1.18**

### Enhancements
- [PR #11](https://github.com/zscaler/zscaler-sdk-go/pull/11) Added support to getAll method for LSS config ctl & policy ctl
- [PR #15](https://github.com/zscaler/zscaler-sdk-go/pull/15) Added support for ZPA arbitrary clouds @hfinucane
- [PR #16](https://github.com/zscaler/zscaler-sdk-go/pull/16) Added support to ZPA API response with html double escaping

## 0.0.6 (August, 29 2022)

### Notes

- Golang Version: **v1.18.x**

### Enhancements

- [PR #11](https://github.com/zscaler/zscaler-sdk-go/pull/11) Added support to getAll method for LSS config ctl & policy ctl
- [PR #15](https://github.com/zscaler/zscaler-sdk-go/pull/15) Added support for ZPA arbitrary clouds @hfinucane
- [PR #16](https://github.com/zscaler/zscaler-sdk-go/pull/16) Added support to ZPA API response with html double escaping


## 0.0.5 (July, 30 2022)

### Notes

- Golang Version: **v1.18.x**

### Bug Fixes

- Fixed typo

## 0.0.4 (July, 28 2022)

### Notes

- Golang Version: **v1.18.x**

### Ehancements

- Added support to GetAll methods to all resources

## 0.0.1 (July, 24 2022)

### Notes

- Golang Version: **v1.18.x**

🎉 **Initial Release** 🎉