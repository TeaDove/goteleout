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

tag:
	git tag $(VERSION)
	git push origin --tags
