[![release](https://github.com/zscaler/zscaler-sdk-go/actions/workflows/release.yml/badge.svg)](https://github.com/zscaler/zscaler-sdk-go/actions/workflows/release.yml)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/zscaler/zscaler-sdk-go)](https://github.com/zscaler/zscaler-sdk-go/blob/master/.go-version)
[![License](https://img.shields.io/github/license/zscaler/zscaler-sdk-go?color=blue)](https://github.com/zscaler/zscaler-sdk-go/blob/master/LICENSE)
[![Zscaler Community](https://img.shields.io/badge/zscaler-community-blue)](https://community.zscaler.com/)

# Zscaler GO SDK

## Aim and Scope

Zscaler GO SDK aims to access ZIA/ZPA API through HTTPS calls
from a client application purely written in Go language.

For more information about ZIA/ZPA APIs visit:

- [ZIA API](https://help.zscaler.com/zia/getting-started-zia-api).
- [ZPA API](https://help.zscaler.com/zpa/zpa-api/api-developer-reference-guide).

## Prerequisites

- The SDK is built using Go 1.18. Some features may not be
available or supported unless you have installed a relevant version of Go.
Please click [https://golang.org/dl/](https://golang.org/dl/) to download and
get more information about installing Go on your computer.

- Make sure you have properly set both `GOROOT` and `GOPATH`
environment variables.

- Before you begin, make sure you have an administrator account and API Keys in the ZIA and/or ZPA portals.

- For more information on to create API Keys for ZIA and/or ZPA see the following help guides:

  - [Getting Started ZIA API](https://help.zscaler.com/zpa/zpa-api/api-developer-reference-guide).
  - [Getting Started ZPA API](https://help.zscaler.com/zpa/getting-started-zpa-api)

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

	"github.com/zscaler/zscaler-sdk-go/zpa"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/appconnectorgroup"
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

License
=========

MIT License

=======

Copyright (c) 2022 [Zscaler BD Solutions Architect team](https://github.com/zscaler)

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
