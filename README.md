[![release](https://github.com/zscaler/zscaler-sdk-go/actions/workflows/release.yml/badge.svg)](https://github.com/zscaler/zscaler-sdk-go/actions/workflows/release.yml)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/zscaler/zscaler-sdk-go)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/zscaler/zscaler-sdk-go)](https://github.com/zscaler/zscaler-sdk-go/blob/master/.go-version)
[![License](https://img.shields.io/github/license/zscaler/zscaler-sdk-go?color=blue)](https://github.com/zscaler/zscaler-sdk-go/blob/master/LICENSE)

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
