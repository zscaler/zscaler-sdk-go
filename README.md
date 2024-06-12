[![release](https://github.com/zscaler/zscaler-sdk-go/actions/workflows/release.yml/badge.svg)](https://github.com/zscaler/zscaler-sdk-go/actions/workflows/release.yml)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/zscaler/zscaler-sdk-go)](https://github.com/zscaler/zscaler-sdk-go/v2/blob/master/.go-version)
[![Go Report Card](https://goreportcard.com/badge/github.com/zscaler/zscaler-sdk-go)](https://goreportcard.com/report/github.com/zscaler/zscaler-sdk-go)
[![codecov](https://codecov.io/github/zscaler/zscaler-sdk-go/graph/badge.svg?token=OVX3UWIWSK)](https://codecov.io/github/zscaler/zscaler-sdk-go)
[![License](https://img.shields.io/github/license/zscaler/zscaler-sdk-go?color=blue)](https://github.com/zscaler/zscaler-sdk-go/v2/blob/master/LICENSE)
[![Zscaler Community](https://img.shields.io/badge/zscaler-community-blue)](https://community.zscaler.com/)

## Support Disclaimer

-> **Disclaimer:** Please refer to our [General Support Statement](docs/guides/support.md) before proceeding with the use of this provider. You can also refer to our [troubleshooting guide](docs/guides/troubleshooting.md) for guidance on typical problems.

# Zscaler SDK GO

* [Need help?](#need-help)
* [Getting started](#getting-started)
* [Usage guide](#usage-guide)
* [Contributing](#contributing)

This repository contains the ZIA/ZPA/ZDX/ZCC SDK for Golang. This SDK can be
used in your server-side code to interact with the Zscaler platform

For more information about ZIA/ZPA APIs visit:

* [ZPA API](https://help.zscaler.com/zpa/zpa-api/api-developer-reference-guide)
* [ZIA API](https://help.zscaler.com/zia/getting-started-zia-api)
* [ZDX API](https://help.zscaler.com/zdx/understanding-zdx-api)
* [ZCC API](https://help.zscaler.com/client-connector/getting-started-client-connector-api)
* [ZCON API](https://help.zscaler.com/cloud-branch-connector/getting-started-cloud-branch-connector-api)

## Need help?

If you run into problems using the SDK, you can

* Refer to our [General Support Statement](/docs/guides/support.md)
* Post [issues][github-issues] here on GitHub (for code errors)

The latest release can always be found on the [releases page][github-releases].

## Getting started

To install the Zscaler GO SDK in your project:

  - Create a module file by running `go mod init`
  - You can skip this step if you already use `go mod`
  - Run `go get github.com/zscaler/zscaler-sdk-go/v2@latest`. This will add
    the SDK to your `go.mod` file.
  - Import the package in your project with `import "github.com/zscaler/zscaler-sdk-go/v2/zpa"` or
  - Before you begin, make sure you have an administrator account and API Keys in the ZIA and/or ZPA portals.
  - For more information on to create API Keys for ZIA and/or ZPA see the following help guides:

  - [Getting Started ZIA API](https://help.zscaler.com/zpa/zpa-api/api-developer-reference-guide).
  - [Getting Started ZPA API](https://help.zscaler.com/zpa/getting-started-zpa-api)
  - [Getting Started ZDX API](https://help.zscaler.com/zdx/about-zdx-api)

## Installation

To download all packages in the repo with their dependencies, simply run

`go get github.com/zscaler/zscaler-sdk-go`

## Getting Started

One can start using Zscaler Go SDK by initializing client and making a request.
Here is an example of creating a ZPA App Connector Group.

```go
package main

import (
 "log"
 "os"

 "github.com/zscaler/zscaler-sdk-go/v2/zpa"
 "github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appconnectorgroup"
)

func main() {
 /*
  If you set one of the value of the parameters to empty string, the client will fallback to:
   - The env variables: ZPA_CLIENT_ID, ZPA_CLIENT_SECRET, ZPA_CUSTOMER_ID, ZPA_CLOUD
   - Or if the env vars are not set, the client will try to use the config file which should be placed at  $HOME/.zpa/credentials.json on Linux and OS X, or "%USERPROFILE%\.zpa/credentials.json" on windows
    with the following format:
   {
    "zpa_client_id": "",
    "zpa_client_secret": "",
    "zpa_customer_id": "",
    "zpa_cloud": "https://config.private.zscaler.com"
   }
 */
 zpa_client_id := os.Getenv("ZPA_CLIENT_ID")
 zpa_client_secret := os.Getenv("ZPA_CLIENT_SECRET")
 zpa_customer_id := os.Getenv("ZPA_CUSTOMER_ID")
 zpa_cloud := os.Getenv("ZPA_CLOUD")
 config, err := zpa.NewConfig(zpa_client_id, zpa_client_secret, zpa_customer_id, zpa_cloud, "userAgent")
 if err != nil {
  log.Printf("[ERROR] creating config failed: %v\n", err)
  return
 }
 zpaClient := zpa.NewClient(config)
 appConnectorGroupService := appconnectorgroup.New(zpaClient)
 app := appconnectorgroup.AppConnectorGroup{
  Name:                   "Example app connector group",
  Description:            "Example  app connector group",
  Enabled:                true,
  CityCountry:            "California, US",
  CountryCode:            "US",
  Latitude:               "37.3382082",
  Longitude:              "-121.8863286",
  Location:               "San Jose, CA, USA",
  UpgradeDay:             "SUNDAY",
  UpgradeTimeInSecs:      "66600",
  OverrideVersionProfile: true,
  VersionProfileID:       "0",
  DNSQueryType:           "IPV4",
 }
 // Create new app connector group
 createdResource, _, err := appConnectorGroupService.Create(app)
 if err != nil {
  log.Printf("[ERROR] creating app connector group failed: %v\n", err)
  return
 }
 // Update app connector group
 createdResource.Description = "New description"
 _, err = appConnectorGroupService.Update(createdResource.ID, createdResource)
 if err != nil {
  log.Printf("[ERROR] updating app connector group failed: %v\n", err)
  return
 }
 // Delete app connector group
 _, err = appConnectorGroupService.Delete(createdResource.ID)
 if err != nil {
  log.Printf("[ERROR] deleting app connector group failed: %v\n", err)
  return
 }
}
```

!> **WARNING:** Hard-coding the **ANY** Zscaler credentials such as API Keys, client ID, and client Secrets,
works for quick tests, but for real projects you should use a more secure ways of storing these values
(such as environment variables).

## Usage guide

These examples will help you understand how to use this library.

Once you initialize a `client`, you can call methods to make requests to the
ZPA and/or ZIA APIs. Most methods are grouped by the API endpoint they belong to. For
example, methods that call the [ZPA Application Segment
API](https://help.zscaler.com/zpa/application-controller) are organized under
`Application Controller`.

## Caching

In the default configuration the ZPA and ZIA client utilizes a memory cache that has a time to live on its cached values.

See [Configuration Setter Object](#configuration-setter-object)  `WithCache(cache bool)`, `WithCacheTtl(int32`, and `WithCacheCleanWindow(i int32)`.

This helps to keep HTTP requests to the ZPA and ZIA API at a minimum. In the case where the client needs to be certain it is accessing recent data; for instance, list items, delete an item, then list items again; be sure to make use of the refresh next facility to clear the request cache. To completely disable the request
memory cache configure the client with `WithCache(false)` or set the following environment variable ``ZSCALER_SDK_CACHE_DISABLED`` to `true`.

### Configuration Setter Object

The client is configured with a configuration setter object passed to the `NewClient` function.

| function | description |
|----------|-------------|
| WithCache(cache bool) | Use request memory cache |
| WithCacheTtl(i int32) | Cache time to live in seconds |
| WithCacheCleanWindow(i int32) | Cache clean up interval in seconds

## Contributing

We're happy to accept contributions and PRs! Please see the [contribution
guide](https://github.com/zscaler/zscaler-sdk-go/blob/master/CONTRIBUTING.md) to understand how to
structure a contribution.

[github-issues]: https://github.com/zscaler/zscaler-sdk-go/issues
[github-releases]: https://github.com/zscaler/zscaler-sdk-go/releases

License
=========

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
