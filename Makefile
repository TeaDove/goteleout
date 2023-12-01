.PHONY: ckeck install upload

PKG_VERSION ?= $(shell cat VERSION)
PKG_OUTPUT ?= goteleout
GO ?= GO111MODULE=on CGO_ENABLED=0 go
GOOS ?= $(shell $(GO) version | cut -d' ' -f4 | cut -d'/' -f1)
GOARCH ?= $(shell $(GO) version | cut -d' ' -f4 | cut -d'/' -f2)

lint:
	pre-commit run -a

test-unit:
	go test ./... -run "TestUnit_*"

test-integration:
	go test ./... -run "TestIntegration_*"

test: test-unit lint test-integration

run:
	go run main.go

install:
	go install


clean:
	@echo ">> CLEAN"
	@$(GO) clean -i ./...
	@rm -f goteleout-*-*
	@rm -rf dist/*
	@echo ">> OK"


crosscompile:
	@echo ">> CROSSCOMPILE linux/amd64"
	@GOOS=linux GOARCH=amd64 $(GO) build -o $(PKG_OUTPUT)-$(PKG_VERSION)-linux-amd64
	@echo ">> OK"
	@echo ">> CROSSCOMPILE darwin/amd64"
	@GOOS=darwin GOARCH=amd64 $(GO) build -o $(PKG_OUTPUT)-$(PKG_VERSION)-darwin-amd64
	@echo ">> OK"
	@echo ">> CROSSCOMPILE windows/amd64"
	@GOOS=windows GOARCH=amd64 $(GO) build -o $(PKG_OUTPUT)-$(PKG_VERSION)-windows-amd64
	@echo ">> OK"

	@echo ">> CROSSCOMPILE linux/arm64"
	@GOOS=linux GOARCH=arm64 $(GO) build -o $(PKG_OUTPUT)-$(PKG_VERSION)-linux-arm64
	@echo ">> OK"
	@echo ">> CROSSCOMPILE darwin/arm64"
	@GOOS=darwin GOARCH=arm64 $(GO) build -o $(PKG_OUTPUT)-$(PKG_VERSION)-darwin-arm64
	@echo ">> OK"
	@echo ">> CROSSCOMPILE windows/arm64"
	@GOOS=windows GOARCH=arm64 $(GO) build -o $(PKG_OUTPUT)-$(PKG_VERSION)-windows-arm64
	@echo ">> OK"
