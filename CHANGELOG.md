# Changelog

# 0.5.8 (January, 11 2022)

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

# 0.5.7 (January, 4 2022)

## Notes
- Golang: **v1.19**

### Enhancements

- [PR #63](https://github.com/zscaler/zscaler-sdk-go/pull/63) Added ``omitempty`` bool parameters in the ZIA URL Firewall Filtering resource ``enable_full_logging``

# 0.5.6 (January, 4 2022)

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
