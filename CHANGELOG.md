# Changelog

# 1.3.5 (April, 7 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #102](https://github.com/zscaler/zscaler-sdk-go/pull/102): Log request ID and API call duration for each request
- [PR #104](https://github.com/zscaler/zscaler-sdk-go/pull/104): Removed lock client function on ZPA client package

### Bug Fixes

- [PR #105](https://github.com/zscaler/zscaler-sdk-go/pull/105): Added function to temporarily handle ZPA upstream bad request errors.

# 1.3.4 (March, 29 2023)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #99](https://github.com/zscaler/zscaler-sdk-go/pull/99): Fixed ZIA API Client to log the user-agent information during debug

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
