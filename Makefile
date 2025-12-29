.PHONY: ckeck install upload

PKG_VERSION ?= $(shell cat VERSION)
PKG_OUTPUT ?= build/goteleout
GO ?= GO111MODULE=on CGO_ENABLED=0 go
GOOS ?= $(shell $(GO) version | cut -d' ' -f4 | cut -d'/' -f1)
GOARCH ?= $(shell $(GO) version | cut -d' ' -f4 | cut -d'/' -f2)

test:
	gotestsum --format-hide-empty-pkg -- ./... --race

install:
	go install

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

git-check-pushed:
	git status -s | xargs --null test -z

release: test git-check-pushed crosscompile
	gh release create $(PKG_VERSION) ./build/* -t="$(PKG_VERSION)" -p=false -n="new release!!!"
