VERSION ?= $(shell cat VERSION)

test:
	gotestsum --format-hide-empty-pkg -- ./... --race

install:
	go install

tag:
	git tag $(VERSION)
	git push origin --tags
