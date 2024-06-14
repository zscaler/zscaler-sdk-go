[![release](https://github.com/zscaler/zscaler-sdk-go/actions/workflows/release.yml/badge.svg)](https://github.com/zscaler/zscaler-sdk-go/actions/workflows/release.yml)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/zscaler/zscaler-sdk-go)](https://github.com/zscaler/zscaler-sdk-go/v2/blob/master/.go-version)
[![Go Report Card](https://goreportcard.com/badge/github.com/zscaler/zscaler-sdk-go)](https://goreportcard.com/report/github.com/zscaler/zscaler-sdk-go)
[![codecov](https://codecov.io/github/zscaler/zscaler-sdk-go/graph/badge.svg?token=OVX3UWIWSK)](https://codecov.io/github/zscaler/zscaler-sdk-go)
[![License](https://img.shields.io/github/license/zscaler/zscaler-sdk-go?color=blue)](https://github.com/zscaler/zscaler-sdk-go/v2/blob/master/LICENSE)
[![Zscaler Community](https://img.shields.io/badge/zscaler-community-blue)](https://community.zscaler.com/)

## Support Disclaimer

-> **Disclaimer:** Please refer to our [General Support Statement](docs/guides/support.md) before proceeding with the use of this provider. You can also refer to our [troubleshooting guide](docs/guides/troubleshooting.md) for guidance on typical problems.

# Official Zscaler SDK GO Overview

- [Getting Started](#getting-started)
* [Need help?](#need-help)
* [Getting started](#getting-started)
* [Usage guide](#usage-guide)
* [Contributing](#contributing)

This repository contains the ZIA/ZPA/ZDX/ZCC/ZCON SDK for Golang. This SDK can be
used in your server-side code to interact with the Zscaler platform

Each package is supportedd by an individual and robust HTTP client designed to handle failures on different levels by performing intelligent retries.

## Getting started

The SDK is compatible with Go version 1.18.x and up. You must use [Go Modules](https://blog.golang.org/using-go-modules) to install the SDK.

To install the Zscaler GO SDK in your project:

  - Create a module file by running `go mod init`
  - You can skip this step if you already use `go mod`
  - Run `go get github.com/zscaler/zscaler-sdk-go/v2@latest`. This will add
    the SDK to your `go.mod` file.
  - Import the package in your project with `import "github.com/zscaler/zscaler-sdk-go/v2/zpa"`.

### You'll also need

*  An administrator account in whichiever one of the Zscaler products you want to interact with.
* API Keys in the the respective Zscaler cloud products.
* For more information on getting started with Zscaler APIs visit one of the following links:

* [ZPA API](https://help.zscaler.com/zpa/zpa-api/api-developer-reference-guide)
* [ZIA API](https://help.zscaler.com/zia/getting-started-zia-api)
* [ZDX API](https://help.zscaler.com/zdx/understanding-zdx-api)
* [ZCC API](https://help.zscaler.com/client-connector/getting-started-client-connector-api)
* [ZCON API](https://help.zscaler.com/cloud-branch-connector/getting-started-cloud-branch-connector-api)

## Authentication<a id="authentication"></a>

Each Zscaler product has separate developer documentation and authentication methods. In this section you will find

1. Credentials that are hard-coded into configuration arguments.

   :warning: **Caution**: Zscaler does not recommend hard-coding credentials into arguments, as they can be exposed in plain text in version control systems. Use environment variables instead.

### ZIA native authentication

- For authentication via Zscaler Internet Access, you must provide `username`, `password`, `api_key` and `cloud`

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

### ZPA native authentication

For authentication via Zscaler Private Access, you must provide `client_id`, `client_secret`, `customer_id` and `cloud`

The ZPA Cloud is identified by several cloud name prefixes, which determines which API endpoint the requests should be sent to. The following cloud environments are supported:

* `PRODUCTION`
* `ZPATWO`
* `BETA`
* `GOV`
* `GOVUS`

### Environment variables

You can provide credentials via the `ZPA_CLIENT_ID`, `ZPA_CLIENT_SECRET`, `ZPA_CUSTOMER_ID`, `ZPA_CLOUD` environment variables, representing your ZPA `client_id`, `client_secret`, `customer_id` and `cloud` of your ZPA account, respectively.

~> **NOTE** `ZPA_CLOUD` environment variable is required, and is used to identify the correct API gateway where the API requests should be forwarded to.

| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `client_id`       | _(String)_ The ZPA API client ID generated from the ZPA console.| `ZPA_CLIENT_ID` |    
| `client_secret`       | _(String)_ The ZPA API client secret generated from the ZPA console.| `ZPA_CLIENT_SECRET` |
| `customer_id`       | _(String)_ The ZPA tenant ID found in the Administration > Company menu in the ZPA console.| `ZPA_CUSTOMER_ID` |   
| `cloud`       | _(String)_ The Zscaler cloud for your tenancy.| `ZPA_CLOUD` |

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

### Environment variables

You can provide credentials via the `ZCC_CLIENT_ID`, `ZCC_CLIENT_SECRET`, `ZCC_CLOUD` environment variables, representing your ZCC `APIKey`, `SecretKey`, and `cloudEnv` of your ZCC account, respectively.

~> **NOTE** `ZCC_CLOUD` environment variable is required, and is used to identify the correct API gateway where the API requests should be forwarded to.

| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `APIKey`       | _(String)_ A string that contains the apiKey for the Mobile Portal.| `ZCC_CLIENT_ID` |    
| `SecretKey`       | _(String)_ A string that contains the secret key for the Mobile Portal.| `ZCC_CLIENT_SECRET` | 
| `cloudEnv`       | _(String)_ The host and basePath for the ZCC cloud services API is `$mobileadmin.<Zscaler Cloud Name>/papi`.| `ZCC_CLOUD` |

### ZCON native authentication

- For authentication via Zscaler Cloud Connector, you must provide `username`, `password`, `api_key` and `cloud`

The ZCON Cloud is identified by several cloud name prefixes, which determines which API endpoint the requests should be sent to. The following cloud environments are supported:

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

You can provide credentials via the `ZCON_USERNAME`, `ZCON_PASSWORD`, `ZCON_API_KEY`, `ZCON_CLOUD` environment variables, representing your ZCON `username`, `password`, `api_key` and `cloud` respectively.

| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `username`       | _(String)_ A string that contains the email ID of the API admin.| `ZCON_USERNAME` |    
| `password`       | _(String)_ A string that contains the password for the API admin.| `ZCON_PASSWORD` |
| `api_key`       | _(String)_ A string that contains the obfuscated API key (i.e., the return value of the obfuscateApiKey() method).| `ZCON_API_KEY` |   
| `cloud`       | _(String)_ The host and basePath for the cloud services API is `$connector.<Zscaler Cloud Name>/api/v1`.| `ZCON_CLOUD` |

### ZDX native authentication

For authentication via Zscaler Digital Experience (ZDX), you must provide `APIKeyID`, `SecretKey`

The ZDX Cloud is identified by several cloud name prefixes, which determines which API endpoint the requests should be sent to.

### Environment variables

You can provide credentials via the `ZDX_API_KEY_ID`, `ZDX_API_KEY_ID` environment variables, representing your ZDX `APIKey`, `SecretKey` of your ZDX account, respectively.

| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `APIKey`       | _(String)_ A string that contains the apiKey for the ZDX Portal.| `ZDX_API_KEY_ID` |    
| `SecretKey`       | _(String)_ A string that contains the secret key for the ZDX Portal.| `ZDX_API_KEY_ID` | 

## Initialize a Client

### ZPA Client Initialization
```go
import (
	"fmt"
	"context"
	"github.com/zscaler/zscaler-sdk-golang/v2/zpa"
)

func main() {
	clientID      := ""
	clientSecret  := ""
	customerID    := ""
	cloudEnv      := ""
	userAgent     := ""

	config, err := zpa.NewConfig(clientID, clientSecret, customerID, cloudEnv, userAgent)
	if err != nil {
		log.Fatalf("Error creating configuration: %v\n", err)
	}
	client := zpa.NewClient(config)
	ctx := context.TODO()

	fmt.Printf("Context: %+v\nClient: %+v\n", ctx, client)
}
```
### ZIA Client Initialization
```go
import (
	"fmt"
	"context"
	"github.com/zscaler/zscaler-sdk-golang/v2/zia"
)

func main() {
	username  := ""
	password  := ""
	apiKey    := ""
	cloudEnv  := "" 

	client, err := zia.NewClient(username, password, apiKey, cloudEnv, userAgent)
	if err != nil {
		log.Fatalf("Error creating client: %v\n", err)
	}
	ctx := context.TODO()
	fmt.Printf("Context: %+v\nClient: %+v\n", ctx, client)
}
```

### ZCC Client Initialization
```go
import (
	"fmt"
	"context"
	"github.com/zscaler/zscaler-sdk-golang/v2/zcc"
)

func main() {

	APIKey    :=  "" 
	SecretKey :=  ""
	cloudEnv  :=  ""
	userAgent :=  ""

	config, err := zcc.NewConfig(APIKey, SecretKey, cloudEnv, userAgent)
	if err != nil {
		log.Fatalf("Error creating configuration: %v\n", err)
	}
	client := zcc.NewClient(config)
	ctx := context.TODO()

	fmt.Printf("Context: %+v\nClient: %+v\n", ctx, client)
}
```
### ZDX Client Initialization
```go
import (
	"fmt"
	"context"
	"github.com/zscaler/zscaler-sdk-golang/v2/zdx"
)

func main() {

	APIKey    := ""
	SecretKey := ""
	cloudEnv  := ""

	config, err := zdx.NewConfig(APIKey, SecretKey, cloudEnv)
	if err != nil {
		log.Fatalf("Error creating configuration: %v\n", err)
	}
	client := zdx.NewClient(config)
	ctx := context.TODO()

	fmt.Printf("Context: %+v\nClient: %+v\n", ctx, client)
}
```

### ZCON Client Initialization
```go
import (
	"fmt"
	"context"
	"github.com/zscaler/zscaler-sdk-golang/v2/zcon"
)

func main() {
	username    := ""  
	password    := ""  
	apiKey      := ""    
	cloudEnv    := "" 
	userAgent   := "" 

	client, err := zcon.NewClient(username, password, apiKey, cloudEnv, userAgent)
	if err != nil {
		log.Fatalf("Error creating client: %v\n", err)
	}

	ctx := context.TODO()
	fmt.Printf("Context: %+v\nClient: %+v\n", ctx, client)
}
```

Hard-coding any of the Zscaler API credentials works for quick tests, but for real
projects you should use a more secure way of storing these values (such as
environment variables). This library supports a few different configuration
sources, covered in the [configuration reference](#configuration-reference) section.

## Usage guide

These examples will help you understand how to use this library. You can also
browse the full [API reference documentation][sdkapiref].

Once you initialize a `client`, you can call methods to make requests to the
respective Zscaler API. Most methods are grouped by the API endpoint they belong to. For
example, methods that call the ZPA [Application
Segments](https://help.zscaler.com/zpa/application-controller#/mgmtconfig/v1/admin/customers/{customerId}/application-get) are organized under
`zpa.applicationsegment`.

## Caching

In the default configuration the ZPA and ZIA client utilizes a memory cache that has a time to live on its cached values.

This helps to keep HTTP requests to the ZPA and ZIA API at a minimum. In the case where the client needs to be certain it is accessing recent data; for instance, list items, delete an item, then list items again; be sure to make use of the refresh next facility to clear the request cache. To completely disable the request
memory cache configure the client with `WithCache(false)` or set the following environment variable ``ZSCALER_SDK_CACHE_DISABLED`` to `true`.

The SDK supports caching for GET requests to improve performance and reduce the number of API calls. 
The cache can be configured and enabled/disabled using the following configuration parameters:

- `cacheEnabled`: Enables or disables caching.
- `cacheTtl`: Time-to-live for cached entries.
- `cacheCleanwindow`: Interval for cleaning expired cache entries.
- `cacheMaxSizeMB`: Maximum size of the cache in megabytes.

When a cached response is available and still valid, the SDK serves the response from the cache instead of making an API call. This behavior can be overridden by setting `freshCache` to `true`, which forces the SDK to bypass the cache and fetch a fresh response.

## Connection Retry / Rate Limiting

This SDK is designed to handle connection retries and rate limiting to ensure reliable and efficient API interactions.

### ZPA Retry Logic

By default, this SDK retries requests that receive a `429 Too Many Requests` response from the server. The retry mechanism respects the `Retry-After` header provided in the response. The `Retry-After` header indicates the time required to wait before another call can be made. For example, a value of `13s` in the `Retry-After` header means the SDK should wait 13 seconds before retrying the request.

Additionally, the SDK uses an exponential backoff strategy for other server errors, where the wait time between retries increases exponentially up to a maximum limit. This is managed by the `BackoffConfig` configuration, which specifies the following:

- `Enabled`: Set to `true` to enable the backoff and retry mechanism.
- `RetryWaitMinSeconds`: Minimum time to wait before retrying a request.
- `RetryWaitMaxSeconds`: Maximum time to wait before retrying a request.
- `MaxNumOfRetries`: Maximum number of retries for a request.

To comply with API rate limits, the SDK includes a custom rate limiter. The rate limiter ensures that requests do not exceed the following limits:

- ``GET`` requests: Maximum 20 requests in a 10-second interval.
- ``POST``, ``PUT``, ``DELETE`` requests: Maximum 10 requests in a 10-second interval.

If the request rate exceeds these limits, the SDK waits for an appropriate amount of time before proceeding with the request. The rate limiter tracks the number of requests and enforces these limits to avoid exceeding the allowed rate.

### ZIA Retry Logic

The ZIA API client in this SDK is designed to handle retries and rate limiting to ensure reliable and efficient interactions with the ZIA API.

The retry mechanism for the ZIA API client works as follows:

- The SDK retries requests that receive a `429 Too Many Requests` response or other recoverable errors.
- The primary mechanism for retries leverages the `Retry-After` header included in the response from the server. This header indicates the amount of time to wait before retrying the request. If the `Retry-After` header is present in the response body, it is also respected.
- If the `Retry-After` header is not provided, the SDK uses an exponential backoff strategy with configurable parameters:
  - `RetryWaitMinSeconds`: Minimum time to wait before retrying a request.
  - `RetryWaitMaxSeconds`: Maximum time to wait before retrying a request.
  - `MaxNumOfRetries`: Maximum number of retries for a request.
- The SDK also includes custom handling for specific error codes and messages to decide if a retry should be attempted.

## Contributing

We're happy to accept contributions and PRs! Please see the [contribution
guide](https://github.com/zscaler/zscaler-sdk-go/blob/master/CONTRIBUTING.md) to understand how to
structure a contribution.

[sdkapiref]: https://pkg.go.dev/github.com/zscaler/zscaler-sdk-go/v2
[github-issues]: https://github.com/zscaler/zscaler-sdk-go/issues
[github-releases]: https://github.com/zscaler/zscaler-sdk-go/releases

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
