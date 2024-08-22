VERSION := $(shell git describe --tags --always)
GIT_SHA := $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -X 'github.com/zscaler/zscaler-sdk-go/zscaler-cli/cmd.Version=$(VERSION)' \
           -X 'github.com/zscaler/zscaler-sdk-go/zscaler-cli/cmd.GitSHA=$(GIT_SHA)' \
           -X 'github.com/zscaler/zscaler-sdk-go/zscaler-cli/cmd.BuildTime=$(BUILD_TIME)'

build:
	go build -ldflags "$(LDFLAGS)" -o bin/zscaler-cli ./zscaler-cli

clean:
	rm -f bin/zscaler-cli
