[![release](https://github.com/zscaler/zscaler-sdk-go/actions/workflows/release.yml/badge.svg)](https://github.com/zscaler/zscaler-sdk-go/actions/workflows/release.yml)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/zscaler/zscaler-sdk-go)](https://github.com/zscaler/zscaler-sdk-go/v3/blob/master/.go-version)
[![Go Report Card](https://goreportcard.com/badge/github.com/zscaler/zscaler-sdk-go)](https://goreportcard.com/report/github.com/zscaler/zscaler-sdk-go)
[![codecov](https://codecov.io/gh/zscaler/zscaler-sdk-go/graph/badge.svg?token=0VX3UWIWSK)](https://codecov.io/gh/zscaler/zscaler-sdk-go)
[![License](https://img.shields.io/github/license/zscaler/zscaler-sdk-go?color=blue)](https://github.com/zscaler/zscaler-sdk-go/v3/blob/master/LICENSE)
[![Zscaler Support](https://img.shields.io/badge/zscaler-support-blue)](https://zscaler.my.site.com/customers/s/)
[![Zscaler Community](https://img.shields.io/badge/zscaler-community-blue)](https://community.zscaler.com/)

<img src="https://raw.githubusercontent.com/zscaler/zscaler-terraformer/master/images/zscaler_terraformer-logo.svg" width="400">

## Support Disclaimer

-> **Disclaimer:** Please refer to our [General Support Statement](docs/guides/support.md) before proceeding with the use of this provider. You can also refer to our [troubleshooting guide](docs/guides/troubleshooting.md) for guidance on typical problems.

# Official Zscaler SDK GO Overview

* [Release status](#release-status)
* [Need help?](#need-help)
* [Getting Started](#getting-started)
* [Authentication](#authentication)
* [OneAPI New Framework](#oneapi-new-framework)
* [Legacy API Framework](#legacy-api-framework)
* [Usage guide](#usage-guide)
* [Configuration reference](#configuration-reference)
* [Pagination](#pagination)
* [Contributing](#contributing)

This repository contains the ZIA/ZPA/ZDX/ZCC/ZTW SDK for Golang. This SDK can be
used in your server-side code to interact with the Zscaler platform

This SDK is designed to support the new Zscaler API framework [OneAPI](https://help.zscaler.com/oneapi/understanding-oneapi)
via a single OAuth2 HTTP client. The SDK is also backwards compatible with the previous
Zscaler API framework, and each package is supported by an individual and robust HTTP client
designed to handle failures on different levels by performing intelligent retries.

## Release status

This library uses semantic versioning and updates are posted in ([release notes](/docs/guides/release-notes.md)) |

| Version | Status                             |
| ------- | ---------------------------------- |
| 1.x     |  :warning: (Retired)  |
| 2.x     |  :warning: Retiring  |
| 3.x     |  :heavy_check_mark: Release ([migration guide](MIGRATING.md)) |

The latest release can always be found on the ([releases page](github-releases))

## Need help?

If you run into problems, please refer to our [General Support Statement](docs/guides/support.md) before proceeding with the use of this SDK. You can also refer to our [troubleshooting guide](docs/guides/troubleshooting.md) for guidance on typical problems. You can also raise an issue via ([github issues page](https://github.com/zscaler/zscaler-sdk-go/issues))

## Getting started

The SDK is compatible with Go version 1.18.x and up. You must use [Go Modules](https://blog.golang.org/using-go-modules) to install the SDK.

### Install current release

To install the Zscaler GO SDK in your project:

* Create a module file by running `go mod init`
* You can skip this step if you already use `go mod`
* Run `go get github.com/zscaler/zscaler-sdk-go/v3@latest`. This will add
    the SDK to your `go.mod` file.
* Import the package in your project with `import "github.com/zscaler/zscaler-sdk-go/v3/zscaler"`.

### You'll also need

* An administrator account in the Zscaler products you want to interact with.
* [OneAPI](https://help.zscaler.com/oneapi/understanding-oneapi): If using the OneAPI entrypoint you must have a API Client created in the [Zidentity platform](https://help.zscaler.com/zidentity/about-api-clients)
* Legacy Framework: If using the legacy API framework you must have API Keys credentials in the respective Zscaler cloud products.
* For more information on getting started with Zscaler APIs visit one of the following links:

* [OneAPI](https://help.zscaler.com/oneapi/understanding-oneapi)
* [ZPA API](https://help.zscaler.com/zpa/zpa-api/api-developer-reference-guide)
* [ZIA API](https://help.zscaler.com/zia/getting-started-zia-api)
* [ZDX API](https://help.zscaler.com/zdx/understanding-zdx-api)
* [ZCC API](https://help.zscaler.com/client-connector/getting-started-client-connector-api)
* [ZTW API](https://help.zscaler.com/cloud-branch-connector/getting-started-cloud-branch-connector-api)

## Authentication<a id="authentication"></a>

The latest versions => 3.x of this SDK provides dual API client capability and can be used to interact both with new Zscaler [OneAPI](https://help.zscaler.com/oneapi/understanding-oneapi) framework and the legacy API platform.

Versions of this SDK <= v2.x only support the legacy API platform. If your Zscaler tenant has not been migrated to the new Zscaler [Zidentity platform](https://help.zscaler.com/zidentity/what-zidentity).

If your organization is not ready to move into the Zidentity platform, this SDK can be configured for backwards compatibility by leveraging a built-in attribute called `use_legacy_client` or environment variable `ZSCALER_USE_LEGACY_CLIENT`.

   :warning: **Caution**: Zscaler does not recommend hard-coding credentials into arguments, as they can be exposed in plain text in version control systems. Use environment variables instead.

## OneAPI New Framework

As of the publication of SDK version => 3.x, OneAPI is available for programmatic interaction with the following products:

* [ZIA API](https://help.zscaler.com/oneapi/understanding-oneapi#:~:text=managed%20using%20OneAPI.-,ZIA%20API,-Zscaler%20Internet%20Access)
* [ZPA API](https://help.zscaler.com/oneapi/understanding-oneapi#:~:text=Workload%20Groups-,ZPA%20API,-Zscaler%20Private%20Access)
* [Zscaler Client Connector API](https://help.zscaler.com/oneapi/understanding-oneapi#:~:text=Version%20Profiles-,Zscaler%20Client%20Connector%20API,-Zscaler%20Client%20Connector)

**NOTE** Zscaler Workflow Automation (ZWA) is currently supported only via the legacy authentication method described in this README.

### OneAPI (API Client Scope)

OneAPI Resources are automatically created within the ZIdentity Admin UI based on the RBAC Roles
applicable to APIs within the various products. For example, in ZIA, navigate to `Administration -> Role Management` and select `Add API Role`.

Once this role has been saved, return to the ZIdentity Admin UI and from the Integration menu
select API Resources. Click the `View` icon to the right of Zscaler APIs and under the ZIA
dropdown you will see the newly created Role. In the event a newly created role is not seen in the
ZIdentity Admin UI a `Sync Now` button is provided in the API Resources menu which will initiate an
on-demand sync of newly created roles.

### Default Environment variables

You can provide credentials via the `ZSCALER_CLIENT_ID`, `ZSCALER_CLIENT_SECRET`, `ZSCALER_VANITY_DOMAIN`, `ZSCALER_CLOUD` environment variables, representing your Zidentity OneAPI credentials `clientId`, `clientSecret`, `vanityDomain` and `cloud` respectively.

| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `clientId`       | _(String)_ Zscaler API Client ID, used with `clientSecret` or `PrivateKey` OAuth auth mode.| `ZSCALER_CLIENT_ID` |
| `clientSecret`       | _(String)_ A string that contains the password for the API admin.| `ZSCALER_CLIENT_SECRET` |
| `privateKey`       | _(String)_ A string Private key value.| `ZSCALER_PRIVATE_KEY` |
| `vanityDomain`       | _(String)_ Refers to the domain name used by your organization `https://<vanity_domain>.zslogin.net/oauth2/v1/token` | `ZSCALER_VANITY_DOMAIN` |
| `cloud`       | _(String)_ The host and basePath for the cloud services API is `$api.<cloud_name>.zsapi.net`.| `ZSCALER_CLOUD` |

### Alternative OneAPI Cloud Environments

OneAPI supports authentication and can interact with alternative Zscaler environments i.e `beta`, `alpha` etc. To authenticate to these environments you must provide the following values:

| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `vanityDomain`       | _(String)_ Refers to the domain name used by your organization `https://<vanity_domain>.zslogin.net/oauth2/v1/token` | `ZSCALER_VANITY_DOMAIN` |
| `cloud`       | _(String)_ The host and basePath for the cloud services API is `$api.<cloud_name>.zsapi.net`.| `ZSCALER_CLOUD` |

For example: Authenticating to Zscaler Beta environment:

```sh
export ZSCALER_VANITY_DOMAIN="acme"
export ZSCALER_CLOUD="beta"
```

**Note**: By default this SDK will send the authentication request and subsequent API calls to the default base URL.

### Authenticating to Zscaler Private Access (ZPA)

The authentication to Zscaler Private Access (ZPA) via the OneAPI framework, requires the extra attribute called `customerId` and optionally the attribute `microtenantId`.

| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `clientId`       | _(String)_ Zscaler API Client ID, used with `clientSecret` or `PrivateKey` OAuth auth mode.| `ZSCALER_CLIENT_ID` |
| `clientSecret`       | _(String)_ A string that contains the password for the API admin.| `ZSCALER_CLIENT_SECRET` |
| `privateKey`       | _(String)_ A string Private key value.| `ZSCALER_PRIVATE_KEY` |
| `customerId`       | _(String)_ The ZPA tenant ID found under Configuration & Control > Public API > API Keys menu in the ZPA console.| `ZPA_CUSTOMER_ID` |
| `microtenantId`       | _(String)_ The ZPA microtenant ID found in the respective microtenant instance under Configuration & Control > Public API > API Keys menu in the ZPA console.| `ZPA_MICROTENANT_ID` |
| `vanityDomain`       | _(String)_ Refers to the domain name used by your organization `https://<vanity_domain>.zslogin.net/oauth2/v1/token` | `ZSCALER_VANITY_DOMAIN` |
| `cloud`       | _(String)_ The host and basePath for the cloud services API is `$api.<cloud_name>.zsapi.net`.| `ZSCALER_CLOUD` |

### Initialize OneAPI Client

Construct a client instance by passing it your Zidentity ClientID, ClientSecret and VanityDomain:

```go
import (
 "fmt"
 "context"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

func main() {
  config, err := zscaler.NewConfiguration(
    zscaler.WithClientID(""),
    zscaler.WithClientSecret(""),
    zscaler.WithVanityDomain("acme")
  )
  if err != nil {
    fmt.Printf("Error: %v\n", err)
  }
  service, err := zscaler.NewOneAPIClient(config)
}
```

Construct a client instance by passing it your Zidentity ClientID, PrivateKey and VanityDomain:

```go
import (
 "fmt"
 "context"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

func main() {
  config, err := zscaler.NewConfiguration(
    zscaler.WithClientID(""),
    zscaler.WithPrivateKey("private_key.pem"),
    zscaler.WithVanityDomain("acme")
  )
  if err != nil {
    fmt.Printf("Error: %v\n", err)
  }
  service, err := zscaler.NewOneAPIClient(config)
}
```

Hard-coding the Zscaler clientID and clientSecret works for quick tests, but for real
projects you should use a more secure way of storing these values (such as
environment variables). This library supports a few different configuration
sources, covered in the [configuration reference](#configuration-reference)
  section.

## Usage guide

These examples will help you understand how to use this library. You can also
browse the full [API reference documentation][sdkapiref].

Once you initialize a `client`, you can call methods to make requests to the
Zscaler API. Most methods are grouped by the API endpoint they belong to.

## Caching

In the default configuration the client utilizes a memory cache that has a time
to live on its cached values. See [Configuration Setter
Object](#configuration-setter-object)  `WithCache(cache bool)`,
`WithCacheTtl(i int32)`, and `WithCacheTti(i int32)`. This helps to
keep HTTP requests to the Zscaler API at a minimum. In the case where the client
needs to be certain it is accessing recent data; for instance, list items,
delete an item, then list items again; be sure to make use of the refresh next
facility to clear the request cache. To completely disable the request
memory cache configure the client with `WithCache(false)`.

## Connection Retry / Rate Limiting

By default, this SDK retries requests that are returned with a `429` (Too Many Requests) or `503` (Service Unavailable) response. To disable this functionality, set both `ZSCALER_CLIENT_REQUEST_TIMEOUT` and `ZSCALER_CLIENT_RATE_LIMIT_MAX_RETRIES` to `0`.

Setting only one of the values to `0` disables that specific check. For example:

- If you set `ZSCALER_CLIENT_REQUEST_TIMEOUT=45` and `ZSCALER_CLIENT_RATE_LIMIT_MAX_RETRIES=0`, the SDK will retry **indefinitely** for up to 45 seconds.
- If both are set to non-zero values, the SDK will retry until **either** condition is met (timeout or max retries).
- If both are set to `0`, no retries will occur.

### Retry Header Logic

When rate limiting is triggered, the SDK uses the following headers to determine when and how long to wait before retrying:

- `Retry-After` – preferred if available, as it provides an explicit retry duration (in seconds or `duration` format)
- `X-Ratelimit-Reset` – fallback if `Retry-After` is not present; interpreted as **relative seconds until reset**
- `X-Ratelimit-Remaining` – used for **proactive backoff** to avoid hitting the rate limit entirely

We add a 1-second buffer to all retry values to account for possible clock skew.

```go
backoff_seconds = header["X-Ratelimit-Reset"] + 1
```

If the `backoff_seconds` calculation exceeds the request timeout, the initial
429 response will be allowed through without additional attempts.

When creating your client, you can pass in these settings like you would with
any other configuration.

```go
import (
 "fmt"
 "context"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

func main() {
  config, err := zscaler.NewConfiguration(
    zscaler.WithClientID(""),
    zscaler.WithClientSecret(""),
    zscaler.WithVanityDomain("acme")
    zscaler.WithRequestTimeout(45),
    zscaler.WithRateLimitMaxRetries(3),
    zscaler.WithRateLimitRemainingThreshold(5),
  )
  if err != nil {
    fmt.Printf("Error: %v\n", err)
  }
  client := zscaler.NewOneAPIClient(config)
}
```

### Notes

- `Retry-After` and `X-Ratelimit-Reset` are interpreted as relative durations, not epoch timestamps.
- The SDK does not rely on the Date header for timing due to Zscaler’s headers being relative, not absolute.

### ZPA - List All SCIM Groups By IDP

```go
import (
 "context"
 "fmt"
 "log"

 "github.com/zscaler/zscaler-sdk-go/v3/zscaler"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/scimgroup"
)

func main() {
 config, err := zscaler.NewConfiguration(
  zscaler.WithClientID(""),
  zscaler.WithClientSecret(""),
  zscaler.WithVanityDomain("acme"),
  zscaler.WithZPACustomerID("12354547545"),
 )
 if err != nil {
  fmt.Printf("Error: %v\n", err)
 }

 service, err := zscaler.NewOneAPIClient(config)
 if err != nil {
  log.Fatalf("Error creating OneAPI client: %v", err)
 }

 ctx := context.Background()
 idp_Id := "216196257331285825"

 allGroups, resp, err := scimgroup.GetAllByIdpId(ctx, service, idp_Id)
 if err != nil {
  log.Fatalf("Error Getting Groups: %v", err)
 }
 fmt.Printf("Groups: %+v\n Response: %+v\n\n", allGroups, resp)
 for index, user := range allGroups {
  fmt.Printf("Group %d: %+v\n", index, user)
 }
}
```

### ZPA - Get a Scim Groups

```go
import (
 "context"
 "fmt"
 "log"

 "github.com/zscaler/zscaler-sdk-go/v3/zscaler"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/scimgroup"
)

func main() {
 config, err := zscaler.NewConfiguration(
  zscaler.WithClientID(""),
  zscaler.WithClientSecret(""),
  zscaler.WithVanityDomain("acme"),
  zscaler.WithZPACustomerID("12354547545"),
 )
 if err != nil {
  fmt.Printf("Error: %v\n", err)
 }

 service, err := zscaler.NewOneAPIClient(config)
 if err != nil {
  log.Fatalf("Error creating OneAPI client: %v", err)
 }

 ctx := context.Background()
 groupID := "1405240"

 allGroups, _, err := scimgroup.Get(ctx, service, groupID)
 if err != nil {
  log.Fatalf("Error retrieving scim group: %v", err)
 }
 fmt.Printf("Successfully retrieved SCIM group: %+v\n", allGroups)
}
```

### Create a ZPA Segment Group

```go
import (
 "context"
 "fmt"
 "log"

 "github.com/zscaler/zscaler-sdk-go/v3/zscaler"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
 "github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func main() {
 name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
 config, err := zscaler.NewConfiguration(
  zscaler.WithClientID(""),
  zscaler.WithClientSecret(""),
  zscaler.WithVanityDomain("acme"),
  zscaler.WithZPACustomerID("12354547545"),
 )
 if err != nil {
  fmt.Printf("Error: %v\n", err)
 }

 service, err := zscaler.NewOneAPIClient(config)
 if err != nil {
  log.Fatalf("Error creating OneAPI client: %v", err)
 }

 ctx := context.Background()
 newGroup := &segmentgroup.SegmentGroup{
  Name:        name,
  Description: name,
  Enabled:     true,
 }
 createdGroup, _, err := segmentgroup.Create(ctx, service, newGroup)
 if err != nil {
  log.Fatalf("Error creating segment group: %v", err)
 }
 fmt.Printf("Successfully created segment group: ID: %s, Name: %s\n", createdGroup.ID, createdGroup.Name)
}
```

### Update a ZPA Segment Group

```go
func main() {
 name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
 config, err := zscaler.NewConfiguration(
  zscaler.WithClientID(""),
  zscaler.WithClientSecret(""),
  zscaler.WithVanityDomain("acme"),
  zscaler.WithZPACustomerID("12354547545"),
 )
 if err != nil {
  log.Fatalf("Error creating configuration: %v", err)
 }

 service, err := zscaler.NewOneAPIClient(config)
 if err != nil {
  log.Fatalf("Error creating OneAPI client: %v", err)
 }

 ctx := context.TODO()
 groupID := "5448754152554"
 groupToUpdate, resp, err := segmentgroup.Get(ctx, service, groupID)
 if err != nil {
  log.Fatalf("Error fetching group to update: %v", err)
 }
 fmt.Printf("Group to update: %+v\n Response: %+v\n\n", groupToUpdate, resp)

 updateGroup := &segmentgroup.SegmentGroup{
  Name: name + "-updated2",
 }

 updatedGroup, err := segmentgroup.UpdateV2(ctx, service, groupToUpdate.ID, updateGroup)
 if err != nil {
  log.Fatalf("Error updating group: %v", err)
 }
 fmt.Printf("Updated Group: %+v\n Response: %+v\n\n", updatedGroup, resp)
}
```

### Create a ZIA Rule Label

```go
import (
 "context"
 "fmt"
 "log"

 "github.com/zscaler/zscaler-sdk-go/v3/zscaler"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/rule_labels"
 "github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func main() {
 name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
 config, err := zscaler.NewConfiguration(
  zscaler.WithClientID(""),
  zscaler.WithClientSecret(""),
  zscaler.WithVanityDomain("acme"),
 )
 if err != nil {
  fmt.Printf("Error: %v\n", err)
 }

 service, err := zscaler.NewOneAPIClient(config)
 if err != nil {
  log.Fatalf("Error creating OneAPI client: %v", err)
 }

 ctx := context.Background()
 newLabel := &rule_labels.RuleLabels{
  Name:        name,
  Description: name,
 }
 createdLabel, _, err := rule_labels.Create(ctx, service, newLabel)
 if err != nil {
  log.Fatalf("Error creating rule label: %v", err)
 }
 fmt.Printf("Successfully created rule label: ID: %d, Name: %s\n", createdLabel.ID, createdLabel.Name)
}
```

### Update a ZIA Rule Label

```go
import (
 "context"
 "fmt"
 "log"

 "github.com/zscaler/zscaler-sdk-go/v3/zscaler"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/rule_labels"
 "github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func main() {
 name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
 config, err := zscaler.NewConfiguration(
  zscaler.WithClientID(""),
  zscaler.WithClientSecret(""),
  zscaler.WithVanityDomain("acme"),
 )
 if err != nil {
  fmt.Printf("Error: %v\n", err)
 }

 service, err := zscaler.NewOneAPIClient(config)
 if err != nil {
  log.Fatalf("Error creating OneAPI client: %v", err)
 }

 ctx := context.TODO()
 labelID := 2073922
 labelToUpdate, err := rule_labels.Get(ctx, service, labelID)
 if err != nil {
  log.Fatalf("Error fetching label to update: %v", err)
 }
 fmt.Printf("Label to update: %+v\n", labelToUpdate)

 updateLabel := &rule_labels.RuleLabels{
  Name: name + "-updated2",
 }

 updatedLabel, resp, err := rule_labels.Update(ctx, service, labelToUpdate.ID, updateLabel)
 if err != nil {
  log.Fatalf("Error updating label: %v", err)
 }

 fmt.Printf("Updated Label: %+v\n Response: %+v\n\n", updatedLabel, resp)
}
```

### List All ZCC Devices

```go
import (
 "context"
 "fmt"
 "log"

 "github.com/zscaler/zscaler-sdk-go/v3/zscaler"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/devices"
)

func main() {
 config, err := zscaler.NewConfiguration(
  zscaler.WithClientID(""),
  zscaler.WithClientSecret(""),
  zscaler.WithVanityDomain("acme"),
 )
 if err != nil {
  fmt.Printf("Error: %v\n", err)
 }

 service, err := zscaler.NewOneAPIClient(config)
 if err != nil {
  log.Fatalf("Error creating OneAPI client: %v", err)
 }

 ctx := context.Background()
 username := ""
 osType := ""

 listDevices, err := devices.GetAll(ctx, service, username, osType)
 if err != nil {
  log.Fatalf("Error listing devices: %v", err)
 }
 fmt.Printf("Devices: %+v\n", listDevices)

 for index, user := range listDevices {
  fmt.Printf("Device %d: %+v\n", index, user)
 }
}
```

## Configuration reference

This library looks for configuration in the following sources:

0. An `zscaler.yaml` file in a `.zscaler` folder in the current user's home directory
   (`~/.zscaler/zscaler.yaml` or `%userprofile\.zscaler\zscaler.yaml`)
0. A `.zscaler.yaml` file in the application or project's root directory
0. Environment variables
0. Configuration explicitly passed to the constructor (see the example in
   [Getting started](#getting-started))

Higher numbers win. In other words, configuration passed via the constructor
will override configuration found in environment variables, which will override
configuration in `zscaler.yaml` (if any), and so on.

### YAML configuration

When you use OneAPI OAuth 2.0 the full YAML configuration looks like:

```yaml
zscaler:
  client:
    clientId: "{yourClientId}"
    clientSecret: "{yourClientSecret}"
 vanityDomain: "{yourVanityDomain}"
 customerId: "{yourZpaCustomerId}" # Required if interacting with ZPA Service
 connectionTimeout: 30 # seconds
    requestTimeout: 0 # seconds
    rateLimit:
      maxRetries: 4
    proxy:
      port: null
      host: null
      username: null
      password: null
```

When you use OneAPI OAuth 2.0 with private key, the full YAML configuration looks like:

```yaml
zscaler:
  client:
    clientId: "{yourClientId}"
    privateKey: | # THIS IS AN EXAMPLE OF A PARTIAL PRIVATE KEY. NOT USED IN PRODUCTION
        -----BEGIN RSA PRIVATE KEY-----
        MIIEogIBAAKCAQEAl4F5CrP6Wu2kKwH1Z+CNBdo0iteHhVRIXeHdeoqIB1iXvuv4
        THQdM5PIlot6XmeV1KUKuzw2ewDeb5zcasA4QHPcSVh2+KzbttPQ+RUXCUAr5t+r
        0r6gBc5Dy1IPjCFsqsPJXFwqe3RzUb...
        -----END RSA PRIVATE KEY-----
 vanityDomain: "{yourVanityDomain}"
 customerId: "{yourZpaCustomerId}" # Required if interacting with ZPA Service
 connectionTimeout: 30 # seconds
    requestTimeout: 0 # seconds
    rateLimit:
      maxRetries: 4
    proxy:
      port: null
      host: null
      username: null
      password: null
```

### Environment variables

Each one of the configuration values above can be turned into an environment
variable name with the `_` (underscore) character:

* `ZSCALER_CLIENT_ID`
* `ZSCALER_CLIENT_SECRET`
* and so on

### Configuration Setter Object

The client is configured with a configuration setter object passed to the `NewOneAPIClient` function.

| function | description |
|----------|-------------|
| WithClientID(clientId string) | OneAPI Client ID |
| WithClientSecret(clientSecret string) | OneAPI Client Secret  |
| WithPrivateKey(privateKey string) | OneAPI Private key value |
| WithVanityDomain(vanityDomain string) | The domain name used by your organization |
| WithZscalerCloud(cloud string) | The alternative Zscaler cloud name for your organization i.e `beta` |
| WithSandboxToken(sandboxToken string) | The Zscaler Internet Access Sandbox Token |
| WithSandboxCloud(sandboxCloud string) | The Zscaler Internet Access Sandbox cloud name |
| WithZPACustomerID(customerId string) | The ZPA tenant ID |
| WithZPAMicrotenantID(microtenantId string) | The ZPA Microtenant ID  |
| WithUserAgentExtra(userAgent string) | Append additional information to the HTTP User-Agent |
| WithCache(cache bool) | Use request memory cache |
| WithCacheManager(cacheManager cache.Cache) | Use custom cache object that implements the `cache.Cache` interface |
| WithCacheTtl(i int32) | Cache time to live in seconds |
| WithCacheTti(i int32) | Cache clean up interval in seconds |
| WithProxyPort(i int32) | HTTP proxy port |
| WithProxyHost(host string) | HTTP proxy host |
| WithProxyUsername(username string) | HTTP proxy username |
| WithProxyPassword(pass string) | HTTP proxy password |
| WithHttpClient(httpClient http.Client) | Custom net/http client |
| WithHttpClientPtr(httpClient *http.Client) | pointer to custom net/http client |
| WithTestingDisableHttpsCheck(httpsCheck bool) | Disable net/http SSL checks |
| WithRequestTimeout(requestTimeout int64) | HTTP request time out in seconds |
| WithRateLimitMaxRetries(maxRetries int32) | Max number of request retries when http request times out |
| WithRateLimitRemainingThreshold(retryRemainingThreshold int32) | Max number of request retries when http request times out |
| WithRateLimitMaxWait(maxWait int32) | Max wait time to wait before next retry |
| WithRateLimitMinWait(minWait int32) | Min wait time to wait before next retry |
| WithDebug(debug int32) | Enable debug mode for troubleshooting |

### Zscaler Client Base Configuration

The Zscaler Client's base configuration starts at

| config setting |
|----------------|
| WithConnectionTimeout(60) |
| WithCache(true) |
| WithCacheTtl(300) |
| WithCacheTti(300) |
| WithUserAgentExtra("") |
| WithTestingDisableHttpsCheck(false) |
| WithRequestTimeout(0) |
| WithRateLimitMaxRetries(2) |
| WithRateLimitRemainingThreshold (5) |

### Context

Every method that calls the API now has the ability to pass `context.Context`
to it as the first parameter. If you do not have a context or do not know which
context to use, you can pass `context.TODO()` to the methods.

## Legacy API Framework

### ZIA native authentication

* For authentication via Zscaler Internet Access, you must provide `username`, `password`, `api_key` and `cloud`

The ZIA Cloud is identified by several cloud name prefixes, which determines which API endpoint the requests should be sent to. The following cloud environments are supported:

* `zscaler`
* `zscalerone`
* `zscalertwo`
* `zscalerthree`
* `zscloud`
* `zscalerbeta`
* `zscalergov`
* `zscalerten`
* `zspreview`

### Environment variables

You can provide credentials via the `ZIA_USERNAME`, `ZIA_PASSWORD`, `ZIA_API_KEY`, `ZIA_CLOUD` environment variables, representing your ZIA `username`, `password`, `api_key` and `cloud` respectively.

| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `username`       | _(String)_ A string that contains the email ID of the API admin.| `ZIA_USERNAME` |
| `password`       | _(String)_ A string that contains the password for the API admin.| `ZIA_PASSWORD` |
| `api_key`       | _(String)_ A string that contains the obfuscated API key (i.e., the return value of the obfuscateApiKey() method).| `ZIA_API_KEY` |
| `cloud`       | _(String)_ The host and basePath for the cloud services API is `$zsapi.<Zscaler Cloud Name>/api/v1`.| `ZIA_CLOUD` |

### ZIA Client Initialization

```go
import (
 "context"
 "fmt"
 "log"

 "github.com/zscaler/zscaler-sdk-go/v3/zscaler"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/rule_labels"
)

func main() {
	username := os.Getenv("ZIA_USERNAME")
	password := os.Getenv("ZIA_PASSWORD")
	apiKey   := os.Getenv("ZIA_API_KEY")
	ziaCloud := os.Getenv("ZIA_CLOUD")

	ziaCfg, err := zia.NewConfiguration(
		zia.WithZiaUsername(username),
		zia.WithZiaPassword(password),
		zia.WithZiaAPIKey(apiKey),
		zia.WithZiaCloud(ziaCloud),
		zia.WithDebug(true),
	)
	if err != nil {
		log.Fatalf("Error creating ZPA configuration: %v", err)
	}

	// Initialize ZPA client
	service, err := zscaler.NewLegacyZiaClient(ziaCfg)
	if err != nil {
		log.Fatalf("Error creating ZIA client: %v", err)
	}

 // Create a new context
 ctx := context.Background()
 labels, err := rule_labels.GetAll(ctx, service)
 if err != nil {
  log.Fatalf("Error Listing Labels: %v", err)
 }
 fmt.Printf("Labels: %+v\n", labels)
 for index, label := range labels {
  fmt.Printf("Label %d: %+v\n", index, label)
 }
}

```

### ZPA native authentication

For authentication via Zscaler Private Access, you must provide `client_id`, `client_secret`, `customer_id` and `cloud`

The ZPA Cloud is identified by several cloud name prefixes, which determines which API endpoint the requests should be sent to. The following cloud environments are supported:

* `PRODUCTION`
* `ZPATWO`
* `BETA`
* `GOV`
* `GOVUS`

### ZPA Environment variables

You can provide credentials via the `ZPA_CLIENT_ID`, `ZPA_CLIENT_SECRET`, `ZPA_CUSTOMER_ID`, `ZPA_CLOUD` environment variables, representing your ZPA `client_id`, `client_secret`, `customer_id` and `cloud` of your ZPA account, respectively.

~> **NOTE** `ZPA_CLOUD` environment variable is required, and is used to identify the correct API gateway where the API requests should be forwarded to.

| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `client_id`       | _(String)_ The ZPA API client ID generated from the ZPA console.| `ZPA_CLIENT_ID` |
| `client_secret`       | _(String)_ The ZPA API client secret generated from the ZPA console.| `ZPA_CLIENT_SECRET` |
| `customer_id`       | _(String)_ The ZPA tenant ID found in the Administration > Company menu in the ZPA console.| `ZPA_CUSTOMER_ID` |
| `cloud`       | _(String)_ The Zscaler cloud for your tenancy.| `ZPA_CLOUD` |

### ZPA Client Initialization

```go
import (
 "context"
 "fmt"
 "log"

 "github.com/zscaler/zscaler-sdk-go/v3/zscaler"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
)

func main() {
	clientID := os.Getenv("ZPA_CLIENT_ID")
	clientSecret := os.Getenv("ZPA_CLIENT_SECRET")
	customerID   := os.Getenv("ZPA_CUSTOMER_ID")
	cloud := os.Getenv("ZPA_CLOUD")

	zpaCfg, err := zpa.NewConfiguration(
		zpa.WithZPAClientID(clientID),
		zpa.WithZPAClientSecret(clientSecret),
		zpa.WithZPACustomerID(customerID),
		zpa.WithZPACloud(cloud),
	)
	if err != nil {
		log.Fatalf("Error creating ZPA configuration: %v", err)
	}

	// Initialize ZPA client
	service, err := zscaler.NewLegacyZpaClient(zpaCfg)
	if err != nil {
		log.Fatalf("Error creating ZPA client: %v", err)
	}

 ctx := context.Background()

 groups, resp, err := segmentgroup.GetAll(ctx, service)
 if err != nil {
  log.Fatalf("Error list all segment group: %v", err)
 }
 fmt.Printf("Groups: %+v\n Response: %+v\n\n", groups, resp)
 for index, group := range groups {
  fmt.Printf("Group %d: %+v\n", index, group)
 }
}
```

### ZIA SCIM API

This SDK supports direct interaction with the ZIA SCIM API endpoints. The SCIM APIs allow you to use custom SCIM clients to make REST API calls to Zscaler. The same way as the regular ZIA API, all SCIM APIs are rate limited.
For more details [About SCIM APIs](https://help.zscaler.com/zia/scim-api-examples)

**NOTE**: Zscaler SCIM servers have a rate limit of 5 requests per second. In order to avoid retries, configure your application to comply with this.

The ZIA Cloud is identified by several cloud name prefixes, which determines which API endpoint the requests should be sent to. The following cloud environments are supported:

* `zscaler`
* `zscalerone`
* `zscalertwo`
* `zscalerthree`
* `zscloud`
* `zscalerbeta`
* `zscalergov`
* `zscalerten`
* `zspreview`

### ZIA SCIM API Environment variables

You can provide credentials via the `ZIA_SCIM_API_TOKEN`, `ZIA_SCIM_CLOUD`, `ZIA_SCIM_TENANT_ID` environment variables, representing your ZIA `zia_scim_api_token`, `zia_scim_cloud`, and `zia_scim_tenant_id` of your ZIA account, respectively.

~> **NOTE** `ZIA_SCIM_CLOUD` environment variable is only required. This environment variable is used to identify the correct SCIM API gateway where the API requests should be forwarded to.

| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `zia_scim_api_token`       | _(String)_ A string that contains the ZIA SCIM API Token | `ZIA_SCIM_API_TOKEN` |
| `zia_scim_cloud`       | _(String)_ A string that contains the ZIA Identity cloud environment | `ZIA_SCIM_CLOUD` |
| `zia_scim_tenant_id`       | _(String)_ A string that contains the ZIA Tenant ID i.e `24326813/61233` | `ZIA_SCIM_TENANT_ID` |

### ZIA API SCIM Client Initialization

```go
import (
 "context"
 "fmt"
 "log"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/scim_api"
)

func main() {
	scimToken := os.Getenv("ZIA_SCIM_API_TOKEN")
	scimCloud := os.Getenv("ZIA_SCIM_CLOUD")
	tenantID := os.Getenv("ZIA_SCIM_TENANT_ID")

	scimClient, err := zia.NewScimConfig(
		zia.WithScimToken(scimToken),
		zia.WithScimCloud(scimCloud),
		zia.WithTenantID(tenantID),
	)
	if err != nil {
		log.Fatalf("Failed to create SCIM client: %v", err)
	}

	service := zscaler.NewZIAScimService(scimClient) // This now properly initializes the Client

	ctx := context.Background()
	users, _, err := scim_api.GetAllUsers(ctx, service)
	if err != nil {
		log.Fatalf("Error retrieving SCIM users: %v", err)
	}
	log.Printf("Retrieved SCIM Users: %+v\n", users)
}
```

### ZPA SCIM API

This SDK supports direct interaction with the ZPA SCIM API endpoints. The SCIM APIs allow you to use custom SCIM clients to make REST API calls to Zscaler. The same way as the regular ZPA API All SCIM APIs are rate limited.
For more details [About SCIM APIs](https://help.zscaler.com/zpa/about-scim-apis)

The ZPA Cloud is identified by several cloud name prefixes, which determines which API endpoint the requests should be sent to. The following cloud environments are supported:

* `PRODUCTION`
* `ZPATWO`
* `BETA`
* `GOV`
* `GOVUS`

### ZPA SCIM API Environment variables

You can provide credentials via the `ZPA_SCIM_TOKEN`, `ZPA_IDP_ID`, `ZPA_SCIM_CLOUD` environment variables, representing your ZPA `scimToken`, `idpId`, and `baseURL` of your ZPA account, respectively.

~> **NOTE** `ZPA_SCIM_CLOUD` environment variable is only required when required when authenticating to a ZPA environment other than production. This environment variable is used to identify the correct API gateway where the API requests should be forwarded to.

| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `scimToken`       | _(String)_ A string that contains the ZPA SCIM API Token | `ZPA_SCIM_TOKEN` |
| `idpId`       | _(String)_ A string that contains the ZPA Identity Provider ID | `ZPA_IDP_ID` |
| `baseURL`       | _(String)_ A string that contains the ZPA cloud environment name | `ZPA_SCIM_CLOUD` |

### ZPA API SCIM Client Initialization

```go
import (
 "context"
 "fmt"
 "log"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/scim_api"
)

func main() {
	scimToken := os.Getenv("ZPA_SCIM_TOKEN")
	idpId := os.Getenv("ZPA_IDP_ID")
	scimCloud := os.Getenv("ZPA_SCIM_CLOUD")

	scimClient, err := zpa.NewScimConfig(
		zpa.WithScimToken(scimToken),
		zpa.WithIDPId(idpId),
		zpa.WithScimCloud(scimCloud),
	)
	if err != nil {
		log.Fatalf("failed to create SCIM client: %v", err)
	}

	service := zscaler.NewZPAScimService(scimClient)

  ctx := context.Background()
	groups, _, err := scim_api.GetAllGroups(ctx, service)
	if err != nil {
		log.Fatalf("Error retrieving SCIM groups: %v", err)
	}
	log.Printf("Retrieved SCIM Groups: %+v\n", groups)
}
```

### ZCC native authentication

For authentication via Zscaler Client Connector (Mobile Portal), you must provide `APIKey`, `SecretKey`, `cloudEnv`

The ZCC Cloud is identified by several cloud name prefixes, which determines which API endpoint the requests should be sent to. The following cloud environments are supported:

* `zscaler`
* `zscalerone`
* `zscalertwo`
* `zscalerthree`
* `zscloud`
* `zscalerbeta`
* `zscalergov`
* `zscalerten`
* `zspreview`

### ZCC Environment variables

You can provide credentials via the `ZCC_CLIENT_ID`, `ZCC_CLIENT_SECRET`, `ZCC_CLOUD` environment variables, representing your ZCC `APIKey`, `SecretKey`, and `cloudEnv` of your ZCC account, respectively.

~> **NOTE** `ZCC_CLOUD` environment variable is required, and is used to identify the correct API gateway where the API requests should be forwarded to.

| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `APIKey`       | _(String)_ A string that contains the apiKey for the Mobile Portal.| `ZCC_CLIENT_ID` |
| `SecretKey`       | _(String)_ A string that contains the secret key for the Mobile Portal.| `ZCC_CLIENT_SECRET` |
| `cloudEnv`       | _(String)_ The host and basePath for the ZCC cloud services API is `$mobileadmin.<Zscaler Cloud Name>/papi`.| `ZCC_CLOUD` |

### ZCC Client Initialization

```go
import (
 "context"
 "log"

 "github.com/zscaler/zscaler-sdk-go/v3/zscaler"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/devices"
)

func main() {
	clientID := os.Getenv("ZCC_CLIENT_ID")
	clientSecret := os.Getenv("ZCC_CLIENT_SECRET")
	cloud := os.Getenv("ZCC_CLOUD")

 zccCfg, err := zcc.NewConfiguration(
  zcc.WithZCCClientID(clientID),
  zcc.WithZCCClientSecret(clientSecret),
  zcc.WithZCCCloud(cloud),
 )
	if err != nil {
		log.Fatalf("Error creating ZCC configuration: %v", err)
	}

        // Initialize ZCC client
        service, err := zscaler.NewLegacyZccClient(zccCfg)
	if err != nil {
		log.Fatalf("Error creating ZCC client: %v", err)
	}

 ctx := context.TODO()
 username := "adam.ashcroft@acme.com"
 osType := "3"
 listDevices, err := devices.GetAll(ctx, service, username, osType)
 if err != nil {
  log.Fatalf("Error listing devices: %v", err)
 }

 for _, device := range listDevices {
  log.Printf("Device: %+v\n", device)
 }
}
```

### ZTW native authentication

* For authentication via Zscaler Client Connector, you must provide `username`, `password`, `api_key` and `cloud`

The ZTW Cloud is identified by several cloud name prefixes, which determines which API endpoint the requests should be sent to. The following cloud environments are supported:

* `zscaler`
* `zscalerone`
* `zscalertwo`
* `zscalerthree`
* `zscloud`
* `zscalerbeta`
* `zscalergov`
* `zscalerten`
* `zspreview`

### ZTW Environment variables

You can provide credentials via the `ZTW_USERNAME`, `ZTW_PASSWORD`, `ZTW_API_KEY`, `ZTW_CLOUD` environment variables, representing your ZTW `username`, `password`, `api_key` and `cloud` respectively.

| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `username`       | _(String)_ A string that contains the email ID of the API admin.| `ZTW_USERNAME` |
| `password`       | _(String)_ A string that contains the password for the API admin.| `ZTW_PASSWORD` |
| `api_key`       | _(String)_ A string that contains the obfuscated API key (i.e., the return value of the obfuscateApiKey() method).| `ZTW_API_KEY` |
| `cloud`       | _(String)_ The host and basePath for the cloud services API is `$connector.<Zscaler Cloud Name>/api/v1`.| `ZTW_CLOUD` |

**NOTE**: The Zscaler Cloud Connector (ZTW) API Client instantiation DOES NOT require the use of the `useLegacyClient` attribute.

### ZTW Client Initialization

```go
import (
	"context"
	"fmt"
	"log"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/locationmanagement/location"
)

func main() {
	username := os.Getenv("ZTW_USERNAME")
	password := os.Getenv("ZTW_PASSWORD")
	apiKey   := os.Getenv("ZTW_API_KEY")
	ztwCloud := os.Getenv("ZTW_CLOUD")

  ztwCfg, err := ztw.NewConfiguration(
		ztw.WithZtwUsername(username),
		ztw.WithZtwPassword(password),
		ztw.WithZtwAPIKey(apiKey),
		ztw.WithZtwCloud(ztwCloud),
		ztw.WithDebug(true),
	)
	if err != nil {
		log.Fatalf("Error creating ZTW configuration: %v", err)
	}

	service, err := zscaler.NewLegacyZtwClient(ztwCfg)
	if err != nil {
		log.Fatalf("Error creating ZTW client: %v", err)
	}

	ctx := context.TODO()
	locations, err := location.GetAll(ctx, service)
	if err != nil {
		log.Fatalf("Error listing locations: %v", err)
	}

	fmt.Printf("Locations: %+v\n", locations)
	for index, location := range locations {
		fmt.Printf("Location %d: %+v\n", index, location)
	}
}
```

### ZDX native authentication

For authentication via Zscaler Digital Experience (ZDX), you must provide `APIKeyID`, `SecretKey`

The ZDX Cloud is identified by several cloud name prefixes, which determines which API endpoint the requests should be sent to.

### ZDX Environment variables

You can provide credentials via the `ZDX_API_KEY_ID`, `ZDX_API_SECRET` environment variables, representing your ZDX `APIKey`, `SecretKey` of your ZDX account, respectively.

| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `APIKey`       | _(String)_ A string that contains the apiKey for the ZDX Portal.| `ZDX_API_KEY_ID` |
| `SecretKey`       | _(String)_ A string that contains the secret key for the ZDX Portal.| `ZDX_API_SECRET` |

**NOTE**: The Zscaler Digital Experience (ZDX) API Client instantiation DOES NOT require the use of the `useLegacyClient` attribute.

### ZDX Client Initialization

```go
import (
 "log"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/applications"
)

func main() {
	apiKey := os.Getenv("ZDX_API_KEY_ID")
	secretKey := os.Getenv("ZDX_API_SECRET")
	cloud := os.Getenv("ZDX_CLOUD") // Optional

 zdxCfg, err := zdx.NewConfiguration(
  zdx.WithZDXAPIKeyID(apiKey),
  zdx.WithZDXAPISecret(secretKey),
  zdx.WithZDXCloud(cloud),
  zdx.WithDebug(true),
 )
 if err != nil {
  log.Fatalf("Error creating ZDX configuration: %v", err)
 }

 zdxClient, err := zdx.NewClient(zdxCfg)
 if err != nil {
  log.Fatalf("Error creating ZDX client: %v", err)
 }

 service := services.New(zdxClient)

 filters := common.GetFromToFilters{
  From: 0, // Start time in epoch seconds (optional)
  To:   0, // End time in epoch seconds (optional)
 }

 apps, resp, err := applications.GetAllApps(service, filters)
 if err != nil {
  log.Fatalf("Error retrieving applications: %v", err)
 }

 log.Printf("Successfully retrieved %d applications.\n", len(apps))
 log.Printf("HTTP Response: %v\n", resp.Status)
 for _, app := range apps {
  log.Printf("Application Name: %s, ZDX Score: %.2f\n", app.Name, app.Score)
 }
}
```

### ZWA native authentication

For authentication via Zscaler Workflow Automation (ZWA), you must provide `key_id`, `key_secret`

The ZWA Cloud is identified by several cloud name prefixes, which determines which API endpoint the requests should be sent to.

### ZWA Environment variables

You can provide credentials via the `ZWA_API_KEY_ID`, `ZWA_API_SECRET` environment variables, representing your ZDX `key_id`, `key_secret` of your ZWA account, respectively.

| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `key_id`       | _(String)_ The ZWA string that contains the API key ID.| `ZWA_API_KEY_ID` |    
| `key_secret`       | _(String)_ The ZWA string that contains the key secret.| `ZWA_API_SECRET` |
| `cloud`       | _(String)_ The ZWA string containing cloud provisioned for your organization.| `ZWA_CLOUD` |

**NOTE**: The Zscaler Workflow Automation (ZWA) API Client instantiation DOES NOT require the use of the `useLegacyClient` attribute.

### ZWA Client Initialization

```go
import (
 "log"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zwa"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zwa/services"
 "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zwa/services/dlp_incidents"
)

func main() {
	key_id := os.Getenv("ZWA_API_KEY_ID")
	key_secret := os.Getenv("ZWA_API_SECRET")
	cloud := os.Getenv("ZWA_CLOUD") // Optional

 zwaCfg, err := zwa.NewConfiguration(
	zwa.WithZWAAPIKeyID(key_id),
	zwa.WithZWAAPISecret(key_secret),
  zwa.WithZWACloud(cloud),
	zwa.WithDebug(true),
 )
 if err != nil {
  log.Fatalf("Error creating ZWA configuration: %v", err)
 }

 zwaClient, err := zwa.NewClient(zwaCfg)
 if err != nil {
  log.Fatalf("Error creating ZWA client: %v", err)
 }

 service := services.New(zwaClient)

	dlpIncidentID := "167867439920099003"

	evidence, _, err := dlp_incidents.GetDLPIncidentTriggers(ctx, service, dlpIncidentID)
	if err != nil {
		log.Fatalf("Error fetching DLP incident evidence: %v", err)
	}

	fmt.Printf("Evidence Details: %+v\n", evidence)
}
```

Hard-coding any of the Zscaler API credentials works for quick tests, but for real
projects you should use a more secure way of storing these values (such as
environment variables). This library supports a few different configuration
sources, covered in the [configuration reference](#configuration-reference) section.

## Caching

In the default configuration the client utilizes a memory cache that has a time
to live on its cached values. See [Configuration Setter
Object](#configuration-setter-object)  `WithCache(cache bool)`,
`WithCacheTtl(i int32)`, and `WithCacheTti(i int32)`.  This helps to
keep HTTP requests to the Zscaler API at a minimum. In the case where the client
needs to be certain it is accessing recent data; for instance, list items,
delete an item, then list items again; be sure to make use of the refresh next
facility to clear the request cache. To completely disable the request
memory cache configure the client with `WithCache(false)`.

## Connection Retry / Rate Limiting

This SDK is designed to handle connection retries and rate limiting to ensure reliable and efficient API interactions.

### ZIA and ZPA Retry Logic

By default, this SDK retries requests that receive a `429 Too Many Requests` response from the server. The retry mechanism respects the `Retry-After` header provided in the response. The `Retry-After` header indicates the time required to wait before another call can be made. For example, a value of `13s` in the `Retry-After` header means the SDK should wait 13 seconds before retrying the request.

Additionally, the SDK uses an exponential backoff strategy for other server errors, where the wait time between retries increases exponentially up to a maximum limit. This is managed by the `BackoffConfig` configuration, which specifies the following:

* `Enabled`: Set to `true` to enable the backoff and retry mechanism.
* `RetryWaitMinSeconds`: Minimum time to wait before retrying a request.
* `RetryWaitMaxSeconds`: Maximum time to wait before retrying a request.
* `MaxNumOfRetries`: Maximum number of retries for a request.

To comply with API rate limits, the SDK includes a custom rate limiter. The rate limiter ensures that requests do not exceed the following limits:

* ``GET`` requests: Maximum 20 requests in a 10-second interval.
* ``POST``, ``PUT``, ``DELETE`` requests: Maximum 10 requests in a 10-second interval.

If the request rate exceeds these limits, the SDK waits for an appropriate amount of time before proceeding with the request. The rate limiter tracks the number of requests and enforces these limits to avoid exceeding the allowed rate.

### ZIA Retry Logic

The ZIA API client in this SDK is designed to handle retries and rate limiting to ensure reliable and efficient interactions with the ZIA API.

The retry mechanism for the ZIA API client works as follows:

* The SDK retries requests that receive a `429 Too Many Requests` response or other recoverable errors.
* The primary mechanism for retries leverages the `Retry-After` header included in the response from the server. This header indicates the amount of time to wait before retrying the request. If the `Retry-After` header is present in the response body, it is also respected.
* If the `Retry-After` header is not provided, the SDK uses an exponential backoff strategy with configurable parameters:
  * `RetryWaitMinSeconds`: Minimum time to wait before retrying a request.
  * `RetryWaitMaxSeconds`: Maximum time to wait before retrying a request.
  * `MaxNumOfRetries`: Maximum number of retries for a request.
* The SDK also includes custom handling for specific error codes and messages to decide if a retry should be attempted.

## Pagination

Each Zscaler service provides pagination support. The SDK independently provides pagination logic for each individual service based on its unique parameters and requirements.

### ZPA Pagination
The ZPA package robust support for pagination, allowing users to fetch all results for a specific API endpoint even if the data spans multiple pages. The SDK abstracts pagination, so you can fetch all records seamlessly without worrying about page-by-page API requests.

The SDK includes a generic function, `GetAllPagesGenericWithCustomFilters`, which automates fetching all resources across multiple pages. Here's an example of how you can use this functionality to list SCIM groups by an IDP ID.

- *Custom Filters*: Use the Filter struct to refine results with parameters such as `Search`, `SortBy`, and `SortOrder`.

#### Example Usage in a Program

```go
func main() {

ctx := context.Background()
idpId := "your-idp-id"

allGroups, resp, err := scimgroup.GetAllByIdpId(ctx, service, idpId)
if err != nil {
    log.Fatalf("Error fetching SCIM groups: %v", err)
}

fmt.Printf("Fetched %d SCIM Groups\n", len(allGroups))
for _, group := range allGroups {
    fmt.Printf("Group: %+v\n", group)
  }
}
```

#### Example - Direct Use of Pagination Functions

The primary pagination functions, such as `GetAllPagesGenericWithCustomFilters` and `GetAllPagesGeneric`, can be directly invoked to handle paginated API requests. Here’s an example of how a user could utilize these functions:
- *Custom Filters*: Use the Filter struct to refine results with parameters such as `Search`, `SortBy`, and `SortOrder`.

```go
import (
	"context"
	"fmt"
	"log"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/scimgroup"
)

func main() {
    // Initialize ZPA configuration
    config, err := zscaler.NewConfiguration(
        zscaler.WithClientID("your-client-id"),
        zscaler.WithClientSecret("your-client-secret"),
        zscaler.WithVanityDomain("your-vanity-domain"),
        zscaler.WithZPACustomerID("your-customer-id"),
        zscaler.WithDebug(true),
    )
    if err != nil {
        log.Fatalf("Error creating configuration: %v", err)
    }

    service, err := zscaler.NewOneAPIClient(config)
    if err != nil {
        log.Fatalf("Error creating ZPA client: %v", err)
    }
    ctx := context.Background()
    relativeURL := "/zpa/mgmtconfig/v1/admin/customers/21619xxxxxxxxxx/scimgroups/idpId/{idpId}"

    // Define filters for pagination
    filters := common.Filter{
        SortBy:    "name",
        SortOrder: "ASC",
        Search:    "example-query",
    }

    // Use the generic pagination function
    allGroups, resp, err := common.GetAllPagesGenericWithCustomFilters[ScimGroup](ctx, service.Client, relativeURL, filters)
    if err != nil {
        log.Fatalf("Error fetching paginated data: %v", err)
    }

    fmt.Printf("Fetched %d SCIM Groups\n", len(allGroups))
    for index, group := range allGroups {
        fmt.Printf("Group %d: %+v\n", index, group)
    }
    fmt.Printf("Response: %+v\n", resp)
}
```

### ZIA Pagination
The ZIA SDK provides pagination support tailored to its API's unique parameters. The SDK allows you to fetch large datasets across multiple pages seamlessly, using built-in utilities or customizable pagination logic.
Pagination in the ZIA service is powered by the `ReadAllPages` and `ReadPage` functions. These utilities enable efficient data retrieval while handling all necessary API parameters for pagination, including page size and sort options.

- **Using Built-In Pagination Functions**: The `ReadAllPages` function automates fetching all pages of data for a given endpoint and aggregates the results into a single slice. For example, the GetAllUsers function leverages `ReadAllPages` to retrieve all users.

```go
func GetAllUsers(ctx context.Context, service *zscaler.Service) ([]Users, error) {
	var users []Users
	err := common.ReadAllPages(ctx, service.Client, usersEndpoint+"?"+common.GetSortParams(common.SortField(service.SortBy), common.SortOrder(service.SortOrder)), &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
```

#### Example Usage in a Program

```go
func main() {
ctx := context.Background()

// Fetch all users
allUsers, err := GetAllUsers(ctx, service)
if err != nil {
    log.Fatalf("Error fetching users: %v", err)
}

fmt.Printf("Fetched %d users\n", len(allUsers))
for _, user := range allUsers {
    fmt.Printf("User: %+v\n", user)
  }
}
```

## Contributing

We're happy to accept contributions and PRs! Please see the [contribution
guide](https://github.com/zscaler/zscaler-sdk-go/blob/master/CONTRIBUTING.md) to understand how to
structure a contribution.

[sdkapiref]: https://pkg.go.dev/github.com/zscaler/zscaler-sdk-go/v3

MIT License
=======

Copyright (c) 2022 [Zscaler](https://github.com/zscaler)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
