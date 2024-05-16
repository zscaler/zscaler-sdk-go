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
	@echo "$(COLOR_OK)  test:integration:zcon        	Run only zcon integration tests$(COLOR_NONE)"
	@echo "$(COLOR_OK)  test:integration:zdx        	Run only zdx integration tests$(COLOR_NONE)"
	@echo "$(COLOR_OK)  test:integration:zia        	Run only zia integration tests$(COLOR_NONE)"
	@echo "$(COLOR_OK)  test:integration:zpa        	Run only zpa integration tests$(COLOR_NONE)"
	@echo "$(COLOR_OK)  test:unit             			Run only unit tests$(COLOR_NONE)"


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
	ZPA_SDK_TEST_SWEEP=true go test ./zpa/sweep -v -sweep=true

sweep\:zia:
	@echo "$(COLOR_WARNING)WARNING: This will destroy infrastructure. Use only in development accounts.$(COLOR_NONE)"
	ZIA_SDK_TEST_SWEEP=true go test ./zia/sweep -v -sweep=true

test\:all:
	@echo "$(COLOR_ZSCALER)Running all tests...$(COLOR_NONE)"
	@make test:integration:zcc
	@make test:integration:zcon
	@make test:integration:zdx
	@make test:integration:zpa
	@make test:integration:zia

test\:integration\:zcc:
	@echo "$(COLOR_ZSCALER)Running zcc integration tests...$(COLOR_NONE)"
	go test -v -race -cover -coverprofile=zcccoverage.out -covermode=atomic ./zcc/... -parallel 1 -timeout 60m
	go tool cover -html=zcccoverage.out -o zcccoverage.html


test\:integration\:zcon:
	@echo "$(COLOR_ZSCALER)Running zcon integration tests...$(COLOR_NONE)"
	go test -v -race -cover -coverprofile=zconcoverage.out -covermode=atomic ./zcon/... -parallel 20 -timeout 60m
	go tool cover -html=zconcoverage.out -o zconcoverage.html

test\:integration\:zdx:
	@echo "$(COLOR_ZSCALER)Running zcon integration tests...$(COLOR_NONE)"
	go test -v -race -cover -coverprofile=zdxcoverage.out -covermode=atomic ./zdx/... -parallel 4 -timeout 60m
	go tool cover -html=zdxcoverage.out -o zdxcoverage.html

test\:integration\:zpa:
	@echo "$(COLOR_ZSCALER)Running zpa integration tests...$(COLOR_NONE)"
	@go test -v -failfast -race -cover -coverprofile=zpacoverage.out -covermode=atomic ./zpa/... -parallel 20 -timeout 60m
	@go tool cover -html=zpacoverage.out -o zpacoverage.html

test\:integration\:zia:
	@echo "$(COLOR_ZSCALER)Running zia integration tests...$(COLOR_NONE)"
	go test -v -race -cover -coverprofile=ziacoverage.out -covermode=atomic ./zia/... -parallel 10 -timeout 60m
	go tool cover -html=ziacoverage.out -o ziacoverage.html

test\:unit:
	@echo "$(COLOR_OK)Running unit tests...$(COLOR_NONE)"
	go test -failfast -race ./tests/unit -test.v

test\:unit\zcc:
	@echo "$(COLOR_OK)Running unit tests...$(COLOR_NONE)"
	go test -failfast -race ./tests/unit/zcc -test.v

test\:unit\zcon:
	@echo "$(COLOR_OK)Running unit tests...$(COLOR_NONE)"
	go test -failfast -race ./tests/unit/zcon -test.v

test\:unit\zdx:
	@echo "$(COLOR_OK)Running unit tests...$(COLOR_NONE)"
	go test -failfast -race ./tests/unit/zdx -test.v

test\:unit\:zia:
	@echo "$(COLOR_OK)Running unit tests...$(COLOR_NONE)"
	go test -failfast -race ./tests/unit/zia -test.v

test\:unit\:zpa:
	@echo "$(COLOR_OK)Running unit tests...$(COLOR_NONE)"
	go test -failfast -race ./tests/unit/zpa -test.v

test\:unit\all:
	@echo "$(COLOR_OK)Running unit tests...$(COLOR_NONE)"
	go test -race ./tests/unit/zcc -test.v
	go test -race ./tests/unit/zcon -test.v
	go test -race ./tests/unit/zdx -test.v
	go test -race ./tests/unit/zia -test.v
	go test -race ./tests/unit/zpa -test.v

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
	@go build -o $(DESTINATION)/ziaActivator ./zia/activation_cli/ziaActivator.go
	ziaActivator

zconActivator: GOOS=$(shell go env GOOS)
zconActivator: GOARCH=$(shell go env GOARCH)
ifeq ($(OS),Windows_NT)  # is Windows_NT on XP, 2000, 7, Vista, 10...
zconActivator: DESTINATION=C:\Windows\System32
else
zconActivator: DESTINATION=/usr/local/bin
endif
zconActivator:
	@echo "==> Installing zconActivator cli $(DESTINATION)"
	cd ./zcon/services/activation_cli
	go mod vendor && go mod tidy
	@mkdir -p $(DESTINATION)
	@rm -f $(DESTINATION)/ziaActivator
	@go build -o $(DESTINATION)/zconActivator ./zcon/services/activation_cli/zconActivator.go
	zconActivator

check-fmt:
	@which $(GOFMT) > /dev/null || GO111MODULE=on go install mvdan.cc/gofumpt@latest

.PHONY: import
import: check-goimports
	@$(GOIMPORTS) -w $$(find . -path ./vendor -prune -o -name '*.go' -print) > /dev/null

check-goimports:
	@which $(GOIMPORTS) > /dev/null || GO111MODULE=on go install golang.org/x/tools/cmd/goimports@latest

.PHONY: fmt vendor fmt coverage test lint doc