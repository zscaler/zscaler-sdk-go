COLOR_OK=\\x1b[0;32m
COLOR_NONE=\x1b[0m
COLOR_ERROR=\x1b[31;01m
COLOR_DESTROY=\033[31m # Red
COLOR_WARNING=\x1b[33;05m
COLOR_ZSCALER=\x1B[34;01m
GOFMT := gofumpt
GOIMPORTS := goimports

help:
	@echo "$(COLOR_ZSCALER)"
	@echo "  ______              _           "
	@echo " |___  /             | |          "
	@echo "    / / ___  ___ __ _| | ___ _ __ "
	@echo "   / / / __|/ __/ _\` | |/ _ \ '__|"
	@echo "  / /__\__ \ (_| (_| | |  __/ |   "
	@echo " /_____|___/\___\__,_|_|\___|_|   "
	@echo "                                  "
	@echo "                                  "
	@echo "$(COLOR_OK)Zscaler SDK for Golang$(COLOR_NONE)"
	@echo ""
	@echo "$(COLOR_WARNING)Usage:$(COLOR_NONE)"
	@echo "$(COLOR_OK)  make [command]$(COLOR_NONE)"
	@echo ""
	@echo "$(COLOR_WARNING)Available commands:$(COLOR_NONE)"
	@echo "$(COLOR_OK)  build                 	Clean and build the Zscaler Golang SDK generated files$(COLOR_NONE)"
	@echo "$(COLOR_WARNING)test$(COLOR_NONE)"
	@echo "$(COLOR_OK)  test:all              	Run all tests$(COLOR_NONE)"
	@echo "$(COLOR_OK)  test:integration:zcc        	Run only zcc integration tests$(COLOR_NONE)"
	@echo "$(COLOR_OK)  test:integration:ztw        	Run only ztw integration tests$(COLOR_NONE)"
	@echo "$(COLOR_OK)  test:integration:zdx        	Run only zdx integration tests$(COLOR_NONE)"
	@echo "$(COLOR_OK)  test:integration:zia        	Run only zia integration tests$(COLOR_NONE)"
	@echo "$(COLOR_OK)  test:integration:zpa        	Run only zpa integration tests$(COLOR_NONE)"
	@echo "$(COLOR_OK)  test:integration:zidentity      Run only zidentity integration tests$(COLOR_NONE)"
	@echo "$(COLOR_OK)  test:unit             			Run all unit tests$(COLOR_NONE)"
	@echo "$(COLOR_OK)  test:unit:zpa         			Run ZPA unit tests$(COLOR_NONE)"
	@echo "$(COLOR_OK)  test:unit:zia         			Run ZIA unit tests$(COLOR_NONE)"
	@echo "$(COLOR_OK)  test:unit:ztw         			Run ZTW unit tests$(COLOR_NONE)"
	@echo "$(COLOR_OK)  test:unit:zdx         			Run ZDX unit tests$(COLOR_NONE)"
	@echo "$(COLOR_OK)  test:unit:zcc         			Run ZCC unit tests$(COLOR_NONE)"
	@echo "$(COLOR_OK)  test:unit:zwa         			Run ZWA unit tests$(COLOR_NONE)"


default: build

build: vendor
	@echo "$(COLOR_ZSCALER)✓ Building SDK Source Code with Go Build...$(COLOR_NONE)"
	@go build -mod vendor -v

vendor:
	@echo "✓ Filling vendor folder with library code ..."
	@go mod vendor

fmt:
	@echo "✓ Formatting source code with goimports ..."
	@go run golang.org/x/tools/cmd/goimports@latest -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")
	@echo "✓ Formatting source code with gofmt ..."
	@gofmt -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

doc:
	@echo "Open http://localhost:6060"
	@go run golang.org/x/tools/cmd/godoc@latest -http=localhost:6060

lint: vendor
	@echo "✓ Linting source code with https://staticcheck.io/ ..."
	@go run honnef.co/go/tools/cmd/staticcheck@v0.4.0 ./...

test: lint
	@echo "✓ Running tests ..."
	@go run gotest.tools/gotestsum@latest --format pkgname-and-test-fails \
		--no-summary=skipped --raw-command go test -v \
		-json -short -coverprofile=coverage.txt ./...

coverage: test
	@echo "✓ Opening coverage for unit tests ..."
	@go tool cover -html=coverage.txt

sweep\:zpa:
	@echo "$(COLOR_WARNING)WARNING: This will destroy infrastructure. Use only in development accounts.$(COLOR_NONE)"
	ZPA_SDK_TEST_SWEEP=true go test ./zscaler/zpa/sweep -v -sweep=true

sweep\:zia:
	@echo "$(COLOR_WARNING)WARNING: This will destroy infrastructure. Use only in development accounts.$(COLOR_NONE)"
	ZIA_SDK_TEST_SWEEP=true go test ./zscaler/zia/sweep -v -sweep=true

sweep\:zidentity:
	@echo "$(COLOR_WARNING)WARNING: This will destroy infrastructure. Use only in development accounts.$(COLOR_NONE)"
	ZIDENTITY_SDK_TEST_SWEEP=true go test ./zscaler/zidentity/sweep -v -sweep=true

test\:all:
	@echo "$(COLOR_ZSCALER)Running all tests...$(COLOR_NONE)"
	@make test:integration:zcc
	@make test:integration:ztw
	@make test:integration:zdx
	@make test:integration:zpa
	@make test:integration:zia

test\:integration\:zcc:
	@echo "$(COLOR_ZSCALER)Running zcc integration tests...$(COLOR_NONE)"
	go test -v -race -cover -coverprofile=zcccoverage.out -covermode=atomic ./zscaler/zcc/... -parallel 1 -timeout 60m
	go tool cover -html=zcccoverage.out -o zcccoverage.html


test\:integration\:ztw:
	@echo "$(COLOR_ZSCALER)Running ztw integration tests...$(COLOR_NONE)"
	go test -v -race -cover -coverprofile=ztwcoverage.out -covermode=atomic ./zscaler/ztw/... -parallel 20 -timeout 60m
	go tool cover -html=ztwcoverage.out -o ztwcoverage.html

test\:integration\:zdx:
	@echo "$(COLOR_ZSCALER)Running ztw integration tests...$(COLOR_NONE)"
	go test -v -race -cover -coverprofile=zdxcoverage.out -covermode=atomic ./zscaler/zdx/... -parallel 4 -timeout 60m
	go tool cover -html=zdxcoverage.out -o zdxcoverage.html
	@go tool cover -func zdxcoverage.out | grep total:

test\:integration\:zpa:
	@echo "$(COLOR_ZSCALER)Running zpa integration tests...$(COLOR_NONE)"
	go test -v -failfast -race -cover -coverprofile=zpacoverage.out -covermode=atomic ./zscaler/zpa/... -parallel 10 -timeout 60m
	go tool cover -html=zpacoverage.out -o zpacoverage.html
	@go tool cover -func zpacoverage.out | grep total:

test\:integration\:zia:
	@echo "$(COLOR_ZSCALER)Running zia integration tests...$(COLOR_NONE)"
	go test -v -failfast -race -cover -coverprofile=ziacoverage.out -covermode=atomic ./zscaler/zia/... ./zscaler/zia/activation_cli/... -parallel 10 -timeout 60m
	go tool cover -html=ziacoverage.out -o ziacoverage.html
	@go tool cover -func ziacoverage.out | grep total:

test\:integration\:zidentity:
	@echo "$(COLOR_ZSCALER)Running zidentity integration tests...$(COLOR_NONE)"
	go test -v -failfast -race -cover -coverprofile=zidentitycoverage.out -covermode=atomic ./zscaler/zidentity/... -parallel 10 -timeout 60m
	go tool cover -html=zidentitycoverage.out -o zidentitycoverage.html
	@go tool cover -func zidentitycoverage.out | grep total:

test\:unit:
	@echo "$(COLOR_OK)Running all unit tests...$(COLOR_NONE)"
	@go test -v -race ./tests/unit/... -timeout 120s

test\:unit\:coverage:
	@echo "$(COLOR_OK)Running all unit tests with source coverage...$(COLOR_NONE)"
	@go test -v -race -cover -coverprofile=unit-coverage.out -covermode=atomic \
		-coverpkg=github.com/zscaler/zscaler-sdk-go/v3/zscaler/... \
		./tests/unit/... -timeout 180s
	@echo ""
	@echo "=== Coverage Summary ==="
	@go tool cover -func unit-coverage.out 2>/dev/null | grep -E "total:" || echo "Coverage report generated"

test\:unit\:zpa:
	@echo "$(COLOR_OK)Running ZPA unit tests with coverage...$(COLOR_NONE)"
	@go test -v -race -cover -coverprofile=unit-zpa-coverage.out -covermode=atomic ./tests/unit/zpa/... -timeout 60s
	@go tool cover -func unit-zpa-coverage.out 2>/dev/null | grep total: || echo "No coverage data"

test\:unit\:zia:
	@echo "$(COLOR_OK)Running ZIA unit tests with coverage...$(COLOR_NONE)"
	@go test -v -race -cover -coverprofile=unit-zia-coverage.out -covermode=atomic ./tests/unit/zia/... -timeout 60s
	@go tool cover -func unit-zia-coverage.out 2>/dev/null | grep total: || echo "No coverage data"

test\:unit\:ztw:
	@echo "$(COLOR_OK)Running ZTW unit tests with coverage...$(COLOR_NONE)"
	@go test -v -race -cover -coverprofile=unit-ztw-coverage.out -covermode=atomic ./tests/unit/ztw/... -timeout 60s
	@go tool cover -func unit-ztw-coverage.out 2>/dev/null | grep total: || echo "No coverage data"

test\:unit\:zdx:
	@echo "$(COLOR_OK)Running ZDX unit tests with coverage...$(COLOR_NONE)"
	@go test -v -race -cover -coverprofile=unit-zdx-coverage.out -covermode=atomic ./tests/unit/zdx/... -timeout 60s
	@go tool cover -func unit-zdx-coverage.out 2>/dev/null | grep total: || echo "No coverage data"

test\:unit\:zcc:
	@echo "$(COLOR_OK)Running ZCC unit tests with coverage...$(COLOR_NONE)"
	@go test -v -race -cover -coverprofile=unit-zcc-coverage.out -covermode=atomic ./tests/unit/zcc/... -timeout 60s
	@go tool cover -func unit-zcc-coverage.out 2>/dev/null | grep total: || echo "No coverage data"

test\:unit\:zwa:
	@echo "$(COLOR_OK)Running ZWA unit tests with coverage...$(COLOR_NONE)"
	@go test -v -race -cover -coverprofile=unit-zwa-coverage.out -covermode=atomic ./tests/unit/zwa/... -timeout 60s
	@go tool cover -func unit-zwa-coverage.out 2>/dev/null | grep total: || echo "No coverage data"

test\:unit\:oneapi:
	@echo "$(COLOR_OK)Running OneAPI unit tests with coverage...$(COLOR_NONE)"
	@go test -v -race -cover -coverprofile=unit-oneapi-coverage.out -covermode=atomic ./tests/unit/common/... -timeout 60s
	@go tool cover -func unit-oneapi-coverage.out 2>/dev/null | grep total: || echo "No coverage data"

test\:unit\:all:
	@echo "$(COLOR_OK)Running all unit tests with combined coverage...$(COLOR_NONE)"
	@go test -v -race -cover -coverprofile=unit-coverage.out -covermode=atomic ./tests/unit/... -timeout 120s
	@go tool cover -html=unit-coverage.out -o unit-coverage.html
	@go tool cover -func unit-coverage.out | grep total:

ziaActivator: GOOS=$(shell go env GOOS)
ziaActivator: GOARCH=$(shell go env GOARCH)
ifeq ($(OS),Windows_NT)  # is Windows_NT on XP, 2000, 7, Vista, 10...
ziaActivator: DESTINATION=C:\Windows\System32
else
ziaActivator: DESTINATION=/usr/local/bin
endif
ziaActivator:
	@echo "==> Installing ziaActivator cli $(DESTINATION)"
	cd ./zia/activation_cli
	go mod vendor && go mod tidy
	@mkdir -p $(DESTINATION)
	@rm -f $(DESTINATION)/ziaActivator
	@go build -o $(DESTINATION)/ziaActivator ./zscaler/zia/activation_cli/ziaActivator.go

ztwActivator: GOOS=$(shell go env GOOS)
ztwActivator: GOARCH=$(shell go env GOARCH)
ifeq ($(OS),Windows_NT)  # is Windows_NT on XP, 2000, 7, Vista, 10...
ztwActivator: DESTINATION=C:\Windows\System32
else
ztwActivator: DESTINATION=/usr/local/bin
endif
ztwActivator:
	@echo "==> Installing ztwActivator cli $(DESTINATION)"
	cd ./zscaler/ztw/services/activation_cli
	go mod vendor && go mod tidy
	@mkdir -p $(DESTINATION)
	@rm -f $(DESTINATION)/ztwActivator
	@go build -o $(DESTINATION)/ztwActivator ./zscaler/ztw/services/activation_cli/ztwActivator.go


check-fmt:
	@which $(GOFMT) > /dev/null || GO111MODULE=on go install mvdan.cc/gofumpt@latest

.PHONY: import
import: check-goimports
	@$(GOIMPORTS) -w $$(find . -path ./vendor -prune -o -name '*.go' -print) > /dev/null

check-goimports:
	@which $(GOIMPORTS) > /dev/null || GO111MODULE=on go install golang.org/x/tools/cmd/goimports@latest

.PHONY: fmt vendor fmt coverage test lint doc