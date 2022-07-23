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
Here is an example of getting IP List.

```go
package main

import (

)

func main() {

}
```
