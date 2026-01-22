---
layout: "zscaler"
page_title: "Release Notes"
description: |-
  The Zscaler SDK GO Release Notes
---

# Zscaler SDK GO: Release Notes

## USAGE

Track all Zscaler SDK GO releases. New resources, features, and bug fixes will be tracked here.

---

``Last updated: v3.8.14``

---

# 3.8.14 (January 21, 2026)

## Notes
- Golang: **v1.24**

### Enhancements

- [PR #401](https://github.com/zscaler/zscaler-sdk-go/pull/401) - Fixed `GetAll` function in `ssl_inspection_rules` zia package by removing unsupported pagination parameters.

# 3.8.13 (January 19, 2026)

## Notes
- Golang: **v1.24**

### Enhancements

- [PR #400](https://github.com/zscaler/zscaler-sdk-go/pull/400) - Fixed `GetAll` function in forwarding_rules ztw package to support new optional parameters.

# 3.8.12 (January 13, 2026)

## Notes
- Golang: **v1.24**

### Enhancements

- [PR #398](https://github.com/zscaler/zscaler-sdk-go/pull/398) - Fixed ZTW Legacy Client environment variables


# 3.8.11 (December 16, 2025)

## Notes
- Golang: **v1.24**

### Enhancements

- [PR #396](https://github.com/zscaler/zscaler-sdk-go/pull/396) - Added new `tags` field to ZPA `applicationsegment`
- [PR #396](https://github.com/zscaler/zscaler-sdk-go/pull/396) - Improved several unit tests across all packages.

# 3.8.10 (December 10, 2025)

## Notes
- Golang: **v1.24**

### Bug Fixes

- [PR #394](https://github.com/zscaler/zscaler-sdk-go/pull/394) - Added new attribute `approvalReviewers` to ZPA `praportal` package.

# 3.8.9 (December 9, 2025)

## Notes
- Golang: **v1.24**

### Bug Fixes

- [PR #393](https://github.com/zscaler/zscaler-sdk-go/pull/393) - Fixed `UpdateURLCategories` function in ZIA `url_category` package to include optional parameters `action` to support partial updates via `ADD_TO_LIST` and `REMOVE_FROM_LIST`.

# 3.8.8 (December 1, 2025)

## Notes
- Golang: **v1.24**

### Bug Fixes

- [PR #391](https://github.com/zscaler/zscaler-sdk-go/pull/391) - Fixed ZTW `provisioning_url` and `location_template` struct resources.

# 3.8.7 (November 21, 2025)

## Notes
- Golang: **v1.24**

### Enhancements

- [PR #388](https://github.com/zscaler/zscaler-sdk-go/pull/388) - Added automatic `x-partner-id` header injection for all API requests when `partnerId` is provided in configuration across OneAPI and all legacy clients (ZIA, ZPA, ZTW, ZCC, ZDX, ZWA)
- [PR #388](https://github.com/zscaler/zscaler-sdk-go/pull/388) - Added `GetWeightedLoadBalancerConfig` and `UpdateWeightedLoadBalancerConfig` functions for ZPA application segments
- [PR #388](https://github.com/zscaler/zscaler-sdk-go/pull/388) - Added optional filter parameters to ZIA location groups `GetAll` function and `GetLocationGroupCount` function
- [PR #388](https://github.com/zscaler/zscaler-sdk-go/pull/388) - Added `GetLocationSupportedCountries` function to retrieve list of supported countries for location configuration
- [PR #388](https://github.com/zscaler/zscaler-sdk-go/pull/388) - Added optional filter parameters to ZIA location lite `GetAll` function and updated struct with sublocation scope fields
- [PR #388](https://github.com/zscaler/zscaler-sdk-go/pull/388) - Added `GetCustomFileTypeCount` function with optional query filter parameter
- [PR #388](https://github.com/zscaler/zscaler-sdk-go/pull/388) - Added optional filter parameters to `GetFileTypeCategories` function (enums, excludeCustomFileTypes)
- [PR #388](https://github.com/zscaler/zscaler-sdk-go/pull/388) - Added `type` parameter to URL categories `GetAll` and `GetCustomURLCategories` to support filtering by category type (`ALL`, `URL_CATEGORY`, `TLD_CATEGORY`)
- [PR #388](https://github.com/zscaler/zscaler-sdk-go/pull/388) - Added optional filter parameters to traffic capture `GetAll`, `GetByName`, and firewall filtering rules functions
- [PR #388](https://github.com/zscaler/zscaler-sdk-go/pull/388) - Added `GetTrafficCaptureRuleOrder`, `GetTrafficCaptureRuleLabels`, and `GetTrafficCaptureRuleCount` functions
- [PR #388](https://github.com/zscaler/zscaler-sdk-go/pull/388) - Added `GetFirewallFilteringRuleCount` function with support for all optional filter parameters
- [PR #388](https://github.com/zscaler/zscaler-sdk-go/pull/388) - Added `excludeType` parameter to IP destination groups `GetAll` function

### Bug Fixes

- [PR #388](https://github.com/zscaler/zscaler-sdk-go/pull/388) - Fixed ZIA location management and VPN credentials pagination to use 1000 max page size to prevent API errors
- [PR #388](https://github.com/zscaler/zscaler-sdk-go/pull/388) - Updated ZIA `common.ReadAllPages` default page size from 5000 to 1000 with support for custom page sizes
- [PR #388](https://github.com/zscaler/zscaler-sdk-go/pull/388) - Fixed URL categories and IP destination groups to remove pagination since APIs don't support it

# 3.8.6 (November 19, 2025)

## Notes
- Golang: **v1.24**

### Enhancements

- [PR #386](https://github.com/zscaler/zscaler-sdk-go/pull/386) - Added support to the following ZIA Endpoints:
    - Added `GET /customFileTypes` Retrieves the list of custom file types. Custom file types can be configured as rule conditions in different ZIA policies.
    - Added `POST /customFileTypes` Adds a new custom file type. 
    - Added `PUT /customFileTypes` Updates information for a custom file type based on the specified ID
    - Added `DELETE /customFileTypes/{id}` Deletes a custom file type based on the specified ID
    - Added `GET /customFileTypes/count` Retrieves the count of custom file types available
    - Added `GET /fileTypeCategories` Retrieves the list of all file types, including predefined and custom file types

### New ZIA Endpoint - Traffic Capture Policy

- [PR #386](https://github.com/zscaler/zscaler-sdk-go/pull/386) - Added the following new ZIA Endpoints
    - Added `GET /trafficCaptureRules` Retrieves the list of Traffic Capture policy rules
    - Added `GET /trafficCaptureRules/{ruleId}` Retrieves the Traffic Capture policy rule based on the specified rule ID
    - Added `PUT /trafficCaptureRules/{ruleId}` Updates information for the Traffic Capture policy rule based on the specified rule ID
    - Added `DELETE /trafficCaptureRules/{ruleId}` Deletes the Traffic Capture policy rule based on the specified rule ID
    - Added `GET /trafficCaptureRules/count` Retrieves the rule count for Traffic Capture policy based on the specified search criteria
    - Added `GET /trafficCaptureRules/order` Retrieves the rule order information for the Traffic Capture policy
    - Added `GET /trafficCaptureRules/ruleLabels` Retrieves the list of rule labels associated with the Traffic Capture policy rules

# 3.8.5 (November 12, 2025)

## Notes
- Golang: **v1.24**

### Enhancements

- [PR #385](https://github.com/zscaler/zscaler-sdk-go/pull/385) - Added new ZPA `service_edge_group` attributes `exclusiveForBusinessContinuity`, `city` and `nameWithoutTrim`

# 3.8.4 (November 11, 2025)

## Notes
- Golang: **v1.24**

### Bug Fixes

- [PR #383](https://github.com/zscaler/zscaler-sdk-go/pull/383) - Added automatic `x-partner-id` header injection for all API requests when `partnerId` is provided in the configuration. The header is automatically included in all requests across OneAPI and Legacy clients (ZIA, ZPA, ZTW, ZCC, ZDX, ZWA) when `partnerId` is specified via config dictionary or `ZSCALER_PARTNER_ID` environment variable.

- [PR #383](https://github.com/zscaler/zscaler-sdk-go/pull/383)- Added the following ZPA Endpoints:
    - Added `GET /weightedLbConfig` Get Weighted Load Balancer Config for AppSegment
    - Added `PUT /weightedLbConfig` Update Weighted Load Balancer Config for AppSegment

# 3.8.3 (November 6, 2025)

## Notes
- Golang: **v1.24**

### Bug Fixes

- [PR #381](https://github.com/zscaler/zscaler-sdk-go/pull/381) - Fixed SCIM and SAML attribute endpoints to use plain search strings instead of filter format, and improved URL encoding for ZPA endpoints to use `%20` for spaces instead of `+` to match API requirements 

# 3.8.2 (November 5, 2025)

## Notes
- Golang: **v1.24**

### Bug Fixes

- [PR #380](https://github.com/zscaler/zscaler-sdk-go/pull/380) - Fixed ZPA search functionality to automatically convert simple search strings to API filter format (`name+EQ+<value>`) to prevent `filtering.input.invalid.operand` errors when searching for resources with multi-word names


# 3.8.1 (November 5, 2025)

## Notes
- Golang: **v1.24**

### Bug Fixes

- [PR #380](https://github.com/zscaler/zscaler-sdk-go/pull/380) - Fixed ZPA search functionality to automatically convert simple search strings to API filter format (`name+EQ+<value>`) to prevent `filtering.input.invalid.operand` errors when searching for resources with multi-word names

# 3.8.0 (October 31, 2025)

## Notes
- Golang: **v1.24**

### New ZPA Endpoint - Application Server Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /server/summary` Get all the configured application servers Name and IDs

### New ZPA Endpoint - Application Segment Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /application/{applicationId}/mappings` Get the application segment mapping details
    - Added `DELETE /application/{applicationId}/deleteAppByType` Delete a BA/Inspection and PRA Application
    - Added `POST /application/{applicationId}/validate` Validate conflicting wildcard domain names. Expect the applicationID to be populated in the case of update
    - Added `GET /application/configured/count` Returns the count of configured application Segment for the provided customer between the date range passed in request body.
    - Added `GET /application/count/currentAndMaxLimit` get current Applications count of domains and maxLimit configured for a given customer

### New ZPA Endpoint - App Connector Group

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /appConnectorGroup/summary` Get all the configured App Connector Group id and name.

### New ZPA Endpoint - Branch Connector Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /branchConnector` Get all BranchConnectors configured for a given customer.

### New ZPA Endpoint - Branch Connector Group Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /branchConnectorGroup/summary` Get all branch connector group id and names configured for a given customer.
    - Added `GET /branchConnectorGroup` Get all configured Branch Connector Groups.

### New ZPA Endpoint - Browser Protection Profile Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /activeBrowserProtectionProfile` Get the active browser protection profile details for the specified customer.
    - Added `GET /browserProtectionProfile` Gets all configured browser protection profiles for the specified customer.
    - Added `PUT /browserProtectionProfile/setActive/{browserProtectionProfileId}` Updates a specified browser protection profile as active for the specified customer.

### New ZPA Endpoint - Customer Config Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /config/isZiaCloudConfigAvailable` Check if zia cloud config for a given customer is available.
    - Added `GET /config/ziaCloudConfig` Get zia cloud service config for a given customer.
    - Added `POST /config/ziaCloudConfig` Add or update zia cloud service config for a given customer.
    - Added `GET /sessionTerminationOnReauth` Get session termination on reauth for a given customer.
    - Added `PUT /sessionTerminationOnReauth` Add /update boolean value for session termination on reauth.

### New ZPA Endpoint - Customer DR Tool Version Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /customerDRToolVersion` Fetch latest the Customer Support DR Tool Versions sorted by latest filter

### New ZPA Endpoint - Customer Version Profile Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /versionProfiles/{versionProfileId}` Update Version Profile for customer

### New ZPA Endpoint - Cloud Connector Group Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /cloudConnectorGroup/summary` Get all edge connector group id and names configured for a given customer

### New ZPA Endpoint - Extranet Resource Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /extranetResource/partner` Get all extranet resources

### New ZPA Endpoint - Machine Group Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /machineGroup/summary` Get all Machine Group Id and Names configured for a given customer

### New ZPA Endpoint - Managed Browser Profile Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /managedBrowserProfile/search` Gets all the managed browser profiles for a customer

### New ZPA Endpoint - Provisioning Key Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /associationType/{associationType}/zcomponent/{zcomponentId}/provisioningKey` get provisioningKey details by zcomponentId for associationType.

### New ZPA Endpoint - OAuth User Code Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `POST /{associationType}/usercodes` Verifies the provided list of user codes for a given component provisioning.
    - Added `POST /{associationType}/usercodes/status` Adds a new Provisioning Key for the specified customer.

### New ZPA Endpoint - Policy-Set Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /riskScoreValues` Gets values of risk scores for the specified customer.
    - Added `GET /policySet/rules/policyType/{policyType}/count` For a customer, get count of policy rules for a given policy type. Providing only endtime would give cumulative count till the endTime.Providing both startTime and endtime would give count between that time period.Not Providing startTime and endtime would give overall count.
    - Added `GET /policySet/rules/policyType/{policyType/application/{applicationId}` Gets paginated policy rules for the specified policy type by application id

### New ZPA Endpoint - Server Group Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /serverGroup/summary` Get all Server Group id and names configured for a given customer

### New ZPA Endpoint - Step up Auth Level Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /stepupauthlevel/summary` Get a step up auth levels.

### New ZPA Endpoint - Step up Auth Level Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /userportal/aup/{id}` Get user portal aup
    - Added `PUT /userportal/aup/{id}` Update user portal aup
    - Added `DELETE /userportal/aup/{id}` Delete user portal aup
    - Added `GET /userportal/aup` Get all AUPs configured for a given customer
    - Added `POST /userportal/aup` Add a new aup for a given customer.

### New ZPA Endpoint - ZPN Location Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /location/extranetResource/{zpnErId}`
    - Added `PUT /location/summary` Get all Location id and names configured for a given customer.

### New ZPA Endpoint - ZPN Location Group Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /locationGroup/extranetResource/{zpnErId}`

### New ZPA Endpoint - Workload Tag Group Controller

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /workloadTagGroup/summary`

### New ZTW Endpoint - Partner Integrations - Public Account Info

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /publicCloudInfo` - Retrieves the list of AWS accounts with metadata
    - Added `POST /publicCloudInfo` - Creates a new AWS account with the provided account and region details.
    - Added `GET /publicCloudInfo/cloudFormationTemplate` - Retrieves the CloudFormation template URL.
    - Added `GET /publicCloudInfo/count` - Retrieves the total number of AWS accounts.
    - Added `POST /publicCloudInfo/generateExternalId` - Creates an external ID for an AWS account.
    - Added `GET /publicCloudInfo/lite` - Retrieves basic information about the AWS cloud accounts
    - Added `GET /publicCloudInfo/supportedRegions` - Retrieves a list of AWS regions supported for workload discovery settings (WDS).
    - Added `GET /publicCloudInfo/{id}` - Retrieves the existing AWS account details based on the provided ID.
    - Added `PUT /publicCloudInfo/{id}` - Updates the existing AWS account details based on the provided ID.
    - Added `DELETE /publicCloudInfo/{id}` - Removes a specific AWS account based on the provided ID.
    - Added `DELETE /publicCloudInfo/{id}/changeState` - Enables or disables a specific AWS account in all regions based on the provided ID.

### New ZTW Endpoint - Partner Integrations - Workload Discovery Service

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /discoveryService/workloadDiscoverySettings` - Retrieves the workload discovery service settings.
    - Added `PUT /discoveryService/{id}/permissions` - Verifies the specified AWS account permissions using the discovery role and external ID.

### New ZTW Endpoint - Partner Integrations - Account Groups

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added the following new ZPA Endpoints
    - Added `GET /accountGroups` - Retrieves the details of AWS account groups with metadata.
    - Added `POST /accountGroups` - Creates an AWS account group. You can create a maximum of 128 groups in each organization. 
    - Added `GET /accountGroups/count` - Retrieves the total number of AWS account groups.
    - Added `GET /accountGroups/lite` - Retrieves the ID and name of all the AWS account groups.
    - Added `PUT /accountGroups/{id}` - Updates the existing AWS account group details based on the provided ID.
    - Added `DELETE /accountGroups/{id}` - Removes a specific AWS account group based on the provided ID.

### Enhancements

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added support to new ZIA `ipdestinationgroups` parameter `override` of type bool. This parameter indicates whether the IPs must be overridden. When set to false, the IPs are appended; else the existing IPs are overridden. The default value is true.
[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added support to new ZIA `dlp_web_rules` attribute `fileTypeCategories`. This attribute supports the list of file types to which the rule applies. This attribute has replaced the attribute `fileTypes`. Zscaler recommends updating your configurations to use the `fileTypeCategories` attribute in place of `fileTypes`. Both attributes are still supported, but cannot be used concurrently.
[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added support to new ZIA `urlfilteringpolicies` attribute `safeSearchApps` of type list of string.

### Enhanced Error Handling and Retry Logic

- **Added automatic retry for 409 EDIT_LOCK_NOT_AVAILABLE errors**: The SDK now automatically detects and retries 409 Conflict responses when encountering edit lock errors (`EDIT_LOCK_NOT_AVAILABLE`, `Resource Access Blocked`, `Failed during enter Org barrier`). Retries use exponential backoff with configurable `RetryWaitMin` and `RetryWaitMax` settings.
- Improved session invalidation handling: Enhanced `SESSION_NOT_VALID` error detection and token refresh logic. The SDK now properly handles both "SESSION_NOT_VALID" and "getAttribute: Session already invalidated" error messages for automatic token renewal and retry.
- Optimized request timeout calculation: Request timeouts now exclude time spent waiting for rate limits, token refreshes, and server backoff delays. This ensures that rate limiting and authentication delays do not count against the overall request timeout, preventing premature failures in long-running operations.
- Fixed request body buffering: The SDK now properly buffers request bodies to enable retry scenarios, including session invalidation and edit lock conflicts, without losing request data.

# 3.7.6 (October 17, 2025)

## Notes
- Golang: **v1.24**

### Enhancements

[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added support to new ZIA `ipdestinationgroups` parameter `override` of type bool. This parameter indicates whether the IPs must be overridden. When set to false, the IPs are appended; else the existing IPs are overridden. The default value is true.
[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added support to new ZIA `dlp_web_rules` attribute `fileTypeCategories`. This attribute supports the list of file types to which the rule applies. This attribute has replaced the attribute `fileTypes`. Zscaler recommends updating your configurations to use the `fileTypeCategories` attribute in place of `fileTypes`. Both attributes are still supported, but cannot be used concurrently.
[PR #379](https://github.com/zscaler/zscaler-sdk-go/pull/379) - Added support to new ZIA `urlfilteringpolicies` attribute `safeSearchApps` of type list of string.

# 3.7.5 (October 14, 2025)

## Notes
- Golang: **v1.24**

### Enhancements

[PR #378](https://github.com/zscaler/zscaler-sdk-go/pull/378) - Removed mutex locks from ZPA policy GET operations to enable concurrent reads. The `GetPolicyRule` function in both `policysetcontroller` (v1) and `policysetcontrollerv2` now executes in parallel, significantly improving performance for large Terraform configurations. CREATE/UPDATE/DELETE operations remain properly serialized per API requirements.

### Bug Fixes

[PR #378](https://github.com/zscaler/zscaler-sdk-go/pull/378) - Fixed ZPA rate limiting backoff logic that was preventing exponential backoff from being calculated. The rate limiter was immediately returning fixed delays instead of allowing intelligent exponential backoff growth (2s → 4s → 8s → 10s max). This caused severe performance degradation (3-4x slower) for Terraform operations with large resource counts. The fix restores proper exponential backoff behavior while maintaining API rate limit handling via 429 status codes and Retry-After headers.

[PR #378](https://github.com/zscaler/zscaler-sdk-go/pull/378) - Added missing exponential backoff to OneAPI client's `ExecuteRequest` retry loop. Server errors (5xx) now retry with intelligent exponential backoff instead of immediately failing or retrying without delay.

# 3.7.4 (October 3, 2025)

## Notes
- Golang: **v1.24**

### ZTW Log and Control Forwarding

[PR #376](https://github.com/zscaler/zscaler-sdk-go/pull/376) - Added the following new ZTW API Endpoints:
    - Added `GET /ecRules/self` Retrieves the list of Log and Control forwarding rules.
    - Added `GET /ecRules/self/{ruleId}` Retrieves a Log and Control forwarding rule configuration based on the specified ID.
    - Added `POST /ecRules/self` Create a Log and Control forwarding rule.
    - Added `PUT /ecRules/self/{ruleId}` Updates Log and Control forwarding rule.
    - Added `DELETE ecRules/self/{ruleId}` Deletes Log and Control forwarding rule.

### ZTW DNS Control Forwarding Rule

[PR #376](https://github.com/zscaler/zscaler-sdk-go/pull/376) - Added the following new ZTW API Endpoints:
    - Added `GET /ecRules/ecDns` Retrieves the list of DNS forwarding rules.
    - Added `GET /ecRules/ecDns/{ruleId}` Retrieves a DNS forwarding rule configuration based on the specified ID.
    - Added `POST /ecRules/ecDns` Create a DNS forwarding rule.
    - Added `PUT /ecRules/ecDns/{ruleId}` Updates DNS forwarding rule.
    - Added `DELETE ecRules/ecDns/{ruleId}` Deletes DNS forwarding rule.

### ZTW DNS Gateway
[PR #376](https://github.com/zscaler/zscaler-sdk-go/pull/376) - Added the following new ZIA API Endpoints:
    - Added `GET /dnsGateways` Retrieves a list of DNS Gateways.
    - Added `GET /dnsGateways/lite` Retrieves a list of DNS Gateways
    - Added `GET /dnsGateways/{gatewayId}` Retrieves the DNS Gateway based on the specified ID
    - Added `POST /dnsGateways` Adds a new DNS Gateway.
    - Added `PUT /dnsGateways/{gatewayId}` Updates the DNS Gateway based on the specified ID
    - Added `DELETE /dnsGateways/{gatewayId}` Deletes a DNS Gateway based on the specified ID

# 3.7.3 (October 3, 2025)

## Notes
- Golang: **v1.24**

### ZTW Log and Control Forwarding

[PR #376](https://github.com/zscaler/zscaler-sdk-go/pull/376) - Added the following new ZTW API Endpoints:
    - Added `GET /ecRules/self` Retrieves the list of Log and Control forwarding rules.
    - Added `GET /ecRules/self/{ruleId}` Retrieves a Log and Control forwarding rule configuration based on the specified ID.
    - Added `POST /ecRules/self` Create a Log and Control forwarding rule.
    - Added `PUT /ecRules/self/{ruleId}` Updates Log and Control forwarding rule.
    - Added `DELETE ecRules/self/{ruleId}` Deletes Log and Control forwarding rule.

### ZTW DNS Control Forwarding Rule

[PR #376](https://github.com/zscaler/zscaler-sdk-go/pull/376) - Added the following new ZTW API Endpoints:
    - Added `GET /ecRules/ecDns` Retrieves the list of DNS forwarding rules.
    - Added `GET /ecRules/ecDns/{ruleId}` Retrieves a DNS forwarding rule configuration based on the specified ID.
    - Added `POST /ecRules/ecDns` Create a DNS forwarding rule.
    - Added `PUT /ecRules/ecDns/{ruleId}` Updates DNS forwarding rule.
    - Added `DELETE ecRules/ecDns/{ruleId}` Deletes DNS forwarding rule.

### ZTW DNS Gateway
[PR #376](https://github.com/zscaler/zscaler-sdk-go/pull/376) - Added the following new ZIA API Endpoints:
    - Added `GET /dnsGateways` Retrieves a list of DNS Gateways.
    - Added `GET /dnsGateways/lite` Retrieves a list of DNS Gateways
    - Added `GET /dnsGateways/{gatewayId}` Retrieves the DNS Gateway based on the specified ID
    - Added `POST /dnsGateways` Adds a new DNS Gateway.
    - Added `PUT /dnsGateways/{gatewayId}` Updates the DNS Gateway based on the specified ID
    - Added `DELETE /dnsGateways/{gatewayId}` Deletes a DNS Gateway based on the specified ID

# 3.7.2 (September 30, 2025)

## Notes
- Golang: **v1.24**

### Bug Fixes

[PR #375](https://github.com/zscaler/zscaler-sdk-go/pull/375) - Implemented fixes and enhancements to ZTW API endpoint packages.


# 3.7.1 (September 22, 2025)

## Notes
- Golang: **v1.24**

### Bug Fixes

[PR #373](https://github.com/zscaler/zscaler-sdk-go/pull/373) - Enhanced session management for ZIA Legacy client to handle 5-minute idle timeout with proactive session validation and refresh capabilities
Please refer to the [Developer Guide](https://help.zscaler.com/zia/getting-started-zia-api#CreateSession) for more details.

[PR #373](https://github.com/zscaler/zscaler-sdk-go/pull/373) - Enhanced session timeout validation and error handling
- Added centralized session invalidation error detection for "SESSION_NOT_VALID" and "Session already invalidated" messages
- Fixed race condition in OAuth2 token renewal ticker with proper mutex locking
- Improved session management with enhanced debugging and automatic token refresh on 401 errors

### Enhancements

[PR #373](https://github.com/zscaler/zscaler-sdk-go/pull/373) - Included function `GetByName` in the ZPA package `c2c_ip_ranges` to allow search by name.


# 3.7.0 (September 15, 2025)

## Notes
- Golang: **v1.23**

#### NEW ZIA Endpoints

[PR #370](https://github.com/zscaler/zscaler-sdk-go/pull/370) - Added the following new ZIA API Endpoints:
    - Added `GET /virtualZenNodes` Retrieves the ZIA Virtual Service Edge for an organization
    - Added `GET /virtualZenNodes/{id}` Retrieves the ZIA Virtual Service Edge for an organization based on the specified ID
    - Added `POST /virtualZenNodes` Adds a ZIA Virtual Service Edge for an organization
    - Added `PUT /virtualZenNodes/{id}` Updates the ZIA Virtual Service Edge for an organization based on the specified ID
    - Added `DELETE /virtualZenNodes/{id}` Deletes the ZIA Virtual Service Edge for an organization based on the specified ID

[PR #370](https://github.com/zscaler/zscaler-sdk-go/pull/370) - Added the following new ZIA API Endpoints:
    - Added `GET /workloadGroups/{id}` Retrieves the workload group based on the specified ID
    - Added `POST /workloadGroups` Adds a workload group for an organization
    - Added `PUT /workloadGroups/{id}` Updates the workload group for an organization based on the specified ID
    - Added `DELETE /workloadGroups/{id}` Updates the workload group based on the specified ID

[PR #370](https://github.com/zscaler/zscaler-sdk-go/pull/370) - Added the following new ZIA API Endpoints:
    - Added `GET /casbTenant/scanInfo` Retrieves the SaaS Security Scan Configuration information

# 3.6.4 (August 26, 2025)

## Notes
- Golang: **v1.23**

### Bug Fixes
[PR #367](https://github.com/zscaler/zscaler-sdk-go/pull/367) - Added attribute `deviceGroups` to ZIA Forwarding Control Rules.

# 3.6.3 (August 26, 2025)

## Notes
- Golang: **v1.23**

### Bug Fixes
[PR #367](https://github.com/zscaler/zscaler-sdk-go/pull/367) - Fixed ZIA `c2c_incident_receiver` `onboardableEntity` nested attribute.
[PR #367](https://github.com/zscaler/zscaler-sdk-go/pull/367) - Added support to `receiver` attribute within ZIA `casb_dlp_rules` resource.

# 3.6.2 (August 26, 2025)

## Notes
- Golang: **v1.23**

### Bug Fixes
[PR #366](https://github.com/zscaler/zscaler-sdk-go/pull/366) - Fixed ZIA `c2c_incident_receiver` `onboardableEntity` nested attribute.

# 3.6.1 (August 22, 2025)

## Notes
- Golang: **v1.23**

### Bug Fixes
[PR #364](https://github.com/zscaler/zscaler-sdk-go/pull/364) - Fixed `cbiProfile` within the `zia` `urlfilteringpolicies`, update and create functions

# 3.6.0 (August 18, 2025)

## Notes
- Golang: **v1.23**

#### NEW Enhancement - ZIdentity API Support

[PR #363](https://github.com/zscaler/zscaler-sdk-go/pull/363): Zscaler [Zidentity](https://help.zscaler.com/zidentity/what-zidentity) API is now available and is supported by this SDK. See [README](https://github.com/zscaler/zscaler-sdk-go/blob/master/README.md) for authentication instructions.

### New ZPA Endpoint - Admin SSO Configuration Controller
[PR #363](https://github.com/zscaler/zscaler-sdk-go/pull/363) Added the following new ZPA API Endpoints:
    - Added `GET /v2/ssoLoginOptions` Get SSO Login Details
    - Added `POST /v2/ssoLoginOptions` Updates SSO Options for customer

### New ZPA Endpoint - C2C IP Ranges

[PR #363](https://github.com/zscaler/zscaler-sdk-go/pull/363) Added the following new ZPA API Endpoints:
    - Added `POST /v2/ipRanges/search` Get the IP Range by `page` and `pageSize`
    - Added `GET /v2/ipRanges` Get All the IP Range
    - Added `POST /v2/ipRanges` Add new IP Range
    - Added `GET /v2/ipRanges/{ipRangeId}` Get the IP Range Details
    - Added `PUT /v2/ipRanges/{ipRangeId}` Update the IP Range Details
    - Added `DELETE /v2/ipRanges/{ipRangeId}` Delete IP Range

### New ZPA Endpoint - API Keys

[PR #363](https://github.com/zscaler/zscaler-sdk-go/pull/363) Added the following new ZPA API Endpoints:
    - Added `GET /apiKeys` Get all apiKeys details
    - Added `POST /apiKeys` Create api keys for customer
    - Added `GET /apiKeys/{id}` Get apiKeys details by ID
    - Added `PUT /apiKeys/{id}` Update apiKeys by ID
    - Added `DELETE /apiKeys/{id}` Delete apiKeys

### New ZPA Endpoint - Customer Controller

[PR #363](https://github.com/zscaler/zscaler-sdk-go/pull/363) Added the following new ZPA API Endpoints:
    - Added `GET /v2/associationtype/{type}/domains` Get domains for a customer
    - Added `POST /v2/associationtype/{type}/domains` Add or update domains for a customer.

### New ZPA Endpoint - NPClient

[PR #363](https://github.com/zscaler/zscaler-sdk-go/pull/363) Added the following new ZPA API Endpoints:
    - Added `GET /vpnConnectedUsers` Get all applications configuired for a given customer

### New ZPA Endpoint - Private Cloud Controller Group

[PR #363](https://github.com/zscaler/zscaler-sdk-go/pull/363) Added the following new ZPA API Endpoints:
    - Added `GET /privateCloudControllerGroup` Get details of all configured Private Cloud Controller Groups
    - Added `POST /privateCloudControllerGroup` Add a new Private Cloud Controller Groups
    - Added `GET /privateCloudControllerGroup/{privateCloudControllerGroupId}` Get the Private Cloud Controller Group details for the specified ID
    - Added `PUT /privateCloudControllerGroup/{privateCloudControllerGroupId}` Update the Private Cloud Controller Group details for the specified ID
    - Added `DELETE /privateCloudControllerGroup/{privateCloudControllerGroupId}` Delete the Private Cloud Controller Group for the specified ID
    - Added `DELETE /privateCloudControllerGroup/summary` Get all the configured Private Cloud Controller Group ID and Name

### New ZPA Endpoint - Private Cloud Controller Group

[PR #363](https://github.com/zscaler/zscaler-sdk-go/pull/363) Added the following new ZPA API Endpoints:
    - Added `GET /privateCloudController` Get all the configured Private Cloud Controller details
    - Added `PUT /privateCloudController/{privateCloudControllerGroupId}/restart` Trigger restart of the Private Cloud Controller
    - Added `GET /privateCloudController/{privateCloudControllerId}` Gets the Private Cloud Controller details for the specified ID.
    - Added `PUT /privateCloudController/{privateCloudControllerId}` Updates the Private Cloud Controller for the specified ID
    - Added `DELETE /privateCloudController/{privateCloudControllerId}` Delete the Private Cloud Controller for the specified ID

### New ZPA Endpoint - User Portal Controller

[PR #363](https://github.com/zscaler/zscaler-sdk-go/pull/363) Added the following new ZPA API Endpoints:
    - Added `GET /userPortal` Get all configured User Portals
    - Added `GET /userPortal/{id}` Get User Portal for the specified ID
    - Added `PUT /userPortal/{Id}` Update User Portal for the specified ID
    - Added `POST /userPortal` Add a new User Portal
    - Added `DELETE /userPortal/{Id}` Delete a User Portal

### New ZPA Endpoint - User Portal Link Controller

[PR #363](https://github.com/zscaler/zscaler-sdk-go/pull/363) Added the following new ZPA API Endpoints:
    - Added `GET /userPortalLink` Get all configured User Portal Links
    - Added `GET /userPortalLink/{id}` Get User Portal Link for the specified ID
    - Added `GET /userPortalLink/userPortal/{portalId}` Get User Portal Link for a given portal
    - Added `PUT /userPortalLink/{Id}` Update User Portal Link for the specified ID
    - Added `POST /userPortalLink` Add a new User Portal Link
    - Added `POST /userPortalLink/bulk` Add list of User Portal Link
    - Added `DELETE /userPortalLink/{Id}` Delete a User Portal Link for the specified ID

### New ZPA Endpoint - Z-Path Config Override Controller

[PR #363](https://github.com/zscaler/zscaler-sdk-go/pull/363) Added the following new ZPA API Endpoints:
    - Added `GET /configOverrides/{id}` Get config-override details by configId
    - Added `GET /configOverrides` Get all config-override details
    - Added `PUT /configOverrides/{id}` Update config-override for the specified ID
    - Added `POST /configOverrides` Create config-override

### New ZPA Endpoint - Multimatch Domains

[PR #363](https://github.com/zscaler/zscaler-sdk-go/pull/363) Added the following new ZPA API Endpoints:
    - Added `GET /multimatchUnsupportedReferences` Get the unsupported feature references for multimatch for domains
    - Added `GET /bulkUpdateMultiMatch` Update multimatch feature in multiple applications.

### New ZIA Cloud-to-Cloud - Receiver for Web DLP Rule

[PR #363](https://github.com/zscaler/zscaler-sdk-go/pull/363) Added attribute `receiver` to support Cloud-to-Cloud - Receiver for Web DLP Rule configuration.

### Bug Fixes
[PR #363](https://github.com/zscaler/zscaler-sdk-go/pull/363) - Fixed ZIA `urlfilteringpolicies` `GET` function to return the complete payload.

# 3.5.4 (July 23, 2025)

## Notes
- Golang: **v1.23**

### Bug Fixes
[PR #357](https://github.com/zscaler/zscaler-sdk-go/pull/356) - Added `omitempty` to ZIA `cloud_nss` attributes to prevent JSON malformed.

# 3.5.3 (July 22, 2025)

## Notes
- Golang: **v1.23**

### Bug Fixes
[PR #356](https://github.com/zscaler/zscaler-sdk-go/pull/356) - Fixed ZIA `bandwidth_classes` attribute `fileSize`
[PR #356](https://github.com/zscaler/zscaler-sdk-go/pull/356) - Added attribute `browserEunTemplateId` to ZIA resources `filetypecontrol`, `urlfilteringpolicies`
[PR #356](https://github.com/zscaler/zscaler-sdk-go/pull/356) - Fixed zpabeta url entrypoint in ZPA `config_scim` client.

# 3.5.2 (June 27, 2025)

## Notes
- Golang: **v1.23**

### Bug Fixes
[PR #352](https://github.com/zscaler/zscaler-sdk-go/pull/352) - Added `ruleMutex.Lock()` to ZPA `policysetcontroller` function `BulkReorder` to prevent simulteneous API calls when invoking the bulk reordering endpoint.

# 3.5.1 (June 23, 2025)

## Notes
- Golang: **v1.23**

### Bug Fixes
[PR #351](https://github.com/zscaler/zscaler-sdk-go/pull/351) - Fixed `CheckErrorInResponse` function to parse and display API error messages more clearly in ZIA Legacy client.

# 3.5.0 (June 19, 2025)

## Notes
- Golang: **v1.23**

### Zscaler Digital Experience (ZDX) - OneAPI Support
[PR #350](https://github.com/zscaler/zscaler-sdk-go/pull/350) - Zscaler Digital Experience (ZDX) API endpoints are now supported via OneAPI.

### New ZIA Endpoint - Browser Control Policy

[PR #350](https://github.com/zscaler/zscaler-sdk-go/pull/350) Added the following new ZIA API Endpoints:
    - Added `GET /browserControlSettings` Retrieves the Browser Control status and the list of configured browsers in the Browser Control policy
    - Added `PUT /browserControlSettings` Updates the Browser Control settings.

### New ZIA Endpoint - SaaS Security API (Casb DLP Rules)

[PR #350](https://github.com/zscaler/zscaler-sdk-go/pull/350) Added the following new ZIA API Endpoints:
    - Added `GET /casbDlpRules` Retrieves the SaaS Security Data at Rest Scanning Data Loss Prevention (DLP) rules based on the specified rule type.
    - Added `GET /casbDlpRules/{ruleId}` Retrieves the SaaS Security Data at Rest Scanning DLP rule based on the specified ID
    - Added `GET /casbDlpRules/all` Retrieves all the SaaS Security Data at Rest Scanning DLP rules
    - Added `POST /casbDlpRules` Adds a new SaaS Security Data at Rest Scanning DLP rule
    - Added `PUT /casbDlpRules/{ruleId}` Updates the SaaS Security Data at Rest Scanning DLP rule based on the specified ID
    - Added `DELETE /casbDlpRules/{ruleId}` Deletes the SaaS Security Data at Rest Scanning DLP rule based on the specified ID

### New ZIA Endpoint - SaaS Security API (Casb Malware Rules)

[PR #350](https://github.com/zscaler/zscaler-sdk-go/pull/350) Added the following new ZIA API Endpoints:
    - Added `GET /casbMalwareRules` Retrieves the SaaS Security Data at Rest Scanning Malware Detection rules based on the specified rule type.
    - Added `GET /casbMalwareRules/{ruleId}` Retrieves the SaaS Security Data at Rest Scanning Malware Detection rule based on the specified ID
    - Added `GET /casbMalwareRules/all` Retrieves all the SaaS Security Data at Rest Scanning Malware Detection rules
    - Added `POST /casbMalwareRules` Adds a new SaaS Security Data at Rest Scanning Malware Detection rule.
    - Added `PUT /casbMalwareRules/{ruleId}` Updates the SaaS Security Data at Rest Scanning Malware Detection rule based on the specified ID
    - Added `DELETE /casbMalwareRules/{ruleId}` Deletes the SaaS Security Data at Rest Scanning Malware Detection rule based on the specified ID

### New ZIA Endpoint - SaaS Security API

[PR #350](https://github.com/zscaler/zscaler-sdk-go/pull/350) Added the following new ZIA API Endpoints:
    - Added `GET /domainProfiles/lite` Retrieves the domain profile summary.
    - Added `GET /quarantineTombstoneTemplate/lite` Retrieves the templates for the tombstone file created when a file is quarantined
    - Added `GET /casbEmailLabel/lite` Retrieves the email labels generated for the SaaS Security API policies in a user's email account
    - Added `GET /casbTenant/{tenantId}/tags/policy` Retrieves the tags used in the policy rules associated with a tenant, based on the tenant ID.
    - Added `GET /casbTenant/lite` Retrieves information about the SaaS application tenant

### New ZIA Location Management Attributes - Extranet Support

[PR #350](https://github.com/zscaler/zscaler-sdk-go/pull/350) The following attributes have been introduced to support Extranet feature configuration:
  - `extranet` - The ID of the extranet resource that must be assigned to the location
  - `extranetIpPool` - The ID of the traffic selector specified in the extranet
  - `extranetDns` - The ID of the DNS server configuration used in the extranet
  - `defaultExtranetTsPool` - A Boolean value indicating that the traffic selector specified in the extranet is the designated default traffic selector
  - `defaultExtranetDns` - A Boolean value indicating that the DNS server configuration used in the extranet is the designated default DNS server

### Internal Enhancements
* [PR #350](https://github.com/zscaler/zscaler-sdk-python/pull/350) - Enhanced `CheckErrorInResponse` function to parse and display API error messages more clearly.

# 3.4.4 (June 6, 2025)

## Notes
- Golang: **v1.23**

### New ZIA Endpoint - Virtual ZEN Clusters:

[PR #348](https://github.com/zscaler/zscaler-sdk-go/pull/348) Added the following new ZIA API Endpoints:
    - Added `GET /virtualZenClusters` Retrieves a list of ZIA Virtual Service Edge clusters.
    - Added `GET /virtualZenClusters/{cluster_id}` Retrieves the Virtual Service Edge cluster based on the specified ID
    - Added `POST /virtualZenClusters` Adds a new Virtual Service Edge cluster. 
    - Added `PUT /virtualZenClusters/{cluster_id}` Updates the Virtual Service Edge cluster based on the specified ID
    - Added `DELETE /virtualZenClusters/{cluster_id}` Deletes the Virtual Service Edge cluster based on the specified ID

### Enhancements
[PR #348](https://github.com/zscaler/zscaler-sdk-go/pull/348) - Added new Policy Client Types: `zpn_client_type_zapp_partner`, `zpn_client_type_vdi`, `zpn_client_type_zia_inspection`

### Bug Fixes
[PR #348](https://github.com/zscaler/zscaler-sdk-go/pull/348) - The SDK's NewOneAPIClient() function was performing OAuth2 authentication unconditionally, which caused the to hang or fail during legacy client initialization. The logic has been updated to skip authentication when the legacy client is in use.

# 3.4.3 (June 5, 2025)

## Notes
- Golang: **v1.23**

### Enhancements
[PR #345](https://github.com/zscaler/zscaler-sdk-go/pull/345) - Added `DELETE` method for ZIA `alertSubscriptions/{alertSubscriptionId}`

### Bug Fixes
[PR #339](https://github.com/zscaler/zscaler-sdk-go/pull/339) - Fixed an issue in the OneAPI client where initial API requests could fail with 401 Unauthorized if a valid OAuth2 token had not yet been obtained. The client now performs immediate authentication before starting the background token renewal ticker, ensuring the token is always present on first use. Removed incorrect retry logic for 401 responses.
[PR #341](https://github.com/zscaler/zscaler-sdk-go/pull/341) - Corrected `YAML` tag for `DefaultCacheMaxSizeMB` fields to use `defaultSize`
[PR #342](https://github.com/zscaler/zscaler-sdk-go/pull/342) - Fix comment referencing provisioning keys in ZPA service
[PR #343](https://github.com/zscaler/zscaler-sdk-go/pull/343) - Document README typo fixes and YAML tag correction as part of the 3.4.3 entry
[PR #345](https://github.com/zscaler/zscaler-sdk-go/pull/345) - Fix `InsecureSkipVerify` flag so TLS checks can be disabled in testing
[PR #345](https://github.com/zscaler/zscaler-sdk-go/pull/345) - Add `Close()` method to stop token renewal ticker
[PR #346](https://github.com/zscaler/zscaler-sdk-go/pull/346) - Fixed ZPA Pagination encoding for special edge cases.

### Documentation
[PR #340](https://github.com/zscaler/zscaler-sdk-go/pull/340) - Fixed README typos and clarified service detection description.

# 3.4.2 (May 29, 2025)

## Notes
- Golang: **v1.23**

### Bug Fixes

### ZPA Privileged Remote Access Portal
[PR #338](https://github.com/zscaler/zscaler-sdk-go/pull/338) - Added support for Zscaler Managed Certificate to resource `zpa_application_segment_browser_access`

# 3.4.1 (May 28, 2025)

## Notes
- Golang: **v1.23**

### Bug Fixes

### ZPA Privileged Remote Access Portal
[PR #335](https://github.com/zscaler/zscaler-sdk-go/pull/335) - Added support for PRA User Portal with Zscaler Managed Certificate

# 3.4.0 (May 28, 2025) - NEW ZPA ENDPOINT RESOURCES

## Notes
- Golang: **v1.23**

### ZPA Administrator Controller
[PR #335](https://github.com/zscaler/zscaler-sdk-go/pull/335) - Added the following new ZPA API Endpoints:
    - Added `GET /administrators` Retrieves a list of administrators in a tenant. A maximum of 200 administrators are returned per request.
    - Added `GET /administrators/{admin_id}` Retrieves administrator details for a specific `{admin_id}`
    - Added `POST /administrators` Create an local administrator account
    - Added `PUT /administrators/{admin_id}` Update a local administrator account for a specific `{admin_id}`
    - Added `DELETE /administrators/{admin_id}` Delete a local administrator account for a specific `{admin_id}`

### ZPA Role Controller
[PR #335](https://github.com/zscaler/zscaler-sdk-go/pull/335) - Added the following new ZPA API Endpoints:
    - Added `GET /permissionGroups` Retrieves all the default permission groups.
    - Added `GET /roles` Retrieves a list of all configured roles in a tenant.
    - Added `GET /roles/{admin_id}` Retrieves a role details for a specific `{role_id}`
    - Added `POST /roles` Adds a new role for a tenant.
    - Added `PUT /roles/{admin_id}` Update a role for a specific `{role_id}`
    - Added `DELETE /roles/{role_id}` Delete a role for a specific `{role_id}`

### ZPA Enrollment Certificate Controller
[PR #335](https://github.com/zscaler/zscaler-sdk-go/pull/335) - Added the following new ZPA API Endpoints:
    - Added `POST /enrollmentCert/csr/generate` Creates a CSR for a new enrollment Certificate
    - Added `POST /enrollmentCert/selfsigned/generate` Creates a self signed Enrollment Certificate
    - Added `POST /enrollmentCert` Creates a enrollment Certificate
    - Added `PUT /enrollmentCert/{cert_id}` Update an existing enrollment Certificate
    - Added `DELETE /enrollmentCert/{cert_id}` Delete an existing enrollment Certificate

### ZPA SAML Attribute Controller
[PR #335](https://github.com/zscaler/zscaler-sdk-go/pull/335) - Added the following new ZPA API Endpoints:
    - Added `POST /samlAttribute` Adds a new `SamlAttribute` for a given tenant
    - Added `PUT /samlAttribute/{attr_id}` Update an existing `SamlAttribute` for a given tenant
    - Added `DELETE /samlAttribute/{attr_id}` Delete an existing `SamlAttribute` for a given tenant

### ZPA Client-Settings Controller
[PR #335](https://github.com/zscaler/zscaler-sdk-go/pull/335) - Added the following new ZPA API Endpoints:
    - Added `GET /clientSetting` Retrieves `clientSetting` details. `ClientCertType` defaults to `CLIENT_CONNECTOR`
    - Added `POST /clientSetting` Create or update `clientSetting` for a customer. `ClientCertType` defaults to `CLIENT_CONNECTOR`
    - Added `DELETE /clientSetting` Delete an existing `clientSetting`. `ClientCertType` defaults to `CLIENT_CONNECTOR`
    - Added `GET /clientSetting/all` Retrieves all `clientSetting` details.

### ZPA Privileged Remote Access Portal
[PR #335](https://github.com/zscaler/zscaler-sdk-go/pull/335) - Added support for PRA User Portal with Zscaler Managed Certificate

#### ZIA SCIM API
[PR #335](https://github.com/zscaler/zscaler-sdk-go/pull/335) - This SDK now supports direct interaction with the ZIA SCIM API Endpoint for user and group management. See [README](https://github.com/zscaler/zscaler-sdk-go/blob/master/README.md)

#### ZPA SCIM API
[PR #335](https://github.com/zscaler/zscaler-sdk-go/pull/335) - Enhanced interaction with ZPA SCIM API via ConfigSetter support and easier client instantiation. See [README](https://github.com/zscaler/zscaler-sdk-go/blob/master/README.md)

### Bug Fixes

[PR #335](https://github.com/zscaler/zscaler-sdk-go/pull/335) – Fixed ZPA central pagination engine `getAllPagesGenericWithCustomFilters` to appropriately parse the response on empty lists. Change was done due to upstream api changes.

# 3.3.1 (May 20, 2025)

## Notes
- Golang: **v1.22**

### Bug Fixes

[PR #333](https://github.com/zscaler/zscaler-sdk-go/pull/333) – Fixed ZIA parameter values `sortOrder` and `sortBy` for pagination.
 - `sortOrder` - Supported Values: `asc`, `desc`, `ruleExecution`
 - `sortBy` - Supported Values: `id`, `name`, `expiry`, `status`, `externalId`, `rank`
 
# 3.3.0 (April 30, 2025) - NEW ZIA ENDPOINT RESOURCES

## Notes
- Golang: **v1.22**

### Bug Fixes

[PR #326](https://github.com/zscaler/zscaler-sdk-go/pull/326) – Added `context.Context` to `startSessionTicker` on both legacy `zia` and `ztw` API Clients to ensure go routine terminarion. This will allow the context to be explicitly cancelled.

### ZIA Password Expiry Settings
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /passwordExpiry/settings` Retrieves the password expiration information for all the admins
    - Added `PUT /passwordExpiry/settings` Updates the password expiration information for all the admins.

### ZIA Alerts
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /alertSubscriptions` Retrieves a list of all alert subscriptions
    - Added `GET /alertSubscriptions/{subscription_id}` Retrieves the alert subscription information based on the specified ID
    - Added `POST /alertSubscriptions` Adds a new alert subscription.
    - Added `PUT /alertSubscriptions/{subscription_id}` Updates an existing alert subscription based on the specified ID

### ZIA NSS Servers
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /nssServers` Retrieves a list of registered NSS servers. 
    - Added `GET /nssServers/{nss_id}` Retrieves the registered NSS server based on the specified ID
    - Added `POST /nssServers` AddsAdds a new NSS server.
    - Added `PUT /nssServers/{nss_id}` Updates an NSS server based on the specified ID
    - Added `DELETE /nssServers/{nss_id}` Deletes an NSS server based on the specified ID

### ZIA Bandwidth Classes
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /bandwidthClasses` Retrieves a list of bandwidth classes for an organization.
    - Added `GET /bandwidthClasses/lite` Retrieves a list of bandwidth classes for an organization
    - Added `GET /bandwidthClasses/{class_id}` Retrieves the alert subscription information based on the specified ID
    - Added `POST /bandwidthClasses` Adds a new bandwidth class.
    - Added `PUT /bandwidthClasses/{class_id}` Updates a bandwidth class based on the specified ID
    - Added `DELETE /bandwidthClasses/{class_id}` Deletes a bandwidth class based on the specified ID

### ZIA Bandwidth Control Rules
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /bandwidthControlRules` Retrieves all the rules in the Bandwidth Control policy.
    - Added `GET /bandwidthControlRules/lite` Retrieves all the rules in the Bandwidth Control policy
    - Added `GET /bandwidthControlRules/{rule_id}` Retrieves the Bandwidth Control policy rule based on the specified ID
    - Added `POST /bandwidthControlRules` Adds a new Bandwidth Control policy rule.
    - Added `PUT /bandwidthControlRules/{rule_id}` Updates the Bandwidth Control policy rule based on the specified ID
    - Added `DELETE /bandwidthControlRules/{rule_id}` Deletes a Bandwidth Control policy rule based on the specified ID

### ZIA NAT Control Policy
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /dnatRules` Retrieves a list of all configured and predefined DNAT Control policies. 
    - Added `GET /dnatRules/{rule_id}` Retrieves the DNAT Control policy rule information based on the specified ID
    - Added `POST /dnatRules` Adds a new DNAT Control policy rule.
    - Added `PUT /dnatRules/{rule_id}` Updates the DNAT Control policy rule information based on the specified ID
    - Added `DELETE /dnatRules/{rule_id}` Deletes the DNAT Control policy rule information based on the specified ID

### ZIA Risk Profiles
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /riskProfiles` Retrieves the cloud application risk profile.
    - Added `GET /riskProfiles/lite` Retrieves the cloud application risk profile
    - Added `GET /riskProfiles/{profile_id}` Retrieves the cloud application risk profile based on the specified ID
    - Added `POST /riskProfiles` Adds a new cloud application risk profile. 
    - Added `PUT /riskProfiles/{profile_id}` Updates the cloud application risk profile based on the specified ID
    - Added `DELETE /riskProfiles/{profile_id}` Deletes the cloud application risk profile based on the specified ID

### ZIA Cloud Application Instances
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /cloudApplicationInstances` Retrieves the list of cloud application instances configured in the ZIA Admin Portal.
    - Added `GET /cloudApplicationInstances/{instance_id}` Retrieves information about a cloud application instance based on the specified ID
    - Added `POST /cloudApplicationInstances` Add a new cloud application instance. 
    - Added `PUT /cloudApplicationInstances/{instance_id}` Updates information about a cloud application instance based on the specified ID
    - Added `DELETE /cloudApplicationInstances/{instance_id}` Deletes a cloud application instance based on the specified ID

### ZIA Cloud Application Instances
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /cloudApplicationInstances` Retrieves the list of cloud application instances configured in the ZIA Admin Portal.
    - Added `GET /cloudApplicationInstances/{instance_id}` Retrieves information about a cloud application instance based on the specified ID
    - Added `POST /cloudApplicationInstances` Add a new cloud application instance. 
    - Added `PUT /cloudApplicationInstances/{instance_id}` Updates information about a cloud application instance based on the specified ID
    - Added `DELETE /cloudApplicationInstances/{instance_id}` Deletes a cloud application instance based on the specified ID

### ZIA Tenancy Restriction Profile
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /tenancyRestrictionProfile` Retrieves all the restricted tenant profiles.
    - Added `GET /tenancyRestrictionProfile/{profile_id}`Retrieves the restricted tenant profile based on the specified ID
    - Added `POST /tenancyRestrictionProfile` Creates restricted tenant profiles. 
    - Added `PUT /tenancyRestrictionProfile/{profile_id}` Updates the restricted tenant profile based on the specified ID
    - Added `DELETE /tenancyRestrictionProfile/{profile_id}` Deletes the restricted tenant profile based on the specified ID
    - Added `GET /tenancyRestrictionProfile/app-item-count/{app_type}/{item_type}` Retrieves the item count of the specified item type for a given application, excluding any specified profile

### ZIA Tenancy Restriction Profile
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /tenancyRestrictionProfile` Retrieves all the restricted tenant profiles.
    - Added `GET /tenancyRestrictionProfile/{profile_id}`Retrieves the restricted tenant profile based on the specified ID
    - Added `POST /tenancyRestrictionProfile` Creates restricted tenant profiles. 
    - Added `PUT /tenancyRestrictionProfile/{profile_id}` Updates the restricted tenant profile based on the specified ID
    - Added `DELETE /tenancyRestrictionProfile/{profile_id}` Deletes the restricted tenant profile based on the specified ID

### ZIA DNS Gateway
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /dnsGateways` Retrieves a list of DNS Gateways.
    - Added `GET /dnsGateways/lite` Retrieves a list of DNS Gateways
    - Added `GET /dnsGateways/{gateway_id}` Retrieves the DNS Gateway based on the specified ID
    - Added `POST /dnsGateways` Adds a new DNS Gateway.
    - Added `PUT /dnsGateways/{gateway_id}` Updates the DNS Gateway based on the specified ID
    - Added `DELETE /dnsGateways/{gateway_id}` Deletes a DNS Gateway based on the specified ID

### ZIA Proxies
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /proxies` Retrieves a list of all proxies configured for third-party proxy services.
    - Added `GET /proxies/lite` Retrieves a list of all proxies configured for third-party proxy services
    - Added `GET /proxies/{proxy_id}` Retrieves the proxy information based on the specified ID
    - Added `POST /proxies` Adds a new proxy for a third-party proxy service.
    - Added `PUT /proxies/{proxy_id}` Updates an existing proxy based on the specified ID
    - Added `DELETE /proxies/{proxy_id}` Deletes an existing proxy based on the specified ID
    - Added `DELETE /dedicatedIPGateways/lite` Retrieves a list of dedicated IP gateways.

### ZIA FTP Settings
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /ftpSettings` Retrieves the FTP Control status and the list of URL categories for which FTP is allowed.
    - Added `PUT /ftpSettings` Updates the FTP Control settings.

### ZIA Mobile Malware Protection Policy
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /mobileAdvanceThreatSettings` Retrieves all the rules in the Mobile Malware Protection policy
    - Added `PUT /mobileAdvanceThreatSettings` Updates the Mobile Malware Protection rule information. 

### ZIA Mobile Malware Protection Policy
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /configAudit` Retrieves the System Audit Report.
    - Added `GET /configAudit/ipVisibility` Retrieves the IP visibility audit report.
    - Added `GET /configAudit/pacFile` Retrieves the PAC file audit report.
**Note**: This endpoint is accessible via Zscaler OneAPI only.

### ZIA Time Intervals
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /timeIntervals` Retrieves the System Audit Report.
    - Added `GET /timeIntervals/{interval_id}` Retrieves the configured time interval based on the specified ID
    - Added `POST /timeIntervals/{interval_id}` Adds a new time interval.
    - Added `PUT /timeIntervals/{interval_id}` Updates the time interval based on the specified ID
    - Added `DELETE /timeIntervals/{interval_id}` Deletes a time interval based on the specified ID

### ZIA Data Center Exclusions
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /dcExclusions` Retrieves the list of Zscaler data centers (DCs) that are currently excluded from service to your organization based on configured exclusions in the ZIA Admin Portal
    - Added `POST /dcExclusions/{dc_id}` Adds a data center (DC) exclusion to disable the tunnels terminating at a virtual IP address of a Zscaler DC
    - Added `PUT /dcExclusions/{dc_id}` Updates a Zscaler data center (DC) exclusion configuration based on the specified ID.
    - Added `DELETE /dcExclusions/{dc_id}` Deletes a Zscaler data center (DC) exclusion configuration based on the specified ID. 
    - Added `GET /datacenters` Retrieves the list of Zscaler data centers (DCs) that can be excluded from service to your organization

### ZIA SubClouds
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /subclouds` Retrieves all the subclouds and the excluded data centers that are associated with the subcloud
    - Added `GET subclouds/isLastDcInCountry/{cloud_id}` Retrieves the list of all the excluded data centers in a country
    - Added `PUT /subclouds/{cloud_id}` Updates the subcloud and excluded data centers based on the specified ID

### ZIA IPv6 Configuration
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /ipv6config` Gets the IPv6 configuration details for the organization.
    - Added `GET ipv6config/dns64prefix` Gets the list of NAT64 prefixes configured as the DNS64 prefix for the organization.
    - Added `GET /ipv6config/nat64prefix` Gets the list of NAT64 prefixes configured for the organization. 

### ZIA Groups
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /groups` Retrieves a list of groups. The search parameters find matching values in the name or comments attributes.configured exclusions in the ZIA Admin Portal
    - Added `GET /groups/lite` Retrieves a list of group names. The search parameters find matching values in the name or comments attributes.
    - Added `GET /groups/{group_id}` Retrieves the group based on the specified ID
    - Added `POST /groups` Adds a new group.
    - Added `PUT /groups/{group_id}` Updates an existing group based on the specified ID.
    - Added `DELETE /groups/{group_id}` Deletes the group based on the specified ID.

### ZIA Departments
[PR #326](https://github.com/zscaler/zscaler-sdk-python/pull/326) - Added the following new ZIA API Endpoints:
    - Added `GET /departments` Retrieves a list of groups. The search parameters find matching values in the name or comments attributes.configured exclusions in the ZIA Admin Portal
    - Added `GET /departments/lite` Retrieves a list of group names. The search parameters find matching values in the name or comments attributes.Retrieves a list of departments. The search parameters find matching values within the name or comments fields.
    - Added `GET /departments/lite/{department_id}` Retrieves the department based on the specified ID
    - Added `GET /departments/{department_id}` Retrieves the department based on the specified ID
    - Added `POST /departments` Adds a department for an organization. 
    - Added `PUT /departments/{department_id}` Updates the department for an organization based on the specified ID.
    - Added `DELETE /departments/{department_id}` Deletes a department for an organization based on the specified ID.

# 3.2.4 (April 23, 2025)

## Notes
- Golang: **v1.22**

### Bug Fixes

[PR #323](https://github.com/zscaler/zscaler-sdk-go/pull/323) – Fixed ZPA resource `pracredentialpool` missing attribute `credentialMappingCount`
[PR #323](https://github.com/zscaler/zscaler-sdk-go/pull/323) – Fixed ZPA resource `policysetcontrollerv2` attribute `credentialPool` pointer
[PR #323](https://github.com/zscaler/zscaler-sdk-go/pull/323) – Fixed ZDX rate limit handling in the `getRetryAfter` function by correctly parsing and calculating retry delays based on the headers: `X-Ratelimit-Remaining-Second` and `X-Ratelimit-Limit-Second`.

# 3.2.3 (April 22, 2025)

## Notes
- Golang: **v1.22**

### Bug Fixes
[PR #322](https://github.com/zscaler/zscaler-sdk-go/pull/322) – Improved OneAPI client rate limit handling by leveraging all available headers for more accurate retry behavior:

- `X-Ratelimit-Reset`
- `X-Ratelimit-Remaining`
- `X-Ratelimit-Limit`

This enhancement enables proactive throttling and reduces the likelihood of encountering 429 responses by calculating wait times more precisely.

# 3.2.2 (April 17, 2025)

## Notes
- Golang: **v1.22**

### Bug Fixes
[PR #320](https://github.com/zscaler/zscaler-sdk-go/pull/320) - Fixed ZIA Admin Roles endpoint.

# 3.2.1 (April 17, 2025)

## Notes
- Golang: **v1.22**

### Bug Fixes
[PR #320](https://github.com/zscaler/zscaler-sdk-go/pull/320) - Fixed ZIA Admin Roles endpoint.

# 3.2.0 (April 16, 2025)

## Notes
- Golang: **v1.22**

### Cloud & Branch Connector - OneAPI Support
[PR #320](https://github.com/zscaler/zscaler-sdk-go/pull/320) - Cloud & Branch Connector package is now compatible with OneAPI and Legacy API framework. Please refer to README for details.
[PR #320](https://github.com/zscaler/zscaler-sdk-go/pull/320) - Cloud & Branch Connector package has been renamed from `zcon` to `ztw`

### ZTW Policy Management
[PR #320](https://github.com/zscaler/zscaler-sdk-go/pull/320) - Added the following new ZTW API Endpoints:
    - Added `GET /ecRules/ecRdr` Retrieves the list of traffic forwarding rules.
    - Added `PUT /ecRules/ecRdr/{ruleId}` Updates a traffic forwarding rule configuration based on the specified ID.
    - Added `POST /ecRules/ecRdr` Creates a new traffic forwarding rule.
    - Added `GET /ecRules/ecRdr/count` Retrieves the count of traffic forwarding rules available in the Cloud & Branch Connector Admin Portal.

### ZTW Policy Resources
[PR #320](https://github.com/zscaler/zscaler-sdk-go/pull/320) - Added the following new ZTW API Endpoints:
    - Added `GET /ipSourceGroups` Retrieves the list of source IP groups.
    - Added `GET /ipSourceGroups/lite` Retrieves the list of source IP groups. This request retrieves basic information about the source IP groups, such as name and ID. For extensive details, use the GET /ipSourceGroups request.
    - Added `POST /ipSourceGroups` Adds a new custom source IP group.
    - Added `DELETE /ipSourceGroups/{ipGroupId}` Deletes a source IP group based on the specified ID.
    - Added `GET /ipDestinationGroups` Retrieves the list of destination IP groups.
    - Added `GET /ipDestinationGroups/lite` Retrieves the list of destination IP groups. This request retrieves basic information about the destination IP groups, ID, name, and type. For extensive details, use the GET /ipDestinationGroups request.
    - Added `POST /ipDestinationGroups` Adds a new custom destination IP group.
    - Added `DELETE /ipDestinationGroups/{ipGroupId}` Deletes the destination IP group based on the specified ID. Default destination groups that are automatically created cannot be deleted.
    - Added `GET /ipGroups` Retrieves the list of IP pools.
    - Added `GET /ipGroups/lite` Retrieves the list of IP pools. This request retrieves basic information about the IP pools, such as name and ID. For extensive details, use the GET /ipGroups request.
    - Added `POST /ipGroups` Adds a new custom IP pool.
    - Added `DELETE /ipGroups/{ipGroupId}` Deletes an IP pool based on the specified ID.
    - Added `GET /networkServices` Retrieves the list of all network services. The search parameters find matching values within the name or description attributes.
    - Added `POST /networkServices` Creates a new network service.
    - Added `PUT /networkServices/{serviceId}` Updates the network service information for the specified service ID.
    - Added `DELETE /networkServices/{serviceId}` Deletes the network service for the specified ID.
    - Added `GET /networkServicesGroups` Retrieves the list of network service groups.
    - Added `GET /zpaResources/applicationSegments` Retrieves the list of ZPA application segments that can be configured in traffic forwarding rule criteria.

### ZIA Admin Role Endpoints
[PR #320](https://github.com/zscaler/zscaler-sdk-go/pull/320) - Added the following new ZIA API Endpoints:
    - Added `GET /adminRoles/{roleId}` Retrieves the admin role based on the specified ID
    - Added `GET /adminRoles/lite` Retrieves a name and ID dictionary of all admin roles. The list only includes the name and ID for all admin roles. 
    - Added `POST /adminRoles` Adds an admin role.
    - Added `PUT /adminRoles/{roleId}` Updates the admin role based on the specified ID.
    - Added `DELETE /adminRoles/{roleId}` Deletes the admin role based on the specified ID.

### Bug Fixes
[PR #320](https://github.com/zscaler/zscaler-sdk-go/pull/320) - Enhanced Updated function on ZPA `applicationsegmentpra` package to include the attribute `deleteAppsPra` in the payload during PRA application removal.
[PR #320](https://github.com/zscaler/zscaler-sdk-go/pull/320) - Fixed ZTW Cloud Connector ECGroup `ECVM` Struct.

# 3.1.14 (April 14, 2025)

## Notes
- Golang: **v1.22**

### Bug Fixes
[PR #319](https://github.com/zscaler/zscaler-sdk-go/pull/319) - Set pointer in the credential block attribute in the ZPA `policysetcontrollerv2` resource.

# 3.1.13 (March 28, 2025)

## Notes
- Golang: **v1.22**

### Bug Fixes
[PR #316](https://github.com/zscaler/zscaler-sdk-go/pull/316) - Fixed `credential` block attribute in the ZPA `policysetcontrollerv2` resource.
[PR #316](https://github.com/zscaler/zscaler-sdk-go/pull/316) - Fixed `zpa_service_edge_controller` `listen_ips` mismatched attribute type.

# 3.1.12 (March 25, 2025)

## Notes
- Golang: **v1.22**

### Bug Fixes
[PR #315](https://github.com/zscaler/zscaler-sdk-go/pull/315) - Fixed ZPA URL Encoding to support edge cases containing special characters.

[PR #315](https://github.com/zscaler/zscaler-sdk-go/pull/315) - Fixed ZCC API Client instantiation/mapping across all functions.

# 3.1.11 (March 25, 2025)

## Notes
- Golang: **v1.22**

### Bug Fixes
[PR #314](https://github.com/zscaler/zscaler-sdk-go/pull/314) - Fixed ZPA URL Encoding to support edge cases containing special characters.

[PR #314](https://github.com/zscaler/zscaler-sdk-go/pull/314) - Fixed ZCC API Client instantiation/mapping across all functions.

# 3.1.10 (March 17, 2025)

## Notes
- Golang: **v1.22**

### Bug Fixes
[PR #312](https://github.com/zscaler/zscaler-sdk-go/pull/312) - Fixed ZPA URL Encoding. The ZPA API Client now supports partial searches as supported by the API engine itself as well as dash separation names.

# 3.1.9 (March 17, 2025)

## Notes
- Golang: **v1.22**

### Bug Fixes
[PR #311](https://github.com/zscaler/zscaler-sdk-go/pull/1) - Fixed ZPA URL Encoding. The ZPA API Client now supports partial searches as supported by the API engine itself as well as dash separation names.

# 3.1.8 (March 14, 2025)

## Notes
- Golang: **v1.22**

### Bug Fixes
[PR #310](https://github.com/zscaler/zscaler-sdk-go/pull/310) - Fixed ZPA Pagination Encoding. The ZPA API Client now supports partial searches as supported by the API engine itself.

# 3.1.7 (March 5, 2025)

## Notes
- Golang: **v1.22**

### Bug Fixes
[PR #308](https://github.com/zscaler/zscaler-sdk-go/pull/308) - Fixed ZPA `customerversionprofile` resource with new attributes

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

### ZIA Policy Export
[PR #302](https://github.com/zscaler/zscaler-sdk-go/pull/302) - Implemented fix on the legacy API clients for `ZCC`, `ZIA` and `ZPA` to prevent rate limit override after client instantiation.

# 3.1.3 (February 5, 2025)

## Notes
- Golang: **v1.22**

### ZIA SSL Inspection Rules
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