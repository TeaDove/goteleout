.PHONY: ckeck install upload

GOTELEOUT_VERSION ?= $(shell cat VERSION)
GOTELEOUT_OUTPUT ?= goteleout
GOTELEOUT_VERSION ?= $(shell cat VERSION)
GO ?= GO111MODULE=on CGO_ENABLED=0 go
GOOS ?= $(shell $(GO) version | cut -d' ' -f4 | cut -d'/' -f1)
GOARCH ?= $(shell $(GO) version | cut -d' ' -f4 | cut -d'/' -f2)

check:
	pre-commit run -a
	go test -v ./...

install:
	go install

TAG := "${VERSION}.dev"
upload: # ONLY FOR DEV PURPOSES
	git add . || true
	git commit -m 'auto: uploading updates for tagging' || true
	git push
	git tag -d ${TAG} || true
	git tag -a ${TAG} -m "auto: development release"
	git push --delete origin ${TAG} || true
	git push origin ${TAG}

build:
	go build -o $(GOTELEOUT_OUTPUT)

clean:
	@echo -n ">> CLEAN"
	@$(GO) clean -i ./...
	@rm -f goteleout-*-*
	@rm -rf dist/*
	@printf '%s\n' '$(OK)'


crosscompile:
	@echo -n ">> CROSSCOMPILE linux/amd64"
	@GOOS=linux GOARCH=amd64 $(GO) build -o $(GOTELEOUT_OUTPUT)-$(GOTELEOUT_VERSION)-linux-amd64
	@printf '%s\n' '$(OK)'
	@echo -n ">> CROSSCOMPILE darwin/amd64"
	@GOOS=darwin GOARCH=amd64 $(GO) build -o $(GOTELEOUT_OUTPUT)-$(GOTELEOUT_VERSION)-darwin-amd64
	@printf '%s\n' '$(OK)'
	@echo -n ">> CROSSCOMPILE windows/amd64"
	@GOOS=windows GOARCH=amd64 $(GO) build -o $(GOTELEOUT_OUTPUT)-$(GOTELEOUT_VERSION)-windows-amd64
	@printf '%s\n' '$(OK)'

	@echo -n ">> CROSSCOMPILE linux/arm64"
	@GOOS=linux GOARCH=arm64 $(GO) build -o $(GOTELEOUT_OUTPUT)-$(GOTELEOUT_VERSION)-linux-arm64
	@printf '%s\n' '$(OK)'
	@echo -n ">> CROSSCOMPILE darwin/arm64"
	@GOOS=darwin GOARCH=arm64 $(GO) build -o $(GOTELEOUT_OUTPUT)-$(GOTELEOUT_VERSION)-darwin-arm64
	@printf '%s\n' '$(OK)'
	@echo -n ">> CROSSCOMPILE windows/arm64"
	@GOOS=windows GOARCH=arm64 $(GO) build -o $(GOTELEOUT_OUTPUT)-$(GOTELEOUT_VERSION)-windows-arm64
	@printf '%s\n' '$(OK)'
