-include Makefile.local.mk

build:
	go build

test:
	go test -v ./... | grep -v "^\?"
