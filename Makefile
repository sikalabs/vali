-include Makefile.local.mk

build:
	go build

test:
	go test -v ./... | grep -v "^\?"

release-snapshot:
	rm -rf dist
	goreleaser build --snapshot --skip before
